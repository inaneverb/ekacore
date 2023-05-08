// Copyright © 2020-2023. All rights reserved.
// Author: Ilya Stroy.
// Contacts: iyuryevich@pm.me, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekaunsafe

import (
	"encoding/binary"
	"unsafe"
)

var (
	goRuntimeAddrEnc binary.ByteOrder // uses in GoRuntimeAddressEncoder()
)

////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////

func init() {
	var x uint64 = 0x6B617479616C7675
	var b = make([]byte, unsafe.Sizeof(x))

	for i, n := 0, len(b); i < n; i++ {
		b[i] = *(*byte)(unsafe.Add(unsafe.Pointer(&x), i))
	}

	switch {
	case x == binary.LittleEndian.Uint64(b):
		goRuntimeAddrEnc = binary.LittleEndian

	case x == binary.BigEndian.Uint64(b):
		goRuntimeAddrEnc = binary.BigEndian

	default:
		panic("Init: Unknown byte order of storing Go addresses in runtime")
	}
}

////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////

// GoRuntimeAddressEncoder returns binary.ByteOrder, that is binary.BigEndian
// or binary.LittleEndian, representing how exactly Go saves addresses
// (and other integers) in RAM for current OS.
func GoRuntimeAddressEncoder() binary.ByteOrder {
	return goRuntimeAddrEnc
}
