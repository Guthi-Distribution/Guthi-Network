package platform

import (
	"GuthiNetwork/utility"
	"sync"
)

var site_mutex sync.Mutex
var token_mutex sync.Mutex

/*
Docs:

	https://mittaltutorials.wordpress.com/2016/01/03/exp7-suzuki-kasami-algorithm/
*/
type Message int8

const (
	REQUEST_MESSAGE = 0
	REPLY_MESSAGE   = 1
)

func CreateToken(net_platform *NetworkPlatform) tokenInfo {
	token := tokenInfo{
		Id:             net_platform.Self_node.NodeID,
		Waiting_queue:  make([]uint64, 0),
		Token_sequence: make(map[uint64]uint64),
	}

	return token
}

var token tokenInfo
var site siteInfo

func Lock(net_platform *NetworkPlatform) {
	net_platform.code_execution_mutex.Lock()
	defer site.setExecuting(true)

	_, found := site.Request_messages[net_platform.Self_node.NodeID]

	// just making sure that it exist, and does not create any runtime error
	// TODO: Maybe @ppok knows some other way
	if !found {
		site.Request_messages[net_platform.Self_node.NodeID] = 0
	}
	site.Request_messages[net_platform.Self_node.NodeID] += 1

	if site.doesHaveToken() {
		// it already has the token so just return
		return
	}

	// send request to other node
	SendTokenRequest(net_platform)
	// wait until it has the token
	for !site.doesHaveToken() {
		// TODO:
	}
	site.IsExecuting = true
}

func (net_platform *NetworkPlatform) ClaimToken() {
	token.Id = net_platform.Self_node.NodeID

	site.setHasToken(true)
}

// called by the user
func Unlock(net_platform *NetworkPlatform) {
	defer net_platform.code_execution_mutex.Unlock()
	site.IsExecuting = false
	token.Token_sequence[net_platform.Self_node.NodeID] = uint64(site.Request_messages[net_platform.Self_node.NodeID])

	if !site.doesHaveToken() {
		return
	}
	// For every site Sj whose id is not in the token queue, it appends its id to the token queue if RNi [j]=LN[j]+1.
	for key, value := range site.Request_messages {
		if ((token.Token_sequence[key] + 1) == value) && (utility.FindInArray(token.Waiting_queue, key) == -1) {
			token.Waiting_queue = utility.Enqueue(token.Waiting_queue, key)
		}
	}

	if !site.doesHaveToken() {
		return
	}

	// If the token queue is nonempty after the above update, Si deletes the top site id from the token queue and sends the token to the site indicated by the id.
	id, err := utility.TopQueue(token.Waiting_queue)
	if err != nil {
		// handle error
		return
	}
	token.Waiting_queue, _ = utility.Dequeue(token.Waiting_queue)

	if !site.doesHaveToken() {
		return
	}
	idx := net_platform.GetNodeFromId(id)
	if idx == -1 {
		// node is not connected anymore
	}

	SendToken(net_platform, net_platform.Connected_nodes[idx].GetAddressString())
}
