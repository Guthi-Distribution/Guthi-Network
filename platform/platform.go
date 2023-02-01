package platform

import (
	"GuthiNetwork/core"
	"GuthiNetwork/lib"
	"errors"
	"net"
	"sync"
	"time"
)

/*
Site struct for suzuki kasami synchronization
*/
type siteInfo struct {
	IsExecuting      bool
	HasToken         bool
	Request_messages map[uint64]uint64
}

func (_site *siteInfo) doesHaveToken() bool {
	site_mutex.Lock()
	defer site_mutex.Unlock()
	return site.HasToken
}

func (_site *siteInfo) setHasToken(has_token bool) {
	site_mutex.Lock()
	defer site_mutex.Unlock()
	_site.HasToken = has_token
}
func (_site *siteInfo) setExecuting(is_executing bool) {
	site_mutex.Lock()
	defer site_mutex.Unlock()
	_site.IsExecuting = is_executing
}

type tokenInfo struct {
	Id             uint64 // id of the site having the token
	Waiting_queue  []uint64
	Token_sequence map[uint64]uint64
}

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
	Self_node            *NetworkNode `json:"self_node"`
	symbol_table         lib.SymbolTable
	listener             net.Listener
	Connected_nodes      []NetworkNode `json:"connected_nodes"` // nodes that are connected right noe
	Connection_History   []string      `json:"history"`         // nodes information that are prevoisly connected
	Connection_caches    []CacheEntry  `json:"cache_entry"`
	symbol_table_mutex   sync.Mutex
	code_execution_mutex sync.Mutex
}

func CreateNetworkPlatform(name string, address string, port int) (*NetworkPlatform, error) {
	platform := &NetworkPlatform{}

	var err error
	if address == "" {
		address = "127.0.0.1"
	} else if address == "f" {
		address = GetNodeAddress()
	}
	platform.Self_node, err = CreateNetworkNode(name, address, port)
	platform.symbol_table = make(lib.SymbolTable)
	platform.symbol_table_mutex = sync.Mutex{}
	platform.code_execution_mutex = sync.Mutex{}

	if err != nil {
		return nil, err
	}
	platform.listener, err = net.Listen("tcp", platform.Self_node.Socket.String())
	if err != nil {
		return platform, err
	}

	// initialize sem lock variables
	token.Token_sequence = make(map[uint64]uint64)
	site_mutex = sync.Mutex{}
	site.setHasToken(false)
	site.IsExecuting = false
	site.Request_messages = make(map[uint64]uint64)

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
	SendVariableToNodes(net_platform.symbol_table[id], net_platform)
	if err != nil {
		return err
	}

	return nil
}

func (net_platform *NetworkPlatform) CreateOrSetValue(id string, data any) error {
	err := lib.CreateOrSetValue(id, data, &net_platform.symbol_table)
	SendVariableToNodes(net_platform.symbol_table[id], net_platform)
	if err != nil {
		return err
	}

	return nil
}

func (net_platform *NetworkPlatform) SetValue(id string, _value *lib.Variable) error {
	net_platform.symbol_table_mutex.Lock()
	value := net_platform.symbol_table[id]
	value.SetVariable(_value)
	net_platform.symbol_table_mutex.Unlock()
	value.UnLock()
	// SendVariableToNodes(value, net_platform)
	sendVariableInvalidation(value, net_platform)
	return nil
}

func (net_platform *NetworkPlatform) setReceivedValue(id string, _value *lib.Variable) {
	net_platform.symbol_table_mutex.Lock()
	value := net_platform.symbol_table[id]
	value.SetVariable(_value)
	net_platform.symbol_table_mutex.Unlock()
}

func (net_platform *NetworkPlatform) SetData(id string, data interface{}) error {
	net_platform.symbol_table_mutex.Lock()
	value := net_platform.symbol_table[id]
	value.SetValue(data)
	net_platform.symbol_table_mutex.Unlock()

	sendVariableInvalidation(value, net_platform)
	return nil
}

func (net_platform *NetworkPlatform) GetValue(id string) (*lib.Variable, error) {
	net_platform.symbol_table_mutex.Lock()
	value, exists := net_platform.symbol_table[id]
	net_platform.symbol_table_mutex.Unlock()
	if !exists {
		return nil, errors.New("Variable not found")
	}
	if !value.IsValid() {
		sendGetVariable(net_platform, value)
	}

	// wait until the value is valid
	for !value.IsValid() {

	}

	return value, nil
}

/*
@internal
Don't care if the validate or not
*/
func (net_platform *NetworkPlatform) getValueInvalidated(id string) (*lib.Variable, error) {
	net_platform.symbol_table_mutex.Lock()
	value, exists := net_platform.symbol_table[id]
	net_platform.symbol_table_mutex.Unlock()
	if !exists {
		return nil, errors.New("Variable not found")
	}

	return value, nil
}

func (net_platform *NetworkPlatform) GetNodeFromId(id uint64) int16 {
	for idx, node := range net_platform.Connected_nodes {
		if node.NodeID == id {
			return int16(idx)
		}
	}

	return -1
}
