package bpks

// FreeSpace represents a continuous area of free space
type FreeSpace struct {
	Min uint64
	Max uint64
}

// NewFreeSpaceFromBuffer returns a new FreeSpace parsed from the supplied buffer
func NewFreeSpaceFromBuffer(buffer []byte) FreeSpace {
	x := FreeSpace{
		Min: sliceToUint64(buffer[0:8]),
		Max: sliceToUint64(buffer[8:16]),
	}
	return x
}

// Cmp compares the Min of this FreeSpace to the Min of another FreeSpace and returns:
// * -1 If this Min is less than the other Min
// * 0 If this Min is equal to the other Min
// * +1 If this Min is more than the other Min
func (fs FreeSpace) Cmp(other FreeSpace) int {
	if fs.Min < other.Min {
		return -1
	}
	if fs.Min > other.Min {
		return 1
	}
	return 0
}

// AsSlice returns this FreeSpace seriaised as a []byte of length 16.
func (fs FreeSpace) AsSlice() []byte {
	buf := uint64ToSlice(fs.Min)
	buf = append(buf, uint64ToSlice(fs.Max)...)
	return buf
}
