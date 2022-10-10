package main

// There should be one univeral listening port

import (
	"fmt"
	"time"

	"./nodes"
)
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
	node := nodes.CreateNetworkNode("localhost", "localhost", 8080)
	net_platform.Connected_nodes = append(net_platform.Connected_nodes, *node)

	net_platform.Self_node = node
	// The initiating host is itself a client in the net platform.
	net_platform.Lock.Unlock()
	return true
}

func main() {
	InitializePlatform()
	nodes.ListenForTCPConnection(&net_platform)

	// Check the connection by launching multiple instance of nc 
	
}
