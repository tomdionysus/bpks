package bpks

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
	return x
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
