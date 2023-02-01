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

	"github.com/mitchellh/go-ps"
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
	range_sum := 1000
	minimum := 1 + range_sum*range_number
	maximum := range_sum + range_sum*range_number
	for i := minimum; i <= maximum; i++ {

		// LOG: works when sleep is large because synchronization is not needed that much
		// but when data is propagating quickly, it really fucks up
		platform.Lock(net_platform)
		variable, _ := net_platform.GetValue(total_sum.Id)

		prev_sum := 0

		prev_sum = variable.GetData().(int)
		fmt.Printf("\n\nPrevious Value: %d\n", prev_sum)
		fmt.Printf("Adding value: %d\n", i)
		prev_sum += i

		net_platform.SetData(total_sum.Id, prev_sum)
		fmt.Printf("Updated Value: %d\n", total_sum.GetData().(int))
		platform.Unlock(net_platform)
		time.Sleep(time.Millisecond * 50)
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

func main() {
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

	net_platform.CreateVariable("total_sum", int(0))
	total_sum, err := net_platform.GetValue("total_sum")

	if !isLaunchedByDebugger() {
		fmt.Println("Not Debugging process")
		if true {
			reader := bufio.NewReader(os.Stdin)
			reader.ReadString('\n')
		}
	} else {
		for len(net_platform.Connected_nodes) == 0 {

		}
	}
	sum(total_sum, net_platform)
	fmt.Printf("Total sum: %d\n", total_sum.GetData())
	sg.Wait()
}
