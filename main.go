package main

// There should be one univeral listening port

import (
	"GuthiNetwork/api"
	"GuthiNetwork/nodes"
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

func main() {
	net_platform.Self_node = nodes.CreateNetworkNode("localhost", "127.0.0.1", 8000)
	api.InitlializeServer(&net_platform)
}
