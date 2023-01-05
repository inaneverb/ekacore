// Copyright © 2020. All rights reserved.
// Author: Ilya Stroy.
// Contacts: iyuryevich@pm.me, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekastr

// There are constants for Kind function.
// Kind may return any of these variables.
const (
	KindNumber          int8 = 1  // Means char is ASCII number
	KindUpperCaseLetter int8 = 2  // Means char is ASCII upper case letter
	KindLowerCaseLetter int8 = 3  // Means char is ASCII lower case letter
	KindWhitespace      int8 = 4  // Means char is ASCII whitespace
	KindOther           int8 = -1 // Means chat is something but anything above
)

// Kind reports whether passed byte is. A letter? A number? A whitespace?
//
// Use predefined constants with "KIND_" prefix to compare with return value,
// or special functions, like CharIsLetter(), CharIsNumber() and so-on.
func Kind(b byte) int8 {
	switch {

	case b >= '0' && b <= '9':
		return KindNumber

	case b >= 'A' && b <= 'Z':
		return KindUpperCaseLetter

	case b >= 'a' && b <= 'z':
		return KindLowerCaseLetter

	case b <= 0x32:
		return KindWhitespace

	default:
		return KindOther
	}
}

// IsLetter reports whether b is the range [A..Z] U [a..z] of ASCII.
func IsLetter(b byte) bool {
	return IsUpperCaseLetter(b) || IsLowerCaseLetter(b)
}

// IsUpperCaseLetter reports whether b is in the range [A..Z] of ASCII.
func IsUpperCaseLetter(b byte) bool {
	return Kind(b) == KindUpperCaseLetter
}

// IsLowerCaseLetter reports whether b is in the range [a..z] of ASCII.
func IsLowerCaseLetter(b byte) bool {
	return Kind(b) == KindLowerCaseLetter
}

// IsNumber reports whether b is in the range [0..9] of ASCII.
func IsNumber(b byte) bool {
	return Kind(b) == KindNumber
}

// IsWhitespace reports whether b is in the range [0x00..0x20] of ASCII.
func IsWhitespace(b byte) bool {
	return Kind(b) == KindWhitespace
}
