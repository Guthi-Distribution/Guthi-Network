package nodes

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
	"net"
)

// payload to send when we will be sending the data
type ConnectData struct {
	AddrFrom   string
	SourceNode NetworkNode
}

func intiateTCPConnection(node *NetworkNode) (*net.Conn, error) {
	tcp_con, err := net.Dial("tcp", node.Socket.String())
	if err != nil {
		fmt.Println("Failed to initiate tcp connection with the host : ", node)
		return nil, err
	}
	return &tcp_con, err
}

func (net_platform *NetworkPlatform) ConnectToNode(address string) error {
	payload := ConnectData{
		AddrFrom:   net_platform.Self_node.Socket.String(),
		SourceNode: *net_platform.Self_node,
	}

	// connect to the network
	data := GobEncode(payload)
	data = append(CommandStringToBytes("connect"), data...)
	err := sendDataToAddress(address, data, net_platform)
	if err != nil {
		return err
	}

	return nil
}

// For connecting to the network, at least one node need to be known
func ConnectToNetwork(node *NetworkNode, net_platform *NetworkPlatform) bool {
	// Connect as a client to the network
	// Maybe implement something like OSPF routing algorithm to create map of the network ??
	tcp_connection, err := intiateTCPConnection(node)
	if err != nil {
		log.Printf("Connection setup error: %s", err)
		return false
	}

	// TODO :: Perform other necessary actions to get in sync with the network
	entry := CreateCacheEntry(tcp_connection, node, node.NodeID)
	net_platform.Connection_caches = append(net_platform.Connection_caches, entry)

	SyncWithNetwork(net_platform)
	return true
}

func HandleConnectionInitiation(request []byte, net_platform *NetworkPlatform) {
	var payload RequestMessage
	gob.NewDecoder(bytes.NewBuffer(request)).Decode(&payload)
	send_payload := NodesMessage{
		AddrFrom: net_platform.GetNodeAddress(),
		Nodes:    []NetworkNode{*net_platform.Self_node},
	}

	sendDataToAddress(payload.AddrFrom, append(CommandStringToBytes("node"), GobEncode(send_payload)...), net_platform)
	SendGetNode(payload.AddrFrom, net_platform) // request complete node information fo
}
