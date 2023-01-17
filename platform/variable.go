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

type VariableInfo struct {
	AddrFrom string
	Value    *lib.Variable
}

type TableInfo struct {
	AddrFrom string
	Table    lib.SymbolTable
}

func SendVariableToNodes(value *lib.Variable, net_platfrom *NetworkPlatform) error {
	variable := VariableInfo{
		net_platfrom.Self_node.GetAddressString(),
		value,
	}
	data := append(CommandStringToBytes("variable"), GobEncode(variable)...)
	var err error
	for _, node := range net_platfrom.Connected_nodes {
		err = sendDataToAddress(node.GetAddressString(), data, net_platfrom)
		if err != nil {
			return err
		}
	}

	return nil
}

/*
FIXME: This should not be occuring when the token is locked as other operation are being carried and can cause invalid or out of date data access
can cause data races
final solution can be to invalidate the data
sending data on each update can be problematic and is very resource consuming

NOTE: Possible solution
  - invalidate the data
  - when the data is accessed, and is found to be invalidated, then the variable is requested from the node that changed the value
  - this solution can also fix the problem for out of the order data acess
  - if the data is an entire structure, then if only some value is changed, we need to send/receive entire structure (TODO: Work on this)
  - Currently, workin on primitive type only, how can we handle complex structure (Certainly json is not the way, may be Gob encode??)
    TODO: Consult punpun and Babu rd sir in this issue
*/
func HandleReceiveVariable(request []byte, net_platform *NetworkPlatform) error {
	// LOG: Test this stuff maynot be valid
	// TODO: Maybe have another mutex, but that can lead to deadlock?
	for site.IsExecuting {

	}
	var payload VariableInfo
	err := gob.NewDecoder(bytes.NewBuffer(request)).Decode(&payload)
	if err != nil {
		return errors.New(fmt.Sprintf("Gob decoder error:%s", err))
	}

	value, err := net_platform.GetValue(payload.Value.Id)
	if err != nil {
		if value.Dtype != payload.Value.Dtype {
			// Send Error to the node as differnet data type
			log.Panic("Type mismatch for received variable")
		}
	}
	fmt.Println("Received value: ", payload.Value.Data)
	net_platform.setReceivedValue(payload.Value.Id, payload.Value)
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
			if !net_platform.symbol_table[id].Timestamp.Before(value.Timestamp) {
				err = SendVariableToNodes(net_platform.symbol_table[id], net_platform)
			}
		} else {
			net_platform.symbol_table[id] = value
		}
	}

	return err
}
