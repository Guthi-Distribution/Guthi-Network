package platform

import (
	"GuthiNetwork/lib"
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"log"
)

/*
-----------------------------------------------------------------------------------------------
-----------------------------------Network Part of symbol table--------------------------------
-----------------------------------------------------------------------------------------------
*/
type GetVariable struct {
	AddrFrom string
	Id       uint32
}

type VariableInfo struct {
	AddrFrom string
	Value    *lib.Variable
}

type ArrayInfo struct {
	AddrFrom string
	Value    []*lib.Variable
}

type IndexedArrayInfo struct {
	AddrFrom     string
	Value        []*lib.Variable
	InitialIndex int
	Count        int
}

type TableInfo struct {
	AddrFrom string
	Table    lib.SymbolTable
}

func SendVariableToNodes(value *lib.Variable, net_platform *NetworkPlatform) error {
	value.Lock()
	defer value.UnLock()
	variable := VariableInfo{
		net_platform.Self_node.GetAddressString(),
		value,
	}
	data := append(CommandStringToBytes("variable"), GobEncode(variable)...)
	var err error
	for _, node := range net_platform.Connected_nodes {
		err = sendDataToAddress(node.GetAddressString(), data, net_platform)
		if err != nil {
			return err
		}
	}

	return nil
}

func Send_array_to_nodes(id string, net_platform *NetworkPlatform) error {
	var values []*lib.Variable
	array_value, err := net_platform.GetData(id)

	if err != nil {
		panic(err)
	}

	size := array_value.(array).Size
	values = make([]*lib.Variable, size+1)
	values[0] = net_platform.symbol_table[lib.GetHashValue(id)]
	for i := 1; i <= size; i++ {
		value, err := net_platform.getValueInvalidated(lib.GetHashValue(get_array_id(id, i)))
		if err != nil {
			return err
		}
		values[i] = value
	}
	variable := ArrayInfo{
		net_platform.Self_node.GetAddressString(),
		values,
	}
	data := append(CommandStringToBytes("array"), GobEncode(variable)...)

	for _, node := range net_platform.Connected_nodes {
		err := sendDataToAddress(node.GetAddressString(), data, net_platform)
		if err != nil {
			return err
		}
	}

	return nil
}

func SendIndexedArray(id string, initial_index int, count int, net_platform *NetworkPlatform) error {
	var values []*lib.Variable
	array_value, err := net_platform.GetData(id)

	if err != nil {
		panic(err)
	}

	if initial_index+count > array_value.(array).Size {
		return errors.New("Index out of bounds\n")
	}

	size := count
	values = make([]*lib.Variable, size+1)
	values[0] = net_platform.symbol_table[lib.GetHashValue(id)]

	index := 1
	for i := initial_index + 1; i <= initial_index+count; i++ {
		value, err := net_platform.getValueInvalidated(lib.GetHashValue(get_array_id(id, i)))
		if err != nil {
			return err
		}
		values[index] = value
		index++
	}

	variable := IndexedArrayInfo{
		net_platform.Self_node.GetAddressString(),
		values,
		initial_index,
		count,
	}
	data := append(CommandStringToBytes("indexed_array"), GobEncode(variable)...)

	for _, node := range net_platform.Connected_nodes {
		err := sendDataToAddress(node.GetAddressString(), data, net_platform)
		if err != nil {
			return err
		}
	}

	return nil
}

func sendGetVariable(net_platform *NetworkPlatform, value_id string, address string) {
	payload := GetVariable{
		net_platform.GetNodeAddress(),
		lib.GetHashValue(value_id),
	}
	// log.Printf("Sending request of id: %s\n", value_id)
	data := append(CommandStringToBytes("get_var"), GobEncode(payload)...)
	sendDataToAddress(address, data, net_platform)
}

func handleGetVariableRequest(request []byte, net_platform *NetworkPlatform) {
	var payload GetVariable
	gob.NewDecoder(bytes.NewReader(request)).Decode(&payload)

	value, err := net_platform.getValueInvalidated(payload.Id)
	if err != nil {
		log.Printf("Id: %d does not exists\n", payload.Id)
		return
	}
	variable := VariableInfo{
		net_platform.Self_node.GetAddressString(),
		value,
	}
	value.Lock()
	defer value.UnLock()
	value.SetSourceNode(payload.AddrFrom)
	value.SetValid(false)
	data := append(CommandStringToBytes("variable"), GobEncode(variable)...)
	sendDataToAddress(payload.AddrFrom, data, net_platform)
}

func handleReceiveVariable(request []byte, net_platform *NetworkPlatform) error {
	var payload VariableInfo
	err := gob.NewDecoder(bytes.NewBuffer(request)).Decode(&payload)
	if err != nil {
		return errors.New(fmt.Sprintf("Gob decoder error:%s", err))
	}

	value, err := net_platform.getValueInvalidated(payload.Value.Id)
	if err != nil {
		return err
	}
	value.Lock()
	defer value.UnLock()
	if err != nil {
		// if value.Dtype != payload.Value.Dtype {
		// 	// Send Error to the node as differnet data type
		// 	log.Panic("Type mismatch for received variable")
		// }
	}
	// fmt.Println("Received value: ", payload.Value.Data)
	net_platform.setReceivedValue(payload.Value.Id, payload.Value)
	value.SetValid(true)
	value.SetSourceNode("")
	return nil
}

func HandleReceiveArray(request []byte, net_platform *NetworkPlatform) error {
	var payload ArrayInfo
	err := gob.NewDecoder(bytes.NewBuffer(request)).Decode(&payload)
	if err != nil {
		return errors.New(fmt.Sprintf("Gob decoder error:%s", err))
	}

	value, err := net_platform.getValueInvalidated(payload.Value[0].Id)
	if err != nil {
		return err
	}
	net_platform.setReceivedValue(payload.Value[0].Id, payload.Value[0])
	for i := 0; i < payload.Value[0].Data.(array).Size; i++ {
		value, _ := net_platform.getValueInvalidated(payload.Value[i].Id)
		if value == nil || value.GetSourceNode() == payload.AddrFrom {
			net_platform.setReceivedValue(payload.Value[i].Id, payload.Value[i])
		}
	}
	value.SetValid(true)
	value.SetSourceNode("")
	log.Println("Received entire array")
	return nil
}

func HandleReceiveIndexedArray(request []byte, net_platform *NetworkPlatform) error {
	var payload IndexedArrayInfo
	err := gob.NewDecoder(bytes.NewBuffer(request)).Decode(&payload)
	if err != nil {
		return errors.New(fmt.Sprintf("Gob decoder error:%s", err))
	}

	value, err := net_platform.getValueInvalidated(payload.Value[0].Id)
	if err != nil {
		return err
	}
	net_platform.setReceivedValue(payload.Value[0].Id, payload.Value[0])

	index := 1
	for i := payload.InitialIndex + 1; i < payload.Count+payload.InitialIndex; i++ {
		value, _ := net_platform.getValueInvalidated(payload.Value[index].Id)
		if value == nil || value.GetSourceNode() == payload.AddrFrom {
			net_platform.setReceivedValue(payload.Value[index].Id, payload.Value[index])
		}
		index++
	}
	payload.Value = nil
	value.SetValid(true)
	value.SetSourceNode("")
	log.Println("Received entire array")
	return nil
}

func SendTableToNode(net_platform *NetworkPlatform, address string) error {
	net_platform.symbol_table_mutex.RLock()
	defer net_platform.symbol_table_mutex.RUnlock()
	variables := TableInfo{
		net_platform.Self_node.GetAddressString(),
		net_platform.symbol_table,
	}
	data := GobEncode(variables)
	return sendDataToAddress(address, append(CommandStringToBytes("symbol_table"), data...), net_platform)
}

func HandleReceiveSymbolTable(request []byte, net_platform *NetworkPlatform) error {
	var payload TableInfo
	err := gob.NewDecoder(bytes.NewBuffer(request)).Decode(&payload)
	if err != nil {
		return errors.New(fmt.Sprintf("Gob decoder error:%s", err))
	}
	fmt.Printf("INFO: Table Size: %d\n", len(payload.Table))

	for id, value := range payload.Table {
		net_platform.symbol_table_mutex.RLock()
		value_node, found := net_platform.symbol_table[id]
		net_platform.symbol_table_mutex.RUnlock()
		if found {
			if value_node.Timestamp.Before(value.Timestamp) {
				err = SendVariableToNodes(value_node, net_platform)
				continue
			}
		}
		value.SetValid(true)
		// value.SetSourceNode(payload.AddrFrom)
		net_platform.symbol_table_mutex.Lock()
		net_platform.symbol_table[id] = value
		net_platform.symbol_table_mutex.Unlock()

	}

	// sendTableReceiveAcknowledgement(net_platform, payload.AddrFrom)

	return err
}

func sendTableReceiveAcknowledgement(net_platform *NetworkPlatform, address string) error {
	variables := GetInformation{
		AddrFrom: net_platform.GetNodeAddress(),
	}
	data := GobEncode(variables)
	return sendDataToAddress(address, append(CommandStringToBytes("symbol_table_ack"), data...), net_platform)
}

func handleReceiveSymbolTableAck(request []byte, net_platform *NetworkPlatform) error {
	var payload GetInformation
	err := gob.NewDecoder(bytes.NewBuffer(request)).Decode(&payload)
	if err != nil {
		return errors.New(fmt.Sprintf("Gob decoder error:%s", err))
	}
	dispatch_pending_call(payload.AddrFrom)

	return err
}

type ValidityInfo struct {
	AddrFrom string
	VarId    uint32
}

func sendVariableInvalidation(value *lib.Variable, net_platform *NetworkPlatform) interface{} {
	value.SetValid(true)
	payload := ValidityInfo{
		net_platform.GetNodeAddress(),
		value.Id,
	}
	data := append(CommandStringToBytes("validity_info"), GobEncode(payload)...)
	for i := range net_platform.Connected_nodes {
		sendDataToNode(&net_platform.Connected_nodes[i], data, net_platform)
	}
	return data
}

func handleVariableInvalidation(request []byte, net_platform *NetworkPlatform) {
	var payload ValidityInfo
	gob.NewDecoder(bytes.NewReader(request)).Decode(&payload)

	value, err := net_platform.getValueInvalidated(payload.VarId)
	if err != nil {
		// log.Println(err)
		return
	}
	value.SetValid(false)
	value.SetSourceNode(payload.AddrFrom)
}
