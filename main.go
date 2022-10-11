package main

// There should be one univeral listening port

import (
	"GuthiNetwork/nodes"
	"GuthiNetwork/shm"
	"fmt"
	"log"
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
	mem, err := shm.CreateSharedMemory()
	if err != nil {
		log.Fatal(err)
	}
	s := "Hello"
	mem.WriteSharedMemory([]byte(s))
	defer mem.RemoveSharedMemory()
	// var n int
	mem.WriteSharedMemory([]byte("Hello there again mfers"))
	mem.ReadSharedMemory()
	defer mem.RemoveSharedMemory()
}
