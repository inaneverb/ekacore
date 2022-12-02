// Copyright © 2020-2022. All rights reserved.
// Author: Ilya Stroy.
// Contacts: iyuryevich@pm.me, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekamath

import (
	"github.com/qioalice/ekago/v3/ekaext"
)

func Min[T ekaext.Ordered](a, b T) T {
	return ekaext.If(a < b, a, b)
}

func Max[T ekaext.Ordered](a, b T) T {
	return ekaext.If(a > b, a, b)
}

func Clamp[T ekaext.Ordered](v, a, b T) T {
	a, b = Min(a, b), Max(a, b)
	return Min(Max(v, a), b)
}
