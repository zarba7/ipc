//go:build windows
// +build windows

package sharedMemery

import (
	"fmt"
	"os"
	"syscall"
)

// https://docs.microsoft.com/en-us/windows/win32/memory/creating-named-shared-memory
// Create used to create or open existing
func Create(key uint, size uint16) (ISharedMemory, error) {
	name, err := syscall.UTF16PtrFromString(fmt.Sprint(key))
	if err != nil {
		return nil, err
	}

	// args are:
	// syscall.InvalidHandle - use paging file
	// nil - default security
	// syscall.PAGE_READWRITE - read/write access
	// 0 - maximum object size (high-order DWORD)
	// uint32(size) - maximum object size (low-order DWORD)
	// name - name of mapping object
	hMapFile, err := syscall.CreateFileMapping(syscall.InvalidHandle, nil, syscall.PAGE_READWRITE, 0, uint32(size), name)
	if err != nil {
		return nil, os.NewSyscallError("CreateFileMapping", err)
	}

	addr, err := syscall.MapViewOfFile(hMapFile, syscall.FILE_MAP_WRITE|syscall.FILE_MAP_READ, 0, 0, uintptr(size))
	if err != nil {
		return nil, os.NewSyscallError("MapViewOfFile", err)
	}

	the := &shm{
		//name:     name,
		handler: &hMapFile,
		//addr:    addr,
	}
	the.init(addr, size)
	return the, nil
}

type shm struct {
	segment
	//name  *uint16
	handler *syscall.Handle
	//addr    uintptr
}

// Detach used to detach from memory segment
func (the *shm) Detach() error {
	err := syscall.UnmapViewOfFile(the.addr())
	if err != nil {
		return os.NewSyscallError("UnmapViewOfFile", err)
	}

	err = syscall.CloseHandle(*the.handler)
	if err != nil {
		return os.NewSyscallError("CloseHandle", err)
	}
	return nil
}
