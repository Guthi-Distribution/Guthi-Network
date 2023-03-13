package lib

import (
	"errors"
	"hash/fnv"
	"reflect"
	"sync"
	"time"
)

type State int
type Type int

const (
	Modified  State = 0
	Shared    State = 1
	Exclusive State = 2
	Invalid   State = 3
)

type SymbolTable map[uint32]*Variable

/*
Variable that is communicated in distributed network
//TODO: Add some metadata so that it can be extended further
  - one can be the owner of the variable,
  - if it is the owner then the owner can change once the owner has failed, differnet failure handling can be implemented for this case too
  - ughhh to much work
*/
type Variable struct {
	Id      uint32 // id is basically variable name
	Dtype   string
	IsConst bool // if it is constant, we don't need to request for it to another node, we can just retrieve iit locally

	Data        interface{}
	Timestamp   time.Time
	mutex       sync.Mutex // for locally accessing variable by multiple goroutine @internallly reading and writng
	access_lock sync.Mutex // to prevent race condition for receiveing the value and exclusive region

	is_valid    bool
	source_node string // ip of the source node, is acessed only when is_vallid is true
}

func CreateVariable(id string, data any, symbol_table *SymbolTable) error {
	value := Variable{}
	hash_value := GetHashValue(id)
	if _, found := (*symbol_table)[hash_value]; found {
		return errors.New("Variable already exist")
	}
	value.Dtype = reflect.TypeOf(data).String()
	value.Data = data
	value.IsConst = false
	value.is_valid = true

	value.Id = hash_value
	value.Timestamp = time.Now()
	(*symbol_table)[value.Id] = &value

	return nil
}

func CreateConstantVariable(id string, data any, symbol_table *SymbolTable) error {
	value := Variable{}
	hash_value := GetHashValue(id)
	if _, found := (*symbol_table)[hash_value]; found {
		return errors.New("Variable already exist")
	}

	value.Dtype = reflect.TypeOf(data).String()
	value.Data = data
	value.IsConst = true
	value.is_valid = true

	value.Id = hash_value
	value.Timestamp = time.Now()
	(*symbol_table)[value.Id] = &value

	return nil
}

func CreateOrSetValue(id string, data any, symbol_table *SymbolTable) error {
	hash_value := GetHashValue(id)
	value := Variable{}
	if variable, exists := (*symbol_table)[hash_value]; exists {
		if reflect.TypeOf((*symbol_table)[hash_value].Data) != reflect.TypeOf(data) {
			return errors.New("Type mismatch for previous and new value")
		}

		variable.SetValue(data)
		return nil
	}
	value.Dtype = reflect.TypeOf(data).String()
	value.Data = data
	value.is_valid = true
	value.Id = hash_value
	value.IsConst = false

	(*symbol_table)[hash_value] = &value
	// value.Timestamp = time.Now()
	return nil
}

func (value *Variable) SetValue(data any) error {
	value.mutex.Lock()
	defer value.mutex.Unlock()
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
	value.mutex.Lock()
	defer value.mutex.Unlock()
	return value.Data
}

func (value *Variable) SetVariable(variable *Variable) {
	value.mutex.Lock()
	defer value.mutex.Unlock()
	value.Data = variable.Data
	value.Timestamp = variable.Timestamp
	// value.Dtype = variable.Dtype
	value.IsConst = variable.IsConst
}

func (value *Variable) IsValid() bool {
	return value.is_valid
}

func (value *Variable) Lock() {
	value.access_lock.Lock()
}

func (value *Variable) UnLock() {
	value.access_lock.Unlock()
}

func (value *Variable) SetValid(validity bool) {
	value.is_valid = validity
}

func (value *Variable) GetSourceNode() string {
	return value.source_node
}

func (value *Variable) SetSourceNode(source string) {
	value.source_node = source
}

func GetHashValue(id string) uint32 {
	h := fnv.New32()
	h.Write([]byte(id))
	return h.Sum32()
}
