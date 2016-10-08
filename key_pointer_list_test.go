package bpks

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestKeyPointerListAdd(t *testing.T) {
	lst := &KeyPointerList{}

	k := KeyPointer{
		Key:          Key{0, 0, 1, 2, 3, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		BlockAddress: 81723123,
	}
	lst.Add(k)

	assert.Equal(t, 1, len(*lst))

	k2 := KeyPointer{
		Key:          Key{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 2, 3, 0, 0, 0},
		BlockAddress: 189247182,
	}
	lst.Add(k2)

	assert.Equal(t, 2, len(*lst))

	k3 := KeyPointer{
		Key:          Key{1, 2, 3, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		BlockAddress: 982374223,
	}
	lst.Add(k3)

	assert.Equal(t, 3, len(*lst))
}

func TestKeyPointerFind(t *testing.T) {
	lst := &KeyPointerList{}
	key := Key{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 2, 3, 0, 0, 0}

	k := KeyPointer{
		Key:          key,
		BlockAddress: 81723123,
	}
	lst.Add(k)
	assert.Equal(t, 1, len(*lst))

	kp, found := lst.Find(key)
	assert.True(t, found)
	assert.Equal(t, k, kp)
}

func TestKeyPointerFindNotFound(t *testing.T) {
	lst := &KeyPointerList{}
	key := Key{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 2, 3, 0, 0, 0}

	k := KeyPointer{
		Key:          key,
		BlockAddress: 81723123,
	}
	lst.Add(k)
	assert.Equal(t, 1, len(*lst))

	_, found := lst.Find(Key{0, 0, 1, 0xFF, 0, 0, 0, 0, 0, 0, 1, 2, 3, 0, 0, 0})
	assert.False(t, found)
}

func TestKeyPointerListRemove(t *testing.T) {
	lst := &KeyPointerList{}
	key := Key{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 2, 3, 0, 0, 0}

	k := KeyPointer{
		Key:          key,
		BlockAddress: 81723123,
	}
	lst.Add(k)
	assert.Equal(t, 1, len(*lst))

	kp, found := lst.Remove(key)
	assert.True(t, found)
	assert.Equal(t, 0, len(*lst))
	assert.Equal(t, kp.Key, key)
}

func TestKeyPointerListRemoveMidSet(t *testing.T) {
	lst := &KeyPointerList{}

	k := KeyPointer{
		Key:          Key{0, 0, 1, 2, 3, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		BlockAddress: 81723123,
	}
	lst.Add(k)

	assert.Equal(t, 1, len(*lst))

	k2 := KeyPointer{
		Key:          Key{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 2, 3, 0, 0, 0},
		BlockAddress: 189247182,
	}
	lst.Add(k2)

	assert.Equal(t, 2, len(*lst))

	k3 := KeyPointer{
		Key:          Key{1, 2, 3, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		BlockAddress: 982374223,
	}
	lst.Add(k3)

	assert.Equal(t, 3, len(*lst))

	kp, found := lst.Remove(k2.Key)
	assert.True(t, found)
	assert.Equal(t, 2, len(*lst))
	assert.Equal(t, k2, kp)

	kp, found = lst.Find(k.Key)
	assert.True(t, found)
	assert.Equal(t, k, kp)

	kp, found = lst.Find(k3.Key)
	assert.True(t, found)
	assert.Equal(t, k3, kp)
}

func TestKeyPointerListRemoveFirst(t *testing.T) {
	lst := &KeyPointerList{}

	k := KeyPointer{
		Key:          Key{0, 0, 1, 2, 3, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		BlockAddress: 81723123,
	}
	lst.Add(k)

	assert.Equal(t, 1, len(*lst))

	k2 := KeyPointer{
		Key:          Key{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 2, 3, 0, 0, 0},
		BlockAddress: 189247182,
	}
	lst.Add(k2)

	assert.Equal(t, 2, len(*lst))

	k3 := KeyPointer{
		Key:          Key{1, 2, 3, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		BlockAddress: 982374223,
	}
	lst.Add(k3)

	assert.Equal(t, 3, len(*lst))

	kp, found := lst.Remove(k.Key)
	assert.True(t, found)
	assert.Equal(t, 2, len(*lst))
	assert.Equal(t, k, kp)

	kp, found = lst.Find(k2.Key)
	assert.True(t, found)
	assert.Equal(t, k2, kp)

	kp, found = lst.Find(k3.Key)
	assert.True(t, found)
	assert.Equal(t, k3, kp)
}
