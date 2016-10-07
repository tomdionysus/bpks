package bpks

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUint64ToSlice(t *testing.T) {
	b := uint64ToSlice(uint64(0x123456789ABCDEF0))

	assert.Equal(t, []byte{0x12, 0x34, 0x56, 0x78, 0x9a, 0xbc, 0xde, 0xf0}, b)
}

func TestUint16ToSlice(t *testing.T) {
	b := uint16ToSlice(uint16(0x1234))

	assert.Equal(t, []byte{0x12, 0x34}, b)
}
