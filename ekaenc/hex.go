// Copyright © 2020-2023. All rights reserved.
// Author: Ilya Stroy.
// Contacts: iyuryevich@pm.me, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekaenc

import (
	"github.com/inaneverb/ekacore/ekaext/v4"
	"github.com/inaneverb/ekacore/ekastr/v4"
)

// KindHex is the same as ekastr.Kind(), but returns ekastr.KindOther,
// if 'b' doesn't belong [0..9] U [a..f] U [A..F].
func KindHex(b byte) int8 {
	var q = (b > 'f' && b <= 'z') || (b > 'F' && b < 'Z')
	return ekaext.If(q, ekastr.KindOther, ekastr.Kind(b))
}

// DecodeHex decodes hex byte and returns. If 'b' is valid hex,
// then returned value belongs [0x00..0x0F]. Otherwise 0xFF (255) is returned.
func DecodeHex(b byte) byte {
	switch KindHex(b) {
	case ekastr.KindNumber:
		return b - '0'
	case ekastr.KindLowerCaseLetter:
		return b - 'a' + 10
	case ekastr.KindUpperCaseLetter:
		return b - 'A' + 10
	default:
		return 0xFF
	}
}

// DecodeHexFull decodes two hex items into a single byte. 'l' stands for left,
// 'r' - right. Returns 0 if any of them is overflowing hex table.
// It allows to use both of lower and upper case for letters.
//
// Examples:
//   - L: '9', R: '1'. Out: 1001 0001
//   - L: 'F', R: 'a': Out: 1111 1010
//   - L: 'x', R: '2': Out: 0000 0000 // 0, because 'x' is not hex.
//
// You can easily check, whether decoding with an error (of incorrect hex)
// or not using binary OR:
//
//	var a, b = 'r', 1
//	if res := DecodeHexFull(a, b); res == 0 && a | b != 0 {
//	    panic("error; any of 'a' or 'b' is not hex")
//	}
func DecodeHexFull(l, r byte) (out byte) {
	var l1, r1 = DecodeHex(l), DecodeHex(r)
	if l1|r1 <= 0x0F {
		out = (l1 << 4) | r1
	}
	return out
}

// DecodeHex decodes hex byte and returns correspondent octet.
// If 'c' is valid hex, then returned value belongs [0x00..0x0F].
// Otherwise 0xFF (255) is returned.

// EncodeHex encodes lowest 4 bits of byte as hex value and returns it.
// If 'b' <= 0x0F, then valid hex value is returned; 0xFF (255) otherwise.
// Lowercase is used to represent [0x0A..0x0F] values.
func EncodeHex(b byte) byte {
	const hexTable = "0123456789abcdef"
	if b <= 0x0F {
		return hexTable[b]
	} else {
		return 0xFF
	}
}

// EncodeHexFull encodes full byte into two hex values.
// For example, 0xAB produces 1st returned arg is 'a', and 2nd is 'b'.
// Lowercase is used for letters.
//
// Examples:
//   - B: 92.  Out: '5', 'c' (0x5C).
//   - B: 234. Out: 'e', 'a' (0xEA).
//   - B: 0.   Out: '0', '0' (0x00).
//
// There's no way to get an exception during decoding.
// So, all returned values are always valid.
func EncodeHexFull(b byte) (byte, byte) {
	return EncodeHex(b >> 4), EncodeHex(b & 0x0F)
}
