package nodes

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"hash/crc64"
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
	address  []string
}

type NodesRequest struct {
	AddrFrom string
	nodes    []NetworkNode
}

type RequestMessage struct {
	AddrFrom string
}

type AckMessage struct {
	AddrFrom string
}

type NetworkNode struct {
	NodeID uint64
	Name   string
	// TCP Addr is akin to socket. So, its only used when its listening for connection, right?
	Socket *net.TCPAddr
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

func sendData(node *NetworkNode, data []byte, net_platform *NetworkPlatform) {
	conn, err := net.Dial(node.Socket.Network(), node.Socket.String())

	if err != nil {
		log.Printf("Connection Failed, for node %s", node.Name)
		net_platform.RemoveNode(*node)
		return
	}
	defer conn.Close()

	_, err = conn.Write(data)
	if err != nil {
		log.Printf("Sending data failed, error: %s", err.Error())
	}
}

func SendGetNode(addr string, net_platform *NetworkPlatform) {
	conn, err := net.Dial("tcp", addr)

	address := make([]string, 1)
	address[0] = addr
	payload := GetNodes{
		AddrFrom: net_platform.Self_node.Socket.String(),
		address:  address,
	}
	data := GobEncode(payload)
	data = append(CommandToBytes("getnodes"), data...)
	if err != nil {
		log.Fatalf("Sending data failed, error: %s", err.Error())
	}

	_, err = conn.Write(data)
}

// sends all the node address
func HandleAddr(request []byte) {

}

func HandleUnknownCommand() {

}

func HandleGetNodes(request []byte, net_platform *NetworkPlatform) {
	var payload GetNodes
	gob.NewDecoder(bytes.NewBuffer(request)).Decode(&payload)

	if payload.address[0] == net_platform.Self_node.Socket.Network() {
		send_payload := NodesRequest{
			AddrFrom: net_platform.GetNodeAddress(),
			nodes:    []NetworkNode{*net_platform.Self_node},
		}

		GobEncode(send_payload)
		// sendData(payload.AddrFrom, GobEncode(send_payload), net_platform)
	}
}

func HandleGetNode(request []byte) {

}

func HandleTCPConnection(tcp_connection *net.TCPConn, net_platform *NetworkPlatform) {
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

	request, err := ioutil.ReadAll(tcp_connection)

	defer tcp_connection.Close()
	if err != nil {
		log.Printf(err.Error())
	}

	// first 32 bytes to hold the commnd
	// TODO: Format the header data
	command := string(request[:COMMAND_LENGTH])
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
	listener, err := net.ListenTCP(net_platform.Self_node.Socket.Network(), net_platform.Self_node.Socket)

	if err != nil {
		log.Fatalf("Listener error: %s", err.Error())
	}

	if net_platform.Self_node.Socket.Port == 6969 {
		SendGetNode("127.0.0.1:7000", net_platform)
	}

	// The call to listen always blocks
	// There's no way to get notified when there is a pending connection in Go?
	// go ProcessConnections(net_platform)
	for {
		conn, _ := listener.AcceptTCP()
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
