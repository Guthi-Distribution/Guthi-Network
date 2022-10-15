package nodes

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"hash/crc64"
	"io"
	"log"
	"net"
	"time"
)

const (
	COMMAND_LENGTH = 16
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

// if you want nodes to hold additional information, use this
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
		sendDataToNode(cached.node_ref, []byte(msg), net_platform)

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

	// what does this table thing do?
	table := crc64.MakeTable(100)
	networkNode.Socket, _ = net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", address, port))
	networkNode.NodeID = crc64.Checksum([]byte(id), table)

	return networkNode
}

func CommandStringToBytes(cmd string) []byte {
	var bytes [COMMAND_LENGTH]byte

	for i, c := range cmd {
		bytes[i] = byte(c)
	}

	return bytes[:]
}

func BytesToCommandString(bytes []byte) string {
	var cmd []byte

	for _, b := range bytes {
		if b != 0x0 {
			cmd = append(cmd, b)
		}
	}

	return fmt.Sprintf("%s", cmd)
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

	_, err = conn.Write(data)
	if err != nil {
		log.Printf("Sending data failed, error: %s\n", err.Error())
	}
}

func sendDataToAddress(addr string, data []byte, net_platform *NetworkPlatform) {
	conn, err := net.Dial("tcp", addr)

	if err != nil {
		log.Printf("Connection Failed, for node with address: %s\n", addr)
		net_platform.RemoveNodeWithAddress(addr)
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
	data = append(CommandStringToBytes("getnodes"), data...)
	sendDataToAddress(addr, data, net_platform)
}

func SendEcho(addr string, net_platform *NetworkPlatform) {
	payload := "Hello World!"
	data := GobEncode(payload)
	data = append(CommandStringToBytes("echo"), data...)
	sendDataToAddress(addr, data, net_platform)
}

// sends all the node address
func HandleAddr(request []byte) {

}

func HandleUnknownCommand() {

}

func ReplyBack(msg []byte, conn net.Conn, net_platform *NetworkPlatform) {
	wr, err := conn.Write(msg)
	if err != nil {
		fmt.Printf("Failed to reply back to connection %v", conn)
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
		send_payload := NodesRequest{
			AddrFrom: net_platform.GetNodeAddress(),
			Nodes:    []NetworkNode{*net_platform.Self_node},
		}

		sendDataToAddress(payload.AddrFrom, append(CommandStringToBytes("node"), GobEncode(send_payload)...), net_platform)

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

// TODO :: Should read be handled concurrently via go routines?
func HandleTCPConnection(conn net.Conn, net_platform *NetworkPlatform) bool {
	// request, err := io.ReadAll(conn)
	request := make([]byte, 2048)
	conn.SetReadDeadline(time.Now().Add(100 * time.Millisecond))
	len, err := io.ReadAtLeast(conn, request, COMMAND_LENGTH)
	if len == 0 {
		return true
	}

	// defer tcp_connection.Close()
	if err != nil {
		// Close the connection
		if errors.Is(err, net.ErrClosed) {
			log.Printf("Connection closed by the peer")
			return false
		}
		log.Printf(err.Error())
		return false
	}

	// first 32 bytes to hold the command
	// TODO: Format the header data
	l := 0
	for i := range request {
		if i == 32 || l == len {
			break
		}
		l++
	}
	command := BytesToCommandString(request[:l])
	log.Printf("Command: %s %d", command, l)

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

	case "echo":
		ReplyBack(request[l+1:len], conn, net_platform)
		break
	case "getresources":
		HandleResources(request[l+1:len], conn, net_platform)
		break
	}
	return true
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
		time.Sleep(100 * time.Millisecond)
		for i, caches := range net_platform.Connection_caches {
			result := HandleTCPConnection(caches.connection, net_platform)
			if !result {
				// Remove the connection from the cache,bruh bruh. WTF. Khai element remove garney function
				temp := net_platform.Connection_caches
				temp[i] = temp[len(net_platform.Connection_caches)-1]
				net_platform.Connection_caches = temp[:len(temp)-1]
			}
		}
	}
}

func AcceptConnection(conn net.Conn, net_platform *NetworkPlatform) {
	// store the information about the newly connected node into the net_platform struct
	new_node := CreateNetworkNode("unknown", "127.0.0.1", 8000)
	net_platform.Connected_nodes = append(net_platform.Connected_nodes, *new_node)

	// Operation on connection caches are omitted for now
	cache_entry := CreateCacheEntry(conn, nil, new_node.NodeID)
	net_platform.Connection_caches = append(net_platform.Connection_caches, cache_entry)

	log.Print("Connection accepted", conn)
}

// For self
func ListenForTCPConnection(net_platform *NetworkPlatform) {
	listener, err := net.Listen("tcp", net_platform.Self_node.Socket.String())
	if err != nil {
		log.Fatalf("Listener error: %s\n", err.Error())
	}
	defer listener.Close()

	// if net_platform.Self_node.Socket.Port == 6969 {
	// 	log.Printf("Sending get nodes request")
	// 	SendGetNode("127.0.0.1:7000", net_platform)
	// }

	// The call to listen always blocks
	// There's no way to get notified when there is a pending connection in Go?
	go ProcessConnections(net_platform)
	log.Printf("Localhost is listening ... \n")
	for {
		conn, _ := listener.Accept()
		if err != nil {
			fmt.Printf("Failed to Accept the incoming connection.  Error: %s\n", err.Error())
			break
		}
		go AcceptConnection(conn, net_platform)
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
