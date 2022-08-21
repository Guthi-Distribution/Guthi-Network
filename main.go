package main

// There should be one univeral listening port

import (
	"fmt"
	"math/rand"
	"net"
	"os"
	"sync"
	"time"

	"./src/nodes"
)

// Go is such a stupid language, All hail C++

// tf is Rune ... lol
func wait_loop(elapsed time.Duration) {
	for {
		fmt.Printf("\r")
		for _, r := range "-\\|/" {
			fmt.Printf("%c", r)
			time.Sleep(elapsed)
		}
	}
}

type NetworkPlatform struct {
	// Well, there's just a single writer but multiple readers. So RWMutex sounds better choice
	lock              sync.Mutex
	self_node         nodes.NetworkNode
	connected_nodes   []nodes.NetworkNode
	connection_caches []CacheEntry
}

type CacheEntry struct {
	connection *net.TCPConn
	node_ref   *nodes.NetworkNode
}

type NodeConnectionCache struct {
	cache []CacheEntry
	// Need some reference to the network node
}

var net_platform NetworkPlatform

// For connecting to the network, at least one node need to be known
func ConnectToNetwork(node *nodes.NetworkNode) bool {
	// Connect as a client to the network
	// Maybe implement something like OSPF routing algorithm to create map of the network ??
	tcp_connection := nodes.IntiateTCPConnection(node)
	if tcp_connection == nil {
		return false
	}

	// TODO :: Perform other necessary actions to get in sync with the network
	entry := CacheEntry{tcp_connection, node}
	net_platform.connection_caches = append(net_platform.connection_caches, entry)

	SyncWithNetwork()
	return true
}

// Onto nodes discovery
// How to decide if networks are in sync? ans -> After certain time lol
// Should this function be called on regular basis? On certain interval or not?

func SyncWithNetwork() uint16 {
	// Receive information about connected nodes from its neighbor nodes
	msg := "Send nudes"
	var discovered_nodes uint16
	for _, cached := range net_platform.connection_caches {

		// This should be bidirectional
		_, err := cached.connection.Write([]byte(msg))

		if err != nil {
			fmt.Println("Failed to recieve response from one of the connected nodes.")
			continue
		}
		// TODO :: Send and receive msg and interpret it
		// Wait for it them to recieve message and compare to them
		// Return the IP Address and port number of other nodes which are listening for p2p connection
		discovered_nodes++
	}
	return discovered_nodes
}

func BroadcastMessage(msg string) {
	for _, cached := range net_platform.connection_caches {
		go cached.connection.Write([]byte(msg))
	}
}

func InitializePlatform() bool {
	net_platform.lock.Lock()

	// Add localhost to the net_platform here
	// TODO :: Maybe choose a different port for local connection and retry to find new unassigned port
	// TODO :: Use public IP

	localhost := "localhost:8080"
	tcp_addr, err := net.ResolveTCPAddr("tcp", localhost)
	if err != nil {
		return false
	}

	node := nodes.NetworkNode{1, "localhost", tcp_addr}
	net_platform.connected_nodes = append(net_platform.connected_nodes, node)

	net_platform.self_node = node
	// The initiating host is itself a client in the net platform.
	net_platform.lock.Unlock()
	return true
}

func main() {
	// Allocating a thread/goroutine whatever they are called, for actively listening all the times would be such a waste
	// For small scale, there are few nodes that are likely to connect to the network, so it could be event driven along with other tasks performed simultaneously
	fmt.Println("Hello from Gulang")
	fmt.Println("All Hail C++ annnndddd ... Rust")

	// During start phase, create a map or dynamic array of Nodes
	// Pass this to a go routine that actively looks for the node connection from the outside
	// Race condition here we go

	if !InitializePlatform() {
		fmt.Println("Failed to initialize the platform")
		os.Exit(-1)
	}

	go ListenForTCPConnection(&net_platform.self_node)
}

// For self
func ListenForTCPConnection(node *nodes.NetworkNode) {
	listener, err := net.ListenTCP("tcp", node.Socket)

	// The call to listen always blocks
	// There's no way to get notified when there is a pending connection in Go?
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

func HandleTCPConnection(tcp_connection *net.TCPConn) {
	// store the information about the newly connected node into the net_platform struct
	// So connection established, now retrieve information about the host

	// TODO :: Test this implementation, left for Go experts
	// Assuming that Garbage collected language can handle anything, literally anything
	// Like some memory allocated by another runtime too.. lol

	tcp_addr := tcp_connection.RemoteAddr().(*net.TCPAddr)
	new_node := nodes.NetworkNode{rand.Intn(1000), "unknown", tcp_addr}
	net_platform.connected_nodes = append(net_platform.connected_nodes, new_node)
}
