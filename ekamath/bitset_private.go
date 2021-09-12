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
	_BITSET_MINIMUM_CAPACITY       = 128
	_BITSET_MINIMUM_CAPACITY_BYTES = _BITSET_MINIMUM_CAPACITY >> 3 // 16

	_BITSET_BITS_PER_CHUNK = bits.UintSize // 32 or 64

	_BITSET_CHUNK_OFFSET = 4 + (_BITSET_BITS_PER_CHUNK >> 5) // 5 or 6
	_BITSET_CHUNK_MASK   = _BITSET_BITS_PER_CHUNK - 1        // 31 or 63

	_BITSET_MASK_FULL = ^uint(0)
)

// Returns a chunk number of the current BitSet.
func (bs *BitSet) chunkSize() uint {
	return uint(len(bs.bs))
}

// Returns a chunk capacity of the current BitSet.
func (bs *BitSet) chunkCapacity() uint {
	return uint(cap(bs.bs))
}

// Reports whether BitSet can contain a bit with provided index.
// It includes IsValid() call, so you don't need to call it explicitly.
func (bs *BitSet) isValidIdx(idx uint, lowerBound uint, skipUpperBoundCheck bool) bool {
	return bs.IsValid() && idx >= lowerBound &&
		(skipUpperBoundCheck || bsChunksForBits(idx+1) <= bs.chunkSize())
}

// Returns a next upped or downed bit index depends on `f`.
func (bs *BitSet) nextGeneric(idx uint, isDown bool) (uint, bool) {

	chunk, offset := uint(0), uint(0)
	if idx != 0 {
		chunk, offset = bsFromIdx(idx)
	}

	v := bs.bs[chunk] >> offset
	if isDown {
		v = ^v
	}
	if n := bs1stUp(v); n < _BITSET_BITS_PER_CHUNK - offset {
		return bsToIdx(chunk, offset+n) + 1, true
	}

	for i, n := chunk+1, uint(len(bs.bs)); i < n; i++ {
		v := bs.bs[i]
		if isDown {
			v = ^v
		}
		if n := bs1stUp(v); n != _BITSET_BITS_PER_CHUNK {
			return bsToIdx(i, n) + 1, true
		}
	}

	return idx, false
}

// Returns a prev upped or downed bit index depends on `f1`, `f2`.
func (bs *BitSet) prevGeneric(idx uint, isDown bool) (uint, bool) {

	chunk, offset := bsFromIdx(idx)

	v := bs.bs[chunk] << (_BITSET_BITS_PER_CHUNK - offset + 1)
	if isDown {
		v = ^v
	}
	if n := bsLastUp(v); n < _BITSET_BITS_PER_CHUNK-offset {
		return bsToIdx(chunk, offset-n-1), true
	}

	for i := int(chunk) - 1; i >= 0; i-- {
		v := bs.bs[i]
		if isDown {
			v = ^v
		}
		if n := bsLastUp(v); n != _BITSET_BITS_PER_CHUNK {
			return bsToIdx(uint(i), _BITSET_BITS_PER_CHUNK-n), true
		}
	}

	return idx, false
}

// Returns a number of bits that are upped (set to 1).
func bsCountOnes(n uint) uint {
	return uint(bits.OnesCount(n))
}

// Returns a minimum number of chunks that is required to store `n` bits.
func bsChunksForBits(n uint) uint {

	c := n >> _BITSET_CHUNK_OFFSET

	if n & _BITSET_CHUNK_MASK != 0 {
		c += 1
	}

	return c
}

// Returns a byte number (starting from 0) and bit offset (starting from 0)
// for the provided index.
func bsFromIdx(n uint) (chunk, offset uint) {
	return n >> _BITSET_CHUNK_OFFSET, n & _BITSET_CHUNK_MASK
}

// Returns an index inside a bitset (starting from 1)
// for the provided chunk (byte number, starting from 0)
// and offset (bit offset, starting from 0).
func bsToIdx(chunk, offset uint) uint {
	return (chunk << _BITSET_CHUNK_OFFSET) | (offset & _BITSET_CHUNK_MASK)
}

// Returns an index (starting from 0) of 1st upped (set to 1) bit.
// If there's no such bits, _BITSET_BITS_PER_CHUNK is returned.
func bs1stUp(n uint) uint {
	return uint(bits.TrailingZeros(n))
}

// Returns an index (starting from 0) of 1st downed (set to 0) bit.
// If there's no such bits, _BITSET_BITS_PER_CHUNK is returned.
func bsLastUp(n uint) uint {
	return uint(bits.LeadingZeros(n))
}
