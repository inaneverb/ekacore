// Copyright © 2020. All rights reserved.
// Author: Ilya Yuryevich.
// Contacts: iyuryevich@pm.me, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekastr

// RequiredForI32 reports how many bytes are required
// to represent provided int32 as a string. It also covers negative values.
func RequiredForI32(i int32) int {
	return RequiredForI64(int64(i))
}

// RequiredForI64 reports how many bytes are required
// to represent provided int64 as a string. It also covers negative values.
func RequiredForI64(i int64) int {

	var requiredBytes int

	if i < 0 {
		i = -i
		requiredBytes += 1
	}

	switch {

	case i < 10:
		requiredBytes += 1

	case i < 100:
		requiredBytes += 2

	case i < 1_000:
		requiredBytes += 3

	case i < 10_000:
		requiredBytes += 4

	case i < 100_000:
		requiredBytes += 5

	case i < 1_000_000:
		requiredBytes += 6

	case i < 10_000_000:
		requiredBytes += 7

	case i < 100_000_000:
		requiredBytes += 8

	case i < 1_000_000_000:
		requiredBytes += 9

	case i < 10_000_000_000:
		requiredBytes += 10

	case i < 100_000_000_000:
		requiredBytes += 11

	case i < 1_000_000_000_000:
		requiredBytes += 12

	case i < 10_000_000_000_000:
		requiredBytes += 13

	case i < 100_000_000_000_000:
		requiredBytes += 14

	case i < 1_000_000_000_000_000:
		requiredBytes += 15

	case i < 10_000_000_000_000_000:
		requiredBytes += 16

	case i < 100_000_000_000_000_000:
		requiredBytes += 17

	case i < 1_000_000_000_000_000_000:
		requiredBytes += 18

	default:
		// i is int64 and max(int64) == 9_223_372_036_854_775_807
		requiredBytes += 19
	}

	return requiredBytes
}

// BItoa32 is the same as just itoa for int32, but instead returning string,
// it writes int32 value to the provided []byte, and reports, how many bytes
// were written.
//
// WARNING!
// If there's not enough space in 'to' to write provided int32,
// -1 is returned and provided []byte remains intact.
func BItoa32(to []byte, i int32) (n int) {
	return BItoa64(to, int64(i))
}

// BItoa64 is the same as BItoa32 but for int64.
func BItoa64(to []byte, i int64) (n int) {

	var idx = RequiredForI64(i)
	if len(to) < idx {
		return -1
	}

	if i < 0 {
		i = -i
		to[0] = '-'
	}

	n = idx
	for idx--; i != 0; idx-- {
		to[idx] = byte(i%10) + '0'
		i /= 10
	}

	return n
}
