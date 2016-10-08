package bpks

import (
	"fmt"
)

// KeyPointer associates a Key with a uint64 block address.
type KeyPointer struct {
	Key          Key
	BlockAddress uint64
}

// NewKeyPointerFromBuffer returns a new KeyPointer parsed from the supplied buffer
func NewKeyPointerFromBuffer(buffer []byte) KeyPointer {
	key := [16]byte{}
	copy(key[:], buffer[0:16])
	x := KeyPointer{
		Key:          Key(key),
		BlockAddress: sliceToUint64(buffer[16:24]),
	}

	fmt.Printf("-- Init KeyPointer from buffer %s -> %d\n", x.Key, x.BlockAddress)

	return x
}

// Nil returns true if the Key of this KeyPointer is 'nil' (0x00000000000000000000000000000000)
func (kp KeyPointer) Nil() bool {
	return kp.Key.Nil()
}

// Cmp compares the Key of this KeyPointer to the Key of another KeyPointer and returns:
// * -1 If this Key is less than the other Key
// * 0 If this Key is equal to the other Key
// * +1 If this Key is more than the other Key
func (kp KeyPointer) Cmp(other KeyPointer) int {
	return kp.Key.Cmp(other.Key)
}

// AsSlice returns this KeyPointer seriaised as a []byte of length 24.
func (kp KeyPointer) AsSlice() []byte {
	buf := kp.Key.AsSlice()
	buf = append(buf, uint64ToSlice(kp.BlockAddress)...)
	return buf
}
