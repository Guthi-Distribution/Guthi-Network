package platform

/*
TODO:
- Check if the node is alive or not
- Implement checkpointing
*/
import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"hash/crc64"
	"io"
	"io/ioutil"
	"log"
	"net"
	"time"
)

const (
	COMMAND_LENGTH = 16
)

// Onto nodes discovery
// How to decide if networks are in sync? ans -> After certain time lol
// Should this function be called on regular basis? On certain interval or not?
func SyncWithNetwork(net_platform *NetworkPlatform) uint16 {
	// Receive information about connected nodes from its neighbor nodes
	msg := "Hello there"
	var discovered_nodes uint16

	for _, cached := range net_platform.Connection_caches {

		// This should be bidirectional
		sendDataToNode(cached.node_ref, []byte(msg), net_platform)

		// TODO :: Send and receive msg and interpret it
		// Wait for it them to recieve message and compare to them
		// Return the IP Address and port number of other nodes which are listening for p2p connection
		// Read the message and identify new nodes in the network
		discovered_nodes++
	}
	return discovered_nodes
}

func CreateNetworkNode(name string, address string, port int) (*NetworkNode, error) {
	networkNode := &NetworkNode{}
	networkNode.Name = name

	//TODO: Implement hashing
	id := fmt.Sprintf("%s %s %d", name, address, port)

	// what does this table thing do?
	table := crc64.MakeTable(100)

	var err error
	networkNode.Socket, err = net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", address, port))
	if err != nil {
		return nil, err
	}
	networkNode.NodeID = crc64.Checksum([]byte(id), table)
	return networkNode, nil
}

func sendDataToNode(node *NetworkNode, data []byte, net_platform *NetworkPlatform) {
	// connect to a network
	conn, err := net.Dial(node.Socket.Network(), node.Socket.String())

	if err != nil {
		log.Printf("Connection Failed, for node %s\n", node.Name)
		net_platform.RemoveNode(*node)
		return
	}
	defer conn.Close()

	_, err = io.Copy(conn, bytes.NewReader(data))
	if err != nil {
		log.Printf("Sending data failed, error: %s\n", err.Error())
	}
}

func sendDataToAddress(addr string, data []byte, net_platform *NetworkPlatform) error {
	conn, err := net.Dial("tcp", addr)

	if err != nil {
		log.Printf("Connection Failed, for node with address: %s\n", addr)
		net_platform.RemoveNodeWithAddress(addr)
		return err
	}
	defer conn.Close()

	_, err = io.Copy(conn, bytes.NewReader(data)) // write into connection i.e send data

	if err != nil {
		return err
	}

	return err
}

// send the get node request to a particular node
func SendGetNode(addr string, net_platform *NetworkPlatform) error {
	payload := GetNodes{
		AddrFrom: net_platform.Self_node.Socket.String(),
		Address: []string{
			addr,
		},
	}
	data := GobEncode(payload)
	data = append(CommandStringToBytes("getnodes"), data...)
	return sendDataToAddress(addr, data, net_platform)
}

func SendEcho(addr string, net_platform *NetworkPlatform) {
	payload := "Hello World!"
	data := GobEncode(payload)
	data = append(CommandStringToBytes("echo"), data...)
	sendDataToAddress(addr, data, net_platform)
}

func HandleUnknownCommand() {

}

func ReplyBack(msg []byte, conn net.Conn, net_platform *NetworkPlatform) {
	wr, err := conn.Write(msg)
	if err != nil {
		fmt.Println("Failed to reply back to connection", conn)
	}
	fmt.Printf("Bytes written back %d.\n", wr)
}

func HandleResources(msg []byte, conn net.Conn, net_platform *NetworkPlatform) {
	// Files are handled by the C++ runtime
	//  Interface with C++
	log.Println("Resource requested ", msg)
}

func HandleGetNodes(request []byte, net_platform *NetworkPlatform) {
	var payload GetNodes
	gob.NewDecoder(bytes.NewBuffer(request)).Decode(&payload)

	fmt.Printf("Sending nodes to address: %s\n", payload.AddrFrom)
	if payload.Address[0] == net_platform.Self_node.Socket.String() {
		// if the receiving address is the self address, then it is send
		send_payload := NodesMessage{
			AddrFrom: net_platform.GetNodeAddress(),
			Nodes:    []NetworkNode{*net_platform.Self_node},
		}

		sendDataToAddress(payload.AddrFrom, append(CommandStringToBytes("node"), GobEncode(send_payload)...), net_platform)
	}
}

func HandleNode(request []byte, net_platform *NetworkPlatform) {
	var payload NodesMessage
	gob.NewDecoder(bytes.NewBuffer(request)).Decode(&payload)
	fmt.Printf("Address received from %s\n", payload.AddrFrom)
	net_platform.Connected_nodes = append(net_platform.Connected_nodes, payload.Nodes...)

	if len(payload.Nodes) == 0 {
		log.Printf("Nodes received length is zero")
		return
	}
	entry := CreateCacheEntry(nil, &payload.Nodes[0], payload.Nodes[0].NodeID)
	net_platform.Connection_caches = append(net_platform.Connection_caches, entry)

	SyncWithNetwork(net_platform)
}

// TODO: Should read be handled concurrently via go routines?
func HandleTCPConnection(conn net.Conn, net_platform *NetworkPlatform) error {
	// request, err := io.ReadAll(conn)
	request, err := ioutil.ReadAll(conn)
	conn.SetReadDeadline(time.Now().Add(100 * time.Millisecond))

	if err != nil {
		// Close the connection
		if errors.Is(err, net.ErrClosed) {
			log.Printf("Connection closed by the peer")
			return err
		}
		return err
	}

	// first 32 bytes to hold the command
	// TODO: Format the header data
	command := BytesToCommandString(request[:COMMAND_LENGTH])
	log.Printf("Command: %s", command)

	switch command {
	default:
		HandleUnknownCommand()
		break

	case "getnodes":
		HandleGetNodes(request[COMMAND_LENGTH:], net_platform)
		break

	// for receiving a node
	case "node":
		HandleNode(request[COMMAND_LENGTH:], net_platform)
		break

	case "echo":
		ReplyBack(request[COMMAND_LENGTH:], conn, net_platform)
		break

	case "connect":
		HandleConnectionInitiation(request[COMMAND_LENGTH:], net_platform)
		break

	case "connection_reply":
		HandleConnectionReply(request[COMMAND_LENGTH:], net_platform)
		break

	case "getresources":
		HandleResources(request[COMMAND_LENGTH:], conn, net_platform)
		break
	}
	return nil
}

// Gob Encode
// Details: https://pkg.go.dev/encoding/gob
func GobEncode(data interface{}) []byte {
	var buff bytes.Buffer

	// the encoded data is stored in buff and the data to be encoded is `data`
	err := gob.NewEncoder(&buff).Encode(data)

	if err != nil {
		log.Panic(err)
	}

	return buff.Bytes()
}

// For self
func ListenForTCPConnection(net_platform *NetworkPlatform) {
	// create a listener that is used to listen to other connection (somethong like that)
	// Listen announces on the local network address. @docs
	listener, err := net.Listen("tcp", net_platform.Self_node.Socket.String())
	if err != nil {
		log.Fatalf("Listener error: %s\n", err.Error())
	}
	defer listener.Close()

	// The call to listen always blocks
	// There's no way to get notified when there is a pending connection in Go?
	log.Printf("Localhost is listening ... \n")
	for {
		conn, _ := listener.Accept()
		if err != nil {
			fmt.Printf("Failed to Accept the incoming connection.  Error: %s\n", err.Error())
			break
		}
		go HandleTCPConnection(conn, net_platform)
	}
}

// To be implemented later on
