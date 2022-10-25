package shm

const (
	IPC_RMID   int = 0
	SEM_UNDO   int = 0x1000
	IPC_CREAT      = 01000 /* Create key if key does not exist. */
	IPC_EXCL       = 02000 /* Fail if key exists.  */
	IPC_NOWAIT     = 04000 /* Return error on wait.  */
	GETPID         = 11    /* get sempid */
	GETVAL         = 12    /* get semval */
	GETALL         = 13    /* get all semval's */
	GETNCNT        = 14    /* get semncnt */
	GETZCNT        = 15    /* get semzcnt */
	SETVAL         = 16    /* set semval */
	SETALL         = 17    /* set all semval's */
)
