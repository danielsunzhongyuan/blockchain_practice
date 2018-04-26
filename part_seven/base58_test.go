package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBase58Encode(t *testing.T) {
	data := []byte{0x00, 0x00, 0x00}
	encoded := Base58Encode(data)
	assert.Equal(t, []byte{'1', '1', '1'}, encoded)

	data = []byte{'a'}
	encoded = Base58Encode(data)
	assert.Equal(t, []byte{'2', 'g'}, encoded)

	data = []byte{'1', 0x00}
	encoded = Base58Encode(data)
	assert.Equal(t, []byte{'4', 'j', 'H'}, encoded)
}

func TestBase58Decode(t *testing.T) {
	data := []byte{0x00, 0x00, 0x00}
	assert.Equal(t, data, Base58Decode(Base58Encode(data)))

	data = []byte{'a'}
	assert.Equal(t, data, Base58Decode(Base58Encode(data)))

	data = []byte{'1', 0x00}
	assert.Equal(t, data, Base58Decode(Base58Encode(data)))
}
