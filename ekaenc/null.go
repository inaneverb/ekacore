// Copyright © 2020-2023. All rights reserved.
// Author: Ilya Stroy.
// Contacts: iyuryevich@pm.me, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekaenc

import (
	"github.com/inaneverb/ekacore/ekastr/v4"
)

func NullAsStringLowerCase() string { return "null" }

func NullAsStringUpperCase() string { return "NULL" }

func NullAsBytesLowerCase() []byte {
	return ekastr.ToBytes(NullAsStringLowerCase())
}

func NullAsBytesUpperCase() []byte {
	return ekastr.ToBytes(NullAsStringUpperCase())
}

func IsNullAsString(s string) bool {
	return s == NullAsStringLowerCase() || s == NullAsStringUpperCase()
}

func IsNullAsBytes(b []byte) bool {
	return IsNullAsString(ekastr.FromBytes(b))
}
