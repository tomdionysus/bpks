// Package bpks implements a B+Tree Key Store that stores key/value pairs on an underlying io.ReadWriteSeeker device.
package bpks

import (
	"bytes"
	"errors"
	"fmt"
	"io"
)

const BLOCK_SIZE = 4096

// BPKS (B+Tree Key Store) is a key-value store based around a B+Tree
type BPKS struct {
	Device io.ReadWriteSeeker
	Root   *IndexBlock
}

// BPKSHeader is the byte array "BPKS" plus a major version (0x00, 0x01)
var BPKSHeader = []byte{0x42, 0x50, 0x4b, 0x53, 0x0, 0x1}

var firstFreeBlock uint64 = 2

// New returns a new BPKS attached to the specified io.ReadWriteSeeker
func New(device io.ReadWriteSeeker) *BPKS {
	return &BPKS{
		Device: device,
	}
}

// Mount mounts the BPKS keystore on the attached device. An error is returned if the
// device does not contain a formatted BPKS keystore.
func (me *BPKS) Mount() error {
	// Check Header
	_, err := me.Device.Seek(0, 0)
	if err != nil {
		return err
	}
	var buf = make([]byte, 6)
	_, err = me.Device.Read(buf)
	if err != nil {
		return err
	}
	if bytes.Compare(buf, BPKSHeader) != 0 {
		return errors.New("Not a BPKS device")
	}

	// Load Index Block
	root, err := me.ReadIndexBlock(2)
	if err != nil {
		return err
	}
	me.Root = root
	return nil
}

// Format initialises a new BPKS keystore on the attached ReadWriteSeeker. This
// will erase all keys and values from an existing keystore.
func (me *BPKS) Format() error {
	// Header
	_, err := me.Device.Seek(0, 0)
	if err != nil {
		return err
	}
	_, err = me.Device.Write(BPKSHeader)
	if err != nil {
		return err
	}

	// TODO: SpaceBPKS

	// Root Index Block
	me.Root = NewIndexBlock(me, 2)
	return me.WriteIndexBlock(me.Root)
}

// Allocate gets the block address of the first free block on the device and marks it used.
func (me *BPKS) Allocate() uint64 {
	firstFreeBlock++
	return firstFreeBlock
}

// Deallocate frees the specified block address for reuse.
func (me *BPKS) Deallocate(blockAddress uint64) {
	panic("Not Implemented")
}

// Low Level KeyPointer Funcs

// Add writes the specified KeyPointer to the keystore
func (me *BPKS) Add(kp KeyPointer) error {
	return me.Root.Add(kp)
}

// Find finds the specified Key in the keystore, returning its KeyPointer if found, or
// an empty KeyPointer and false if not
func (me *BPKS) Find(key Key) (KeyPointer, bool, error) {
	return me.Root.Find(key)
}

// High Level Storage Functions
func (me *BPKS) Set(key string, data []byte) error {
	// TODO: detect replace

	// Write Key
	firstDataBlockAddress := me.Allocate()
	kp := KeyPointer{
		Key:          NewKeyFromStringMD5(key),
		BlockAddress: firstDataBlockAddress,
	}
	// TODO: Write multi-block data
	db := &DataBlock{
		BPKS:         me,
		BlockAddress: firstDataBlockAddress,
		Data:         data,
	}
	err := me.Add(kp)
	if err != nil {
		return err
	}
	err = me.WriteDataBlock(db)
	if err != nil {
		return err
	}
	return nil
}

func (me *BPKS) Get(key string) ([]byte, bool, error) {
	kp, found, err := me.Find(NewKeyFromStringMD5(key))
	if err != nil {
		return nil, false, err
	}
	if !found {
		return nil, false, nil
	}
	db, err := me.ReadDataBlock(kp.BlockAddress)
	if err != nil {
		return nil, false, err
	}
	// TODO: Read multi-block data
	return db.Data, true, nil
}

// IO Funcs

func (me *BPKS) ReadIndexBlock(blockAddress uint64) (*IndexBlock, error) {
	fmt.Printf("Reading Index Block at address %d (offset %d)\n", blockAddress, blockAddress*BLOCK_SIZE)
	_, err := me.Device.Seek(int64(blockAddress*BLOCK_SIZE), 0)
	if err != nil {
		return nil, err
	}
	buffer := [BLOCK_SIZE]byte{}
	fmt.Printf("- Reading BLOCK_SIZE Bytes\n")
	c, err := me.Device.Read(buffer[:])
	if err != nil {
		return nil, err
	}
	fmt.Printf("- Read %d bytes\n", c)
	return NewIndexBlockFromBuffer(me, blockAddress, buffer[:]), nil
}

func (me *BPKS) WriteIndexBlock(block *IndexBlock) error {
	fmt.Printf("Writing Index Block at address %d (offset %d)\n", block.BlockAddress, block.BlockAddress*BLOCK_SIZE)
	_, err := me.Device.Seek(int64(block.BlockAddress*BLOCK_SIZE), 0)
	if err != nil {
		return err
	}
	buffer := block.AsSlice()
	c, err := me.Device.Write(buffer[:])
	if err != nil {
		return err
	}
	fmt.Printf("- Wrote %d bytes\n", c)
	return nil
}

func (me *BPKS) ReadDataBlock(blockAddress uint64) (*DataBlock, error) {
	fmt.Printf("Reading Data Block at address %d (offset %d)\n", blockAddress, blockAddress*BLOCK_SIZE)
	_, err := me.Device.Seek(int64(blockAddress*BLOCK_SIZE), 0)
	if err != nil {
		return nil, err
	}
	buffer := [BLOCK_SIZE]byte{}
	fmt.Printf("- Reading %d Bytes\n", BLOCK_SIZE)
	c, err := me.Device.Read(buffer[:])
	if err != nil {
		return nil, err
	}
	fmt.Printf("- Read %d bytes\n", c)
	return NewDataBlockFromBuffer(me, blockAddress, buffer[:]), nil
}

func (me *BPKS) WriteDataBlock(block *DataBlock) error {
	fmt.Printf("Writing Data Block at address %d (offset %d)\n", block.BlockAddress, block.BlockAddress*BLOCK_SIZE)
	_, err := me.Device.Seek(int64(block.BlockAddress*BLOCK_SIZE), 0)
	if err != nil {
		return err
	}
	buffer := block.AsSlice()
	c, err := me.Device.Write(buffer[:])
	if err != nil {
		return err
	}
	fmt.Printf("- Wrote %d bytes\n", c)
	return nil
}
