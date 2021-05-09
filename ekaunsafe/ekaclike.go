// Copyright © 2020. All rights reserved.
// Author: Ilya Stroy.
// Contacts: qioalice@gmail.com, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekaunsafe

import (
	"unsafe"

	"github.com/qioalice/ekago/v3/internal/ekaclike"
)

// To see docs and comments,
// navigate to the origin package.

type (
	Interface = ekaclike.Interface
)

func UnpackInterface(i interface{}) Interface {
	return ekaclike.UnpackInterface(i)
}

func TakeRealAddr(i interface{}) unsafe.Pointer {
	return ekaclike.TakeRealAddr(i)
}

func TakeCallableAddr(i interface{}) unsafe.Pointer {
	return ekaclike.TakeCallableAddr(i)
}

func Addr2Callable(realPtr unsafe.Pointer) (callablePtr unsafe.Pointer) {
	return ekaclike.Addr2Callable(realPtr)
}

func Addr2Real(callablePtr unsafe.Pointer) (realPtr unsafe.Pointer) {
	return ekaclike.Addr2Real(realPtr)
}
