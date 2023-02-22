package main

// There should be one univeral listening port

import (
	"GuthiNetwork/platform"
	"bufio"
	"encoding/gob"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/mitchellh/go-ps"
)

/*
	TODO:
		- State Management
		- Creation of variable in single node
*/

type Config struct {
	Name    string `json:"name"`
	Address string `json:"address"`
}

var range_number int // 1 for 100 to 200 and false for 0 to 100

func render_mandelbrot(range_number int) {
	diff := (256 / 2)
	min := 0 + range_number*diff
	max := diff + range_number*diff
	fmt.Println(min, max)
	net_platform := platform.GetPlatform()
	for i := 0; i < 256; i++ {
		for j := min; j < max; j++ {
			_, err := net_platform.GetDataOfArray("mandelbrot", 256*i+j)
			if err != nil {
				log.Printf("Index: %d\n", 256*i+j)
				panic(err)
			}
			// data := _data.(Color)
		}
	}

	fmt.Println("Completed")
}

func LoadConfiguration(file string) Config {
	var config Config
	configFile, err := os.Open(file)
	defer configFile.Close()
	if err != nil {
		config.Name = ""
		config.Address = ""
		return config
	}
	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&config)
	return config
}

func sum(range_number int) {
	range_sum := 1000
	minimum := 1 + range_sum*range_number
	maximum := range_sum + range_sum*range_number
	net_platform := platform.GetPlatform()
	for i := minimum; i <= maximum; i++ {
		platform.Lock(net_platform)
		prev_sum := 0
		prev_sum_interface, _ := net_platform.GetData("total_sum")
		prev_sum = prev_sum_interface.(int)
		prev_sum += i
		net_platform.SetData("total_sum", prev_sum)
		platform.Unlock(net_platform)
		time.Sleep(time.Millisecond)
	}

	sum, _ := net_platform.GetData("total_sum")
	fmt.Println("Total sum: ", sum)
}

func isLaunchedByDebugger() bool {
	pid := os.Getppid()

	// We loop in case there were intermediary processes like the gopls language server.
	for pid != 0 {
		switch p, err := ps.FindProcess(pid); {
		case err != nil:
			return false
		case p.Executable() == "dlv":
			return true
		default:
			pid = p.PPid()
		}
	}
	return false
}

type Color struct {
	R uint16
	G uint16
	B uint16
}

func main() {
	port := flag.Int("port", 6969, "Port for the network") // send port using command line argument (-port 6969)
	sum_type := flag.Int("range", 0, "Type of range")

	flag.Parse()
	range_number = *sum_type
	fmt.Println(range_number, *sum_type)

	config := LoadConfiguration("configy.json")
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
	// net_platform.CreateFile("test", "test_contents")
	c := Color{}
	gob.Register(c)

	net_platform.RegisterFunction(render_mandelbrot)
	if *port == 6969 {
		curr_time := time.Now().UnixMilli()
		net_platform.CreateArray("mandelbrot", 256*256, c)
		fmt.Println(time.Now().UnixMilli() - curr_time)
		fmt.Println("Not Debugging process")

		if !isLaunchedByDebugger() {
			reader := bufio.NewReader(os.Stdin)
			reader.ReadString('\n')
		} else {
			for len(net_platform.Connected_nodes) == 0 {

			}
		}
		net_platform.CallFunction(platform.GetFunctionName(render_mandelbrot), 0, "")
		net_platform.CallFunction(platform.GetFunctionName(render_mandelbrot), 1, net_platform.Connected_nodes[0].GetAddressString())
	}

	sg.Wait()
}
