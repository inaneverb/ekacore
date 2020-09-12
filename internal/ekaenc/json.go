// Copyright © 2020. All rights reserved.
// Author: Ilya Stroy.
// Contacts: qioalice@gmail.com, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekaenc

//goland:noinspection GoSnakeCaseUsage
const (
	NULL_JSON = `null`
)

//goland:noinspection GoSnakeCaseUsage
var (
	NULL_JSON_BYTES_SLICE = []byte(NULL_JSON)
)

// IsNullJSON returns true if b == nil or b[:4] == "null" (case insensitive).
func IsNullJSON(b []byte) bool {

	if b == nil {
		return true
	}

	z := len(b) == 4
	z = z && b[0] == 'N' || b[0] == 'n'
	z = z && (b[1] == 'U' || b[1] == 'u')
	z = z && (b[2] == 'L' || b[2] == 'l')
	z = z && (b[3] == 'L' || b[3] == 'l')

	return z
}
