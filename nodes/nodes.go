package nodes

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"hash/crc64"
	"io"
	"io/ioutil"
	"log"
	"net"
	"time"
)

const (
	COMMAND_LENGTH = 32
)

// var net_platform NetworkPlatform

// message for get address request
// request for all the known address
// maybe we could handle nodes directly???
type GetAddress struct {
	AddrFrom   string
	message_id uint64
}

type GetNodes struct {
	AddrFrom string
	Address  []string
}

type NodesRequest struct {
	AddrFrom string
	Nodes    []NetworkNode
}

type RequestMessage struct {
	AddrFrom string
}

type AckMessage struct {
	AddrFrom string
}

// Onto nodes discovery
// How to decide if networks are in sync? ans -> After certain time lol
// Should this function be called on regular basis? On certain interval or not?
func SyncWithNetwork(net_platform *NetworkPlatform) uint16 {
	// Receive information about connected nodes from its neighbor nodes
	msg := "Hello there"
	var discovered_nodes uint16

	for _, cached := range net_platform.Connection_caches {

		// This should be bidirectional
		sendData(cached.node_ref, []byte(msg), net_platform)

		// TODO :: Send and receive msg and interpret it
		// Wait for it them to recieve message and compare to them
		// Return the IP Address and port number of other nodes which are listening for p2p connection
		// Read the message and identify new nodes in the network
		discovered_nodes++
	}
	return discovered_nodes
}

// For connecting to the network, at least one node need to be known
func ConnectToNetwork(node *NetworkNode, net_platform *NetworkPlatform) bool {
	// Connect as a client to the network
	// Maybe implement something like OSPF routing algorithm to create map of the network ??
	tcp_connection := IntiateTCPConnection(node)
	if tcp_connection == nil {
		return false
	}

	// TODO :: Perform other necessary actions to get in sync with the network
	entry := CreateCacheEntry(tcp_connection, node, node.NodeID)
	net_platform.Connection_caches = append(net_platform.Connection_caches, entry)

	SyncWithNetwork(net_platform)
	return true
}

func CreateNetworkNode(name string, address string, port int) *NetworkNode {
	networkNode := &NetworkNode{}
	networkNode.Name = name

	//TODO: Implement hashing
	id := fmt.Sprintf("%s %s %d", name, address, port)
	table := crc64.MakeTable(100)
	networkNode.Socket, _ = net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", address, port))
	networkNode.NodeID = crc64.Checksum([]byte(id), table)
	return networkNode
}

func CommandToBytes(cmd string) []byte {
	var bytes [COMMAND_LENGTH]byte

	for i, c := range cmd {
		bytes[i] = byte(c)
	}

	return bytes[:]
}

func BytesToCommand(bytes []byte) string {
	var cmd []byte

	for _, b := range bytes {
		if b != 0x0 {
			cmd = append(cmd, b)
		}
	}

	return fmt.Sprintf("%s", cmd)
}

func sendData(node *NetworkNode, data []byte, net_platform *NetworkPlatform) {
	conn, err := net.Dial(node.Socket.Network(), node.Socket.String())

	if err != nil {
		log.Printf("Connection Failed, for node %s\n", node.Name)
		net_platform.RemoveNode(*node)
		return
	}
	defer conn.Close()

	_, err = conn.Write(data)
	if err != nil {
		log.Printf("Sending data failed, error: %s\n", err.Error())
	}
}

func sendDataToAddress(addr string, data []byte, net_platform *NetworkPlatform) {
	conn, err := net.Dial("tcp", addr)

	if err != nil {
		log.Printf("Connection Failed, for node with address: %s\n", addr)
		//TODO: implement remove node for address too
		// Difficulty: Easy
		// net_platform.RemoveNode(*node)
		return
	}
	defer conn.Close()

	_, err = io.Copy(conn, bytes.NewReader(data))

	if err != nil {
		log.Panic(err)
	}
}

// send the get node request to a particular node
func SendGetNode(addr string, net_platform *NetworkPlatform) {
	payload := GetNodes{
		AddrFrom: net_platform.Self_node.Socket.String(),
		Address: []string{
			addr,
		},
	}
	data := GobEncode(payload)
	data = append(CommandToBytes("getnodes"), data...)
	sendDataToAddress(addr, data, net_platform)
}

// sends all the node address
func HandleAddr(request []byte) {

}

func HandleUnknownCommand() {

}

func HandleGetNodes(request []byte, net_platform *NetworkPlatform) {
	var payload GetNodes
	gob.NewDecoder(bytes.NewBuffer(request)).Decode(&payload)

	fmt.Printf("Sending nodes to address: %s\n", payload.AddrFrom)
	if payload.Address[0] == net_platform.Self_node.Socket.String() {
		// if the receiving address is the self address, then it is send
		send_payload := NodesRequest{
			AddrFrom: net_platform.GetNodeAddress(),
			Nodes:    []NetworkNode{*net_platform.Self_node},
		}

		sendDataToAddress(payload.AddrFrom, append(CommandToBytes("node"), GobEncode(send_payload)...), net_platform)

		// if the address is not known for this node, it is fetched
		if net_platform.knows(payload.AddrFrom) {
			SendGetNode(payload.AddrFrom, net_platform)
		}
	}
}

func HandleNode(request []byte, net_platform *NetworkPlatform) {
	var payload NodesRequest
	gob.NewDecoder(bytes.NewBuffer(request)).Decode(&payload)
	fmt.Printf("Address received from %s\n", payload.AddrFrom)
	net_platform.Connected_nodes = append(net_platform.Connected_nodes, payload.Nodes...)
}

func HandleTCPConnection(conn net.Conn, net_platform *NetworkPlatform) {
	// store the information about the newly connected node into the net_platform struct
	// So connection established, now retrieve information about the host
	// TODO :: Test this implementation, left for Go experts
	// Assuming that Garbage collected language can handle anything, literally anything
	// Like some memory allocated by another runtime too.. lol
	// new_node := CreateNetworkNode("unknown", "127.0.0.1", 8000)
	// net_platform.Connected_nodes = append(net_platform.Connected_nodes, *new_node)

	// // Operation on connection caches are omitted for now
	// cache_entry := CreateCacheEntry(tcp_connection, nil, new_node.NodeID)
	// net_platform.Connection_caches = append(net_platform.Connection_caches, cache_entry)

	request, err := ioutil.ReadAll(conn)

	// defer tcp_connection.Close()
	if err != nil {
		log.Printf(err.Error())
	}
	defer conn.Close()

	// first 32 bytes to hold the commnd
	// TODO: Format the header data
	command := BytesToCommand(request[:COMMAND_LENGTH])
	log.Printf("Command: %s", command)

	switch command {
	default:
		HandleUnknownCommand()
		break

	case "addr":
		HandleAddr(request)
		break

	case "getnodes":
		HandleGetNodes(request[COMMAND_LENGTH:], net_platform)
		break

	// for receiving a node
	case "node":
		HandleNode(request[COMMAND_LENGTH:], net_platform)
		break
	}
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

func ProcessConnections(net_platform *NetworkPlatform) {
	// Process other currently running connections
	for {
		time.Sleep(1000 * time.Millisecond)
		for _, caches := range net_platform.Connection_caches {
			// If there's message to be sent to that node send it here.
			// Else, receive message here

			msg := make([]byte, 2048)
			// Echo back the same message to client
			sendData(caches.node_ref, msg, net_platform)
		}
	}
}

// For self
func ListenForTCPConnection(net_platform *NetworkPlatform) {
	listener, err := net.Listen("tcp", net_platform.Self_node.Socket.String())
	if err != nil {
		log.Fatalf("Listener error: %s\n", err.Error())
	}
	defer listener.Close()

	if net_platform.Self_node.Socket.Port == 6969 {
		log.Printf("Sending get nodes request")
		SendGetNode("127.0.0.1:7000", net_platform)
	}

	// The call to listen always blocks
	// There's no way to get notified when there is a pending connection in Go?
	// go ProcessConnections(net_platform)
	for {
		conn, _ := listener.Accept()
		if err != nil {
			fmt.Printf("Failed to Accept the incoming connection.  Error: %s\n", err.Error())
			break
		}
		go HandleTCPConnection(conn, net_platform)
	}
}

func IntiateTCPConnection(node *NetworkNode) *net.TCPConn {
	tcp_con, err := net.DialTCP("tcp", nil, node.Socket)
	if err != nil {
		fmt.Println("Failed to initiate tcp connection with the host : ", node)
		return nil
	}
	return tcp_con
}

// To be implemented later on
