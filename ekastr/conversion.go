// Copyright © 2020-2023. All rights reserved.
// Author: Ilya Stroy.
// Contacts: iyuryevich@pm.me, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekastr

import (
	"github.com/qioalice/ekago/v4/ekaunsafe"
)

// FromBytes converts byte slice to a string without memory allocation.
func FromBytes(b []byte) string {
	return ekaunsafe.BytesToString(b)
}

// ToBytes converts string to a byte slice without memory allocation.
// String literals should be used only as read-only. Panic otherwise.
func ToBytes(s string) []byte {
	return ekaunsafe.StringToBytes(s)
}
