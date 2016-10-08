// Package bpks implements a B+Tree-like Key Store that stores key/value pairs on an underlying io.ReadWriteSeeker device.
//
// bpks is currently ALPHA and should not be used in production.
package bpks

import (
	"bytes"
	"errors"
	// "fmt"
	"io"
)

// BlockSize is the block size for the B+Tree in bytes.
const BlockSize = 4096

// BPKS (B+Tree Key Store) is a key-value store based around a B+Tree.
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
func (bp *BPKS) Mount() error {
	// Check Header
	_, err := bp.Device.Seek(0, 0)
	if err != nil {
		return err
	}
	var buf = make([]byte, 6)
	_, err = bp.Device.Read(buf)
	if err != nil {
		return err
	}
	if bytes.Compare(buf, BPKSHeader) != 0 {
		return errors.New("Not a BPKS device")
	}

	// Load Index Block
	root, err := bp.ReadIndexBlock(2)
	if err != nil {
		return err
	}
	bp.Root = root
	return nil
}

// Format initialises a new BPKS keystore on the attached ReadWriteSeeker. This
// will erase all keys and values from an existing keystore.
func (bp *BPKS) Format() error {
	// Header
	_, err := bp.Device.Seek(0, 0)
	if err != nil {
		return err
	}
	_, err = bp.Device.Write(BPKSHeader)
	if err != nil {
		return err
	}

	// TODO: SpaceBPKS

	// Root Index Block
	bp.Root = NewIndexBlock(bp, 2)
	return bp.WriteIndexBlock(bp.Root)
}

// Allocate gets the block address of the first free block on the device and marks it used.
func (bp *BPKS) Allocate() uint64 {
	firstFreeBlock++
	return firstFreeBlock
}

// Deallocate frees the specified block address for reuse.
func (bp *BPKS) Deallocate(blockAddress uint64) {
}

// Set writes a key/value pair of the MD5 of the supplied string, and data, to the key store,
// returning nil on success or an error.
func (bp *BPKS) Set(key string, data []byte) error {
	// TODO: detect replace

	// Write Key
	firstDataBlockAddress := bp.Allocate()
	kp := KeyPointer{
		Key:          NewKeyFromStringMD5(key),
		BlockAddress: firstDataBlockAddress,
	}
	// TODO: Write multi-block data
	db := &DataBlock{
		BPKS:         bp,
		BlockAddress: firstDataBlockAddress,
		Data:         data,
	}
	err := bp.Root.Add(kp)
	if err != nil {
		return err
	}
	err = bp.WriteDataBlock(db)
	if err != nil {
		return err
	}
	return nil
}

// Get finds and reads the value of the Key which is the MD5 of the given string, returning
// the data, whether they key was found, and/or an error if any.
func (bp *BPKS) Get(key string) ([]byte, bool, error) {
	kp, found, err := bp.Root.Find(NewKeyFromStringMD5(key))
	if err != nil {
		return nil, false, err
	}
	if !found {
		return nil, false, nil
	}
	db, err := bp.ReadDataBlock(kp.BlockAddress)
	if err != nil {
		return nil, false, err
	}
	// TODO: Read multi-block data
	return db.Data, true, nil
}

// Delete finds and deletes the Key which is the MD5 of the given string and its value,
// returning whether they key was found and removed, and/or an error if any.
func (bp *BPKS) Delete(key string) (bool, error) {
	kp, found, err := bp.Root.Remove(NewKeyFromStringMD5(key))
	if err != nil {
		return false, err
	}
	if !found {
		return false, nil
	}
	bp.Deallocate(kp.BlockAddress)
	// TODO: Read multi-block data
	return true, nil
}

// IO Funcs

// ReadIndexBlock reads and returns the IndexBlock at the specified block address, returning
// a pointer to the parsed IndexBlock and/or an error if any.
func (bp *BPKS) ReadIndexBlock(blockAddress uint64) (*IndexBlock, error) {
	// fmt.Printf("Reading Index Block at address %d (offset %d)\n", blockAddress, blockAddress*BlockSize)
	_, err := bp.Device.Seek(int64(blockAddress*BlockSize), 0)
	if err != nil {
		return nil, err
	}
	buffer := [BlockSize]byte{}
	// fmt.Printf("- Reading BlockSize Bytes\n")
	_, err = bp.Device.Read(buffer[:])
	if err != nil {
		return nil, err
	}
	// fmt.Printf("- Read %d bytes\n", c)
	return NewIndexBlockFromBuffer(bp, blockAddress, buffer[:]), nil
}

// WriteIndexBlock writes the specified IndexBlock to its block address, returning
// nil on success or an error.
func (bp *BPKS) WriteIndexBlock(block *IndexBlock) error {
	// fmt.Printf("Writing Index Block at address %d (offset %d)\n", block.BlockAddress, block.BlockAddress*BlockSize)
	_, err := bp.Device.Seek(int64(block.BlockAddress*BlockSize), 0)
	if err != nil {
		return err
	}
	buffer := block.AsSlice()
	_, err = bp.Device.Write(buffer[:])
	if err != nil {
		return err
	}
	// fmt.Printf("- Wrote %d bytes\n", c)
	return nil
}

// ReadDataBlock reads and returns the DataBlock at the specified block address, returning
// a pointer to the parsed DataBlock and/or an error if any.
func (bp *BPKS) ReadDataBlock(blockAddress uint64) (*DataBlock, error) {
	// fmt.Printf("Reading Data Block at address %d (offset %d)\n", blockAddress, blockAddress*BlockSize)
	_, err := bp.Device.Seek(int64(blockAddress*BlockSize), 0)
	if err != nil {
		return nil, err
	}
	buffer := [BlockSize]byte{}
	// fmt.Printf("- Reading %d Bytes\n", BlockSize)
	_, err = bp.Device.Read(buffer[:])
	if err != nil {
		return nil, err
	}
	// fmt.Printf("- Read %d bytes\n", c)
	return NewDataBlockFromBuffer(bp, blockAddress, buffer[:]), nil
}

// WriteDataBlock writes the specified DataBlock to its block address, returning
// nil on success or an error.
func (bp *BPKS) WriteDataBlock(block *DataBlock) error {
	// fmt.Printf("Writing Data Block at address %d (offset %d)\n", block.BlockAddress, block.BlockAddress*BlockSize)
	_, err := bp.Device.Seek(int64(block.BlockAddress*BlockSize), 0)
	if err != nil {
		return err
	}
	buffer := block.AsSlice()
	_, err = bp.Device.Write(buffer[:])
	if err != nil {
		return err
	}
	// fmt.Printf("- Wrote %d bytes\n", c)
	return nil
}
