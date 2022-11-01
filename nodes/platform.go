package nodes

import (
	"sync"
)

type NetworkPlatform struct {
	// Well, there's just a single writer but multiple readers. So RWMutex sounds better choice
	Lock              sync.Mutex
	Self_node         *NetworkNode
	Connected_nodes   []NetworkNode
	Connection_caches []CacheEntry
}

func CreateNetworkPlatform(name string, address string, port int) (*NetworkPlatform, error) {
	platform := &NetworkPlatform{}

	var err error
	platform.Self_node, err = CreateNetworkNode(name, address, port)
	if err != nil {
		return nil, err
	}

	return platform, nil
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

func (self *NetworkPlatform) RemoveNodeWithAddress(addr string) {
	new_arr := make([]NetworkNode, len(self.Connected_nodes))
	j := 0

	for _, elem := range self.Connected_nodes {
		if elem.Socket.String() != addr {
			new_arr[j] = elem
			j++
		}
	}

	self.Connected_nodes = new_arr
}

func (self *NetworkPlatform) GetNodeAddress() string {
	return self.Self_node.Socket.String()
}

// see if the node knows a node with address
func (self *NetworkPlatform) knows(addr string) bool {
	for _, node := range self.Connected_nodes {
		if node.Socket.String() == addr {
			return true
		}
	}
	return false
}
