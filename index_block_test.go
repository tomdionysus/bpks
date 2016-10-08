package bpks

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func xTestIndexBlockAdd(t *testing.T) {
	disk := NewRAMDisk(4 * 1024 * 1024)
	kvs := New(disk)
	err := kvs.Format()
	assert.Nil(t, err)
	kvs.Mount()

	k := KeyPointer{
		Key:          Key{0, 0, 1, 2, 3, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		BlockAddress: 81723123,
	}
	err = kvs.Root.Add(k)
	assert.Nil(t, err)

	assert.Equal(t, 1, len(*kvs.Root.KeyPointerList))
	assert.Equal(t, k.Key, (*kvs.Root.KeyPointerList)[0].Key)
	assert.Equal(t, k.BlockAddress, (*kvs.Root.KeyPointerList)[0].BlockAddress)

	k2 := KeyPointer{
		Key:          Key{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 2, 3, 0, 0, 0},
		BlockAddress: 81723123,
	}
	err = kvs.Root.Add(k2)
	assert.Nil(t, err)

	assert.Equal(t, 2, len(*kvs.Root.KeyPointerList))
	assert.Equal(t, k2.Key, (*kvs.Root.KeyPointerList)[0].Key)
	assert.Equal(t, k2.BlockAddress, (*kvs.Root.KeyPointerList)[0].BlockAddress)
	assert.Equal(t, k.Key, (*kvs.Root.KeyPointerList)[1].Key)
	assert.Equal(t, k.BlockAddress, (*kvs.Root.KeyPointerList)[1].BlockAddress)

	k3 := KeyPointer{
		Key:          Key{1, 1, 1, 2, 3, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		BlockAddress: 81723123,
	}
	err = kvs.Root.Add(k3)
	assert.Nil(t, err)

	assert.Equal(t, 3, len(*kvs.Root.KeyPointerList))
	assert.Equal(t, k2.Key, (*kvs.Root.KeyPointerList)[0].Key)
	assert.Equal(t, k2.BlockAddress, (*kvs.Root.KeyPointerList)[0].BlockAddress)
	assert.Equal(t, k.Key, (*kvs.Root.KeyPointerList)[1].Key)
	assert.Equal(t, k.BlockAddress, (*kvs.Root.KeyPointerList)[1].BlockAddress)
	assert.Equal(t, k3.Key, (*kvs.Root.KeyPointerList)[2].Key)
	assert.Equal(t, k3.BlockAddress, (*kvs.Root.KeyPointerList)[2].BlockAddress)
}

func TestIndexBlockAddSplit(t *testing.T) {
	fmt.Printf("\n\n------------------\n\n")
	disk := NewRAMDisk(4 * 1024 * 1024)
	kvs := New(disk)
	err := kvs.Format()
	assert.Nil(t, err)
	kvs.Mount()

	for i := 0; i < 256; i++ {
		k := KeyPointer{
			Key:          NewKeyRandom(),
			BlockAddress: 81723123,
		}
		err = kvs.Root.Add(k)
		assert.Nil(t, err)
	}

	kv3 := Key{1, 1, 1, 2, 3, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	k3 := KeyPointer{
		Key:          kv3,
		BlockAddress: 81723123,
	}
	err = kvs.Root.Add(k3)

	assert.Equal(t, 84, len(*kvs.Root.KeyPointerList))

	kp, found, err := kvs.Root.Find(kv3)
	assert.Nil(t, err)
	assert.True(t, found)
	assert.Equal(t, kv3, kp.Key)
	fmt.Printf("\n\n------------------\n\n")
}
