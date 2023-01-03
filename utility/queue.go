package utility

import "errors"

func enqueue[T interface{}](queue []T, element T) []T {
	queue = append(queue, element)
	return queue
}

func dequeue[T interface{}](queue []T) ([]T, error) {
	if len(queue) == 0 {
		return nil, errors.New("Attempting to deque from a empty list")
	}

	if len(queue) == 1 {
		return []T{}, nil
	}

	return queue[1:], nil
}

func find[T comparable](list []T, data T) int16 {
	for idx, _data := range list {
		if _data == data {
			return int16(idx)
		}
	}

	return -1
}
