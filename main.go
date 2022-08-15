package main

import (
	"Guthi/node"
	"flag"
)

func main() {
	port := flag.Int("port", 6969, "Enter port number")
	protocol := flag.String("protocol", "tcp", "Enter port number")
	node.StartServer(*port, *protocol)
}
