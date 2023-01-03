package main

// There should be one univeral listening port

import (
	"GuthiNetwork/api"
	"GuthiNetwork/lib"
	"GuthiNetwork/platform"
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"sync"
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

type Config struct {
	Name    string `json:"name"`
	Address string `json:"address"`
}

var range_number int // 1 for 100 to 200 and false for 0 to 100

func sum(total_sum *lib.Variable, net_platform *platform.NetworkPlatform) {
	minimum := 1 + 99*range_number
	maximum := 101 + 99*range_number

	prev_sum := 0
	variable, _ := net_platform.GetValue(total_sum.Id)
	prev_sum = variable.GetData().(int)
	for i := minimum; i <= maximum; i++ {
		fmt.Printf("Previous sum: %d\n", prev_sum)
		prev_sum += i
		total_sum = variable
		fmt.Printf("Total sum: %d\n", variable.GetData())
		time.Sleep(time.Millisecond * 10)
	}

	net_platform.SetValue(total_sum.Id, prev_sum)
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

	// send request to the central node
	if net_platform.Self_node.Socket.Port != 6969 {
		net_platform.ConnectToNode("127.0.0.1:6969") // one of the way to connect to a particular node, request all the nodes information it has
	}

	if *port == 6969 {
		go api.StartServer(net_platform)
	}
	var sg sync.WaitGroup
	sg.Add(1)
	go platform.ListenForTCPConnection(net_platform) // listen for connection

	net_platform.CreateVariable("total_sum", int(0))
	total_sum, err := net_platform.GetValue("total_sum")
	if range_number == 1 || true {
		reader := bufio.NewReader(os.Stdin)
		reader.ReadString('\n')
	}
	sum(total_sum, net_platform)
	fmt.Printf("Total sum: %d\n", total_sum.GetData())
	sg.Wait()
}
