package platform

import (
	"GuthiNetwork/core"
	"GuthiNetwork/lib"
	"encoding/gob"
	"errors"
	"fmt"
	"net"
	"sync"
	"time"
)

type NodeFailureEventHandler func(*NetworkPlatform, string) // interface so that use can pass it's own structures

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
	mutex          sync.Mutex
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

	// events
	node_failure_event_handler NodeFailureEventHandler
}

var network_platform *NetworkPlatform

func CreateNetworkPlatform(name string, address string, port int) (*NetworkPlatform, error) {
	// only one struct is possible, need only one value
	if network_platform != nil {
		return network_platform, nil
	}
	network_platform = &NetworkPlatform{}

	var err error
	if address == "" {
		address = "127.0.0.1"
	} else if address == "f" {
		address = GetNodeAddress()
	}
	network_platform.Self_node, err = CreateNetworkNode(name, address, port)
	network_platform.symbol_table = make(lib.SymbolTable)
	network_platform.symbol_table_mutex = sync.Mutex{}
	network_platform.code_execution_mutex = sync.Mutex{}
	network_platform.node_failure_event_handler = nil

	if err != nil {
		return nil, err
	}
	network_platform.listener, err = net.Listen("tcp", network_platform.Self_node.Socket.String())
	if err != nil {
		return network_platform, err
	}

	// initialize sem lock variables
	token.Token_sequence = make(map[uint64]uint64)
	token.mutex = sync.Mutex{}
	site_mutex = sync.Mutex{}
	site.setHasToken(false)
	site.IsExecuting = false
	site.Request_messages = make(map[uint64]uint64)

	return network_platform, nil
}

func GetPlatform() *NetworkPlatform {
	return network_platform
}

func (net_platform *NetworkPlatform) BindNodeFailureEventHandler(handler NodeFailureEventHandler) {
	net_platform.node_failure_event_handler = handler
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
	net_platform.symbol_table_mutex.Lock()
	defer net_platform.symbol_table_mutex.Unlock()
	err := lib.CreateVariable(id, data, &net_platform.symbol_table)
	SendVariableToNodes(net_platform.symbol_table[id], net_platform)
	if err != nil {
		return err
	}

	return nil
}

func (net_platform *NetworkPlatform) CreateOrSetValue(id string, data any) error {
	net_platform.symbol_table_mutex.Lock()
	defer net_platform.symbol_table_mutex.Unlock()
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
	defer value.UnLock()
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
		return nil, errors.New(fmt.Sprintf("Variable %s not found", id))
	}
	if !value.IsValid() {
		sendGetVariable(net_platform, value)
	}

	// wait until the value is valid
	for !value.IsValid() {

	}

	return value, nil
}

func (net_platform *NetworkPlatform) GetData(id string) (interface{}, error) {
	net_platform.symbol_table_mutex.Lock()
	value, exists := net_platform.symbol_table[id]
	net_platform.symbol_table_mutex.Unlock()
	if !exists {
		return nil, errors.New(fmt.Sprintf("Variable %s not found", id))
	}
	if !value.IsValid() {
		sendGetVariable(net_platform, value)
	}

	// wait until the value is valid
	for !value.IsValid() {

	}

	return value.GetData(), nil
}

/*
@internal
Don't care if the validate or not
*/
func (net_platform *NetworkPlatform) getValueInvalidated(id string) (*lib.Variable, error) {
	net_platform.symbol_table_mutex.Lock()
	defer net_platform.symbol_table_mutex.Unlock()

	//FIXME: CRD
	value, exists := net_platform.symbol_table[id]
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

type Array struct {
	Size int
}

func get_array_id(id string, index int) string {
	return fmt.Sprintf("__%s-%d", id, index)
}

func (net_platform *NetworkPlatform) CreateArray(id string, size int, data interface{}) error {
	if size < 0 {
		return errors.New("Size cannot be less than 0")
	}
	var err error
	net_platform.symbol_table_mutex.Lock()
	_, exists := net_platform.symbol_table[id]
	// if the array exists, then it is being created by another node
	// FIXME: Do something else
	if exists {
		return nil
	}
	net_platform.symbol_table_mutex.Unlock()

	var arr Array
	arr.Size = size
	net_platform.CreateVariable(id, arr)

	//TODO: Multithread this piece of shit
	for i := 0; i < size; i++ {
		err = net_platform.CreateOrSetValue(get_array_id(id, i), data)
		if err != nil {
			// TODO: Cleanup all the created array
			return err
		}
	}

	return nil
}

func (net_platform *NetworkPlatform) SetDataArray(id string, index int, data interface{}) error {
	value, err := net_platform.GetValue(id)
	if err != nil {
		return errors.New("Variable does not exist")
	}

	array := value.GetData().(Array)
	if array.Size <= index {
		err_str := fmt.Sprintf("Index out of range, size: %d, index: %d", array.Size, index)
		return errors.New(err_str)
	}
	value.SetValue(data)

	sendVariableInvalidation(value, net_platform)
	return nil
}

func (net_platform *NetworkPlatform) GetValueArray(id string, index int) (*lib.Variable, error) {
	value, err := net_platform.GetValue(id)
	if err != nil {
		return nil, errors.New("Variable does not exist")
	}

	array := value.GetData().(Array)
	if array.Size <= index {
		err_str := fmt.Sprintf("Index out of range, size: %d, index: %d", array.Size, index)
		return nil, errors.New(err_str)
	}
	value, err = net_platform.GetValue(get_array_id(id, index))
	if err != nil {
		return nil, err
	}
	return value, nil
}

func (net_platform *NetworkPlatform) GetDataArray(id string, index int) (interface{}, error) {
	value, err := net_platform.GetValue(id)
	if err != nil {
		return nil, errors.New("Variable does not exist")
	}

	array := value.GetData().(Array)
	if array.Size <= index {
		err_str := fmt.Sprintf("Index out of range, size: %d, index: %d", array.Size, index)
		return nil, errors.New(err_str)
	}
	value, err = net_platform.GetValue(get_array_id(id, index))
	if err != nil {
		return nil, err
	}
	return value.GetData(), nil
}

func init() {
	gob.Register(Array{})
}
