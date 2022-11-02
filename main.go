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
	sem, err := shm.CreateSemaphore()
	defer sem.RemoveSemaphore()
	if err != nil {
		log.Fatalf("Semaphore creation error: %s", err)
	}
	err = sem.Lock()
	if err != nil {
		log.Fatalf("Lock error: %s", err)
	}
	fmt.Printf("Locked\n")
	err = sem.Unlock()
	if err != nil {
		log.Fatalf("Unlock error: %s", err)
	}
}
