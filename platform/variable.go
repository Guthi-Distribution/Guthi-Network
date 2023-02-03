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
	Id       string
}

type VariableInfo struct {
	AddrFrom string
	Value    *lib.Variable
}

type TableInfo struct {
	AddrFrom string
	Table    lib.SymbolTable
}

func SendVariableToNodes(value *lib.Variable, net_platform *NetworkPlatform) error {
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

func sendGetVariable(net_platform *NetworkPlatform, value *lib.Variable) {
	payload := GetVariable{
		net_platform.GetNodeAddress(),
		value.Id,
	}

	data := append(CommandStringToBytes("get_var"), GobEncode(payload)...)
	sendDataToAddress(value.GetSourceNode(), data, net_platform)
}

func handleGetVariableRequest(request []byte, net_platform *NetworkPlatform) {
	var payload GetVariable
	gob.NewDecoder(bytes.NewReader(request)).Decode(&payload)

	value, _ := net_platform.getValueInvalidated(payload.Id)
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

/*
FIXME: This should not be occuring when the token is locked as other operation are being carried and can cause invalid or out of date data access
can cause data races
final solution can be to invalidate the data
sending data on each update can be problematic and is very resource consuming

	sendGetVariable(net_platform *NetworkPlatform, value *lib.Variable) {
		payload := GetVariable{
			net_platform.GetNodeAddress(),
			value.Id,
		}

		data := append(CommandStringToBytes("get_var"), GobEncode(payload)...)
		sendDataToAddress(value.GetSourceNode(), data, net_platform)

NOTE: Possible solution
  - invalidate the data
  - when the data is accessed, and is found to be invalidated, then the variable is requested from the node that changed the value
  - this solution can also fix the problem for out of the order data acess
  - if the data is an entire structure, then if only some value is changed, we need to send/receive entire structure (TODO: Work on this)
  - Currently, workin on primitive type only, how can we handle complex structure (Certainly json is not the way, may be Gob encode??)
    TODO: Consult punpun and Babu rd sir in this issue
*/
func HandleReceiveVariable(request []byte, net_platform *NetworkPlatform) error {
	// TODO: Maybe have another mutex, but that can lead to deadlock?
	var payload VariableInfo
	err := gob.NewDecoder(bytes.NewBuffer(request)).Decode(&payload)
	if err != nil {
		return errors.New(fmt.Sprintf("Gob decoder error:%s", err))
	}

	value, err := net_platform.getValueInvalidated(payload.Value.Id)
	value.Lock()
	defer value.UnLock()
	if err != nil {
		if value.Dtype != payload.Value.Dtype {
			// Send Error to the node as differnet data type
			log.Panic("Type mismatch for received variable")
		}
	}
	fmt.Println("Received value: ", payload.Value.Data)
	net_platform.setReceivedValue(payload.Value.Id, payload.Value)
	value.SetValid(true)
	value.SetSourceNode("")
	return nil
}

func SendTableToNode(net_platform *NetworkPlatform, address string) error {
	variables := TableInfo{
		net_platform.Self_node.GetAddressString(),
		net_platform.symbol_table,
	}
	data := GobEncode(variables)
	return sendDataToAddress(address, append(CommandStringToBytes("symbol_table"), GobEncode(data)...), net_platform)
}

func HandleReceiveSymbolTable(request []byte, net_platform *NetworkPlatform) error {
	var payload TableInfo
	err := gob.NewDecoder(bytes.NewBuffer(request)).Decode(&payload)
	if err != nil {
		return errors.New(fmt.Sprintf("Gob decoder error:%s", err))
	}

	for id, value := range payload.Table {
		if _, found := net_platform.symbol_table[id]; found {
			if net_platform.symbol_table[id].Timestamp.Before(value.Timestamp) {
				err = SendVariableToNodes(net_platform.symbol_table[id], net_platform)
				continue
			}
		}
		log.Println("Setiing validity as false")
		value.SetValid(false)
		value.SetSourceNode(payload.AddrFrom)
		net_platform.symbol_table[id] = value
	}

	return err
}

type ValidityInfo struct {
	AddrFrom string
	VarId    string
	Validity bool
}

func sendVariableInvalidation(value *lib.Variable, net_platform *NetworkPlatform) {
	value.SetValid(true)
	payload := ValidityInfo{
		net_platform.GetNodeAddress(),
		value.Id,
		false,
	}
	data := append(CommandStringToBytes("validity_info"), GobEncode(payload)...)
	for _, node := range net_platform.Connected_nodes {
		sendDataToAddress(node.GetAddressString(), data, net_platform)
	}
}

func handleVariableInvalidation(request []byte, net_platform *NetworkPlatform) {
	var payload ValidityInfo
	gob.NewDecoder(bytes.NewReader(request)).Decode(&payload)

	value, _ := net_platform.getValueInvalidated(payload.VarId)
	fmt.Println("Received Value invalidation: ", payload.VarId)
	value.SetValid(false)
	value.SetSourceNode(payload.AddrFrom)
}
