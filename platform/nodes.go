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
	"io/ioutil"
	"log"
	"net"
	"time"
)

const (
	COMMAND_LENGTH = 16
)

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

/*
Handle Get node request from the node
@param request: request byte
@net_platform: pointer to network platform
*/
func HandleGetNodes(request []byte, net_platform *NetworkPlatform) error {
	var payload GetNodes
	gob.NewDecoder(bytes.NewBuffer(request)).Decode(&payload)

	fmt.Printf("Sending nodes to address: %s\n", payload.AddrFrom)
	if payload.Address[0] == net_platform.Self_node.Socket.String() {
		// if the receiving address is the self address, then it is send
		send_payload := NodesMessage{
			AddrFrom: net_platform.GetNodeAddress(),
			Nodes:    []NetworkNode{*net_platform.Self_node},
		}

		return sendDataToAddress(payload.AddrFrom, append(CommandStringToBytes("node"), GobEncode(send_payload)...), net_platform)
	}

	return nil
}

/*
Handles when nodes information is received
- Adds the nodes to connected nodes
*/
func HandleNodeResponse(request []byte, net_platform *NetworkPlatform) {
	var payload NodesMessage
	gob.NewDecoder(bytes.NewBuffer(request)).Decode(&payload)
	fmt.Printf("Address received from %s\n", payload.AddrFrom)
	net_platform.Connected_nodes = append(net_platform.Connected_nodes, payload.Nodes...)

	if len(payload.Nodes) == 0 {
		log.Printf("Nodes received length is zero")
		return
	}
	entry := CreateCacheEntry(&payload.Nodes[0], payload.Nodes[0].NodeID)
	net_platform.Connection_caches = append(net_platform.Connection_caches, entry)
}

/*
Wrapper function for all the handling of various request and response
*/
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
		HandleNodeResponse(request[COMMAND_LENGTH:], net_platform)
		break

	case "echo":
		HandleEchoMessage(request[COMMAND_LENGTH:], net_platform)
		break

	case "echo_reply":
		HandleEchoReply(request[COMMAND_LENGTH:], net_platform)
		break

	case "connect":
		HandleConnectionInitiation(request[COMMAND_LENGTH:], net_platform)
		break

	case "connection_reply":
		HandleConnectionReply(request[COMMAND_LENGTH:], net_platform)
		break

	case "get_mem_info":
		HandleGetMemoryInformation(request[COMMAND_LENGTH:], net_platform)
		break

	case "get_cpu_info":
		HandleGetCpuInformation(request[COMMAND_LENGTH:], net_platform)
		break

	case "cpuinfo":
		HandleReceiveCpuInformation(request[COMMAND_LENGTH:], net_platform)
		break

	case "meminfo":
		HandleReceiveMemoryInformation(request[COMMAND_LENGTH:], net_platform)
		break

	case "get_fs":
		HandleGetFileSystem(request[COMMAND_LENGTH:], net_platform)
		break

	case "filesystem":
		HandleReceiveFileSystem(request[COMMAND_LENGTH:], net_platform)
		break

	case "variable":
		HandleReceiveVariable(request[COMMAND_LENGTH:], net_platform)
		break

	case "symbol_table":
		HandleReceiveSymbolTable(request[COMMAND_LENGTH:], net_platform)
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
		log.Panic("Gob Encode error:" + err.Error())
	}

	return buff.Bytes()
}

// For self
func ListenForTCPConnection(net_platform *NetworkPlatform) {
	// create a listener that is used to listen to other connection (somethong like that)
	// Listen announces on the local network address. @docs

	// The call to listen always blocks
	// There's no way to get notified when there is a pending connection in Go?
	log.Printf("Localhost is listening ... \n")
	// go RequestInfomation(net_platform)
	// go CommunicateFileSystem(net_platform)
	// go Synchronize(net_platform)
	for {
		conn, err := net_platform.listener.Accept()
		if err != nil {
			fmt.Printf("Failed to Accept the incoming connection.  Error: %s\n", err.Error())
			break
		}
		go HandleTCPConnection(conn, net_platform)
	}
}
