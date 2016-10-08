package bpks

import (
	"fmt"
)

// BlockSize Bytes

// 0 - Prev 8 bytes
// 8 - Next 8 bytes
// 16 - Length 2 bytes
// 18 - Data max 4078 bytes

// DataBlock represents a 4078 byte data block with block addresses of the previous and next blocks.
type DataBlock struct {
	// BPKS is a pointer to the BPKS that owns this DataBlock
	BPKS *BPKS
	// BlockAddress is the uint64 block address of this DataBlock
	BlockAddress uint64

	// Prev is the uint64 block address of the previous DataBlock in this chain, or nil
	Prev uint64
	// Next is the uint64 block address of the prevnextious DataBlock in this chain, or nil
	Next uint64
	// Data is data contained in this DataBlock, of maximum length BlockSize - 18
	Data []byte
}

// NewDataBlockFromBuffer returns a pointer to a new DataBlock with the specified BPKS owner and
// block address, parsed from the specified buffer.
func NewDataBlockFromBuffer(bpks *BPKS, blockAddress uint64, buffer []byte) *DataBlock {
	fmt.Printf("-- Init Data Block from buffer len %d\n", len(buffer))
	ln := sliceToUint16(buffer[16:18])
	if ln > BlockSize {
		ln = BlockSize
	}
	x := &DataBlock{
		BPKS:         bpks,
		BlockAddress: blockAddress,
		Prev:         sliceToUint64(buffer[0:8]),
		Next:         sliceToUint64(buffer[8:16]),
	}
	fmt.Printf("-- Data is %d bytes long\n", ln)
	x.Data = buffer[18 : 18+ln]
	return x
}

// AsSlice serialises and returns the DataBlock as a []byte, padded to BlockSize.
func (db *DataBlock) AsSlice() []byte {
	buf := uint64ToSlice(db.Prev)
	buf = append(buf, uint64ToSlice(db.Next)...)
	buf = append(buf, uint16ToSlice(uint16(len(db.Data)))...)
	buf = append(buf, db.Data...)
	if len(buf) < BlockSize {
		x := make([]byte, BlockSize-len(buf))
		buf = append(buf, x...)
	}
	return buf
}
