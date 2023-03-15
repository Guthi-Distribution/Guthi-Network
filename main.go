package main

import (
	"GuthiNetwork/platform"
	"flag"
	"fmt"
	"sync"
)

// There should be one univeral listening port

func main() {
	count = 0
	width = 512
	height = 512
	port := flag.Int("port", 6969, "Port for the network") // send port using command line argument (-port 6969)
	sum_type := flag.Int("range", 0, "Type of range")

	flag.Parse()
	range_number = *sum_type
	fmt.Println(range_number, *sum_type)

	config := LoadConfiguration("config.json")

	net_platform, err := platform.CreateNetworkPlatform(config.Name, config.Address, *port, true)
	if err != nil {
		panic(err)
	}

	if net_platform.Self_node.Socket.Port != 6969 {
		net_platform.ConnectToNode("127.0.0.1:6969") // one of the way to connect to a particular node, request all the nodes information it has
	}
	go platform.ListenForTCPConnection(net_platform)

	net_platform.TrackFile("test.txt")
	var sg sync.WaitGroup

	sg.Add(1)

	sg.Wait()
}
