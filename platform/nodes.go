package platform

/*
TODO:
- Implement checkpointing
*/
import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"hash/crc64"
	"io"
	"log"
	"net"
)

const (
	COMMAND_LENGTH = 24
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
	networkNode.conn = nil
	networkNode.function_state = map[string]interface{}{}
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
	fmt.Println("Received Node info")
	var payload NodesMessage
	gob.NewDecoder(bytes.NewBuffer(request)).Decode(&payload)
	fmt.Printf("Address received from %s\n", payload.AddrFrom)
	net_platform.Connected_nodes = append(net_platform.Connected_nodes, payload.Nodes...)

	if len(payload.Nodes) == 0 {
		fmt.Printf("Nodes received length is zero")
		return
	}
	entry := CreateCacheEntry(&payload.Nodes[0], payload.Nodes[0].NodeID)
	net_platform.Connection_caches = append(net_platform.Connection_caches, entry)
}

/*
Wrapper function for all the handling of various request and response
first 32 bytes commandf, rest payload
*/
func handleTCPConnection(request []byte, net_platform *NetworkPlatform) error {

	// first 24 bytes to hold the command
	command := BytesToCommandString(request[:COMMAND_LENGTH])

	// TODO: Log this into file
	// log.Printf("Command: %s\n", command)
	switch command {
	default:
		HandleUnknownCommand()
		break

	case "getnodes":
		fmt.Printf("Command: %s\n", command)
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
		fmt.Printf("Command: %s\n", command)
		HandleConnectionInitiation(request[COMMAND_LENGTH:], net_platform)
		break

	case "connection_reply":
		fmt.Printf("Command: %s\n", command)
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
		handleReceiveVariable(request[COMMAND_LENGTH:], net_platform)
		break

	case "array":
		log.Printf("Command: %s\n", command)
		HandleReceiveArray(request[COMMAND_LENGTH:], net_platform)
		break

	case "indexed_array":
		HandleReceiveIndexedArray(request[COMMAND_LENGTH:], net_platform)
		break

	case "symbol_table":
		fmt.Printf("Command: %s\n", command)
		HandleReceiveSymbolTable(request[COMMAND_LENGTH:], net_platform)
		break

	case "symbol_table_ack":
		fmt.Printf("Command: %s\n", command)
		handleReceiveSymbolTableAck(request[COMMAND_LENGTH:], net_platform)
		break

	case "token_request_sk":
		HandleTokenRequest(request[COMMAND_LENGTH:], net_platform)
		break

	case "token":
		HandleReceiveToken(request[COMMAND_LENGTH:], net_platform)
		break

	case "get_var":
		handleGetVariableRequest(request[COMMAND_LENGTH:], net_platform)
		break

	case "validity_info":
		handleVariableInvalidation(request[COMMAND_LENGTH:], net_platform)
		break

	case "function_dispatch":
		fmt.Printf("Command: %s\n", command)
		handleFunctionDispatch(request[COMMAND_LENGTH:], net_platform)
		break

	case "func_state":
		handleFunctionState(request[COMMAND_LENGTH:])
		break

	case "func_completed":
		handleFunctionCompletion(request[COMMAND_LENGTH:])
		break

	}

	request = nil

	return nil
}

// Gob Encode
// Details: https://pkg.go.dev/encoding/gob
func GobEncode(data any) []byte {
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
	go RequestInfomation(net_platform)
	// go CommunicateFileSystem(net_platform)
	go Synchronize(net_platform)
	for {
		conn, err := net_platform.listener.AcceptTCP()
		if err != nil {
			fmt.Printf("Failed to Accept the incoming connection.  Error: %s\n", err.Error())
			continue
		}
		// conn.SetReadDeadline(time.Now().Add(100 * time.Millisecond))
		go func(conn *net.TCPConn) {
			defer conn.Close()
			for {

				var request []byte
				request = make([]byte, 4)
				_, err := io.ReadAtLeast(conn, request, 4)
				if err != nil {
					log.Printf("Connection Reading error while reading length: %s\n", err)
					return
				}

				length := getLengthFromBytes(request)
				if length <= 0 {
					continue
				}

				request = make([]byte, length)
				_, err = io.ReadAtLeast(conn, request, length)
				if err != nil {
					request = nil
					// Close the connection
					if errors.Is(err, net.ErrClosed) {
						fmt.Printf("Connection closed by the peer")
						return
					}
					log.Printf("Connection Reading error while reading data: %s\n", err)
					return
				}
				if len(request) != 0 {
					go handleTCPConnection(request, net_platform)
				}
			}
		}(conn)
	}
}
