// Copyright © 2019. All rights reserved.
// Author: Ilya Yuryevich.
// Contacts: qioalice@gmail.com, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package dangerous

import "unsafe"

// Interface represents what "interface{}" means in internal Golang parts.
type Interface struct {
	Type uintptr        // pointer to the type definition struct
	Word unsafe.Pointer // pointer to the value
}

//
func TypedInterface(i interface{}) Interface {
	return *(*Interface)(unsafe.Pointer(&i))
}
