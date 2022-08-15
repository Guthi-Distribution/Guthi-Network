package node

import (
	"Guthi/utility"
	"fmt"
	"io/ioutil"
	"log"
	"net"
)

//TODO: Interface to the C++ filesystem
//TODO: Time Implementation
//TODO: Implement runtime avaialble resources for every node
type Node struct {
	id string
	//TODO: Implement to update it dynamically to get the runtime information
	avaiable_nodes []string // address of all the connected nodes
	time           Time     // time of the node
}

func CreateNode(address string) Node {
	node := Node{
		id: address,
	}

	return node
}

func (node *Node) GetAddress() string {
	return node.id
}

func (node *Node) connect(address string, portNumber int32) {
	nodeAddress := fmt.Sprintf("%s: %d", address, portNumber)

	if node.id == nodeAddress {
		panic("Destination same as the receiver")
	}
}

func HandleConnection(conn net.Conn) {
	//TODO: Get request command
	// TODO: Get Node Information
	//TODO: Decide on the overall format of the message
	// TODO: Get payload

	req, err := ioutil.ReadAll(conn)
	defer conn.Close()
	if err != nil {
		utility.ErrThenLogPanic(err)
	}

	command := string(req[:20])
	fmt.Println(command)
}

func StartServer(port int, protocol string) {
	address := fmt.Sprintf("%s:%d", utility.GetNodeAddress(), port)
	fmt.Println(address)
	ln, err := net.Listen(protocol, address)

	if err != nil {
		log.Panic(err)
	}
	defer ln.Close()

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Panic(err)
		}

		go HandleConnection(conn)
	}
}
