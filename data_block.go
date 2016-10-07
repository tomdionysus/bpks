package bpks

import (
	"fmt"
)

// BLOCK_SIZE Bytes

// 0 - Prev 8 bytes
// 8 - Next 8 bytes
// 16 - Length 2 bytes
// 18 - Data max 4078 bytes

// DataBlock represents a 4078 byte data block with block addresses of the previous and next blocks.
type DataBlock struct {
	// Not Serialized
	BPKS         *BPKS
	BlockAddress uint64

	// Serialized
	Prev   uint64
	Next   uint64
	Length uint16
	Data   []byte
}

func NewDataBlockFromBuffer(bpks *BPKS, blockAddress uint64, buffer []byte) *DataBlock {
	fmt.Printf("-- Init Data Block from buffer len %d\n", len(buffer))
	ln := sliceToUint16(buffer[16:18])
	if ln > BLOCK_SIZE {
		ln = BLOCK_SIZE
	}
	x := &DataBlock{
		BPKS:         bpks,
		BlockAddress: blockAddress,
		Prev:         sliceToUint64(buffer[0:8]),
		Next:         sliceToUint64(buffer[8:16]),
		Length:       ln,
	}
	fmt.Printf("-- Data is %d bytes long\n", ln)
	x.Data = buffer[18 : 18+ln]
	return x
}

func (me *DataBlock) AsSlice() []byte {
	me.Length = uint16(len(me.Data))
	buf := uint64ToSlice(me.Prev)
	buf = append(buf, uint64ToSlice(me.Next)...)
	buf = append(buf, uint16ToSlice(me.Length)...)
	buf = append(buf, me.Data...)
	if len(buf) < BLOCK_SIZE {
		x := make([]byte, BLOCK_SIZE-len(buf))
		buf = append(buf, x...)
	}
	return buf
}
