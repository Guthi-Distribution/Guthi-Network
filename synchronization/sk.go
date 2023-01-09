package synchronization

import (
	"GuthiNetwork/platform"
	"GuthiNetwork/utility"
)

type Message int8

const (
	REQUEST_MESSAGE = 0
	REPLY_MESSAGE   = 1
)

type Token platform.Token

func CreateToken(net_platform *platform.NetworkPlatform) Token {
	token := Token{
		Id:             net_platform.Self_node.NodeID,
		Waiting_queue:  make([]uint64, 0),
		Token_sequence: make(map[uint64]uint64),
	}

	return token
}

var site platform.SiteInfo

func (token *Token) Lock(net_platform *platform.NetworkPlatform) {
	requestToken(net_platform)
}

func requestToken(net_platform *platform.NetworkPlatform) {
	if site.HasToken {
		// it already has the token so do nothing
		return
	}
	_, found := site.Request_messages[net_platform.Self_node.NodeID]
	if !found {
		site.Request_messages[net_platform.Self_node.NodeID] = make(map[uint64]uint64)
	}
	site.Request_messages[net_platform.Self_node.NodeID][net_platform.Self_node.NodeID] += 1

	// TODO: Implement sending code
	// send request to other node
	platform.SendRequestToken(net_platform, site)

	// listen for receiving token
	go platform.ListenForToken(net_platform)
	// wait until it has the token
	for site.HasToken {

	}
}

// called by the user
func (token *Token) ReleaseToken(net_platform *platform.NetworkPlatform) {
	token.Token_sequence[net_platform.Self_node.NodeID] = uint64(site.Request_messages[net_platform.Self_node.NodeID][net_platform.Self_node.NodeID])
	rn_i := site.Request_messages[net_platform.Self_node.NodeID]
	for key, value := range rn_i {
		if ((token.Token_sequence[key] + 1) == value) && (utility.FindInArray(token.Waiting_queue, key) == -1) {
			token.Waiting_queue = utility.Enqueue(token.Waiting_queue, key)
		}
	}

	id, err := utility.TopQueue(token.Waiting_queue)
	if err != nil {
		// handle error
		return
	}
	utility.Dequeue(token.Waiting_queue)
	idx := net_platform.GetNodeFromId(id)
	if idx == -1 {
		// node is not connected anymore
	}
	platform.SendToken(net_platform.Connected_nodes[idx].GetAddressString())
}
