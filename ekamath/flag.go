// Copyright © 2020-2021. All rights reserved.
// Author: Ilya Stroy.
// Contacts: iyuryevich@pm.me, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekamath

import (
	"github.com/inaneverb/ekacore/ekaext/v4"
)

// FlagAny reports whether at least one flag from 'f2' is set in 'f'.
// Returns false if there are no flags in 'f' from 'f2'.
func FlagAny[T ekaext.Unsigned](f, f2 T) bool {
	return f&f2 != 0
}

// FlagAll reports whether ALL flags from 'f2' are set in 'f'.
// Returns false if at least one flag from 'f2' is NOT set in 'f'.
func FlagAll[T ekaext.Unsigned](f, f2 T) bool {
	return f&f2 == f2
}

// FlagClear downs all flags from 'f2' in 'f'. Returns changed 'f'.
func FlagClear[T ekaext.Unsigned](f, f2 T) T {
	return f &^ f2
}

// FlagApply changes all bits in 'f', that are set in 'f2' to 0 or 1,
// depending on provided 'up' value: 0 for false, 1 for true.
func FlagApply[T ekaext.Unsigned](f, f2 T, up bool) T {
	return f & ekaext.If(up, f2, ^f2)
}
