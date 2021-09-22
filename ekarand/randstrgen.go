// Copyright © 2020. All rights reserved.
// Author: Ilya Stroy.
// Contacts: iyuryevich@pm.me, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekarand

import (
	mrand "math/rand"

	"github.com/qioalice/ekago/v3/ekastr"
)

const (
	charSetLetters = `abcdefghijklmnopqrstuvwxyz`
	charSetDigits  = `1234567890`
	charSetAll     = charSetLetters + charSetDigits
)

// WithLen generates a random sequence with n length that contains both of
// english letters and arabic digits. Returns "" if n <= 0.
func WithLen(n int) string {
	return genWithLenFrom(charSetAll, n)
}

// WithLenOnlyLetters generates a random sequence with n length that contains
// only english letters. Returns "" if n <= 0.
func WithLenOnlyLetters(n int) string {
	return genWithLenFrom(charSetLetters, n)
}

// WithLenOnlyNumbers generates a random sequence with n length that contains
// only arabic digits. Returns "" if n <= 0.
func WithLenOnlyNumbers(n int) string {
	return genWithLenFrom(charSetDigits, n)
}

func genWithLenFrom(charSet string, n int) string {
	if n <= 0 {
		return ""
	}
	res := make([]byte, n)
	for i := 0; i < n; i++ {
		res[i] = charSet[mrand.Intn(len(charSet))]
	}
	return ekastr.B2S(res)
}
