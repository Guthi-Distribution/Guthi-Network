package platform

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"log"
	"os"
	"reflect"
	"runtime"
	"strings"
	"sync"
	"time"
)

// global store/map for storing function
var globalFuncStore = make(map[string]interface{}, 100)
var pending_function_dispatch []remoteFunctionInvoke
var pending_dispatch_mutex sync.Mutex

type FunctionInformation struct {
	Key      string
	Function interface{}
}

type GobEncodedBytes []byte

type RemoteFuncReturn struct {
	CallId    int // will be the hash function
	RetsCount int
	Err       string
	Returns   []interface{}
}

type remoteFunctionInvoke struct {
	FName string
	Value interface{}
}

type function_execution_completed struct {
	AddrFrom    string
	FuncName    string
	Param       interface{}
	ReturnValue interface{}
}

func GetFunctionName(temp interface{}) string {
	strs := strings.Split((runtime.FuncForPC(reflect.ValueOf(temp).Pointer()).Name()), ".")
	return strs[len(strs)-1]
}

func (net_platform *NetworkPlatform) RegisterFunction(f interface{}) error {
	interfaceKind := reflect.ValueOf(f).Kind()
	if interfaceKind != reflect.Func {
		return errors.New("Invalid addition: only functions and their corresponding key (string identifier) can be added to store")
	}
	gob.Register(f)
	key := GetFunctionName(f)
	globalFuncStore[key] = f

	// TODO: Implement this
	// sendFunctionKeyToNodes(key, net_platform)

	return nil
}

func RetrieveFromMap(key string) interface{} {
	return globalFuncStore[key]
}

func GetFunctionSignature(fName string) string {
	f, ok := globalFuncStore[fName]
	if ok {
		return reflect.TypeOf(f).String()
	}
	return ""
}

func CallInterfaceFunction(inArgs GobEncodedBytes) GobEncodedBytes {
	remoteData := remoteFunctionInvoke{}

	err := gob.NewDecoder(bytes.NewReader(inArgs)).Decode(&remoteData)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error while encoding interface:")
		panic(err)
	}

	for k, v := range globalFuncStore {
		if k == remoteData.FName {
			fType := reflect.TypeOf(v)
			inArgsCount := fType.NumIn()
			fValue := reflect.ValueOf(v)

			if inArgsCount != 1 {
				// return error
				errReturn := RemoteFuncReturn{
					RetsCount: -1,
					Returns:   nil,
					Err:       "Can have function with a single argument",
				}

				var encoded bytes.Buffer
				err := gob.NewEncoder(&encoded).Encode(errReturn)
				if err != nil {
					panic(err)
				}
				return encoded.Bytes()
			}

			var in []reflect.Value = make([]reflect.Value, inArgsCount)

			in[0] = reflect.ValueOf(remoteData.Value)
			out := fValue.Call(in)
			var return_value interface{}
			return_value = nil
			if len(out) != 0 {
				return_value = reflect.ValueOf(out[0])
			}
			if handler, exists := network_platform.function_completed[remoteData.FName]; exists {
				handler(remoteData.FName, remoteData.Value, return_value)
			}

			out = nil
		}
	}
	return nil
}

type function_dispatch_status struct {
	Args_count      int
	Agrs            []interface{}
	dispatch_count  int
	completed_count int
}

func (net_platform *NetworkPlatform) CallFunction(func_name string, args interface{}, address string) {
	input := remoteFunctionInvoke{FName: func_name, Value: args}

	// if it is a remote addres, send function call
	if address != "" && address != network_platform.GetNodeAddress() {
		sendFunctionDispatch(func_name, args, net_platform, address)
		return
	}

	// if it is the same node then call the function natively
	for k, v := range globalFuncStore {
		if k == input.FName {
			fType := reflect.TypeOf(v)
			inArgsCount := fType.NumIn()
			fValue := reflect.ValueOf(v)

			if inArgsCount != 1 {
				// return error
				errReturn := RemoteFuncReturn{
					RetsCount: -1,
					Returns:   nil,
					Err:       "Can have function with a single argument",
				}

				var encoded bytes.Buffer
				err := gob.NewEncoder(&encoded).Encode(errReturn)
				if err != nil {
					panic(err)
				}
			}

			var in []reflect.Value = make([]reflect.Value, inArgsCount)

			in[0] = reflect.ValueOf(input.Value)
			out := fValue.Call(in)
			var return_value interface{}
			return_value = nil
			if len(out) != 0 {
				return_value = out[0].Interface()
			}
			if handler, exists := network_platform.function_completed[input.FName]; exists {
				handler(input.FName, input.Value, return_value)
			}

			dispatch_pending_call("")
		}
	}
}

/*
Dispatches a function into multiple nodes
Args:

	func_name: Function that needs to be called
	args: Argument that needs to be provided to different nodes
*/

func (net_platform *NetworkPlatform) DispatchFunction(func_name string, args []interface{}) {
	if len(args) == 0 {
		return
	}

	go net_platform.CallFunction(func_name, args[0], "")

	fmt.Println("Calling function")
	length := len(args)
	args_index := 1

	// dispatch the function call to all the connected nodes
	for index := range net_platform.Connected_nodes {
		if index >= length-1 || args_index >= length {
			break
		}
		args_index++
		net_platform.CallFunction(func_name, args[index+1], net_platform.Connected_nodes[index].GetAddressString())
	}

	// add all the pending function call

	for args_index < length {
		input := remoteFunctionInvoke{FName: func_name, Value: args[args_index]}
		pending_function_dispatch = append(pending_function_dispatch, input)
		args_index++
	}
}

func AddPendingDispatch(func_name string, param interface{}) {
	input := remoteFunctionInvoke{FName: func_name, Value: param}
	pending_dispatch_mutex.Lock()
	pending_function_dispatch = append(pending_function_dispatch, input)
	pending_dispatch_mutex.Unlock()
}

type functionDispatchInfo struct {
	AddrFrom string
	FuncName string
	Param    interface{}
}

func sendFunctionDispatch(func_name string, param interface{}, network_platform *NetworkPlatform, address string) {
	payload := functionDispatchInfo{
		AddrFrom: network_platform.GetNodeAddress(),
		FuncName: func_name,
		Param:    param,
	}
	data := append(CommandStringToBytes("function_dispatch"), GobEncode(payload)...)
	sendDataToAddress(address, data, network_platform)
}

func handleFunctionDispatch(data []byte, net_platform *NetworkPlatform) {
	var payload functionDispatchInfo
	gob.NewDecoder(bytes.NewReader(data)).Decode(&payload)
	log.Println("received function dispatch")

	for k, v := range globalFuncStore {
		if k == payload.FuncName {
			fValue := reflect.ValueOf(v)

			log.Println("Calling Function")
			var args []reflect.Value = make([]reflect.Value, 1)
			args[0] = reflect.ValueOf(payload.Param)

			out := fValue.Call(args)
			// time.Sleep(time.Second)

			var return_value interface{}
			return_value = nil
			if len(out) != 0 {
				return_value = out[0].Interface()
			}

			out = nil

			payload := function_execution_completed{
				network_platform.Self_node.GetAddressString(),
				payload.FuncName,
				args[0].Interface(),
				return_value,
			}
			data := append(CommandStringToBytes("func_completed"), GobEncode(payload)...)
			for i := range network_platform.Connected_nodes {
				sendDataToNode(&network_platform.Connected_nodes[i], data, network_platform)
			}
		}
	}
}

func handleFunctionCompletion(request []byte) {
	var payload function_execution_completed
	gob.NewDecoder(bytes.NewBuffer(request)).Decode(&payload)
	log.Printf("Received completion status\n")
	go dispatch_pending_call(payload.AddrFrom)
	if handler, exists := network_platform.function_completed[payload.FuncName]; exists {
		log.Printf("Calling handler\n")
		handler(payload.FuncName, payload.Param, payload.ReturnValue)
	}
}

var Start_time time.Time

func dispatch_pending_call(addr string) {
	net_platform := GetPlatform()

	pending_dispatch_mutex.Lock()
	length := len(pending_function_dispatch)
	if length > 0 {
		log.Println("Calling pending function")
		dispatch_info := pending_function_dispatch[0]
		pending_function_dispatch = pending_function_dispatch[1:]
		pending_dispatch_mutex.Unlock()
		net_platform.CallFunction(dispatch_info.FName, dispatch_info.Value, addr)
		return
	} else {
		log.Println("No pending function")
		log.Println(time.Now().Sub(Start_time).Seconds())
	}
	pending_dispatch_mutex.Unlock()
}
