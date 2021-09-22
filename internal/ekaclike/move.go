// Copyright © 2021. All rights reserved.
// Author: Ilya Stroy.
// Contacts: iyuryevich@pm.me, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekaclike

import (
	_ "reflect"
	"unsafe"
)

//go:linkname Typedmemmove reflect.typedmemmove
func Typedmemmove(rtype unsafe.Pointer, dst, src unsafe.Pointer)
