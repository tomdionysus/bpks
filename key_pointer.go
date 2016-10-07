package bpks

import (
	"fmt"
)

// KeyPointer associates a Key with a block address
// Length 24 bytes.

type KeyPointer struct {
	Key          Key
	BlockAddress uint64
}

func NewKeyPointerFromBuffer(buffer []byte) KeyPointer {
	keyarr := [16]byte{}
	copy(keyarr[:], buffer[0:16])
	x := KeyPointer{
		Key:          Key(keyarr),
		BlockAddress: sliceToUint64(buffer[16:24]),
	}

	fmt.Printf("-- Init KeyPointer from buffer %s -> %d\n", x.Key, x.BlockAddress)

	return x
}

func (me KeyPointer) Nil() bool {
	return me.Key.Nil()
}

func (me KeyPointer) Cmp(other KeyPointer) int {
	return me.Key.Cmp(other.Key)
}
