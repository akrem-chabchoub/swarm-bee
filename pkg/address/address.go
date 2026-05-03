package address

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"math/bits"
)

// Size is the number of bytes in an address.
const Size = 32

// ZeoroAddress is the zero value of an address.
var ZeroAddress = Address{}

// Address is the 32-byte dientifier in the swarm address space
// XOR distance between two addresses determines routing proximity:
//   - Short XOR distance = logically close = same "neighborhood
//   - Long XOR distance = logically far = differ
type Address struct {
	b []byte
}

func NewAddress(b []byte) Address {
	if len(b) != Size {
		return ZeroAddress
	}

	cp := make([]byte, Size)
	copy(cp, b)

	return Address{b: cp}
}

// Bytes returns the nraw bytes of the address:
func (a Address) Bytes() []byte {
	if len(a.b) == 0 {
		return make([]byte, Size)
	}

	return a.b
}

// Equal returns true if both addresses have the same bytes.
func (a Address) Equal(b Address) bool {
	return bytes.Equal(a.b, b.b)
}

// IsZero returns true if the address is the zero value.
func (a Address) IsZero() bool {
	return bytes.Equal(a.b, ZeroAddress.b)
}


// String returns the hex-encoded address.
func (a Address) String() string {
	return hex.EncodeToString(a.Bytes())
}

// MarshalJSON encodes the address as a hex string in JSON.
func (a Address) MarshalJSON() ([]byte, error) {
	return json.Marshal(a.String())
}

// UnmarshalJSON decodes a hex string from JSON into an address.
func (a *Address) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	decoded, err := hex.DecodeString(s)
	if err != nil {
		return err
	}
	if len(decoded) != Size {
		return errors.New("invalid address length: must be 32 bytes")
	}
	a.b = decoded
	return nil
}

// ProximityOrder returns the number of common leading bits between a and b.
//
// This is the "bin" number in Kademlia:
//   - 0 = no common leading bits  (very far apart)
//   - 8 = first byte is identical (moderately close)
//   - 255 = identical addresses   (same node)
//
// Example:
//
//	a = 10110101 00110010 ...
//	b = 10110011 00101100 ...
//	     ^^^^^^ common (6 bits) → ProximityOrder = 6
func ProximityOrder(a, b Address) int {
	aBytes := a.Bytes()
	bBytes := b.Bytes()

	for i := 0; i < Size; i++ {
		xor := aBytes[i] ^ bBytes[i]
		if xor != 0 {
			// Found the first differing byte.
			// Count leading zeros in the XOR to get common bits in this byte.
			return i*8 + bits.LeadingZeros8(xor)
		}
	}
	return Size * 8 // all 256 bits match
}

// XORBytes computes XOR of two addresses, byte by byte.
// Useful for debugging and visualizing distance.
func XORBytes(a, b Address) []byte {
	aBytes := a.Bytes()
	bBytes := b.Bytes()
	result := make([]byte, Size)
	for i := range result {
		result[i] = aBytes[i] ^ bBytes[i]
	}
	return result
}