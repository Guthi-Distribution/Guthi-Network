package nodes

import (
	"net"
	"time"

	"GuthiNetwork/events"
)

// Implement a event queue here
type Queue struct {
	capacity uint
	len      uint
	data     []events.Events
}

// single cache entry
type CacheEntry struct {
	// the network nodes are stored in array statically, so using ID as ref
	node_ref    *NetworkNode
	node_ref_id uint64
	connection  net.Conn
	time        time.Time // timestamp for when the cache was written
}

func CreateCacheEntry(connection net.Conn, node_ref *NetworkNode, node_ref_id uint64) CacheEntry {
	cache_entry := CacheEntry{
		node_ref,
		node_ref_id,
		connection,
		time.Now(), // might need to consider a distrubted time system
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
