package shm

const (
	IPC_RMID   int = 0
	SETVAL     int = 16
	GETVAL     int = 12
	SEM_UNDO   int = 0x1000
	IPC_CREAT      = 01000 /* Create key if key does not exist. */
	IPC_EXCL       = 02000 /* Fail if key exists.  */
	IPC_NOWAIT     = 04000 /* Return error on wait.  */
)
