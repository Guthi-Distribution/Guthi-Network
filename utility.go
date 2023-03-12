package main

import (
	"GuthiNetwork/platform"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"
)

type Config struct {
	Name    string `json:"name"`
	Address string `json:"address"`
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

type Color struct {
	R uint16
	G uint16
	B uint16
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

func node_failure_handler(node platform.NetworkNode) {
	state, err := node.GetFunctionState("render_mandelbrot")
	if err != nil {
		log.Println(err)
		return
	}
	platform.AddPendingDispatch("render_mandelbrot", state)
}
