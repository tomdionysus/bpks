package bpks

import (
	"fmt"
	"sort"
)

type KeyPointerList []KeyPointer

func NewKeyPointerListFromBuffer(buffer []byte) *KeyPointerList {
	fmt.Printf("-- Init KeyPointerList from buffer len %d\n", len(buffer))
	ln := int(sliceToUint16(buffer[0:2]))
	fmt.Printf("-- Init KeyPointerList from buffer %d entries\n", ln)
	x := KeyPointerList{}
	for i := 0; i < ln; i++ {
		x = append(x, NewKeyPointerFromBuffer(buffer[2+(i*24):2+((i+1)*24)]))
	}
	return &x
}

func (kpl *KeyPointerList) Add(kp KeyPointer) {
	fmt.Printf("KeyPointerList.Add %s -> %d\n", kp.Key, kp.BlockAddress)
	*kpl = append(*kpl, kp)
	sort.Sort(kpl)
}

func (kpl *KeyPointerList) Find(key Key) (KeyPointer, bool) {
	fmt.Printf("KeyPointerList.Find %s\n", key)
	l := len(*kpl)
	i := sort.Search(l, func(i int) bool { return (*kpl)[i].Key.Cmp(key) != -1 })
	if i < l {
		return (*kpl)[i], true
	}

	return KeyPointer{}, false
}

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

func (kpl *KeyPointerList) MinKey() Key {
	if len(*kpl) == 0 {
		return MinKey
	}
	return (*kpl)[0].Key
}

func (kpl *KeyPointerList) MaxKey() Key {
	x := len(*kpl)
	if x == 0 {
		return MaxKey
	}
	return (*kpl)[x-1].Key
}

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
func (kpl *KeyPointerList) Len() int {
	return len(*kpl)
}

func (kpl *KeyPointerList) Less(i, j int) bool {
	return (*kpl)[i].Cmp((*kpl)[j]) == -1
}

func (kpl *KeyPointerList) Swap(i, j int) {
	tp := (*kpl)[i]
	(*kpl)[i] = (*kpl)[j]
	(*kpl)[j] = tp
}
