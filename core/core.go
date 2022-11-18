package core

import "GuthiNetwork/shm"

//TODO: Need better name
/*
	Initializes the the basic structure needed for core communication
*/
func Initialize() error {
	var err error
	shared_memory, err = shm.CreateSharedMemory()
	if err != nil {
		return err
	}
	semaphore, err = shm.CreateSemaphore()
	if err != nil {
		return err
	}

	return nil
}
