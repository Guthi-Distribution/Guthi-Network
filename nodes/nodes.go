package nodes

import (
	"fmt"
	"hash/crc64"
	"io/ioutil"
	"log"
	"net"
	"time"
)

var net_platform NetworkPlatform

// message for get address request
// request for all the known address
// maybe we could handle nodes directly???
type GetAddress struct {
	AddrFrom string
}

type GetNode struct {
	node []NetworkNode
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
func SyncWithNetwork() uint16 {
	// Receive information about connected nodes from its neighbor nodes
	msg := "Hello there"
	var discovered_nodes uint16

	for _, cached := range net_platform.Connection_caches {

		// This should be bidirectional
		_, err := cached.Connection.Write([]byte(msg))

		if err != nil {
			log.Printf("Failed to recieve response from one of the connected nodes. Error %s", err)
			continue
		}
		// TODO :: Send and receive msg and interpret it
		// Wait for it them to recieve message and compare to them
		// Return the IP Address and port number of other nodes which are listening for p2p connection
		// Read the message and identify new nodes in the network
		discovered_nodes++
	}
	return discovered_nodes
}

// For connecting to the network, at least one node need to be known
func ConnectToNetwork(node *NetworkNode) bool {
	// Connect as a client to the network
	// Maybe implement something like OSPF routing algorithm to create map of the network ??
	tcp_connection := IntiateTCPConnection(node)
	if tcp_connection == nil {
		return false
	}

	// TODO :: Perform other necessary actions to get in sync with the network
	entry := CreateCacheEntry(tcp_connection, node, node.NodeID)
	net_platform.Connection_caches = append(net_platform.Connection_caches, entry)

	SyncWithNetwork()
	return true
}

func CreateNetworkNode(name string, address string, port int) *NetworkNode {
	networkNode := &NetworkNode{}
	networkNode.Name = name

	//TODO: Implement hashing
	id := fmt.Sprintf("%s %s %d", name, address, port)
	table := crc64.MakeTable(100)
	networkNode.Socket = &net.TCPAddr{
		IP:   net.IP(address),
		Port: port,
	}
	networkNode.NodeID = crc64.Checksum([]byte(id), table)
	return networkNode
}

// sends all the node address
func HandleAddr(request []byte) {

}

func HandleUnknownCommand() {

}

func HandleTCPConnection(tcp_connection *net.TCPConn) {
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
	command := string(request[:32])
	log.Printf("Command: %s", command)

	switch command {
	default:
		HandleUnknownCommand()
		break

	case "addr":
		HandleAddr(request)
		break
	}
}

func ProcessConnections() {
	// Process other currently running connections
	for {
		time.Sleep(100 * time.Millisecond)
		for _, conn := range net_platform.Connection_caches {
			// If there's message to be sent to that node send it here.
			// Else, receive message here
			msg := make([]byte, 2048)
			len, err := conn.Connection.Read(msg)
			if err != nil {
				// Connecton has been closed from client side, so drop the connection from the cache list and possibly known nodes (after certain time has elapsed)
				// where's the remove method in slice ??????
			} else {
				if len != 0 {
					fmt.Println("Message recevied from : ", conn.GetNodeRef().NodeID, " :-> ", string(msg[:len]))
					// Echo back the same message to client
					conn.Connection.Write(msg[:len])
				}
			}
		}
	}
}

// For self
func ListenForTCPConnection(node *NetworkNode) {
	listener, err := net.ListenTCP("tcp", node.Socket)

	// The call to listen always blocks
	// There's no way to get notified when there is a pending connection in Go?
	go ProcessConnections()
	for {
		conn, _ := listener.AcceptTCP()
		if err != nil {
			fmt.Println("Failed to Accept the incoming connection")
			break
		}
		go HandleTCPConnection(conn)
	}
	listener.Close()
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
