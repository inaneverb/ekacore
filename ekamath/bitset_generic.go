// Copyright © 2021. All rights reserved.
// Author: Ilya Stroy.
// Contacts: qioalice@gmail.com, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekamath

type (
	// BitSetGeneric is a bitset with variate capacity.
	// It can be grown, depends on your cases.
	//
	// It's strongly recommend to instantiate BitSetGeneric using NewBitSet() constructor,
	// but just creating a BitSetGeneric is also possible and ready-to-use
	// (it will be with 0 capacity and will grow when you will try to set any bit).
	BitSetGeneric struct {
		bs []uint
	}
)

// IsValid reports whether current BitSetGeneric is valid.
func (bs *BitSetGeneric) IsValid() bool {
	return bs != nil
}

// IsEmpty reports whether current BitSetGeneric is empty bitset or not.
// Empty bitset is a bitset with all downed (zeroed) bits.
// Returns true if BitSetGeneric is invalid.
func (bs *BitSetGeneric) IsEmpty() bool {
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
// inside current BitSetGeneric.
// Returns 0 if current BitSetGeneric is invalid.
func (bs *BitSetGeneric) Capacity() uint {
	if !bs.IsValid() {
		return 0
	}
	return uint(len(bs.bs)) * 8
}

// ClearUnchecked downs (zeroes) ALL bits in the current BitSetGeneric.
// Panics if current BitSetGeneric is invalid.
func (bs *BitSetGeneric) ClearUnchecked() *BitSetGeneric {
	for i, n := 0, len(bs.bs); i < n; i++ {
		bs.bs[i] = 0
	}
	return bs
}

// GrowUncheckedUpTo grows current BitSetGeneric to be able operate with bits
// up to requested index.
// Panics if BitSetGeneric is invalid.
func (bs *BitSetGeneric) GrowUncheckedUpTo(idx uint) *BitSetGeneric {

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
// growing bitset if it's too small. Does nothing if BitSetGeneric is invalid.
func (bs *BitSetGeneric) Up(idx uint) *BitSetGeneric {
	if bs.IsValid() && idx > 0 {
		bs.GrowUncheckedUpTo(idx).UpUnchecked(idx)
	}
	return bs
}

// UpUnchecked sets bit to 1 with requested index without any check.
// Panics if BitSetGeneric is invalid or if an index is out of bounds.
func (bs *BitSetGeneric) UpUnchecked(idx uint) *BitSetGeneric {
	chunk, offset := bsFromIdx(idx)
	bs.bs[chunk] |= 1 << offset
	return bs
}

// Down sets bit to 0 with requested index checking bounds,
// growing bitset if it's too small. Does nothing if BitSetGeneric is invalid.
func (bs *BitSetGeneric) Down(idx uint) *BitSetGeneric {
	if bs.IsValid() && idx > 0 {
		bs.GrowUncheckedUpTo(idx).DownUnchecked(idx)
	}
	return bs
}

// DownUnchecked sets bit to 0 with requested index without any check.
// Panics if BitSetGeneric is invalid or if an index is out of bounds.
func (bs *BitSetGeneric) DownUnchecked(idx uint) *BitSetGeneric {
	chunk, offset := bsFromIdx(idx)
	bs.bs[chunk] &^= 1 << offset
	return bs
}

// Set calls Up() or Down() with provided index depends on `b`.
func (bs *BitSetGeneric) Set(idx uint, b bool) *BitSetGeneric {
	if bs.IsValid() && idx > 0 {
		bs.GrowUncheckedUpTo(idx).SetUnchecked(idx, b)
	}
	return bs
}

// SetUnchecked calls UpUnchecked() or DownUnchecked() with provided index
// depends on `b`.
func (bs *BitSetGeneric) SetUnchecked(idx uint, b bool) *BitSetGeneric {
	if b {
		return bs.UpUnchecked(idx)
	} else {
		return bs.DownUnchecked(idx)
	}
}

// IsSet reports whether a bit with requested index is set or not.
// Returns false either if bit isn't set, BitSetGeneric is invalid or index is out of bound.
func (bs *BitSetGeneric) IsSet(idx uint) bool {
	return bs.IsValid() && bs.isValidIdx(idx) && bs.IsSetUnchecked(idx)
}

// IsSetUnchecked reports whether a bit with requested index is set or not.
// Panics if BitSetGeneric is invalid or if an index is out of bound.
func (bs *BitSetGeneric) IsSetUnchecked(idx uint) bool {
	chunk, offset := bsFromIdx(idx)
	return bs.bs[chunk] & (1 << offset) != 0
}

// NewBitSet creates a new BitSetGeneric with desired initial capacity.
// If capacity is too small, it will be overwritten with default minimum capacity.
func NewBitSet(capacity uint) *BitSetGeneric {
	return new(BitSetGeneric).GrowUncheckedUpTo(capacity)
}
