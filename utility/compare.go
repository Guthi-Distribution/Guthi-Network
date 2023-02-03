package utility

// copied from https://pkg.go.dev/golang.org/x/exp/constraints
type Ordered interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr | ~float32 | ~float64 | ~string
}

func Max[T Ordered](elem_1 T, elem_2 T) T {
	if elem_1 > elem_2 {
		return elem_1
	} else {
		return elem_2
	}
}
