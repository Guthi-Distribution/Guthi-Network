package platform

import (
	"GuthiNetwork/core"
	"GuthiNetwork/lib"
	"errors"
	"net"
	"time"
)

/*
Network Node, and platform struct and methods
*/
type NetworkNode struct {
	NodeID uint64 `json:"id"`
	Name   string `json:"name"`
	// TCP Addr is akin to socket. So, its only used when its listening for connection, right?
	Socket *net.TCPAddr `json:"address"`
}

func (node *NetworkNode) GetAddressString() string {
	return node.Socket.String()
}

type NetworkPlatform struct {
	// Well, there's just a single writer but multiple readers. So RWMutex sounds better choice
	Self_node          *NetworkNode `json:"self_node"`
	symbol_table       lib.SymbolTable
	listener           net.Listener
	Connected_nodes    []NetworkNode `json:"connected_nodes"` // nodes that are connected right noe
	Connection_History []string      `json:"history"`         // nodes information that are prevoisly connected
	Connection_caches  []CacheEntry  `json:"cache_entry"`
}

func CreateNetworkPlatform(name string, address string, port int) (*NetworkPlatform, error) {
	platform := &NetworkPlatform{}

	var err error
	if address == "" {
		address = GetNodeAddress()
	}
	platform.Self_node, err = CreateNetworkNode(name, address, port)
	platform.symbol_table = make(lib.SymbolTable)
	if err != nil {
		return nil, err
	}
	platform.listener, err = net.Listen("tcp", platform.Self_node.Socket.String())
	if err != nil {
		return platform, err
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

func (self *NetworkPlatform) AddToPreviousNodes(addr string) {
	for _, node := range self.Connection_History {
		if node == addr {
			return
		}
	}
	self.Connection_History = append(self.Connection_History, addr)
}

func (self *NetworkPlatform) AddNode(node NetworkNode) {
	if !self.knows(node.Socket.String()) {
		self.Connected_nodes = append(self.Connected_nodes, node)
		// when adding a node, create a cache entry too
		self.Connection_caches = append(self.Connection_caches, CacheEntry{
			&self.Connected_nodes[len(self.Connected_nodes)-1],
			node.NodeID,
			time.Now(),
			core.ProcessorStatus{},
			core.MemoryStatus{},
		})
	}
}

// TODO: Implement this for cache entry
func (self *NetworkPlatform) RemoveNodeWithAddress(addr string) {
	length := len(self.Connected_nodes)
	if length == 0 {
		return
	}
	new_arr := make([]NetworkNode, length-1)
	j := 0

	for _, elem := range self.Connected_nodes {
		if elem.Socket.String() != addr {
			if length >= j {
				new_arr = append(new_arr, elem)
			} else {
				new_arr[j] = elem
			}
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

func (self *NetworkPlatform) get_node_from_string(addr string) int {
	for i, node := range self.Connected_nodes {
		if node.Socket.String() == addr {
			return i
		}
	}
	return -1
}

func (net_platform *NetworkPlatform) CreateVariable(id string, data any) error {
	err := lib.CreateVariable(id, data, &net_platform.symbol_table)
	if err != nil {
		return err
	}

	return nil
}

func (net_platform *NetworkPlatform) CreateOrSetValue(id string, data any) error {
	err := lib.CreateOrSetValue(id, data, &net_platform.symbol_table)
	if err != nil {
		return err
	}

	return nil
}

func (net_platform *NetworkPlatform) GetValue(id string) (lib.Variable, error) {
	if _, exists := net_platform.symbol_table[id]; !exists {
		return lib.Variable{}, errors.New("Identifier not found")
	}

	return net_platform.symbol_table[id], nil
}
