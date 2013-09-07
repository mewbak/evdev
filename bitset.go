// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package evdev

import (
	"math"
	"unsafe"
)

const (
	WordBitSize  = 64
	WordByteSize = 8
)

// A word is part of a bitset.
type Word uint64

// Bitset defines a set of bit values.
type Bitset []Word

// NewBitset creates a new bitset of the given size.
func NewBitset(bits int) Bitset {
	size := int(math.Ceil((float64(bits) / WordBitSize)))
	return make(Bitset, size)
}

// Len returns the number of bits in the set.
func (b Bitset) Len() int {
	return len(b) * WordBitSize
}

// Bytes returns the bitset as a byte slice.
// This is the same memory, so any changes to the returned slice,
// will affect the bitset.
func (b Bitset) Bytes() []byte {
	if len(b) == 0 {
		return nil
	}

	size := len(b) * WordByteSize
	return (*(*[1<<31 - 1]byte)(unsafe.Pointer(&b[0])))[:size]
}

// Set sets the bit at the given index.
func (b Bitset) Set(i int) {
	w := i / WordBitSize
	if w >= len(b) {
		return
	}

	bit := Word(1 << uint(i%WordBitSize))
	b[w] &^= bit
	b[w] ^= bit
}

// Unset clears the bit at the given index.
func (b Bitset) Unset(i int) {
	w := i / WordBitSize
	if w < len(b) {
		b[w] &^= 1 << uint(i%WordBitSize)
	}
}

// Test returns true if the bit at the given index is set.
func (b Bitset) Test(i int) bool {
	w := i / WordBitSize
	return w < len(b) && ((b[w]>>uint(i%WordBitSize))&1) == 1
}
