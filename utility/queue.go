package utility

import (
	"errors"
)

func Enqueue[T interface{}](queue []T, element T) []T {
	queue = append(queue, element)
	return queue
}

func TopQueue[T any](queue []T) (T, error) {
	var data T
	size := len(queue)
	if size == 0 {
		return data, errors.New("Queue is empy")
	}
	data = queue[size-1]
	return data, nil
}

func Dequeue[T interface{}](queue []T) ([]T, error) {
	if len(queue) == 0 {
		return nil, errors.New("Attempting to deque from a empty list")
	}

	if len(queue) == 1 {
		return []T{}, nil
	}

	return queue[1:], nil
}

func FindInArray[T comparable](list []T, data T) int16 {
	for idx, _data := range list {
		if _data == data {
			return int16(idx)
		}
	}

	return -1
}
