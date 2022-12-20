// Copyright © 2021. All rights reserved.
// Author: Ilya Stroy.
// Contacts: iyuryevich@pm.me, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

// This package is re-write of
// https://github.com/bits-and-blooms/bitset
// that is distributed (for now: 2021 Sep 12)
// by BSD 3-Clause "New" or "Revised" License.
//
// Here's some small changes, like introducing PrevUp(), PrevDown() methods,
// using `uint` instead of `uint64` (thus it can be uint32 on 32bit platforms)
// and some other improvements and changes.

package ekamath

import (
	"encoding/base64"
	"errors"
)

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
type BitSet struct {
	bs []uint
}

var (
	ErrBitSetInvalid             = errors.New("invalid BitSet")
	ErrBitSetInvalidDataToDecode = errors.New("invalid data to decode to BitSet")
)

// ---------------------------------------------------------------------------- //

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

// ---------------------------------------------------------------------------- //

// Capacity returns the number of values (starting from 1) that can be stored
// inside current BitSet.
// Returns 0 if current BitSet is invalid.
func (bs *BitSet) Capacity() uint {
	if !bs.IsValid() {
		return 0
	}
	return uint(len(bs.bs)) * _BITSET_BITS_PER_CHUNK
}

// Count returns number of bits that are upped (set to 1).
// Returns 0 if current BitSet is invalid.
func (bs *BitSet) Count() uint {

	if !bs.IsValid() {
		return 0
	}

	var c uint
	for i, n := uint(0), bs.chunkSize(); i < n; i++ {
		c += bsCountOnes(bs.bs[i])
	}

	return c
}

// CountBetween returns number of bits that are upped (set to 1),
// between range [a..b]. Note: `b` IS IN the range.
// Returns 0 if either current BitSet is invalid, `a` >= `b`
// or any part of that range is out of bound of the BitSet.
func (bs *BitSet) CountBetween(a, b uint) uint {

	if !bs.IsValid() || a >= b {
		return 0
	}

	c1, off1 := bsFromIdx(a - 1)
	c2, off2 := bsFromIdx(b - 1)

	if bs1size := bs.chunkSize(); c1 > bs1size || c2 > bs1size {
		return 0
	}

	// The case when chunk is the same.
	if c1 == c2 {
		mask := (_BITSET_MASK_FULL >> (_BITSET_BITS_PER_CHUNK - off2 + off1 - 1)) << off1
		return bsCountOnes(bs.bs[c1] & mask)
	}

	// There are 3 parts of counting:
	// 1. Discard unnecessary bits and count remains of the chunk of `A` bound,
	// 2. Discard unnecessary bits and count remains of the chunk of `B` bound,
	// 3. Count all in the chunks between those that belongs either `A` or `B` bound.
	var c uint

	// Part 1.
	mask := _BITSET_MASK_FULL << off1
	c += bsCountOnes(bs.bs[c1] & mask)

	// Part 2.
	mask = _BITSET_MASK_FULL >> (_BITSET_BITS_PER_CHUNK - off2 - 1)
	c += bsCountOnes(bs.bs[c2] & mask)

	// Part 3.
	for c1++; c1 < c2; c1++ {
		c += bsCountOnes(bs.bs[c1])
	}

	return c
}

// ---------------------------------------------------------------------------- //

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

	n := bsChunksForBits(idx)

	if l, c := bs.chunkSize(), bs.chunkCapacity(); c == 0 {
		bs.bs = make([]uint, n)

	} else if l <= n {
		if c >= n {
			bs.bs = bs.bs[:n]

		} else {
			old := bs.bs
			bs.bs = make([]uint, n)
			copy(bs.bs, old)
		}
	}

	return bs
}

// ShrinkUpTo shrinks current BitSet to be able operate with bits
// up to requested index.
// Does nothing if current BitSet is less than the requested index.
// Panics if BitSet is invalid.
func (bs *BitSet) ShrinkUpTo(idx uint) *BitSet {

	chunk, offset := bsFromIdx(idx - 1)

	if bs1size := bs.chunkSize(); chunk <= bs1size {

		for i := chunk + 1; i < bs1size; i++ {
			bs.bs[i] = 0
		}

		bs.bs[chunk] &= _BITSET_MASK_FULL >> (_BITSET_BITS_PER_CHUNK - offset - 1)
		bs.bs = bs.bs[:chunk+1]
	}

	return bs
}

// ---------------------------------------------------------------------------- //

// Up sets bit to 1 with requested index checking bounds,
// growing bitset if it's too small. Does nothing if BitSet is invalid.
func (bs *BitSet) Up(idx uint) *BitSet {
	if bs.isValidIdx(idx, 1, true) {
		bs.GrowUnsafeUpTo(idx).UpUnsafe(idx)
	}
	return bs
}

// UpUnsafe sets bit to 1 with requested index without any check.
// Panics if BitSet is invalid or if an index is out of bounds.
func (bs *BitSet) UpUnsafe(idx uint) *BitSet {
	chunk, offset := bsFromIdx(idx - 1)
	bs.bs[chunk] |= 1 << offset
	return bs
}

// Down sets bit to 0 with requested index checking bounds,
// growing bitset if it's too small. Does nothing if BitSet is invalid.
func (bs *BitSet) Down(idx uint) *BitSet {
	if bs.isValidIdx(idx, 1, true) {
		bs.GrowUnsafeUpTo(idx).DownUnsafe(idx)
	}
	return bs
}

// DownUnsafe sets bit to 0 with requested index without any check.
// Panics if BitSet is invalid or if an index is out of bounds.
func (bs *BitSet) DownUnsafe(idx uint) *BitSet {
	chunk, offset := bsFromIdx(idx - 1)
	bs.bs[chunk] &^= 1 << offset
	return bs
}

// Set calls Up() or Down() with provided index depends on `b`.
func (bs *BitSet) Set(idx uint, b bool) *BitSet {
	if bs.isValidIdx(idx, 1, true) {
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

// Invert changes bit to against value with requested index checking bounds,
// growing bitset if it's too small. Does nothing if BitSet is invalid.
func (bs *BitSet) Invert(idx uint) *BitSet {
	if bs.isValidIdx(idx, 1, true) {
		bs.GrowUnsafeUpTo(idx).InvertUnsafe(idx)
	}
	return bs
}

// InvertUnsafe changes bit to against value with requested index without any check.
// Panics if BitSet is invalid or if an index is out of bounds.
func (bs *BitSet) InvertUnsafe(idx uint) *BitSet {
	chunk, offset := bsFromIdx(idx - 1)
	bs.bs[chunk] ^= 1 << offset
	return bs
}

// IsSet reports whether a bit with requested index is set or not.
// Returns false either if bit isn't set, BitSet is invalid or index is out of bound.
func (bs *BitSet) IsSet(idx uint) bool {
	return bs.isValidIdx(idx, 1, false) && bs.IsSetUnsafe(idx)
}

// IsSetUnsafe reports whether a bit with requested index is set or not.
// Panics if BitSet is invalid or if an index is out of bound.
func (bs *BitSet) IsSetUnsafe(idx uint) bool {
	chunk, offset := bsFromIdx(idx - 1)
	return bs.bs[chunk]&(1<<offset) != 0
}

// ---------------------------------------------------------------------------- //

// NextUp returns an index of next upped (set to 1) bit.
// It's safe to use 0 as index because this is the only way to get 1st bit.
//
// You can use this method (and the similar methods, like NextDown(), PrevUp(), PrevDown())
// to iterate over BitSet:
//
//	for v, e := bs.NextUp(0); e; v, e = bs.NextUp(v) {
//	    fmt.Printf("Elem: %d\n", v)
//	}
func (bs *BitSet) NextUp(idx uint) (uint, bool) {
	if bs.isValidIdx(idx, 0, false) {
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
	if bs.isValidIdx(idx, 0, false) {
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
	if bs.isValidIdx(idx, 2, false) {
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
	if bs.isValidIdx(idx, 2, false) {
		return bs.PrevDownUnsafe(idx)
	}
	return idx, false
}

// PrevDownUnsafe is the same as PrevDown but without any bound checks.
// It will lead to UB or panic if you will use an incorrect index.
func (bs *BitSet) PrevDownUnsafe(idx uint) (uint, bool) {
	return bs.prevGeneric(idx, true)
}

// ---------------------------------------------------------------------------- //

// Complement makes a complement operation (invert all bits),
// saving result to the current BitSet and returns it.
// If you want to use current BitSet after this operation, just make a copy before.
// Read more: https://en.wikipedia.org/wiki/Complement_(set_theory)
//
// Does nothing if either current BitSet is invalid or has 0 capacity.
func (bs *BitSet) Complement() *BitSet {

	if bs.IsValid() {
		for i, n := uint(0), bs.chunkSize(); i < n; i++ {
			bs.bs[i] ^= _BITSET_MASK_FULL
		}
	}

	return bs
}

// Union makes a union operation, saving result to the current BitSet and returns it.
// If you want to use current BitSet after this operation, just make a copy before.
// Read more: https://en.wikipedia.org/wiki/Union_(set_theory)
//
// If `bs2` BitSet's capacity > current's one,
// the current BitSet will be grown up to that BitSet capacity.
//
// Does nothing if either current BitSet or provided one is invalid.
func (bs *BitSet) Union(bs2 *BitSet) *BitSet {

	// Capacity() includes IsValid() call
	if bs2cap := bs2.Capacity(); bs.IsValid() && bs2cap > 0 {

		bs.GrowUnsafeUpTo(bs2cap)
		for i, n := uint(0), bs2.chunkSize(); i < n; i++ {
			bs.bs[i] |= bs2.bs[i]
		}
	}

	return bs
}

// Intersection makes an intersection operation, saving result to the current BitSet
// and returns it.
// If you want to use current BitSet after this operation, just make a copy before.
// Read more: https://en.wikipedia.org/wiki/Intersection_(set_theory)
//
// If `bs2` BitSet has bits out of the upper bound of the current BitSet,
// they will be ignored.
// If current BitSet has bits out of the upper bound of `bs2` BitSet,
// they will be zeroed.
//
// Does nothing if either current BitSet or provided one is invalid.
func (bs *BitSet) Intersection(bs2 *BitSet) *BitSet {

	if bs.IsValid() && bs2.IsValid() {

		bs1size := bs.chunkSize()

		i := uint(0)
		for n := Min(bs1size, bs2.chunkSize()); i < n; i++ {
			bs.bs[i] &= bs2.bs[i]
		}

		for ; i < bs1size; i++ {
			bs.bs[i] = 0
		}
	}

	return bs
}

// Difference performs a difference operation, saving result to the current BitSet
// and returns it.
// If you want to use current BitSet after this operation, just make a copy before.
// Read more: https://en.wikipedia.org/wiki/Complement_(set_theory)#Relative_complement
//
// If either `bs2` BitSet or current BitSet has bits out of the upper bound
// of the other BitSet, they will be ignored.
//
// Does nothing if either current BitSet or provided one is invalid.
func (bs *BitSet) Difference(bs2 *BitSet) *BitSet {

	if bs.IsValid() && bs2.IsValid() {

		for i, n := uint(0), Min(bs.chunkSize(), bs2.chunkSize()); i < n; i++ {
			bs.bs[i] &^= bs2.bs[i]
		}
	}

	return bs
}

// SymmetricDifference performs a symmetric difference (XOR) operation,
// saving result to the current BitSet and returns it.
// If you want to use current BitSet after this operation, just make a copy before.
// Read more: https://en.wikipedia.org/wiki/Symmetric_difference
//
// If `bs2` BitSet's capacity > current's one,
// the current BitSet will be grown up to that BitSet capacity.
//
// Does nothing if either current BitSet or provided one is invalid.
func (bs *BitSet) SymmetricDifference(bs2 *BitSet) *BitSet {

	// Capacity() includes IsValid() call
	if bs2cap := bs2.Capacity(); bs.IsValid() && bs2cap != 0 {

		bs1size := bs.chunkSize()
		bs2size := bs2.chunkSize()

		bs.GrowUnsafeUpTo(bs2cap)

		i := uint(0)
		for n := Min(bs1size, bs2size); i < n; i++ {
			bs.bs[i] ^= bs2.bs[i]
		}

		for ; i < bs2size; i++ {
			bs.bs[i] |= bs2.bs[i]
		}
	}

	return bs
}

// ---------------------------------------------------------------------------- //

// MarshalBinary implements BinaryMarshaler interface encoding current BitSet
// in binary form.
//
// It guarantees that if BitSet is valid, the MarshalBinary() cannot fail.
// There's no guarantees about algorithm that will be used to encode/decode.
// Returns nil if BitSet doesn't have any presented underlying chunks.
//
// WARNING!
// User MUST NOT modify returned data. If you need it, clone it firstly.
//
// Despite of warnings and restrictions, this method has O(1) complexity.
func (bs *BitSet) MarshalBinary() ([]byte, error) {

	if !bs.IsValid() {
		return nil, ErrBitSetInvalid
	}

	return bsUnsafeToBytesSlice(bs.bs), nil
}

// UnmarshalBinary implements BinaryUnmarshaler interface decoding provided `data`
// from binary form.
//
// The current BitSet's data will be overwritten by the decoded one
// if decoding operation has been completed successfully.
//
// Provided `data` MUST BE obtained by calling BitSet.MarshalBinary() method.
// There's no guarantees about algorithm that will be used to encode/decode.
//
// Does nothing (and returns nil) if provided `data` is empty.
// Returns ErrBitSetInvalidDataToDecode if provided data is invalid.
//
// WARNING!
// User MUST NOT use provided `data` after passing to this method. UB otherwise.
//
// Despite of warnings and restrictions, this method has O(1) complexity.
func (bs *BitSet) UnmarshalBinary(data []byte) error {

	switch l := len(data); {
	case l == 0:
		return nil

	case bs == nil:
		return ErrBitSetInvalid

	case l&_BITSET_BYTES_PER_CHUNK != 0:
		return ErrBitSetInvalidDataToDecode
	}

	bs.bs = bsUnsafeFromBytesSlice(data)
	return nil
}

// MarshalText implements TextMarshaler interface encoding current BitSet
// in text form.
//
// It guarantees that if BitSet is valid, the MarshalText() cannot fail.
// Returns nil if BitSet doesn't have any presented underlying chunks.
//
// MarshalText guarantees that output data will be base64 encoded,
// but NOT GUARANTEES that decoded data from base64 is user-friendly
// and user can manually read/construct BitSet from that data.
//
// Provided base64 data is NO URL FRIENDLY!
// It means, that was used something but base64.URLEncoding.
func (bs *BitSet) MarshalText() ([]byte, error) {

	binaryEncodedData, err := bs.MarshalBinary()
	if err != nil {
		return nil, err
	}

	buf := make([]byte, base64.StdEncoding.EncodedLen(len(binaryEncodedData)))
	base64.StdEncoding.Encode(buf, binaryEncodedData)

	return buf, nil
}

// UnmarshalText implements TextUnmarshaler interface decoding provided `data`
// from text form.
//
// The current BitSet's data will be overwritten by the decoded one
// if decoding operation has been completed successfully.
//
// Provided `data` MUST BE obtained by calling BitSet.MarshalText() method.
//
// Does nothing (and returns nil) if provided `data` is empty.
// Returns ErrBitSetInvalidDataToDecode if provided data is empty or invalid.
//
// WARNING!
// User MUST NOT use provided `data` after passing to this method. UB otherwise.
func (bs *BitSet) UnmarshalText(data []byte) error {

	switch {
	case len(data) == 0:
		return nil

	case bs == nil:
		return ErrBitSetInvalid
	}

	buf := make([]byte, base64.StdEncoding.DecodedLen(len(data)))
	n, err := base64.StdEncoding.Decode(buf, data)
	if err != nil {
		return err
	}

	buf = buf[:n]
	if len(buf)&_BITSET_BYTES_PER_CHUNK != 0 {
		return ErrBitSetInvalidDataToDecode
	}

	bs.bs = bsUnsafeFromBytesSlice(buf)
	return nil
}

// ---------------------------------------------------------------------------- //

// NewBitSet creates a new BitSet with desired initial capacity.
// If capacity is too small, it will be overwritten with default minimum capacity.
func NewBitSet(capacity uint) *BitSet {
	return new(BitSet).GrowUnsafeUpTo(capacity)
}
