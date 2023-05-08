// Copyright © 2020-2023. All rights reserved.
// Author: Ilya Stroy.
// Contacts: iyuryevich@pm.me, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekastr

// ToUpper is the same as bytes.ToUpper(), but it doesn't make a copy.
// Just works in-place and works only with ASCII.
func ToUpper(b []byte) {
	for i, n := 0, len(b); i < n; i++ {
		if b[i] >= 'a' && b[i] <= 'z' {
			b[i] -= 32
		}
	}
}

// ToLower is the same as bytes.ToLower(), but it doesn't make a copy.
// Just works in-place and works only with ASCII.
func ToLower(b []byte) {
	for i, n := 0, len(b); i < n; i++ {
		if b[i] >= 'A' && b[i] <= 'Z' {
			b[i] += 32
		}
	}
}
