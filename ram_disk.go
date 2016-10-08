package bpks

import (
	"errors"
)

// RAMDisk is a in-memory buffer than implements os.ReadWriteSeeker
type RAMDisk struct {
	buffer []byte
	pos    uint64
}

// NewRAMDisk creates and returns a new RAMDisk of the supplied size in bytes
func NewRAMDisk(sizeInBytes int64) *RAMDisk {
	inst := &RAMDisk{
		buffer: make([]byte, sizeInBytes),
		pos:    0,
	}
	return inst
}

// Read reads from the current position into the supplied slice of bytes, returning the
// number of bytes read and/or any errors.
func (rd *RAMDisk) Read(p []byte) (int, error) {
	var err error
	readlen := uint64(len(p))
	if rd.pos+readlen > uint64(len(rd.buffer)) {
		readlen = uint64(len(rd.buffer)) - rd.pos
		err = errors.New("Read truncated, disk end reached")
	}
	for i := uint64(0); i < readlen; i++ {
		p[i] = rd.buffer[rd.pos+i]
	}
	return int(readlen), err
}

// Read writes at the current position from the supplied slice of bytes, returning the
// number of bytes written and/or any errors.
func (rd *RAMDisk) Write(p []byte) (int, error) {
	var err error
	writelen := uint64(len(p))
	if rd.pos+writelen > uint64(len(rd.buffer)) {
		writelen = uint64(len(rd.buffer)) - rd.pos
		err = errors.New("Write truncated, disk end reached")
	}
	for i := uint64(0); i < writelen; i++ {
		rd.buffer[rd.pos+i] = p[i]
	}
	rd.pos += writelen
	return int(writelen), err
}

// Seek moves the current position to the supplied offset, using the second parameter to
// denote the type of seek: 0 - start of disk, 1 - current position, 2 - end of disk.
// Returns the new position from the start of the disk, and/or any errors.
func (rd *RAMDisk) Seek(offset int64, whence int) (int64, error) {
	p := int64(rd.pos)
	switch whence {
	case 0:
		p = offset
	case 1:
		p += offset
	case 2:
		p = int64(len(rd.buffer)) + offset
	}
	if p < 0 || p > int64(len(rd.buffer)-1) {
		return int64(rd.pos), errors.New("Cannot seek outside of disk limits")
	}
	rd.pos = uint64(p)
	return p, nil
}
