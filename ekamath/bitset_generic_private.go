// Copyright © 2021. All rights reserved.
// Author: Ilya Stroy.
// Contacts: qioalice@gmail.com, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekamath

import (
	"math/bits"
)

//goland:noinspection GoSnakeCaseUsage
const (
	_BITSET_GENERIC_MINIMUM_CAPACITY       = 128
	_BITSET_GENERIC_MINIMUM_CAPACITY_BYTES = _BITSET_GENERIC_MINIMUM_CAPACITY >> 3 // 16

	_BITSET_GENERIC_BITS_PER_CHUNK = bits.UintSize // 32 or 64

	_BITSET_GENERIC_CHUNK_OFFSET = 4 + (_BITSET_GENERIC_BITS_PER_CHUNK >> 5) // 5 or 6
	_BITSET_GENERIC_CHUNK_MASK   = _BITSET_GENERIC_BITS_PER_CHUNK - 1        // 31 or 63
)

// Reports whether BitSetGeneric can contain a bit with provided index.
func (bs *BitSetGeneric) isValidIdx(idx uint) bool {
	return idx > 0 && bsBytesForBits(idx) <= uint(len(bs.bs))
}

// Returns a minimum number of chunks that is required to store `n` bits.
func bsBytesForBits(n uint) uint {

	c := n >> _BITSET_GENERIC_CHUNK_OFFSET

	if n&(n-1) != 0 {
		c += 1
	}

	return c
}

// Returns a byte number (starting from 0) and bit offset (starting from 0)
// for the provided index.
func bsFromIdx(n uint) (chunk, offset uint) {
	return n >> _BITSET_GENERIC_CHUNK_OFFSET, n & _BITSET_GENERIC_CHUNK_MASK
}

// Returns an index inside a bitset (starting from 1)
// for the provided chunk (byte number, starting from 0)
// and offset (bit offset, starting from 0).
func bsToIdx(chunk, offset uint) uint {
	return (chunk << _BITSET_GENERIC_CHUNK_OFFSET) | (offset & _BITSET_GENERIC_CHUNK_MASK)
}
