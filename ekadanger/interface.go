// Copyright © 2020. All rights reserved.
// Author: Ilya Stroy.
// Contacts: qioalice@gmail.com, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekadanger

import (
	"unsafe"
)

type (
	// Interface represents what "interface{}" is in internal Golang parts.
	Interface struct {
		Type uintptr        // pointer to the type definition struct
		Word unsafe.Pointer // pointer to the value
	}
)

// TypedInterface exposes Golang 'i' interface{} internal parts and returns it.
func TypedInterface(i interface{}) Interface {
	return *(*Interface)(unsafe.Pointer(&i))
}
