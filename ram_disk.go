package bpks

import (
	"errors"
)

// RAMDisk is a in-memory buffer than implements os.ReadWriteSeeker
type RAMDisk struct {
	buffer []byte
	pos    uint64
}

func NewRAMDisk(sizeInBytes int64) *RAMDisk {
	inst := &RAMDisk{
		buffer: make([]byte, sizeInBytes),
		pos:    0,
	}
	return inst
}

func (me *RAMDisk) Read(p []byte) (int, error) {
	var err error
	readlen := uint64(len(p))
	if me.pos+readlen > uint64(len(me.buffer)) {
		readlen = uint64(len(me.buffer)) - me.pos
		err = errors.New("Read truncated, disk end reached")
	}
	for i := uint64(0); i < readlen; i++ {
		p[i] = me.buffer[me.pos+i]
	}
	return int(readlen), err
}

func (me *RAMDisk) Write(p []byte) (int, error) {
	var err error
	writelen := uint64(len(p))
	if me.pos+writelen > uint64(len(me.buffer)) {
		writelen = uint64(len(me.buffer)) - me.pos
		err = errors.New("Write truncated, disk end reached")
	}
	for i := uint64(0); i < writelen; i++ {
		me.buffer[me.pos+i] = p[i]
	}
	me.pos += writelen
	return int(writelen), err
}

func (me *RAMDisk) Seek(offset int64, whence int) (int64, error) {
	p := int64(me.pos)
	switch whence {
	case 0:
		p = offset
	case 1:
		p += offset
	case 2:
		p = int64(len(me.buffer)) + offset
	}
	if p < 0 || p > int64(len(me.buffer)-1) {
		return int64(me.pos), errors.New("Cannot seek outside of disk limits")
	}
	me.pos = uint64(p)
	return p, nil
}
