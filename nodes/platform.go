package nodes

import (
	"fmt"
	"sync"
)

type NetworkPlatform struct {
	// Well, there's just a single writer but multiple readers. So RWMutex sounds better choice
	Lock              sync.Mutex
	Self_node         *NetworkNode
	Connected_nodes   []NetworkNode
	Connection_caches []CacheEntry
}

func (self *NetworkPlatform) RemoveNode(node NetworkNode) {
	new_arr := make([]NetworkNode, len(self.Connected_nodes))
	j := 0

	for _, elem := range self.Connected_nodes {
		if elem != node {
			new_arr[j] = elem
			j++
		}
	}

	self.Connected_nodes = new_arr
}

func (self *NetworkPlatform) GetNodeAddress() string {
	fmt.Println(self.Self_node.Socket.String())
	return self.Self_node.Socket.String()
}
