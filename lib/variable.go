package lib

import (
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
	Id      string // id is basically variable name
	Dtype   string
	IsConst bool // if it is constant, we don't need to request for it to another node, we can just retrieve iit locally

	// FIXME: may need something else, because json can be too bloated
	Data      interface{}
	Timestamp time.Time
}

func CreateVariable(id string, data any, symbol_table *SymbolTable) error {
	value := Variable{}
	if _, found := (*symbol_table)[id]; found {
		return errors.New("Variable already exist")
	}
	value.Dtype = reflect.TypeOf(data).String()
	value.Data = data
	value.IsConst = false

	value.Id = id
	value.Timestamp = time.Now()
	(*symbol_table)[id] = value

	return nil
}

func CreateConstantVariable(id string, data any, symbol_table *SymbolTable) error {
	value := Variable{}
	if _, found := (*symbol_table)[id]; found {
		return errors.New("Variable already exist")
	}

	value.Dtype = reflect.TypeOf(data).String()
	value.Data = data
	value.IsConst = true

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
	value.Dtype = reflect.TypeOf(data).String()
	value.Data = data

	value.IsConst = false

	(*symbol_table)[id] = value
	value.Timestamp = time.Now()
	return nil
}

func (value *Variable) SetValue(data any) error {
	if value.IsConst == true {
		return errors.New("Cannot Write to a constant variable")
	}
	if value.Dtype != reflect.TypeOf(data).String() {
		return errors.New("Type mismatch for previous and new value")
	}
	value.Data = data

	return nil
}

func (value *Variable) GetData() interface{} {
	return value.Data
}
