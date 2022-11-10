package core

/*
#cgo CXXFLAGS: "-std=c++20"
#cgo LDFLAGS: -L../Guthi-Core/ -lGuthiCore -lstdc++
#include "../Guthi-Core/src/core/c_api.h"
*/
import "C"
import "fmt"

func Initialize() interface{} {
	a := C.InitFileSystem()
	fmt.Printf("%T\n", a)
	C.PrettyPrintFileSystem()
	return a
}
