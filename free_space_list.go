package bpks

import (
	// "fmt"
	"sort"
)

// FreeSpaceList represents a slice of FreeSpaces
type FreeSpaceList []FreeSpace

// NewFreeSpaceListFromBuffer returns a pointer to a new FreeSpaceList, parsed from the supplied
// buffer.
func NewFreeSpaceListFromBuffer(buffer []byte) *FreeSpaceList {
	ln := int(sliceToUint16(buffer[0:2]))
	// fmt.Printf("-- Init FreeSpaceList from buffer len %d bytes %d entries\n", len(buffer), ln)
	x := FreeSpaceList{}
	for i := 0; i < ln; i++ {
		x = append(x, NewFreeSpaceFromBuffer(buffer[2+(i*16):2+((i+1)*16)]))
	}
	return &x
}

// Add adds the supplied FreeSpace to this list and sorts the list.
func (fsl *FreeSpaceList) Add(kp FreeSpace) {
	// fmt.Printf("FreeSpaceList.Add %s -> %d\n", kp.Key, kp.BlockAddress)
	*fsl = append(*fsl, kp)
	sort.Sort(fsl)
}

// AsSlice returns this FreeSpaceList seriaised as a []byte
func (fsl *FreeSpaceList) AsSlice() []byte {
	buf := []byte{}
	l := fsl.Len()
	buf = append(buf, uint16ToSlice(uint16(l))...)
	for i := 0; i < l; i++ {
		buf = append(buf, (*fsl)[i].AsSlice()...)
	}
	return buf
}

// Implement sort.Interface

// Len returns the current length of this FreeSpaceList
func (fsl *FreeSpaceList) Len() int {
	return len(*fsl)
}

// Less compares the Keys of the FreeSpaces at the indices i and j, and returns true
// if the Key at i is less than the Key at j.
func (fsl *FreeSpaceList) Less(i, j int) bool {
	return (*fsl)[i].Cmp((*fsl)[j]) == -1
}

// Swap swaps the values of the FreeSpaces at the indices i and j.
func (fsl *FreeSpaceList) Swap(i, j int) {
	tp := (*fsl)[i]
	(*fsl)[i] = (*fsl)[j]
	(*fsl)[j] = tp
}
