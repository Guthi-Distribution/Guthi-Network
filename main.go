package main

// There should be one univeral listening port

import (
	"GuthiNetwork/api"
	"GuthiNetwork/platform"
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

func LoadConfiguration(file string) Config {
	var config Config
	configFile, err := os.Open(file)
	defer configFile.Close()
	if err != nil {
		log.Fatal(err.Error())
	}
	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&config)
	return config
}

func main() {
	port := flag.Int("port", 6969, "Port for the network") // send port using command line argument (-port 6969)
	flag.Parse()
	config := LoadConfiguration("config.json")
	net_platform, err := platform.CreateNetworkPlatform(config.Name, config.Address, *port)

	net_platform.CreateVariable("a", 2)

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
	platform.ListenForTCPConnection(net_platform) // listen for connection

	net_platform.CreateVariable("a", "hello there")
	sg.Wait()
}
