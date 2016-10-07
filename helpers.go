package bpks

func sliceToUint64(sl []byte) uint64 {
	return uint64((sl[0] << 56) | (sl[1] << 48) | (sl[2] << 40) | (sl[3] << 32) | (sl[4] << 24) | (sl[5] << 16) | (sl[6] << 8) | sl[7])
}

func sliceToUint16(sl []byte) uint16 {
	return uint16((sl[0] << 8) | sl[1])
}
