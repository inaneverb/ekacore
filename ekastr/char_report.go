// Copyright © 2020. All rights reserved.
// Author: Ilya Stroy.
// Contacts: iyuryevich@pm.me, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekastr

//goland:noinspection GoSnakeCaseUsage
const (

	/*
		There are constants for CharReport function.
		CharReport may return any of these variables.
	*/

	CR_NUMBER            = int8(1)
	CR_UPPER_CASE_LETTER = int8(2)
	CR_LOWER_CASE_LETTER = int8(3)
	CR_WHITESPACE        = int8(4)
)

/*
CharReport reports whether passed byte is. A letter? A number? A whitespace?
Use predefined constants with "CR_" prefix to compare with return value,
or specified methods.
Returns -1 for those bytes that are not covers by "CR_" constants.
*/
func CharReport(b byte) int8 {
	switch {

	case b >= '0' && b <= '9':
		return CR_NUMBER

	case b >= 'A' && b <= 'Z':
		return CR_UPPER_CASE_LETTER

	case b >= 'a' && b <= 'z':
		return CR_LOWER_CASE_LETTER

	case b <= 0x32:
		return CR_WHITESPACE

	default:
		return -1
	}
}

/*
CharIsLetter reports whether b is the range [A..Z] + [a..z] of ASCII table.
*/
func CharIsLetter(b byte) bool {
	return CharIsUpperCaseLetter(b) || CharIsLowerCaseLetter(b)
}

/*
CharIsUpperCaseLetter reports whether b is in the range [A..Z] of ASCII table.
*/
func CharIsUpperCaseLetter(b byte) bool {
	return CharReport(b) == CR_UPPER_CASE_LETTER
}

/*
CharIsLowerCaseLetter reports whether b is in the range [a..z] of ASCII table.
*/
func CharIsLowerCaseLetter(b byte) bool {
	return CharReport(b) == CR_LOWER_CASE_LETTER
}

/*
CharIsNumber reports whether b is in the range [0..9] of ASCII table.
*/
func CharIsNumber(b byte) bool {
	return CharReport(b) == CR_NUMBER
}

/*
CharIsWhitespace reports whether b is in the range [0x00..0x20] of ASCII table.
*/
func CharIsWhitespace(b byte) bool {
	return CharReport(b) == CR_WHITESPACE
}
