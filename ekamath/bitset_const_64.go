// Copyright © 2021. All rights reserved.
// Author: Ilya Stroy.
// Contacts: qioalice@gmail.com, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekamath

type (
	// BitSet64 is a bitset with max capacity == 64.
	// It's ready-to-use after just instantiating an object.
	BitSet64 uint8
)

// IsValid reports whether current BitSet64 is valid.
func (bs *BitSet64) IsValid() bool {
	return bs != nil
}

// IsValidWithIndex extends just IsValid() call with also checking
// whether provided index in the allowed range.
func (bs *BitSet64) IsValidWithIndex(idx uint) bool {
	return bs.IsValid() && idx > 0 && idx <= (_BITSET_64_CHUNK_MASK+1)
}

// IsEmpty reports whether current BitSet64 is empty bitset or not.
// Empty bitset is a bitset with all downed (zeroed) bits.
// Returns true if BitSet64 is invalid.
func (bs *BitSet64) IsEmpty() bool {
	return !bs.IsValid() || *bs == 0
}

// ClearUnchecked downs (zeroes) ALL bits in the current BitSet64.
// Panics if current BitSet64 is invalid.
func (bs *BitSet64) ClearUnchecked() *BitSet64 {
	*bs = 0
	return bs
}

// Up sets bit to 1 with requested index checking bounds.
// Does nothing if BitSetGeneric is invalid or index is out of bounds.
func (bs *BitSet64) Up(idx uint) *BitSet64 {
	if bs.IsValidWithIndex(idx) {
		bs.UpUnchecked(idx)
	}
	return bs
}

// UpUnchecked sets bit to 1 with requested index without any check.
// Panics if BitSetGeneric is invalid or if an index is out of bounds.
func (bs *BitSet64) UpUnchecked(idx uint) *BitSet64 {
	*bs |= 1 << (idx - 1)
	return bs
}

// Down sets bit to 0 with requested index checking bounds.
// Does nothing if BitSetGeneric is invalid or index is out of bounds.
func (bs *BitSet64) Down(idx uint) *BitSet64 {
	if bs.IsValidWithIndex(idx) {
		bs.DownUnchecked(idx)
	}
	return bs
}

// DownUnchecked sets bit to 0 with requested index without any check.
// Panics if BitSetGeneric is invalid or if an index is out of bounds.
func (bs *BitSet64) DownUnchecked(idx uint) *BitSet64 {
	*bs &^= 1 << (idx - 1)
	return bs
}

// Set calls Up() or Down() with provided index depends on `b`.
func (bs *BitSet64) Set(idx uint, b bool) *BitSet64 {
	if bs.IsValidWithIndex(idx) {
		bs.SetUnchecked(idx, b)
	}
	return bs
}

// SetUnchecked calls UpUnchecked() or DownUnchecked() with provided index
// depends on `b`.
func (bs *BitSet64) SetUnchecked(idx uint, b bool) *BitSet64 {
	if b {
		return bs.UpUnchecked(idx)
	} else {
		return bs.DownUnchecked(idx)
	}
}

// IsSet reports whether a bit with requested index is set or not.
// Returns false either if bit isn't set, BitSet64 is invalid or index is out of bounds.
func (bs *BitSet64) IsSet(idx uint) bool {
	return bs.IsValidWithIndex(idx) && bs.IsSetUnchecked(idx)

}

// IsSetUnchecked reports whether a bit with requested index is set or not.
// Panics if BitSet64 is invalid or if an index is out of bounds.
func (bs *BitSet64) IsSetUnchecked(idx uint) bool {
	return *bs&(1<<(idx-1)) != 0
}
