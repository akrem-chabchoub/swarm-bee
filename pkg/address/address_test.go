package address

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

)

func TestNewAddress(t *testing.T) {
	t.Run("valid 32 bytes", func(t *testing.T) {
		b := make([]byte, 32)
		b[0] = 0xAB
		addr := NewAddress(b)
		assert.Equal(t, b, addr.Bytes())
	})

	t.Run("wrong length returns zero address", func(t *testing.T) {
		addr := NewAddress([]byte{1, 2, 3})
		assert.True(t, addr.IsZero())
	})

	t.Run("nil input returns zero address", func(t *testing.T) {
		addr := NewAddress(nil)
		assert.True(t, addr.IsZero())
	})
}

func TestAddressEqual(t *testing.T) {
	b1 := make([]byte, 32)
	b1[0] = 0xFF

	b2 := make([]byte, 32)
	b2[0] = 0xFF

	b3 := make([]byte, 32)
	b3[0] = 0x00

	assert.True(t, NewAddress(b1).Equal(NewAddress(b2)))
	assert.False(t, NewAddress(b1).Equal(NewAddress(b3)))
	assert.True(t, ZeroAddress.Equal(ZeroAddress))
}

func TestAddressString(t *testing.T) {
	b := make([]byte, 32)
	b[0] = 0xDE
	b[1] = 0xAD
	addr := NewAddress(b)
	// First two bytes should be "dead" in hex
	assert.Contains(t, addr.String(), "dead")
}

func TestAddressJSON(t *testing.T) {
	b := make([]byte, 32)
	for i := range b {
		b[i] = byte(i)
	}
	original := NewAddress(b)

	// Marshal to JSON
	data, err := json.Marshal(original)
	require.NoError(t, err)

	// Unmarshal back
	var recovered Address
	err = json.Unmarshal(data, &recovered)
	require.NoError(t, err)

	assert.True(t, original.Equal(recovered))
}

func TestProximityOrder(t *testing.T) {
	// Test: identical addresses have maximum proximity
	b := make([]byte, 32)
	b[0] = 0xFF
	a1 := NewAddress(b)
	a2 := NewAddress(b)
	assert.Equal(t, 256, ProximityOrder(a1, a2))

	// Test: addresses differing only in last bit have proximity 255
	b1 := make([]byte, 32)
	b2 := make([]byte, 32)
	b2[31] = 0x01 // flip last bit
	assert.Equal(t, 255, ProximityOrder(
		NewAddress(b1),
		NewAddress(b2),
	))

	// Test: addresses differing in first bit have proximity 0
	b3 := make([]byte, 32)
	b4 := make([]byte, 32)
	b3[0] = 0b10000000 // starts with 1
	b4[0] = 0b00000000 // starts with 0
	assert.Equal(t, 0, ProximityOrder(
		NewAddress(b3),
		NewAddress(b4),
	))

	// Test: addresses sharing first 3 bits have proximity 3
	// 0b11100000 = 0xE0
	// 0b11011111 = 0xDF
	// XOR = 0b00111111 — leading zeros in XOR = 2... wait
	// Let me recalculate:
	// 0b11100000
	// 0b11010000
	// XOR = 0b00110000 → leading zeros = 2 → proximity = 2
	b5 := make([]byte, 32)
	b6 := make([]byte, 32)
	b5[0] = 0b11100000
	b6[0] = 0b11010000
	assert.Equal(t, 2, ProximityOrder(
		NewAddress(b5),
		NewAddress(b6),
	))
}

func TestProximityOrderBinAssignment(t *testing.T) {
	// Demonstrate how ProximityOrder maps to Kademlia bins.
	// A peer in bin N shares exactly N leading bits with us.
	self := make([]byte, 32)
	self[0] = 0b10000000

	// Peer in bin 0: first bit differs
	peerBin0 := make([]byte, 32)
	peerBin0[0] = 0b00000000
	assert.Equal(t, 0, ProximityOrder(
		NewAddress(self),
		NewAddress(peerBin0),
	))

	// Peer in bin 1: shares first bit, second differs
	peerBin1 := make([]byte, 32)
	peerBin1[0] = 0b11000000 // wait, self starts with 1, peer needs to start with 1 too
	// 10000000 XOR 11000000 = 01000000 → leading zeros = 1 → bin 1 ✓
	assert.Equal(t, 1, ProximityOrder(
		NewAddress(self),
		NewAddress(peerBin1),
	))
}
