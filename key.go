package bpks

import (
	"bytes"
	"fmt"
)

type Key [16]byte

var MinKey Key = Key{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
var MaxKey Key = Key{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}

func (me Key) String() string {
	return fmt.Sprintf("%02x", me[:])
}

func (me Key) Nil() bool {
	return bytes.Compare(me[:], MinKey[:]) == 0
}

func (me Key) Cmp(other Key) int {
	return bytes.Compare(me[:], other[:])
}

func (me Key) AsSlice() []byte {
	x := [16]byte(me)
	return x[:]
}
