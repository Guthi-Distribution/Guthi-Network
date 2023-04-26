package main

// There should be one univeral listening port

import (
	"encoding/gob"
	"flag"
	"fmt"
	"sync"
	"time"

	"github.com/Guthi/guthi_network/platform"
	"github.com/Guthi/guthi_network/renderer"

	"github.com/Guthi/guthi_network/api"
	"github.com/Guthi/guthi_network/utility"
)

var start_time time.Time

func main() {
	count = 0
	port := flag.Int("port", 6969, "Port for the network") // send port using command line argument (-port 6969)
	sum_type := flag.Int("range", 0, "Type of range")

	flag.Parse()
	range_number = *sum_type
	fmt.Println(range_number, *sum_type)

	config := utility.LoadConfiguration("config.json")
	config.Port = *port
	net_platform, err := platform.CreateNetworkPlatform(config)
	if err != nil {
		panic(err)
	}
	if net_platform.Self_node.Socket.Port != 6969 {
		net_platform.ConnectToNode("127.0.0.1:6969") // one of the way to connect to a particular node, request all the nodes information it has
	}
	go platform.ListenForTCPConnection(net_platform)
	go api.StartServer(net_platform)

	var sg sync.WaitGroup
	sg.Add(1)
	c := Color{}

	gob.Register(Color{})
	gob.Register(MandelbrotParam{})

	net_platform.RegisterFunction(render_mandelbrot)
	net_platform.BindNodeFailureEventHandler(node_failure_handler)
	start_time = time.Now()

	if *port == 6969 {
		// Initialize the renderer
		renderer.InitializeRenderer(int32(width), int32(height))
		net_platform.BindFunctionCompletionEventHandler("render_mandelbrot", plot_mandelbrot)
		curr_time := time.Now().UnixMilli()
		net_platform.CreateArray("mandelbrot", width*height*3, c)
		fmt.Println(time.Now().UnixMilli() - curr_time)

		args := []interface{}{}

		for i := 0; i < width/block_size; i++ {
			for j := 0; j < height/block_size; j++ {
				args = append(args, MandelbrotParam{i * block_size, j * block_size})
			}
		}

		fmt.Println(args...)
		net_platform.DispatchFunction("render_mandelbrot", args)
	}

	// Wait for rendering
	renderer.PollEvents()

	sg.Wait()
}
