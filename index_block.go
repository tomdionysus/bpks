package bpks

import (
	"fmt"
)

// Blocks are BLOCK_SIZE Bytes.

// 0 - Minimum KeyPointer
// 24 - Max KeyPointer
// 48 - Length in Keypointers 2 bytes
// 50 - Slice of max 168 KeyPointers

type IndexBlock struct {
	// Not Serialized
	BPKS         *BPKS
	BlockAddress uint64

	// Serialized
	Min            KeyPointer
	Max            KeyPointer
	KeyPointerList *KeyPointerList
}

func NewIndexBlock(bpks *BPKS, blockAddress uint64) *IndexBlock {
	return &IndexBlock{
		BPKS:           bpks,
		BlockAddress:   blockAddress,
		KeyPointerList: &KeyPointerList{},
	}
}

func NewIndexBlockFromBuffer(bpks *BPKS, blockAddress uint64, buffer []byte) *IndexBlock {
	fmt.Printf("-- Init Index Block from buffer len %d\n", len(buffer))
	return &IndexBlock{
		BPKS:           bpks,
		BlockAddress:   blockAddress,
		Min:            NewKeyPointerFromBuffer(buffer[0:24]),
		Max:            NewKeyPointerFromBuffer(buffer[24:48]),
		KeyPointerList: NewKeyPointerListFromBuffer(buffer[48:BLOCK_SIZE]),
	}
}

func (me *IndexBlock) Add(kp KeyPointer) error {
	fmt.Printf("IndexBlock.Add %s\n", kp)
	// If there is a minimum and the key is less than the minimum
	if me.Min.BlockAddress != 0 && me.Min.Cmp(kp) == -1 {
		fmt.Printf("- Min exists and key is less than min %s\n", me.Min.Key)
		left, err := me.BPKS.ReadIndexBlock(me.Min.BlockAddress)
		if err != nil {
			return err
		}
		left.Add(kp)
		return nil
	}

	// If there is a maximum and the key is more than the maximum
	if me.Max.BlockAddress != 0 && me.Max.Cmp(kp) == 1 {
		fmt.Printf("- Max exists and key is more than max %s\n", me.Max.Key)
		right, err := me.BPKS.ReadIndexBlock(me.Max.BlockAddress)
		if err != nil {
			return err
		}
		right.Add(kp)
		return nil
	}

	// If there is space in this block
	if me.KeyPointerList.Len() < 168 {
		me.KeyPointerList.Add(kp)
		me.Min.Key = me.KeyPointerList.MinKey()
		me.Max.Key = me.KeyPointerList.MaxKey()
		return me.BPKS.WriteIndexBlock(me)
	}

	a := KeyPointerList((*me.KeyPointerList)[0:42])
	b := KeyPointerList((*me.KeyPointerList)[126:168])
	c := KeyPointerList((*me.KeyPointerList)[42:126])

	// Split this index block
	left := IndexBlock{
		BPKS:           me.BPKS,
		BlockAddress:   me.BPKS.Allocate(),
		KeyPointerList: &a,
	}
	right := IndexBlock{
		BPKS:           me.BPKS,
		BlockAddress:   me.BPKS.Allocate(),
		KeyPointerList: &b,
	}
	left.Min = me.Min
	right.Max = me.Max
	err := me.BPKS.WriteIndexBlock(&left)
	if err != nil {
		return err
	}
	err = me.BPKS.WriteIndexBlock(&right)
	if err != nil {
		return err
	}
	me.KeyPointerList = &c
	me.Min.Key = me.KeyPointerList.MinKey()
	me.Min.BlockAddress = left.BlockAddress
	me.Max.Key = me.KeyPointerList.MaxKey()
	me.Max.BlockAddress = right.BlockAddress
	return me.BPKS.WriteIndexBlock(me)
}

func (me *IndexBlock) Find(key Key) (KeyPointer, bool, error) {
	// If there is a minimum and the key is less than the minimum
	if !me.Min.Nil() && me.Min.Key.Cmp(key) == -1 {
		left, err := me.BPKS.ReadIndexBlock(me.Min.BlockAddress)
		if err != nil {
			return KeyPointer{}, false, err
		}
		return left.Find(key)
	}

	// If there is a maximum and the key is more than the maximum
	if !me.Max.Nil() && me.Max.Key.Cmp(key) == 1 {
		right, err := me.BPKS.ReadIndexBlock(me.Max.BlockAddress)
		if err != nil {
			return KeyPointer{}, false, err
		}
		return right.Find(key)
	}

	// Find in this indexblock
	kp, found := me.KeyPointerList.Find(key)
	return kp, found, nil
}

func (me *IndexBlock) Remove(key Key) (KeyPointer, bool, error) {
	// If there is a minimum and the key is less than the minimum
	if !me.Min.Nil() && me.Min.Key.Cmp(key) == -1 {
		left, err := me.BPKS.ReadIndexBlock(me.Min.BlockAddress)
		if err != nil {
			return KeyPointer{}, false, err
		}
		return left.Remove(key)
	}

	// If there is a maximum and the key is more than the maximum
	if !me.Max.Nil() && me.Max.Key.Cmp(key) == 1 {
		right, err := me.BPKS.ReadIndexBlock(me.Max.BlockAddress)
		if err != nil {
			return KeyPointer{}, false, err
		}
		return right.Remove(key)
	}

	// Remove in this indexblock
	kp, found := me.KeyPointerList.Remove(key)

	if found {
		me.BPKS.WriteIndexBlock(me)
	}

	return kp, found, nil
}

func (me *IndexBlock) AsSlice() []byte {
	buf := me.Min.AsSlice()
	buf = append(buf, me.Max.AsSlice()...)
	buf = append(buf, me.KeyPointerList.AsSlice()...)
	if len(buf) < BLOCK_SIZE {
		x := make([]byte, BLOCK_SIZE-len(buf))
		buf = append(buf, x...)
	}
	return buf
}
