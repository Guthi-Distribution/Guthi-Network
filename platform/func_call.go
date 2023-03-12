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
)

// global store/map for storing function
var globalFuncStore = make(map[string]interface{}, 100)

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
	AddrFrom string
	FuncName string
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
			// outArgsCount := fType.NumOut()
			fValue := reflect.ValueOf(v)

			if inArgsCount != 1 {
				// return error
				errReturn := RemoteFuncReturn{
					RetsCount: -1,
					Returns:   nil,
					Err:       "Can have function with a single argument",
				}
				//TODO: Add this to platform
				var encoded bytes.Buffer
				err := gob.NewEncoder(&encoded).Encode(errReturn)
				if err != nil {
					panic(err)
				}
				return encoded.Bytes()
			}

			var in []reflect.Value = make([]reflect.Value, inArgsCount)
			// var out []reflect.Value = make([]reflect.Value, outArgsCount)

			in[0] = reflect.ValueOf(remoteData.Value)
			fValue.Call(in)
			if handler, exists := network_platform.function_completed[remoteData.FName]; exists {
				handler()
			}
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

// call id to functionDispatchInfo
// var function_dispatch_status map[int]functionDispatchInfo

/*
TODO:
  - Add node wise parameter binding, along with main node (just like master thread)
  - Add feature for returing function
  - Maybe, add each callback function for every function call
  - Add CallId, so that we can even cache function call with same signature
*/
func (net_platform *NetworkPlatform) CallFunction(func_name string, args interface{}, address string) {
	input := remoteFunctionInvoke{FName: func_name, Value: args}

	if address != "" && address != network_platform.GetNodeAddress() {
		sendFunctionDispatch(func_name, args, net_platform, address)
		return
	}
	var encodedBuffer bytes.Buffer
	err := gob.NewEncoder(&encodedBuffer).Encode(input)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error while encoding interface:")
		panic(err)
	}

	// TODO: Fix this execution issue
	// var retValue RemoteFuncReturn
	CallInterfaceFunction(encodedBuffer.Bytes())
}

/*
TODO: Add pending function dispatch
Dispatches a function into multiple nodes
Args:

	func_name: Function that needs to be called
	args: Argument that needs to be provided to different nodes
*/
func (net_platform *NetworkPlatform) DispatchFunction(func_name string, args []interface{}) {
	input := remoteFunctionInvoke{FName: GetFunctionName(func_name), Value: args}
	if len(args) == 0 {
		return
	}

	var encodedBuffer bytes.Buffer
	err := gob.NewEncoder(&encodedBuffer).Encode(input)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error while encoding interface:")
		panic(err)
	}
	// TODO: Fix this execution issue
	// var retValue RemoteFuncReturn
	CallInterfaceFunction(encodedBuffer.Bytes())
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

	for k, v := range globalFuncStore {
		if k == payload.FuncName {
			fValue := reflect.ValueOf(v)

			var args []reflect.Value = make([]reflect.Value, 1)
			args[0] = reflect.ValueOf(payload.Param)
			fValue.Call(args)

			payload := function_execution_completed{
				network_platform.Self_node.GetAddressString(),
				payload.FuncName,
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
	if handler, exists := network_platform.function_completed[payload.FuncName]; exists {
		log.Printf("Calling handler\n")
		handler()
	}
}
