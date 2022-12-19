package lib

import (
	"bytes"
	"encoding/json"
	"errors"
	"reflect"
	"time"
)

// type uint8 types
type SymbolTable map[string]Variable

var symbol_table = make(map[string]Variable)

/*
Variable that is communicated in distributed network
*/
type Variable struct {
	id    string // id is basically variable name
	dtype reflect.Type

	// FIXME: may need something else, because json can be too bloated
	data      string // json representation of the string
	timestamp time.Time
}

// FIXME: do something with ide here
func CreateVariable(id string, data any) (Variable, error) {
	value := Variable{}
	if _, found := symbol_table[id]; found {
		return value, errors.New("Variable already exist")
	}
	value.dtype = reflect.TypeOf(data)
	buff := bytes.NewBufferString(value.data)
	encoder := json.NewEncoder(buff)
	encoder.Encode(data)
	value.data = buff.String()
	value.id = id
	symbol_table[id] = value
	//TODO: Implement this
	// platform.SendVariableToNodes()
	return value, nil
}

func CreateOrSetValue(id string, data any) (Variable, error) {
	value := Variable{}
	if _, exists := symbol_table[id]; exists {
		if reflect.TypeOf(symbol_table[id].data) != reflect.TypeOf(data) {
			return value, errors.New("Type mismatch for previous and new value")
		}
	}
	value.dtype = reflect.TypeOf(data)
	buff := bytes.NewBufferString(value.data)
	encoder := json.NewEncoder(buff)
	encoder.Encode(data)
	value.data = buff.String()
	symbol_table[id] = value
	return value, nil
}

func (value *Variable) SetValue(data any) error {
	if value.dtype != reflect.TypeOf(data) {
		return errors.New("Type mismatch for previous and new value")
	}
	buff := bytes.NewBufferString(value.data)
	encoder := json.NewEncoder(buff)
	encoder.Encode(data)
	value.data = buff.String()
	return nil
}

func (value *Variable) GetValue() any {
	return value.data
}

func GetSymbolTable() SymbolTable {
	return symbol_table
}
