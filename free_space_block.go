package bpks

// Blocks are BlockSize Bytes.

// 0 - Minimum KeyPointer
// 24 - Max KeyPointer
// 48 - Length in Keypointers 2 bytes
// 50 - Slice of max 168 KeyPointers

// FreeSpaceBlock represents a block containing records of free space.
type FreeSpaceBlock struct {
	// Not Serialized
	BPKS         *BPKS
	BlockAddress uint64

	// Serialized
	Min           uint64
	Max           uint64
	FreeSpaceList *FreeSpaceList
}

func (fsb *FreeSpaceBlock) Allocate() (uint64, error) {
	// TODO: Detect full
	return fsb.FreeSpaceList.Allocate()
}

func (fsb *FreeSpaceBlock) Deallocate(blockaddress uint64) error {
	return fsb.FreeSpaceList.Deallocate(blockaddress)
}

// NewFreeSpaceBlock returns a pointer to a new FreeSpaceBlock with the specified BPKS owner and block address, containing a
// single freespace records
func NewFreeSpaceBlock(bpks *BPKS, blockAddress, min, max uint64) *FreeSpaceBlock {
	return &FreeSpaceBlock{
		BPKS:          bpks,
		BlockAddress:  blockAddress,
		Min:           min,
		Max:           max,
		FreeSpaceList: &FreeSpaceList{FreeSpace{Min: min, Max: max}},
	}
}

// NewFreeSpaceBlockFromBuffer returns a pointer to a new FreeSpaceBlock with the specified BPKS owner and block address,
// parsed from the supplied block buffer.
func NewFreeSpaceBlockFromBuffer(bpks *BPKS, blockAddress uint64, buffer []byte) *FreeSpaceBlock {
	// fmt.Printf("-- Init Index Block from buffer len %d\n", len(buffer))
	return &FreeSpaceBlock{
		BPKS:          bpks,
		BlockAddress:  blockAddress,
		Min:           sliceToUint64(buffer[0:8]),
		Max:           sliceToUint64(buffer[8:16]),
		FreeSpaceList: NewFreeSpaceListFromBuffer(buffer[24:BlockSize]),
	}
}

// AsSlice serialises and returns the IndexBlock as a []byte, padded to BlockSize.
func (fsb *FreeSpaceBlock) AsSlice() []byte {
	buf := uint64ToSlice(fsb.Min)
	buf = append(buf, uint64ToSlice(fsb.Min)...)
	buf = append(buf, fsb.FreeSpaceList.AsSlice()...)
	if len(buf) < BlockSize {
		x := make([]byte, BlockSize-len(buf))
		buf = append(buf, x...)
	}
	return buf
}
