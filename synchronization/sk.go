package synchronization

import (
	"GuthiNetwork/platform"
	"log"
)

type Message int8

const (
	REQUEST_MESSAGE = 0
	REPLY_MESSAGE   = 1
)

type Token struct {
	Id             uint64 // id of the site having the token
	Waiting_queue  map[uint64]uint64
	Token_sequence map[uint64]uint64
}

type Site struct {
	IsExecuting      bool
	HasToken         bool
	request_messages map[uint64]map[uint64]uint16
}

var token Token
var site Site

func RequestToken(net_platform *platform.NetworkPlatform) {
	_, found := site.request_messages[net_platform.Self_node.NodeID]
	if !found {
		site.request_messages[net_platform.Self_node.NodeID] = make(map[uint64]uint16)
	}
	site.request_messages[net_platform.Self_node.NodeID][net_platform.Self_node.NodeID] += 1

	// TODO: Implement sending code
}

func ReceiveRequest(net_platform *platform.NetworkPlatform, node_id uint64, value uint16) {
	if site.request_messages[net_platform.Self_node.NodeID][node_id] < value {
		site.request_messages[net_platform.Self_node.NodeID][node_id] = value
	}

	if token.Id == net_platform.Self_node.NodeID && site.request_messages[net_platform.Self_node.NodeID][node_id] == uint16(token.Token_sequence[node_id]+1) {
		log.Println("Sending token")
		// TODO Send token to the requesting node
	}
}

func ReleaseToken(net_platform *platform.NetworkPlatform) {

}

func GetToken() Token {
	return token
}
