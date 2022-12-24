package platform

import (
	"GuthiNetwork/lib"
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
)

/*
-----------------------------------------------------------------------------------------------
-----------------------------------Network Part of symbol table--------------------------------
-----------------------------------------------------------------------------------------------
*/

type VariableInfo struct {
	AddrFrom string
	Value    lib.Variable
}

type TableInfo struct {
	AddrFrom string
	Table    lib.SymbolTable
}

func SendVariableToNodes(value lib.Variable, net_platfrom *NetworkPlatform) error {
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

func HandleReceiveVariable(request []byte, net_platform *NetworkPlatform) error {
	var payload VariableInfo
	err := gob.NewDecoder(bytes.NewBuffer(request)).Decode(&payload)
	if err != nil {
		return errors.New(fmt.Sprintf("Gob decoder error:%s", err))
	}

	if _, found := net_platform.symbol_table[payload.Value.Id]; found {
		if net_platform.symbol_table[payload.Value.Id].Dtype != payload.Value.Dtype {
			// Send Error to the node as differnet data type
		}
		if !net_platform.symbol_table[payload.Value.Id].Timestamp.Before(payload.Value.Timestamp) {
			err = SendVariableToNodes(net_platform.symbol_table[payload.Value.Id], net_platform)
		}
	}
	err = lib.CreateOrSetValue(payload.Value.Id, payload.Value.Data, &net_platform.symbol_table)
	return nil
}

func SendTableToNodes(net_platform *NetworkPlatform, address string) error {
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
