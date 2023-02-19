package main

// There should be one univeral listening port

import (
	"GuthiNetwork/api"
	"GuthiNetwork/platform"
	"bufio"
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

func render_mandelbrot(net_platform *platform.NetworkPlatform) {
	diff := (720 / 2)
	min := 0 + range_number*diff
	max := diff + range_number*diff
	fmt.Println(min, max)
	for i := 0; i < 1080; i++ {
		for j := min; j < max; j++ {
			_data, err := net_platform.GetDataArray("mandelbrot", 720*i+j)
			if err != nil {
				panic(err)
			}
			data := _data.(Color)
			fmt.Println(data.R)
		}
	}
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

/*
func main() {
	port := flag.Int("port", 6969, "Port for the network") // send port using command line argument (-port 6969)
	sum_type := flag.Int("range", 0, "Type of range")

	flag.Parse()
	range_number = *sum_type
	fmt.Println(range_number, *sum_type)

	config := LoadConfiguration("config.json")
	net_platform, err := platform.CreateNetworkPlatform(config.Name, config.Address, *port)

	fmt.Println(net_platform.Self_node.Socket.IP)
	if err != nil {
		log.Fatalf("Platform Creation error: %s", err)
	}
	gob.Register(Color{})

	// send request to the central node
	if net_platform.Self_node.Socket.Port != 6969 {
		net_platform.ConnectToNode("127.0.0.1:6969") // one of the way to connect to a particular node, request all the nodes information it has

		net_platform.ClaimToken()
		log.Print("Claiming token for this node")
	} else {
	}

	if *port == 6969 {
		go api.StartServer(net_platform)
	}
	var sg sync.WaitGroup
	sg.Add(1)
	go platform.ListenForTCPConnection(net_platform) // listen for connection

	curr_time := time.Now().UnixMilli()
	net_platform.CreateArray("mandelbrot", 1080*720, Color{})
	fmt.Println(time.Now().UnixMilli() - curr_time)
	fmt.Scanln("")

	fmt.Println("hello")
	// render_mandelbrot(net_platform)
	fmt.Println("Repeat hello")

	im := image.NewRGBA(image.Rectangle{image.Point{0, 0}, image.Point{1080, 720}})
	for i := 0; i < 1080; i++ {
		for j := 0; j < 720; j++ {
			im.Set(i, j, color.RGBA{})
		}
	}
	for i := 0; i < 1080; i++ {
		im.SetRGBA(i, 10, color.RGBA{})
	}
	output, _ := os.Create("img.png")
	png.Encode(output, im)

	sg.Wait()
}
*/

func main() {
	// port := flag.Int("port", 6969, "Port for the network") // send port using command line argument (-port 6969)
	// sum_type := flag.Int("range", 0, "Type of range")

	// flag.Parse()
	// range_number = *sum_type
	// fmt.Println(range_number, *sum_type)

	// config := LoadConfiguration("config.json")

	// net_platform, err := platform.CreateNetworkPlatform(config.Name, config.Address, *port)
	// net_platform.CreateVariable("total_sum", int(0))

	// go platform.ListenForTCPConnection(net_platform) // listen for connection
	// if err != nil {
	// 	log.Panic(err)
	// }
	// var sg sync.WaitGroup
	// sg.Add(1)
	// range_number := 1
	// if *port == 6969 {
	// 	net_platform.ClaimToken()
	// 	range_number = 0
	// } else {
	// 	net_platform.ConnectToNode("127.0.0.1:6969")
	// }
	// // core.Initialize()
	// // core.CreateFile("hello_there.txt", "hello_there")

	// net_platform.RegisterFunction(sum)
	// fmt.Printf("Range number: %d\n", range_number)
	// platform.CallFunction(platform.GetFunctionName(sum), range_number, "")
	// sg.Wait()

	port := flag.Int("port", 6969, "Port for the network") // send port using command line argument (-port 6969)
	sum_type := flag.Int("range", 0, "Type of range")

	flag.Parse()
	range_number = *sum_type
	fmt.Println(range_number, *sum_type)

	config := LoadConfiguration("configy.json")
	net_platform, err := platform.CreateNetworkPlatform(config.Name, config.Address, *port)

	fmt.Println(net_platform.Self_node.Socket.IP)
	if err != nil {
		log.Fatalf("Platform Creation error: %s", err)
	}

	// send request to the central node
	if net_platform.Self_node.Socket.Port != 6969 {
		net_platform.ConnectToNode("127.0.0.1:6969") // one of the way to connect to a particular node, request all the nodes information it has
		log.Print("Claiming token for this node")
	} else {
		net_platform.ClaimToken()
	}

	if *port == 6969 {
		go api.StartServer(net_platform)
	}
	var sg sync.WaitGroup
	sg.Add(1)
	go platform.ListenForTCPConnection(net_platform) // listen for connection

	net_platform.CreateVariable("total_sum", int(0))
	total_sum, err := net_platform.GetValue("total_sum")

	if !isLaunchedByDebugger() {
		// fmt.Println("Not Debugging process")
		// reader := bufio.NewReader(os.Stdin)
		// reader.ReadString('\n')
	} else {
		for len(net_platform.Connected_nodes) == 0 {

		}
	}
	net_platform.RegisterFunction(sum)
	// fmt.Printf("Range number: %d\n", range_number)
	if *port == 6969 {
		fmt.Println("Not Debugging process")
		reader := bufio.NewReader(os.Stdin)
		reader.ReadString('\n')
		net_platform.CallFunction(platform.GetFunctionName(sum), 0, "")
		net_platform.CallFunction(platform.GetFunctionName(sum), 1, net_platform.Connected_nodes[0].GetAddressString())
	}
	fmt.Printf("Total sum: %d\n", total_sum.GetData())
	sg.Wait()
}
