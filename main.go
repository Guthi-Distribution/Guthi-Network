package main

// There should be one univeral listening port

import (
	"GuthiNetwork/platform"
	"encoding/gob"
	"flag"
	"fmt"
	"sync"
	"time"
)

func main() {
	count = 0
	width = 256
	height = 256
	port := flag.Int("port", 6969, "Port for the network") // send port using command line argument (-port 6969)
	sum_type := flag.Int("range", 0, "Type of range")

	flag.Parse()
	range_number = *sum_type
	fmt.Println(range_number, *sum_type)

	config := LoadConfiguration("config.json")
	net_platform, err := platform.CreateNetworkPlatform(config.Name, config.Address, *port)
	if err != nil {
		panic(err)
	}
	if net_platform.Self_node.Socket.Port != 6969 {
		net_platform.ConnectToNode("127.0.0.1:6969") // one of the way to connect to a particular node, request all the nodes information it has
	}
	go platform.ListenForTCPConnection(net_platform)

	var sg sync.WaitGroup
	sg.Add(1)
	c := Color{}

	gob.Register(Color{})
	gob.Register(MandelbrotParam{})

	net_platform.RegisterFunction(render_mandelbrot)
	net_platform.BindNodeFailureEventHandler(node_failure_handler)
	if *port == 6969 {
		net_platform.BindFunctionCompletionEventHandler("render_mandelbrot", plot_mandelbrot)
		curr_time := time.Now().UnixMilli()
		net_platform.CreateArray("mandelbrot", width*height, c)
		fmt.Println(time.Now().UnixMilli() - curr_time)
		fmt.Println("Not Debugging process")

		args := []interface{}{
			MandelbrotParam{0, 0},
			MandelbrotParam{1, 0},
		}
		// time.Sleep(time.Second * 2)
		net_platform.DispatchFunction("render_mandelbrot", args)
	}

	sg.Wait()
}
