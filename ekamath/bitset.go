// Copyright © 2021. All rights reserved.
// Author: Ilya Stroy.
// Contacts: qioalice@gmail.com, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

// This package is re-write of
// https://github.com/bits-and-blooms/bitset
// that is distributed (for now: 2021 Sep 12) by BSD 3-Clause "New" or "Revised" License.
//
// Here's some small changes, like introducing PrevUp(), PrevDown() methods,
// using `uint` instead of `uint64` (thus it can be uint32 on 32bit platforms)
// and some other improvements and changes.

package ekamath

type (
	// BitSet is a bitset with variate capacity.
	// It can be grown, depends on your cases.
	//
	// The index of BitSet is starts from 1.
	// In almost all cases it's prohibited to use 0 as index.
	// NextUp(), NextDown() an their unsafe methods are exceptions.
	//
	// It's strongly recommend to instantiate BitSet using NewBitSet() constructor,
	// but just creating a BitSet is also possible and ready-to-use
	// (it will be with 0 capacity and will grow when you will try to set any bit).
	BitSet struct {
		bs []uint
	}
)

// IsValid reports whether current BitSet is valid.
func (bs *BitSet) IsValid() bool {
	return bs != nil
}

// IsEmpty reports whether current BitSet is empty bitset or not.
// Empty bitset is a bitset with all downed (zeroed) bits.
// Returns true if BitSet is invalid.
func (bs *BitSet) IsEmpty() bool {
	if !bs.IsValid() {
		return true
	}
	for i, n := 0, len(bs.bs); i < n; i++ {
		if bs.bs[i] != 0 {
			return false
		}
	}
	return true
}

// Capacity returns the number of values (starting from 1) that can be stored
// inside current BitSet.
// Returns 0 if current BitSet is invalid.
func (bs *BitSet) Capacity() uint {
	if !bs.IsValid() {
		return 0
	}
	return uint(len(bs.bs)) * 8
}

// Count returns number of bits that are upped (set to 1).
// Returns 0 if current BitSet is invalid.
func (bs *BitSet) Count() uint {

	if !bs.IsValid() {
		return 0
	}

	var count uint
	for i, n := 0, len(bs.bs); i < n; i++ {
		count += bsCountOnes(bs.bs[i])
	}

	return count
}

// Clear downs (zeroes) ALL bits in the current BitSet.
// Does nothing if BitSet is invalid.
func (bs *BitSet) Clear() *BitSet {
	if bs.IsValid() {
		for i, n := 0, len(bs.bs); i < n; i++ {
			bs.bs[i] = 0
		}
	}
	return bs
}

// Clone makes a copy of BitSet and returns it.
// If BitSet is invalid, NewBitSet() is called instead.
func (bs *BitSet) Clone() *BitSet {

	if !bs.IsValid() {
		return NewBitSet(0) // make BitSet with default capacity
	}

	cloned := make([]uint, len(bs.bs))
	copy(cloned, bs.bs)

	return &BitSet{
		bs: cloned,
	}
}

// GrowUnsafeUpTo grows current BitSet to be able operate with bits
// up to requested index.
// Panics if BitSet is invalid.
func (bs *BitSet) GrowUnsafeUpTo(idx uint) *BitSet {

	n := bsBytesForBits(idx)

	if l := uint(len(bs.bs)); l == 0 {
		bs.bs = make([]uint, MaxU(_BITSET_GENERIC_MINIMUM_CAPACITY_BYTES, n))

	} else if l <= n {
		old := bs.bs
		bs.bs = make([]uint, n)
		copy(bs.bs, old)
	}

	return bs
}

// Up sets bit to 1 with requested index checking bounds,
// growing bitset if it's too small. Does nothing if BitSet is invalid.
func (bs *BitSet) Up(idx uint) *BitSet {
	if bs.IsValid() && bs.isValidIdx(idx, 1, true) {
		bs.GrowUnsafeUpTo(idx).UpUnsafe(idx)
	}
	return bs
}

// UpUnsafe sets bit to 1 with requested index without any check.
// Panics if BitSet is invalid or if an index is out of bounds.
func (bs *BitSet) UpUnsafe(idx uint) *BitSet {
	chunk, offset := bsFromIdx(idx-1)
	bs.bs[chunk] |= 1 << offset
	return bs
}

// Down sets bit to 0 with requested index checking bounds,
// growing bitset if it's too small. Does nothing if BitSet is invalid.
func (bs *BitSet) Down(idx uint) *BitSet {
	if bs.IsValid() && bs.isValidIdx(idx, 1, true) {
		bs.GrowUnsafeUpTo(idx).DownUnsafe(idx)
	}
	return bs
}

// DownUnsafe sets bit to 0 with requested index without any check.
// Panics if BitSet is invalid or if an index is out of bounds.
func (bs *BitSet) DownUnsafe(idx uint) *BitSet {
	chunk, offset := bsFromIdx(idx-1)
	bs.bs[chunk] &^= 1 << offset
	return bs
}

// Set calls Up() or Down() with provided index depends on `b`.
func (bs *BitSet) Set(idx uint, b bool) *BitSet {
	if bs.IsValid() && bs.isValidIdx(idx, 1, true) {
		bs.GrowUnsafeUpTo(idx).SetUnsafe(idx, b)
	}
	return bs
}

// SetUnsafe calls UpUnsafe() or DownUnsafe() with provided index
// depends on `b`.
func (bs *BitSet) SetUnsafe(idx uint, b bool) *BitSet {
	if b {
		return bs.UpUnsafe(idx)
	} else {
		return bs.DownUnsafe(idx)
	}
}

// IsSet reports whether a bit with requested index is set or not.
// Returns false either if bit isn't set, BitSet is invalid or index is out of bound.
func (bs *BitSet) IsSet(idx uint) bool {
	return bs.IsValid() && bs.isValidIdx(idx, 1, false) && bs.IsSetUnsafe(idx)
}

// IsSetUnsafe reports whether a bit with requested index is set or not.
// Panics if BitSet is invalid or if an index is out of bound.
func (bs *BitSet) IsSetUnsafe(idx uint) bool {
	chunk, offset := bsFromIdx(idx-1)
	return bs.bs[chunk]&(1<<offset) != 0
}

// NextUp returns an index of next upped (set to 1) bit.
// It's safe to use 0 as index because this is the only way to get 1st bit.
//
// You can use this method (and the similar methods, like NextDown(), PrevUp(), PrevDown())
// to iterate over BitSet:
//
//  for v, e := bs.NextUp(0); e; v, e = bs.NextUp(v+1) {
//      fmt.Printf("Elem: %d\n", v)
//  }
//
func (bs *BitSet) NextUp(idx uint) (uint, bool) {
	if bs.IsValid() && bs.isValidIdx(idx, 0, false) {
		return bs.NextUpUnsafe(idx)
	}
	return idx, false
}

// NextUpUnsafe is the same as NextUp but without any bound checks.
// It will lead to UB or panic if you will use an incorrect index.
func (bs *BitSet) NextUpUnsafe(idx uint) (uint, bool) {
	return bs.nextGeneric(idx, false)
}

// NextDown returns an index of next downed (set to 0) bit.
// It's safe to use 0 as index because this is the only way to get 1st bit.
// See NextUp() method to get to know how to use that method to iterate over BitSet.
func (bs *BitSet) NextDown(idx uint) (uint, bool) {
	if bs.IsValid() && bs.isValidIdx(idx, 0, false) {
		return bs.NextDownUnsafe(idx)
	}
	return idx, false
}

// NextDownUnsafe is the same as NextDown but without any bound checks.
// It will lead to UB or panic if you will use an incorrect index.
func (bs *BitSet) NextDownUnsafe(idx uint) (uint, bool) {
	return bs.nextGeneric(idx, true)
}

// PrevUp returns an index of prev upped (set to 1) bit.
// The minimum index you should use to get not false 2nd return argument is 2
// (the minimum index of BitSet is 1 and there's no upped bits before 1 or 0).
// See NextUp() method to get to know how to use that method to iterate over BitSet.
func (bs *BitSet) PrevUp(idx uint) (uint, bool) {
	if bs.IsValid() && bs.isValidIdx(idx, 2, false) {
		return bs.PrevUpUnsafe(idx)
	}
	return idx, false
}

// PrevUpUnsafe is the same as PrevUp but without any bound checks.
// It will lead to UB or panic if you will use an incorrect index.
func (bs *BitSet) PrevUpUnsafe(idx uint) (uint, bool) {
	return bs.prevGeneric(idx, false)
}

// PrevDown returns an index of prev downed (set to 0) bit.
// The minimum index you should use to get not false 2nd return argument is 2
// (the minimum index of BitSet is 1 and there's no downed bits before 1 or 0).
// See NextUp() method to get to know how to use that method to iterate over BitSet.
func (bs *BitSet) PrevDown(idx uint) (uint, bool) {
	if bs.IsValid() && bs.isValidIdx(idx, 2, false) {
		return bs.PrevDownUnsafe(idx)
	}
	return idx, false
}

// PrevDownUnsafe is the same as PrevDown but without any bound checks.
// It will lead to UB or panic if you will use an incorrect index.
func (bs *BitSet) PrevDownUnsafe(idx uint) (uint, bool) {
	return bs.prevGeneric(idx, true)
}

// NewBitSet creates a new BitSet with desired initial capacity.
// If capacity is too small, it will be overwritten with default minimum capacity.
func NewBitSet(capacity uint) *BitSet {
	return new(BitSet).GrowUnsafeUpTo(capacity)
}
