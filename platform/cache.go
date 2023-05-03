package platform

import (
	"time"

	"github.com/Guthi/guthi_network/core"
)

// Implement a event queue here
type Queue struct {
	capacity uint
	len      uint
	// data     []events.Events
}

// single cache entry
type CacheEntry struct {
	// the network nodes are stored in array statically, so using ID as ref
	Node_ref    *NetworkNode `json:"self_node"`
	node_ref_id uint64
	time        time.Time // timestamp for when the cache was written
	Cpu_info    core.ProcessorStatus
	Memory_info core.MemoryStatus
}

func CreateCacheEntry(node_ref *NetworkNode, node_ref_id uint64) CacheEntry {
	cache_entry := CacheEntry{
		node_ref,
		node_ref_id,
		time.Now(), // might need to consider a distrubted time system
		core.ProcessorStatus{},
		core.MemoryStatus{},
	}
	return cache_entry
}

func (cache_entry *CacheEntry) GetNodeRef() *NetworkNode {
	return cache_entry.Node_ref
}

type NodeConnectionCache struct {
	Cache []CacheEntry
	// Need some reference to the network node
}
