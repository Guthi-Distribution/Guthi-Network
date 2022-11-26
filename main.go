package main

// There should be one univeral listening port

import (
	"GuthiNetwork/api"
	"GuthiNetwork/platform"
	"flag"
	"fmt"
	"log"
	"time"

	"GuthiNetwork/core"
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

func main() {
	core.Initialize()
	go core.ReadSharedMemory()
	port := flag.Int("port", 6969, "Port for the network") // send port using command line argument (-port 6969)
	flag.Parse()
	net_platform, err := platform.CreateNetworkPlatform("localhost", "localhost", *port)
	if err != nil {
		log.Fatalf("Platform Creation error: %s", err)
	}

	// send request to the central node
	if net_platform.Self_node.Socket.Port != 6969 {
		net_platform.ConnectToNode("127.0.0.1:6969") // one of the way to connect to a particular node, request all the nodes information it has
	}
	if *port == 6969 {
		go api.StartServer(net_platform)
	}
	platform.ListenForTCPConnection(net_platform) // listen for connection
}
