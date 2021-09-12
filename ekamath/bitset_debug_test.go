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

func (bs *BitSet) DebugDump() {

	fmt.Printf("Bitset dump of [%x]\n", uintptr(unsafe.Pointer(bs)))
	defer fmt.Printf("End of bitset dump of [%x]\n", uintptr(unsafe.Pointer(bs)))

	if bs == nil {
		return
	}

	countOnes := bs.Count()

	fmt.Printf("\tIs valid:       %t\n", bs.IsValid())
	fmt.Printf("\tChunk bit size: %d\n", _BITSET_GENERIC_BITS_PER_CHUNK)
	fmt.Printf("\tChunk len:      %d\n", len(bs.bs))
	fmt.Printf("\tChunk capacity: %d\n", cap(bs.bs))
	fmt.Printf("\tCapacity:       %d\n", bs.Capacity())
	fmt.Printf("\tIs empty:       %t\n", bs.IsEmpty())
	fmt.Printf("\tCount ones:     %d\n", countOnes)
	fmt.Printf("\tAs chunks:      %v\n", bs.bs)

	var bits strings.Builder
	bits.Grow(len(bs.bs) * (8+1))
	for _, chunk := range bs.bs {
		_, _ = fmt.Fprintf(&bits, "%b ", chunk)
	}

	fmt.Printf("\tAs bits:        [%s]\n", strings.TrimSpace(bits.String()))

	ones := make([]uint, 0, countOnes +1)
	for v, e := bs.NextUp(0); e; v, e = bs.NextUp(v+1) {
		ones = append(ones, v)
	}

	fmt.Printf("\tAs elems:       %v\n", ones)
}