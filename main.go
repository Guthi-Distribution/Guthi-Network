package main

// There should be one univeral listening port

import (
	"fmt"
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
	connection_caches []net.TCPConn
}

type NodeConnectionCache struct {
	conn_cache []*net.TCPConn
	// Need some reference to the network node
}

var net_platform NetworkPlatform
var connection_cache NodeConnectionCache

func HandleNetworkConnections() {
	// It have one listener and continuously listens for the incoming connections
	// The problem is call to Accept() blocks until connection has been established

}

// For connecting to the network, at least one node need to be known
func ConnectToNetwork(node *nodes.NetworkNode) bool {
	// Connect as a client to the network
	// Maybe implement something like OSPF routing algorithm to create map of the network ??
	tcp_connection := nodes.IntiateTCPConnection(node)
	if tcp_connection == nil {
		return false
	}

	// TODO :: Perform other necessary actions to get in sync with the network
	connection_cache.conn_cache = append(connection_cache.conn_cache, tcp_connection)
	return true
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

	go nodes.ListenForTCPConnection(&net_platform.self_node)
}
