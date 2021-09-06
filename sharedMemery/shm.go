package sharedMemery

import (
	"encoding/binary"
	"io"
	"math"
	"reflect"
	"unsafe"
)

// ISharedMemory interface represents shared memory segment
type ISharedMemory interface {
	// Read specified amount of data (len)
	Read(one func(data []byte, half bool))
	// Write data to the shared memory segment
	Write(data []byte) (left []byte, err error)
	// Detach from the segment
	Detach() error
}

const HalfFlag = math.MaxUint16

type segment struct {
	data []byte
	len  []byte
}

func (the *segment) addr() uintptr { return (*reflect.SliceHeader)(unsafe.Pointer(&the.len)).Data }
func (the *segment) init(addr uintptr, size uint16) {
	data := &reflect.SliceHeader{
		Data: addr,
		Len:  int(size),
		Cap:  int(size),
	}
	buf := *(*[]byte)(unsafe.Pointer(data))
	the.len, the.data = buf[:2], buf[2:]
	(*reflect.SliceHeader)(unsafe.Pointer(&the.len)).Cap = 2
}
func (the *segment) getLen() uint16  { return binary.BigEndian.Uint16(the.len) }
func (the *segment) setLen(n uint16) { binary.BigEndian.PutUint16(the.len, n) }
func (the *segment) Write(data []byte) (left []byte, err error) {
	n := len(data)
	off := int(the.getLen())
	need := off + n + 2
	if need > len(the.data) {
		if off == 0 {
			the.setLen(HalfFlag)
			n = len(the.data) - 2
			copy(the.data, data[:n])
			return data[n:], nil
		}
		return nil, io.ErrShortWrite
	}
	binary.BigEndian.PutUint16(the.data[off:], uint16(n))
	copy(the.data[off+2:], data)
	the.setLen(uint16(need))
	return nil, nil
}
func (the *segment) Read(fun func(data []byte, half bool)) {
	n := int(the.getLen())
	if n == HalfFlag {
		fun(the.data[:len(the.data)-2], true)
	} else {
		i := 0
		for n > i {
			one := int(binary.BigEndian.Uint16(the.data[i:]))
			i += 2 + one
			fun(the.data[i-one:i], false)
		}
	}
	the.setLen(0)
}
