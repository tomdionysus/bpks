package bpks

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBPKSFormat(t *testing.T) {
	disk := NewRAMDisk(4 * 1024 * 1024)

	BPKS := New(disk)

	err := BPKS.Format()
	assert.Nil(t, err)

	x := make([]byte, 6)

	disk.Seek(0, 0)
	disk.Read(x)

	assert.Equal(t, []byte{0x42, 0x50, 0x4b, 0x53, 0x0, 0x1}, x)
}

func TestBPKSMountBadDevice(t *testing.T) {
	disk := NewRAMDisk(4 * 1024 * 1024)

	BPKS := New(disk)

	err := BPKS.Mount()
	assert.NotNil(t, err)
	assert.Equal(t, "Not a BPKS device", err.Error())
}

func TestBPKSMountGoodDevice(t *testing.T) {
	disk := NewRAMDisk(4 * 1024 * 1024)

	BPKS := New(disk)

	err := BPKS.Format()
	assert.Nil(t, err)
	err = BPKS.Mount()
	assert.Nil(t, err)
}
