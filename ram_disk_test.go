package bpks

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewRAMDisk(t *testing.T) {
	x := NewRAMDisk(1048576)

	assert.NotNil(t, x.buffer)
	assert.Equal(t, 1048576, len(x.buffer))
	assert.Equal(t, uint64(0), x.pos)
}

func TestRAMDiskRead(t *testing.T) {
	x := NewRAMDisk(1048576)
	// Test Standard Read

	x.buffer[0] = 0xFF
	x.buffer[1] = 0xFE
	x.buffer[2] = 0xFD
	x.buffer[3] = 0xFC

	buf := [4]byte{}
	n, err := x.Read(buf[:])
	assert.Nil(t, err)
	assert.Equal(t, 4, n)
	assert.Equal(t, [4]byte{0xFF, 0xFE, 0xFD, 0xFC}, buf)

	// Test Read past end
	x.pos = (1048576) - 2
	n, err = x.Read(buf[:])
	assert.NotNil(t, err)
	assert.Equal(t, 2, n)
	assert.Equal(t, [4]byte{0x00, 0x00, 0xFD, 0xFC}, buf)
}

func TestRAMDiskWrite(t *testing.T) {
	x := NewRAMDisk(1048576)
	// Test Standard Write

	buf := [4]byte{0xFF, 0xFE, 0xFD, 0xFC}

	n, err := x.Write(buf[:])
	assert.Nil(t, err)
	assert.Equal(t, 4, n)
	assert.Equal(t, []byte{0xFF, 0xFE, 0xFD, 0xFC}, x.buffer[0:4])

	// Test write past end
	x.pos = (1048576) - 2
	n, err = x.Write(buf[:])
	assert.NotNil(t, err)
	assert.Equal(t, 2, n)
	assert.Equal(t, []byte{0x0, 0x0, 0xFF, 0xFE}, x.buffer[1048572:1048576])
}

func TestRAMDiskSeek(t *testing.T) {
	x := NewRAMDisk(1048576)

	p, err := x.Seek(5, 0)
	assert.Nil(t, err)
	assert.Equal(t, uint64(5), x.pos)
	assert.Equal(t, int64(5), p)

	p, err = x.Seek(5, 1)
	assert.Nil(t, err)
	assert.Equal(t, uint64(10), x.pos)
	assert.Equal(t, int64(10), p)

	p, err = x.Seek(-1024, 2)
	assert.Nil(t, err)
	assert.Equal(t, uint64(1047552), x.pos)
	assert.Equal(t, int64(1047552), p)

	p, err = x.Seek(-1024, 0)
	assert.NotNil(t, err)
	assert.Equal(t, uint64(1047552), x.pos)
	assert.Equal(t, int64(1047552), p)

	p, err = x.Seek(1025, 1)
	assert.NotNil(t, err)
	assert.Equal(t, uint64(1047552), x.pos)
	assert.Equal(t, int64(1047552), p)

	p, err = x.Seek(1024, 2)
	assert.NotNil(t, err)
	assert.Equal(t, uint64(1047552), x.pos)
	assert.Equal(t, int64(1047552), p)
}
