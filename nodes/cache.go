package nodes

import (
	"GuthiNetwork/events"
	"net"
	"time"
)

// Implement a event queue here
type Queue struct {
	capacity uint
	len      uint
	data     []events.Events
}

// single cache entry
type CacheEntry struct {
	Connection *net.TCPConn
	// the network nodes are stored in array statically, so using ID as ref
	node_ref    *NetworkNode
	node_ref_id uint64
	time        time.Time // timestamp for when the cache was written
}

func CreateCacheEntry(connection *net.TCPConn, node_ref *NetworkNode, node_red_id uint64) CacheEntry {
	cache_entry := CacheEntry{
		Connection:  connection,
		node_ref:    node_ref,
		node_ref_id: node_red_id,
	}

	return cache_entry
}

func (cache_entry *CacheEntry) GetNodeRef() *NetworkNode {
	return cache_entry.node_ref
}

type NodeConnectionCache struct {
	Cache []CacheEntry
	// Need some reference to the network node
}
