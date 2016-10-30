package bpks

type Block interface {
	AsSlice() []byte
	GetBlockAddress() uint64
}
