package shm

import (
	"golang.org/x/sys/windows"
)

const (
	semaphore_name = "guthi_semaphore"
)

type Semaphore struct {
	hnd windows.Handle
}

func CreateSemaphore() (*Semaphore, error) {
	return nil, nil
}

func (s *Semaphore) RemoveSemaphore() error {
	return nil
}

/*
	@Param semNum: semaphore number, just to make compatible with linux
*/
func (s *Semaphore) Lock() error {
	return nil
}
