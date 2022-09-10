package main

// There should be one univeral listening port

import (
	"GuthiNetwork/nodes"
	"fmt"
	"net"
	"time"
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

var net_platform nodes.NetworkPlatform

func BroadcastMessage(msg string) {
	for _, cached := range net_platform.Connection_caches {
		go cached.Connection.Write([]byte(msg))
	}
}

func InitializePlatform() bool {
	net_platform.Lock.Lock()

	// Add localhost to the net_platform here
	// TODO :: Maybe choose a different port for local connection and retry to find new unassigned port
	// TODO :: Use public IP
	localhost := "localhost:8080"
	tcp_addr, err := net.ResolveTCPAddr("tcp", localhost)
	if err != nil {
		return false
	}

	node := nodes.NetworkNode{NodeID: 1, Name: "localhost", Socket: tcp_addr}
	net_platform.Connected_nodes = append(net_platform.Connected_nodes, node)

	net_platform.Self_node = node
	// The initiating host is itself a client in the net platform.
	net_platform.Lock.Unlock()
	return true
}

func main() {
	node := nodes.CreateNetworkNode("Node 1", "127.0.0.1", 8000)
	nodes.ListenForTCPConnection(node)
}
