//go:build windows
// +build windows

package sharedMemery

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var testPayload = "qwertyuioplkjhgfdsazxcvbnm" //26

func TestCreate(t *testing.T) {
	shm, err := Create(2233, 32)
	if err != nil {
		t.Fatal(shm)
	}
	d := []byte(testPayload)
	shm.Write(d)

	shm2, err := Create(2233, 32)
	var d2 []byte
	shm2.Read(func(data []byte, half bool) {
		d2 = data
	})
	assert.Equal(t, d, d2)
	err = shm.Detach()
	if err != nil {
		t.Fatal(err)
	}

	err = shm2.Detach()
	if err != nil {
		t.Fatal(err)
	}
}

func TestSegment_Read(t *testing.T) {
	shm, err := Create(2233, 30)
	if err != nil {
		t.Fatal(shm)
	}

	//d2 := make([]byte, 26, 26)
	shm.Read(func(data []byte, half bool) {
		//d2 = data
	})
	if err != nil {
		t.Fatal(err)
	}

}
