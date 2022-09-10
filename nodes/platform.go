package nodes

import (
	"sync"
)

type NetworkPlatform struct {
	// Well, there's just a single writer but multiple readers. So RWMutex sounds better choice
	Lock              sync.Mutex
	Self_node         NetworkNode
	Connected_nodes   []NetworkNode
	Connection_caches []CacheEntry
}

func (network_platform *NetworkPlatform) RemoveNode(node NetworkNode) {
	new_arr := make([]NetworkNode, len(network_platform.Connected_nodes))
	j := 0

	for _, elem := range net_platform.Connected_nodes {
		if elem != node {
			new_arr[j] = elem
			j++
		}
	}

	network_platform.Connected_nodes = new_arr
}
