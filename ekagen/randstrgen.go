// Copyright © 2020. All rights reserved.
// Author: Ilya Stroy.
// Contacts: qioalice@gmail.com, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekagen

const (
	charSetLetters = `abcdefghijklmnopqrstuvwxyz`
	charSetDigits  = `1234567890`
	charSetAll     = charSetLetters + charSetDigits
)

//
func genWithLenFrom(charSet string, n int) string {

	if n <= 0 {
		return ""
	}

	res := make([]byte, n)

	for i := 0; i < n; i++ {
		res[i] = charSet[r.Intn(len(charSet))]
	}

	return string(res)
}

//
func WithLen(n int) string {
	return genWithLenFrom(charSetAll, n)
}

//
func WithLenOnlyLetters(n int) string {
	return genWithLenFrom(charSetLetters, n)
}

//
func WithLenOnlyNumbers(n int) string {
	return genWithLenFrom(charSetDigits, n)
}
