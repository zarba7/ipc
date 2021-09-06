package semaphore

import (
"flag"
"log"
"syscall"
"time"
"unsafe"
)

/*
#include <sys/sem.h>
typedef struct sembuf sembuf;
typedef struct _semun {
        int val;
} semun;;
*/
import "C"

func semget(key int) int {
	r1, _, _ := syscall.Syscall(syscall.sem, uintptr(key),
		uintptr(1), uintptr(00666))
	if int(r1) < 0 {
		semid, r2, err := syscall.Syscall(syscall.SYS_SEMGET, uintptr(key),
			uintptr(1), uintptr(C.IPC_CREAT|C.IPC_EXCL|00666))
		if int(semid) < 0 {
			log.Printf("error:semget error is %v\n", err)
		}
		var newInit int
		newInit = 1

		r1, r2, err := syscall.Syscall6(syscall.SYS_SEMCTL,
			uintptr(semid),
			uintptr(0),
			uintptr(C.SETVAL),
			uintptr(newInit),
			uintptr(0), uintptr(0))
		if int(r1) < 0 {
			log.Fatal("error:SYS_SEMCTL:", r1, r2, err)
		}
		return int(semid)

	} else {

		semid := r1

		return int(semid)

	}
	return int(r1)
}

func semLock(semid int) int {

	stSemBuf := C.sembuf{
		sem_num: 0,
		sem_op:  -1,
		sem_flg: C.IPC_NOWAIT | C.SEM_UNDO,
	}

	r1, r2, err := syscall.Syscall(syscall.SYS_SEMOP, uintptr(semid), uintptr(unsafe.Pointer(&stSemBuf)), 1)
	if int(r1) < 0 {
		log.Printf("error:semget error is %v,%v,%v\n", r1, r2, err)
	}
	return int(r1)
}
func semShow(semid int) int {

	r1, r2, err := syscall.Syscall(syscall.SYS_SEMCTL,
		uintptr(semid),
		uintptr(0),
		uintptr(C.GETVAL))
	if int(r1) < 0 {
		log.Printf("error:semShow error is %v,%v,%v\n", r1, r2, err)
	}
	return int(r1)
