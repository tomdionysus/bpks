package bpks

// 4096 Bytes

// 0 - Prev 8 bytes
// 8 - Next 8 bytes
// 16 - Length 2 bytes
// 18 - Data max 4078 bytes

// DataBlock represents a 4078 byte data block with block addresses of the previous and next blocks.
type DataBlock struct {
	// Not Serialized
	BPKS         BPKS
	BlockAddress uint64

	// Serialized
	Prev   uint64
	Next   uint64
	Length uint16
	Data   []byte
}
