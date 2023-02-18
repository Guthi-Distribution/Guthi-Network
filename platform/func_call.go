package platform

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"os"
	"reflect"
	"runtime"
	"strings"
)

// global store/map for storing function
var GlobalFuncStore = make(map[string]interface{}, 100)

type GobEncodedBytes []byte

type RemoteFuncReturn struct {
	RetsCount int
	Err       string
	Returns   []interface{}
}

type RemoteFunctionInvoke struct {
	ArgsCount uint
	Values    []interface{}
}

func GetFunctionName(temp interface{}) string {
	strs := strings.Split((runtime.FuncForPC(reflect.ValueOf(temp).Pointer()).Name()), ".")
	return strs[len(strs)-1]
}

func AddToMap(f interface{}) {
	interfaceKind := reflect.ValueOf(f).Kind()
	if interfaceKind != reflect.Func {
		panic("invalid addition: only functions and their corresponding key (string identifier) can be added to store")
	}
	key := GetFunctionName(f)
	GlobalFuncStore[key] = f
}

func RetrieveFromMap(key string) interface{} {
	return GlobalFuncStore[key]
}

func GetFunctionSignature(fName string) string {
	f, ok := GlobalFuncStore[fName]
	if ok {
		return reflect.TypeOf(f).String()
	}
	return ""
}

func interfaceAddNums(a, b int) (int, string) {
	return a + b, ""
}

func CallInterfaceFunction(fName string, inArgs GobEncodedBytes) GobEncodedBytes {
	remoteData := RemoteFunctionInvoke{}

	err := gob.NewDecoder(bytes.NewReader(inArgs)).Decode(&remoteData)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error while encoding interface:")
		panic(err)
	}

	for k, v := range GlobalFuncStore {
		if k == fName {
			fType := reflect.TypeOf(v)
			inArgsCount := fType.NumIn() // in argument is just one, an interface, which the user should cast accordingly
			outArgsCount := fType.NumOut()
			fValue := reflect.ValueOf(v)

			if inArgsCount != int(remoteData.ArgsCount) {
				// return error
				errReturn := RemoteFuncReturn{
					RetsCount: -1,
					Returns:   nil,
					Err:       "provided and required inputs do not match",
				}
				var encoded bytes.Buffer
				err := gob.NewEncoder(&encoded).Encode(errReturn)
				if err != nil {
					panic(err)
				}
				return encoded.Bytes()
			}

			var in []reflect.Value = make([]reflect.Value, inArgsCount)
			var out []reflect.Value = make([]reflect.Value, outArgsCount)

			for i := 0; i < inArgsCount; i++ {
				in[i] = reflect.ValueOf(remoteData.Values[i])
			}

			out = fValue.Call(in)

			retOut := make([]interface{}, outArgsCount-1)

			for i := 0; i < outArgsCount-1; i++ {
				retOut[i] = out[i].Interface()
			}
			errReturn := out[outArgsCount-1].Interface().(string)

			r := RemoteFuncReturn{
				RetsCount: outArgsCount,
				Err:       errReturn,
				Returns:   retOut,
			}

			var encoded bytes.Buffer
			err := gob.NewEncoder(&encoded).Encode(r)
			if err != nil {
				panic(err)
			}
			return encoded.Bytes()
		}
	}
	return nil
}

func callMain() {
	AddToMap(interfaceAddNums)

	inArgs := []interface{}{1, 2}
	argsCount := len(inArgs)

	input := RemoteFunctionInvoke{ArgsCount: uint(argsCount), Values: inArgs}

	var encodedBuffer bytes.Buffer
	err := gob.NewEncoder(&encodedBuffer).Encode(input)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error while encoding interface:")
		panic(err)
	}

	var retValue RemoteFuncReturn
	nbytes := CallInterfaceFunction("interfaceAddNums", encodedBuffer.Bytes())
	err = gob.NewDecoder(bytes.NewReader(nbytes)).Decode(&retValue)
	if err != nil {
		panic(err)
	}
	fmt.Println(retValue.Returns...)
}
