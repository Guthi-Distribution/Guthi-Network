package shm

import (
	"syscall"
	"unsafe"
)

const (
	key     = 69
	sem_num = 0
)

type Semaphore struct {
	semid int
	nsems int
}

type semop struct {
	semNum  uint16
	semOp   int16
	semFlag int16
}

func errnoErr(errno syscall.Errno) error {
	switch errno {
	case syscall.Errno(0):
		return nil
	default:
		return error(errno)
	}
}

func SemGet(key int, nsems int, flags int) (*Semaphore, error) {
	r1, _, errno := syscall.Syscall(syscall.SYS_SEMGET,
		uintptr(key), uintptr(nsems), uintptr(flags))
	if errno == syscall.Errno(0) {
		return &Semaphore{semid: int(r1), nsems: nsems}, nil
	} else {
		return nil, errnoErr(errno)
	}
}

func CreateSemaphore() (*Semaphore, error) {
	sem, err := SemGet(key, 1, IPC_CREAT|(syscall.S_IRUSR|syscall.S_IWUSR|syscall.S_IRGRP|syscall.S_IWGRP))
	if err != nil {
		return nil, err
	}
	return sem, nil
}

func (s *Semaphore) RemoveSemaphore() error {
	_, _, errno := syscall.Syscall(syscall.SYS_SEMCTL, uintptr(s.semid),
		uintptr(0), uintptr(IPC_RMID))
	return errnoErr(errno)
}

func (s *Semaphore) GetVal(semNum int) (int, error) {
	val, _, errno := syscall.Syscall(syscall.SYS_SEMCTL, uintptr(s.semid),
		uintptr(semNum), uintptr(GETVAL))
	return int(val), errnoErr(errno)
}

func (s *Semaphore) Unlock() error {
	post := semop{semNum: uint16(0), semOp: 1, semFlag: 0x1000}
	_, _, errno := syscall.Syscall(syscall.SYS_SEMOP, uintptr(s.semid),
		uintptr(unsafe.Pointer(&post)), uintptr(s.nsems))
	return errnoErr(errno)
}

func (s *Semaphore) Lock() error {
	wait := semop{semNum: uint16(0), semOp: -1, semFlag: 0x1000}
	_, _, errno := syscall.Syscall(syscall.SYS_SEMOP, uintptr(s.semid),
		uintptr(unsafe.Pointer(&wait)), uintptr(s.nsems))
	return errnoErr(errno)
}
