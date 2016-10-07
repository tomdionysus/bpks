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

func (me *KeyPointerList) Add(kp KeyPointer) {
	fmt.Printf("KeyPointerList.Add %s\n", kp)
	*me = append(*me, kp)
	sort.Sort(me)
}

func (me *KeyPointerList) Find(key Key) (KeyPointer, bool) {
	fmt.Printf("KeyPointerList.Find %s\n", key)
	l := len(*me)
	i := sort.Search(l, func(i int) bool { return (*me)[i].Key.Cmp(key) != -1 })
	if i < l {
		return (*me)[i], true
	}

	return KeyPointer{}, false
}

func (me *KeyPointerList) MinKey() Key {
	if len(*me) == 0 {
		return MinKey
	}
	return (*me)[0].Key
}

func (me *KeyPointerList) MaxKey() Key {
	x := len(*me)
	if x == 0 {
		return MaxKey
	}
	return (*me)[x-1].Key
}

func (me *KeyPointerList) AsSlice() []byte {
	buf := []byte{}
	l := me.Len()
	buf = append(buf, uint16ToSlice(uint16(l))...)
	for i := 0; i < l; i++ {
		buf = append(buf, (*me)[i].AsSlice()...)
	}
	return buf
}

// Implement sort.Interface
func (me *KeyPointerList) Len() int {
	return len(*me)
}

func (me *KeyPointerList) Less(i, j int) bool {
	return (*me)[i].Cmp((*me)[j]) == -1
}

func (me *KeyPointerList) Swap(i, j int) {
	tp := (*me)[i]
	(*me)[i] = (*me)[j]
	(*me)[j] = tp
}
