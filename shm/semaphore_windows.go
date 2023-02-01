package shm

/*
	Semaphore operation:
		Lock and unlock
		Semaphore is created by C++ Process and opened by Golang process
*/

/*
#include <Windows.h>
#include <strsafe.h>

char* print_error(const char* message) {
    char *error_message_buff;
    DWORD dw = GetLastError();

    FormatMessage(
        FORMAT_MESSAGE_ALLOCATE_BUFFER |
        FORMAT_MESSAGE_FROM_SYSTEM |
        FORMAT_MESSAGE_IGNORE_INSERTS,
        message,
        dw,
        MAKELANGID(LANG_NEUTRAL, SUBLANG_DEFAULT),
        (LPTSTR)&error_message_buff,
        0, NULL);
	char *error = (char *)malloc(1024);
	sprintf_s(error, 1024, "%s:%s", message, error_message_buff);
	LocalFree(error_message_buff);
	return error;
}

const char* sem_name = "guthi_semaphore";
typedef struct Semaphore {
    void* semHnd;
	char *err;
} Semaphore;

Semaphore* CreateSem() {
	Semaphore *s = (Semaphore *)malloc(sizeof(Semaphore));
	s->semHnd = OpenSemaphore(SEMAPHORE_ALL_ACCESS, 0, sem_name);;
	s->err = NULL;
	if (s->semHnd == NULL) {
		s->err = print_error("Semaphore creation error");
		return s;
	}

	return s;
}

char* lock(void* semHnd, int block_call) {
	DWORD timeout_time = block_call ? INFINITE : 0;
	DWORD wait_result;
	char *err = NULL;
	wait_result = WaitForSingleObject(semHnd, timeout_time);

	if (wait_result == WAIT_OBJECT_0) {
		return 0;
	}
	else if (wait_result == WAIT_TIMEOUT) {
		err = print_error("Semaphore wait timeout");
		return err;
	}
	else if (wait_result == WAIT_FAILED) {
		err = print_error("Semaphore lock error, waiting failed");
		return err;
	}

	return err;
}

char* unlock(void* semHnd) {
	int success = ReleaseSemaphore(semHnd, 1, NULL);
	if (success) {
		return NULL;
	}
	return print_error("Semaphore Release error");
}

void CloseSemaphore(void* semHnd) {
	CloseHandle(semHnd);
}

*/
import "C"
import (
	"errors"
	"unsafe"
)

const (
	semaphore_name = "guthi_semaphore"
)

type Semaphore struct {
	hnd unsafe.Pointer
}

func CreateSemaphore() (*Semaphore, error) {
	sem := &Semaphore{}
	c_sem := C.CreateSem()
	err := C.GoString(c_sem.err)
	if c_sem == nil || c_sem.semHnd == nil {
		return nil, errors.New(err)
	}
	sem.hnd = unsafe.Pointer(c_sem.semHnd)
	return sem, nil
}

func (s *Semaphore) RemoveSemaphore() error {
	C.CloseSemaphore(s.hnd)
	return nil
}

func (s *Semaphore) Lock() error {
	c_err := C.lock(s.hnd, 1)
	err := C.GoString((*C.char)(c_err))
	if len(err) != 0 {
		return errors.New(err)
	}
	return nil
}

func (s *Semaphore) Unlock() error {
	c_err := C.unlock(s.hnd)
	err := C.GoString((*C.char)(c_err))
	if len(err) != 0 {
		return errors.New(err)
	}
	return nil
}
