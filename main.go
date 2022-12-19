package main

// There should be one univeral listening port

import (
	"GuthiNetwork/api"
	"GuthiNetwork/lib"
	"GuthiNetwork/platform"
	"flag"
	"fmt"
	"log"
	"time"
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
	// v, err := lib.CreateVariable("a", 2)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// v.SetValue(3)
	// log.Println(v.GetValue())
	// err := core.Initialize()
	// // if err != nil {
	// // 	log.Fatal(err.Error())
	// // }
	// // go core.ReadSharedMemory()
	lib.CreateVariable("a", 2)
	port := flag.Int("port", 6969, "Port for the network") // send port using command line argument (-port 6969)
	flag.Parse()
	net_platform, err := platform.CreateNetworkPlatform("sanskar", "localhost", *port)
	if err != nil {
		log.Fatalf("Platform Creation error: %s", err)
	}

	// send request to the central node
	if net_platform.Self_node.Socket.Port != 7000 {
		net_platform.ConnectToNode("192.168.45.68:6969") // one of the way to connect to a particular node, request all the nodes information it has
	}
	if *port == 6969 {
		go api.StartServer(net_platform)
	}
	platform.ListenForTCPConnection(net_platform) // listen for connection
}
