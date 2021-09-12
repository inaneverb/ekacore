// Copyright © 2021. All rights reserved.
// Author: Ilya Stroy.
// Contacts: qioalice@gmail.com, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekamath

import (
	"fmt"
	"strings"
	"unsafe"
)

func (bs *BitSet) DebugOnesAsSlice(expectedValues uint) []uint {
	ones := make([]uint, 0, MaxU(expectedValues, _BITSET_MINIMUM_CAPACITY) +1)
	for v, e := bs.NextUp(0); e; v, e = bs.NextUp(v) {
		ones = append(ones, v)
	}
	return ones
}

func (bs *BitSet) DebugFullDump() {

	fmt.Printf("Bitset dump of [%x]\n", uintptr(unsafe.Pointer(bs)))
	defer fmt.Printf("End of bitset dump of [%x]\n", uintptr(unsafe.Pointer(bs)))

	if bs == nil {
		return
	}

	countOnes := bs.Count()

	fmt.Printf("\tIs valid:       %t\n", bs.IsValid())
	fmt.Printf("\tChunk bit size: %d\n", _BITSET_BITS_PER_CHUNK)
	fmt.Printf("\tChunk len:      %d\n", len(bs.bs))
	fmt.Printf("\tChunk capacity: %d\n", cap(bs.bs))
	fmt.Printf("\tCapacity:       %d\n", bs.Capacity())
	fmt.Printf("\tIs empty:       %t\n", bs.IsEmpty())
	fmt.Printf("\tCount ones:     %d\n", countOnes)
	fmt.Printf("\tAs chunks:      %v\n", bs.bs)

	var bits strings.Builder
	bits.Grow(int(bs.chunkSize() * (_BITSET_BITS_PER_CHUNK +1)))
	for _, chunk := range bs.bs {
		_, _ = fmt.Fprintf(&bits, "%b ", chunk)
	}

	fmt.Printf("\tAs bits:        [%s]\n", strings.TrimSpace(bits.String()))
	fmt.Printf("\tAs elems:       %v\n", bs.DebugOnesAsSlice(countOnes))
}