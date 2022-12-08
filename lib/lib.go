package lib

/*
Variable that is communicated in distributed network
*/
type Variable struct {
	dtype uint8
	data  string // json representation of the string
}

// type Variable interface {
// 	GetName() string
// 	GetId() string

// }
