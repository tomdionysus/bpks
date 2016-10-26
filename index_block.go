package bpks

import (
// "fmt"
)

// Blocks are BlockSize Bytes.

// 0 - Minimum KeyPointer
// 24 - Max KeyPointer
// 48 - Length in Keypointers 2 bytes
// 50 - Slice of max 168 KeyPointers

// IndexBlock represents an internal node in the B+Tree, containing a KeyPointerList and optionally
// linking to up two other IndexBlocks representing the left and right child nodes in this tree.
type IndexBlock struct {
	// Not Serialized
	BPKS         *BPKS
	BlockAddress uint64

	// Serialized
	Min            KeyPointer
	Max            KeyPointer
	KeyPointerList *KeyPointerList
}

// NewIndexBlock returns a pointer to a new IndexBlock with the specified BPKS owner and block address.
func NewIndexBlock(bpks *BPKS, blockAddress uint64) *IndexBlock {
	return &IndexBlock{
		BPKS:           bpks,
		BlockAddress:   blockAddress,
		KeyPointerList: &KeyPointerList{},
	}
}

// NewIndexBlockFromBuffer returns a pointer to a new IndexBlock with the specified BPKS owner and block address,
// parsed from the supplied block buffer.
func NewIndexBlockFromBuffer(bpks *BPKS, blockAddress uint64, buffer []byte) *IndexBlock {
	// fmt.Printf("-- Init Index Block from buffer len %d\n", len(buffer))
	return &IndexBlock{
		BPKS:           bpks,
		BlockAddress:   blockAddress,
		Min:            NewKeyPointerFromBuffer(buffer[0:24]),
		Max:            NewKeyPointerFromBuffer(buffer[24:48]),
		KeyPointerList: NewKeyPointerListFromBuffer(buffer[48:BlockSize]),
	}
}

// Add adds a supplied KeyPointer to this IndexBlock or one of its children, splitting the indexblock if necessary
// and writing the modified IndexBlock to the BPKS device, returning nil on success or an error.
func (ib *IndexBlock) Add(kp KeyPointer) error {
	// fmt.Printf("IndexBlock.Add %s -> %d\n", kp.Key, kp.BlockAddress)
	// If there is a minimum and the key is less than the minimum
	if ib.Min.BlockAddress != 0 && ib.Min.Cmp(kp) == -1 {
		// fmt.Printf("- Min exists and key is less than min %s\n", ib.Min.Key)
		left, err := ib.BPKS.ReadIndexBlock(ib.Min.BlockAddress)
		if err != nil {
			return err
		}
		left.Add(kp)
		return nil
	}

	// If there is a maximum and the key is more than the maximum
	if ib.Max.BlockAddress != 0 && ib.Max.Cmp(kp) == 1 {
		// fmt.Printf("- Max exists and key is more than max %s\n", ib.Max.Key)
		right, err := ib.BPKS.ReadIndexBlock(ib.Max.BlockAddress)
		if err != nil {
			return err
		}
		right.Add(kp)
		return nil
	}

	// If there is space in this block
	if ib.KeyPointerList.Len() < 168 {
		ib.KeyPointerList.Add(kp)
		ib.Min.Key = ib.KeyPointerList.MinKey()
		ib.Max.Key = ib.KeyPointerList.MaxKey()
		return ib.BPKS.WriteIndexBlock(ib)
	}

	leftkpl := KeyPointerList((*ib.KeyPointerList)[0:42])
	rightkpl := KeyPointerList((*ib.KeyPointerList)[126:168])
	c := KeyPointerList((*ib.KeyPointerList)[42:126])

	leftblockAddr, err := ib.BPKS.FreeSpace.Allocate()
	if err != nil {
		return err
	}
	rightblockAddr, err := ib.BPKS.FreeSpace.Allocate()
	if err != nil {
		return err
	}

	// Split this index block
	left := IndexBlock{
		BPKS:           ib.BPKS,
		BlockAddress:   leftblockAddr,
		KeyPointerList: &leftkpl,
	}
	right := IndexBlock{
		BPKS:           ib.BPKS,
		BlockAddress:   rightblockAddr,
		KeyPointerList: &rightkpl,
	}
	// fmt.Printf("-- Split Index Block %d -> %d / %d\n", ib.BlockAddress, left.BlockAddress, right.BlockAddress)
	left.Min = ib.Min
	right.Max = ib.Max
	err = ib.BPKS.WriteIndexBlock(&left)
	if err != nil {
		return err
	}
	err = ib.BPKS.WriteIndexBlock(&right)
	if err != nil {
		return err
	}
	ib.KeyPointerList = &c
	ib.Min.Key = ib.KeyPointerList.MinKey()
	ib.Min.BlockAddress = left.BlockAddress
	ib.Max.Key = ib.KeyPointerList.MaxKey()
	ib.Max.BlockAddress = right.BlockAddress
	return ib.BPKS.WriteIndexBlock(ib)
}

// Find finds and returns the KeyPointer associated with the supplied key in this IndexBlock or one of
// its children, returning the KeyPointer, whether the key was found, and an error if any.
func (ib *IndexBlock) Find(key Key) (KeyPointer, bool, error) {
	// fmt.Printf("IndexBlock.Find %d: %s\n", ib.BlockAddress, key)
	// If there is a minimum and the key is less than the minimum
	if ib.Min.BlockAddress != 0 && ib.Min.Key.Cmp(key) == -1 {
		// fmt.Printf("- Min exists and key is less than min %s\n", ib.Min.Key)
		left, err := ib.BPKS.ReadIndexBlock(ib.Min.BlockAddress)
		if err != nil {
			return KeyPointer{}, false, err
		}
		return left.Find(key)
	}

	// If there is a maximum and the key is more than the maximum
	if ib.Max.BlockAddress != 0 && ib.Max.Key.Cmp(key) == 1 {
		// fmt.Printf("- Max exists and key is more than max %s\n", ib.Max.Key)
		right, err := ib.BPKS.ReadIndexBlock(ib.Max.BlockAddress)
		if err != nil {
			return KeyPointer{}, false, err
		}
		return right.Find(key)
	}

	// Find in this indexblock
	kp, found := ib.KeyPointerList.Find(key)
	return kp, found, nil
}

// Remove finds and removes the KeyPointer associated with the supplied key in this IndexBlock or one of
// its children, returning the removed KeyPointer, whether the key was found, and an error if any.
func (ib *IndexBlock) Remove(key Key) (KeyPointer, bool, error) {
	// If there is a minimum and the key is less than the minimum
	if ib.Min.BlockAddress != 0 && ib.Min.Key.Cmp(key) == -1 {
		left, err := ib.BPKS.ReadIndexBlock(ib.Min.BlockAddress)
		if err != nil {
			return KeyPointer{}, false, err
		}
		return left.Remove(key)
	}

	// If there is a maximum and the key is more than the maximum
	if ib.Max.BlockAddress != 0 && ib.Max.Key.Cmp(key) == 1 {
		right, err := ib.BPKS.ReadIndexBlock(ib.Max.BlockAddress)
		if err != nil {
			return KeyPointer{}, false, err
		}
		return right.Remove(key)
	}

	// Remove in this indexblock
	kp, found := ib.KeyPointerList.Remove(key)

	if found {
		ib.BPKS.WriteIndexBlock(ib)
	}

	// TODO: Merge indexblocks if underpopulated

	return kp, found, nil
}

// AsSlice serialises and returns the IndexBlock as a []byte, padded to BlockSize.
func (ib *IndexBlock) AsSlice() []byte {
	buf := ib.Min.AsSlice()
	buf = append(buf, ib.Max.AsSlice()...)
	buf = append(buf, ib.KeyPointerList.AsSlice()...)
	if len(buf) < BlockSize {
		x := make([]byte, BlockSize-len(buf))
		buf = append(buf, x...)
	}
	return buf
}
