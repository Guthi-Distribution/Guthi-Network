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

/*
Variable that is communicated in distributed network
*/
type Variable struct {
	Id    string // id is basically variable name
	Dtype reflect.Type

	// FIXME: may need something else, because json can be too bloated
	Data      string // json representation of the string
	Timestamp time.Time
}

// FIXME: do something with id here
func CreateVariable(id string, data any, symbol_table *SymbolTable) error {
	value := Variable{}
	if _, found := (*symbol_table)[id]; found {
		return errors.New("Variable already exist")
	}

	value.Dtype = reflect.TypeOf(data)

	buff := bytes.NewBufferString(value.Data)
	encoder := json.NewEncoder(buff)
	err := encoder.Encode(data)
	if err != nil {
		return err
	}
	value.Data = buff.String()

	value.Id = id
	value.Timestamp = time.Now()
	(*symbol_table)[id] = value

	return nil
}

func CreateOrSetValue(id string, data any, symbol_table *SymbolTable) error {
	value := Variable{}
	if _, exists := (*symbol_table)[id]; exists {
		if reflect.TypeOf((*symbol_table)[id].Data) != reflect.TypeOf(data) {
			return errors.New("Type mismatch for previous and new value")
		}
	}
	value.Dtype = reflect.TypeOf(data)
	buff := bytes.NewBufferString(value.Data)
	encoder := json.NewEncoder(buff)
	err := encoder.Encode(data)
	if err != nil {
		return err
	}
	value.Data = buff.String()
	(*symbol_table)[id] = value
	value.Timestamp = time.Now()
	return nil
}

func (value *Variable) SetValue(data any) error {
	if value.Dtype != reflect.TypeOf(data) {
		return errors.New("Type mismatch for previous and new value")
	}
	buff := bytes.NewBufferString(value.Data)
	encoder := json.NewEncoder(buff)
	encoder.Encode(data)
	value.Data = buff.String()
	return nil
}

func (value *Variable) GetValue() any {
	return value.Data
}
