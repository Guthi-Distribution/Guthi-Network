package main

// There should be one univeral listening port

import (
	"GuthiNetwork/nodes"
	"flag"
	"fmt"
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

func InitializePlatform() bool {
	net_platform.Lock.Lock()

	// Add localhost to the net_platform here
	// TODO :: Maybe choose a different port for local connection and retry to find new unassigned port
	// TODO :: Use public IP
	node := nodes.CreateNetworkNode("localhost", "localhost", 8000)
	net_platform.Connected_nodes = append(net_platform.Connected_nodes, *node)

	net_platform.Self_node = node
	// The initiating host is itself a client in the net platform.
	net_platform.Lock.Unlock()
	return true
}

func main() {
	var port int
	flag.IntVar(&port, "port", 6969, "-port")
	flag.Parse()
	net_platform.Self_node = nodes.CreateNetworkNode("localhost", "127.0.0.1", port)
	net_platform.GetNodeAddress()
	nodes.ListenForTCPConnection(&net_platform)
}
