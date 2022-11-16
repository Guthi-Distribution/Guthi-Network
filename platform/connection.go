package platform

import (
	"bytes"
	"crypto/rand"
	"encoding/gob"
	"fmt"
	"log"
	"net"
	"time"
)

// TODO: Connection timeout feature added
var pending_connection = make(map[string]uint64)

/*
Initiate TCP Connection, creates connection and returns the connection
*/
func intiateTCPConnection(node *NetworkNode) (*net.Conn, error) {
	tcp_con, err := net.Dial("tcp", node.Socket.String())
	if err != nil {
		fmt.Println("Failed to initiate tcp connection with the host : ", node)
		return nil, err
	}
	return &tcp_con, err
}

/*
Connect to the node, to a specific address
Returns:

	error if any has occured
*/
func (net_platform *NetworkPlatform) ConnectToNode(address string) error {
	rand_num, err := rand.Prime(rand.Reader, 64)
	payload := ConnectionRequest{
		AddrFrom:  net_platform.Self_node.Socket.String(),
		ConnectId: rand_num.Uint64(),
	}
	pending_connection[address] = uint64(time.Now().Unix())
	// connect to the network
	data := GobEncode(payload)
	data = append(CommandStringToBytes("connect"), data...)

	err = sendDataToAddress(address, data, net_platform)
	if err != nil {
		return err
	}

	return nil
}

/*
Connect to the node, to a specific node
Returns:

	error if any has occured
*/
func ConnectToNetwork(node *NetworkNode, net_platform *NetworkPlatform) error {
	// Connect as a client to the network
	// Maybe implement something like OSPF routing algorithm to create map of the network ??
	_, err := intiateTCPConnection(node)
	if err != nil {
		log.Printf("Connection setup error: %s", err)
		return err
	}

	// TODO :: Perform other necessary actions to get in sync with the network
	entry := CreateCacheEntry(node, node.NodeID)
	net_platform.Connection_caches = append(net_platform.Connection_caches, entry)
	SyncWithNetwork(net_platform)
	return nil
}

/*
Respond to connect request from the node
*/
func HandleConnectionInitiation(request []byte, net_platform *NetworkPlatform) error {
	var payload ConnectionRequest
	gob.NewDecoder(bytes.NewBuffer(request)).Decode(&payload)
	send_payload := ConnectionReply{
		AddrFrom:  net_platform.GetNodeAddress(),
		Node:      *net_platform.Self_node,
		ConnectId: payload.ConnectId + 1,
		IsReply:   false,
	}

	err := sendDataToAddress(payload.AddrFrom, append(CommandStringToBytes("connection_reply"), GobEncode(send_payload)...), net_platform)
	if err != nil {
		return err
	}
	if err != nil {
		return err
	}

	return nil
}

func HandleConnectionReply(request []byte, net_platform *NetworkPlatform) error {
	var payload ConnectionReply
	gob.NewDecoder(bytes.NewBuffer(request)).Decode(&payload)
	net_platform.AddNode(payload.Node)

	if !payload.IsReply {
		// then a reply is recieved, reply with the self node information
		send_payload := ConnectionReply{
			AddrFrom:  net_platform.GetNodeAddress(),
			Node:      *net_platform.Self_node,
			ConnectId: payload.ConnectId + 1,
			IsReply:   true,
		}
		err := sendDataToAddress(payload.AddrFrom, append(CommandStringToBytes("connection_reply"), GobEncode(send_payload)...), net_platform)
		if err != nil {
			return err
		}
	}

	return nil
}
