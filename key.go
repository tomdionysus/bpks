package bpks

import (
	"bytes"
	"crypto/md5"
	"crypto/rand"
	"fmt"
)

// Key represents a 128 bit key value, used to uniquely identify a value in the key-value store.
type Key [16]byte

// MinKey is the minimum possible, or 'nil' Key
var MinKey = Key{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}

// MaxKey is the maximum possible Key
var MaxKey = Key{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}

// NewKeyFromBuffer returns a new Key parsed from the supplied buffer
func NewKeyFromBuffer(buffer []byte) Key {
	x := Key{}
	copy(x[:], buffer[0:16])
	return x
}

// NewKeyFromStringMD5 returns a new Key set to the 128 bit MD5 hash of the supplied string
func NewKeyFromStringMD5(str string) Key {
	return Key(md5.Sum([]byte(str)))
}

// NewKeyRandom returns a cryptographically generated random Key.
func NewKeyRandom() Key {
	key := make([]byte, 16)
	_, _ = rand.Read(key)
	return NewKeyFromBuffer(key)
}

// String returns the Key bytes as a string in lowercase hexadecimal format
func (key Key) String() string {
	return fmt.Sprintf("%02x", key[:])
}

// Nil returns true if this Key equals MinKey (nil Key)
func (key Key) Nil() bool {
	return bytes.Compare(key[:], MinKey[:]) == 0
}

// Cmp compares the Key to another Key and returns:
// * -1 If this Key is less than the other Key
// * 0 If this Key is equal to the other Key
// * +1 If this Key is more than the other Key
func (key Key) Cmp(other Key) int {
	return bytes.Compare(key[:], other[:])
}

// AsSlice returns this Key seriaised as a []byte of length 16.
func (key Key) AsSlice() []byte {
	x := [16]byte(key)
	return x[:]
}
