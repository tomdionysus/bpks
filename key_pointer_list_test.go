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
