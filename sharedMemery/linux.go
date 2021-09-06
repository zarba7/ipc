//go:build linux
// +build linux

package sharedMemery

import (
	"errors"
	"os"
	"syscall"
)

type Flag int

// https://github.com/torvalds/linux/blob/master/include/uapi/linux/ipc.h
const (
	/* resource get request flags */
	IPC_CREAT  Flag = 00001000 /* create if key is nonexistent */
	IPC_EXCL   Flag = 00002000 /* fail if key exists */
	IPC_NOWAIT Flag = 00004000 /* return error on wait */

	/* Permission flag for shmget.  */
	SHM_R Flag = 0400 /* or S_IRUGO from <linux/stat.h> */
	SHM_W Flag = 0200 /* or S_IWUGO from <linux/stat.h> */

	/* Flags for `shmat'.  */
	SHM_RDONLY Flag = 010000 /* attach read-only else read-write */
	SHM_RND    Flag = 020000 /* round attach address to SHMLBA */

	/* Commands for `shmctl'.  */
	SHM_REMAP Flag = 040000  /* take-over region on attach */
	SHM_EXEC  Flag = 0100000 /* execution access */

	SHM_LOCK   Flag = 11 /* lock segment (root only) */
	SHM_UNLOCK Flag = 12 /* unlock segment (root only) */
)

const (
	S_IRUSR Flag = 0400         /* Read by owner.  */
	S_IWUSR Flag = 0200         /* Write by owner.  */
	S_IRGRP      = S_IRUSR >> 3 /* Read by group.  */
	S_IWGRP      = S_IWUSR >> 3 /* Write by group.  */
)

type shm struct {
	segment
}

const permission = S_IRUSR | S_IWUSR | S_IRGRP | S_IWGRP
const flag1 = IPC_CREAT

/* The args are:
key - int, used as uniques identifier for the shared memory segment
size - uint, size in bytes to allocate
permission - int, if passed zero, 0600 will be used by default
flags - IPC_CREAT, IPC_EXCL, IPC_NOWAIT. More info can be found here https://github.com/torvalds/linux/blob/master/include/uapi/linux/ipc.h
*/

func Create(key uint, size uint16) (ISharedMemory, error) {
	// OR (bitwise) flags
	var flags = []Flag{flag1}
	var flag Flag
	for i := 0; i < len(flags); i++ {
		flag |= flags[i]
	}

	if permission != 0 {
		flag |= permission
	} else {
		flag |= S_IRUSR | S_IWUSR // default permission
	}

	// second arg could be uintptr(0) - auto
	// third arg - size
	// fourth - shmflg (flags)
	id, _, errno := syscall.RawSyscall(syscall.SYS_SHMGET, uintptr(key), uintptr(size), uintptr(flag))
	if errno != 0 {
		return nil, os.NewSyscallError("SYS_SHMGET", errno)
	}

	addr, _, errno := syscall.RawSyscall(syscall.SYS_SHMAT, id, uintptr(0), uintptr(flag))
	if errno != 0 {
		return nil, errors.New(errno.Error())
	}

	// construct slice from memory segment
	the = &shm{}
	the.init(addr, size)
	return the, nil
}

// Detach used to detach from memory segment
func (the *shm) Detach() error {
	addr = the.addr()
	_, _, errno := syscall.Syscall(syscall.SYS_SHMDT, addr, 0, 0)
	if errno != 0 {
		return errors.New(errno.Error())
	}
	return nil
}
