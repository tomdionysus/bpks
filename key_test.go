package bpks

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestKeyString(t *testing.T) {
	x := Key{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 0xa, 0xb, 0xc, 0xd, 0xe, 0xf}

	assert.Equal(t, "000102030405060708090a0b0c0d0e0f", x.String())
}

func TestKeyNil(t *testing.T) {
	x := Key{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 0xa, 0xb, 0xc, 0xd, 0xe, 0xf}

	assert.False(t, x.Nil())

	x = Key{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}

	assert.True(t, x.Nil())
}

func TestKeyCmp(t *testing.T) {
	x := Key{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1}
	y := Key{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2}
	z := Key{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2}

	assert.Equal(t, 1, y.Cmp(x))
	assert.Equal(t, -1, x.Cmp(y))
	assert.Equal(t, 0, y.Cmp(z))
}
