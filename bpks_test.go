package bpks

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBPKSFormat(t *testing.T) {
	disk := NewRAMDisk(4 * 1024 * 1024)

	bpks := New(disk, 1024)

	err := bpks.Format()
	assert.Nil(t, err)

	x := make([]byte, 6)

	disk.Seek(0, 0)
	disk.Read(x)

	assert.Equal(t, []byte{0x42, 0x50, 0x4b, 0x53, 0x0, 0x1}, x)
}

func TestBPKSMountBadDevice(t *testing.T) {
	disk := NewRAMDisk(4 * 1024 * 1024)

	bpks := New(disk, 1024)

	err := bpks.Mount()
	assert.NotNil(t, err)
	assert.Equal(t, "Not a BPKS device", err.Error())
}

func TestBPKSMountGoodDevice(t *testing.T) {
	disk := NewRAMDisk(4 * 1024 * 1024)

	bpks := New(disk, 1024)

	err := bpks.Format()
	assert.Nil(t, err)
	err = bpks.Mount()
	assert.Nil(t, err)
}

func TestBPKSSetGet(t *testing.T) {
	disk := NewRAMDisk(4 * 1024 * 1024)
	bpks := New(disk, 1024)
	err := bpks.Format()
	assert.Nil(t, err)

	err = bpks.Set("testing!", []byte("Hello World!"))
	assert.Nil(t, err)

	dat, found, err := bpks.Get("testing!")
	assert.Nil(t, err)
	assert.True(t, found)
	assert.Equal(t, "Hello World!", string(dat))
}

func TestBPKSGetNotFound(t *testing.T) {
	disk := NewRAMDisk(4 * 1024 * 1024)
	bpks := New(disk, 1024)
	err := bpks.Format()
	assert.Nil(t, err)

	_, found, err := bpks.Get("testing!")
	assert.Nil(t, err)
	assert.False(t, found)
}

func TestBPKSSetDelete(t *testing.T) {
	disk := NewRAMDisk(4 * 1024 * 1024)
	bpks := New(disk, 1024)
	err := bpks.Format()
	assert.Nil(t, err)

	err = bpks.Set("testing!", []byte("Hello World!"))
	assert.Nil(t, err)

	dat, found, err := bpks.Get("testing!")
	assert.Nil(t, err)
	assert.True(t, found)
	assert.Equal(t, "Hello World!", string(dat))

	found, err = bpks.Delete("testing!")
	assert.Nil(t, err)
	assert.True(t, found)

	_, found, err = bpks.Get("testing!")
	assert.Nil(t, err)
	assert.False(t, found)
}
