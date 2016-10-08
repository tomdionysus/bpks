package bpks

import (
	"fmt"
	"sort"
)

// KeyPointerList represents a slice of KeyPointers
type KeyPointerList []KeyPointer

// NewKeyPointerListFromBuffer returns a pointer to a new KeyPointerList, parsed from the supplied
// buffer.
func NewKeyPointerListFromBuffer(buffer []byte) *KeyPointerList {
	ln := int(sliceToUint16(buffer[0:2]))
	// fmt.Printf("-- Init KeyPointerList from buffer len %d bytes %d entries\n", len(buffer), ln)
	x := KeyPointerList{}
	for i := 0; i < ln; i++ {
		x = append(x, NewKeyPointerFromBuffer(buffer[2+(i*24):2+((i+1)*24)]))
	}
	return &x
}

// Add adds the supplied KeyPointer to this list and sorts the list.
func (kpl *KeyPointerList) Add(kp KeyPointer) {
	// fmt.Printf("KeyPointerList.Add %s -> %d\n", kp.Key, kp.BlockAddress)
	*kpl = append(*kpl, kp)
	sort.Sort(kpl)
}

// Find finds the KeyPointer with the supplied Key in this list and returns that KeyPointer,
// and whether it was found.
func (kpl *KeyPointerList) Find(key Key) (KeyPointer, bool) {
	fmt.Printf("KeyPointerList.Find %s\n", key)
	l := len(*kpl)
	i := sort.Search(l, func(i int) bool { return (*kpl)[i].Key.Cmp(key) != -1 })
	if i < l {
		return (*kpl)[i], true
	}

	return KeyPointer{}, false
}

// Remove finds and removes the KeyPointer with the supplied Key in this list and returns
// that KeyPointer and whether it was found and removed.
func (kpl *KeyPointerList) Remove(key Key) (KeyPointer, bool) {
	fmt.Printf("KeyPointerList.Remove %s\n", key)
	l := len(*kpl)
	i := sort.Search(l, func(i int) bool { return (*kpl)[i].Key.Cmp(key) != -1 })
	if i >= l {
		return KeyPointer{}, false
	}

	kpout := (*kpl)[i]
	nsl := (*kpl)[:i]
	nsl = append(nsl, (*kpl)[i+1:]...)
	*kpl = KeyPointerList(nsl)

	return kpout, true
}

// MinKey returns the Key in the list with the smallest value.
func (kpl *KeyPointerList) MinKey() Key {
	if len(*kpl) == 0 {
		return MinKey
	}
	return (*kpl)[0].Key
}

// MaxKey returns the Key in the list with the largest value.
func (kpl *KeyPointerList) MaxKey() Key {
	x := len(*kpl)
	if x == 0 {
		return MaxKey
	}
	return (*kpl)[x-1].Key
}

// AsSlice returns this KeyPointerList seriaised as a []byte
func (kpl *KeyPointerList) AsSlice() []byte {
	buf := []byte{}
	l := kpl.Len()
	buf = append(buf, uint16ToSlice(uint16(l))...)
	for i := 0; i < l; i++ {
		buf = append(buf, (*kpl)[i].AsSlice()...)
	}
	return buf
}

// Implement sort.Interface

// Len returns the current length of this KeyPointerList
func (kpl *KeyPointerList) Len() int {
	return len(*kpl)
}

// Less compares the Keys of the KeyPointers at the indices i and j, and returns true
// if the Key at i is less than the Key at j.
func (kpl *KeyPointerList) Less(i, j int) bool {
	return (*kpl)[i].Cmp((*kpl)[j]) == -1
}

// Swap swaps the values of the KeyPointers at the indices i and j.
func (kpl *KeyPointerList) Swap(i, j int) {
	tp := (*kpl)[i]
	(*kpl)[i] = (*kpl)[j]
	(*kpl)[j] = tp
}
