// Copyright © 2020. All rights reserved.
// Author: Ilya Stroy.
// Contacts: iyuryevich@pm.me, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekaenc

import (
	"github.com/qioalice/ekago/ekastr/v4"
)

func NullAsStringLowerCase() string { return "null" }

func NullAsStringUpperCase() string { return "NULL" }

func NullAsBytesLowerCase() []byte { return ekastr.S2B(NullAsStringLowerCase()) }

func NullAsBytesUpperCase() []byte { return ekastr.S2B(NullAsStringUpperCase()) }

func IsNullAsString(s string) bool {
	return s == NullAsStringLowerCase() || s == NullAsStringUpperCase()
}

func IsNullAsBytes(b []byte) bool { return IsNullAsString(ekastr.B2S(b)) }
