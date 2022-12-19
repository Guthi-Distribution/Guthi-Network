package main

// There should be one univeral listening port

import (
	"GuthiNetwork/api"
<<<<<<< HEAD
	"GuthiNetwork/lib"
=======
	"GuthiNetwork/core"
>>>>>>> 66cf85f2047b38e2cb00a9951164edaae563c8b6
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
	err := core.Initialize()
	// if err != nil {
	// 	log.Fatal(err.Error())
	// }
	// go core.ReadSharedMemory()
	port := flag.Int("port", 6969, "Port for the network") // send port using command line argument (-port 6969)
	flag.Parse()
	net_platform, err := platform.CreateNetworkPlatform("localhost", "localhost", *port)
	fmt.Println(net_platform.Self_node.Socket.IP)
	if err != nil {
		log.Fatalf("Platform Creation error: %s", err)
	}

	// send request to the central node
	if net_platform.Self_node.Socket.Port != 6969 {
		net_platform.ConnectToNode("192.168.45.68:6969") // one of the way to connect to a particular node, request all the nodes information it has
	}
	if *port == 6969 {
		go api.StartServer(net_platform)
	}
	platform.ListenForTCPConnection(net_platform) // listen for connection
}
