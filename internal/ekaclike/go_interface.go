// Copyright © 2020. All rights reserved.
// Author: Ilya Stroy.
// Contacts: iyuryevich@pm.me, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekaclike

import (
	"unsafe"
)

type (
	/*
		Interface represents what "interface{}" is
		in internal Golang parts.
	*/
	Interface struct {
		Type uintptr        // pointer to the type definition struct
		Word unsafe.Pointer // pointer to the value
	}
)

/*
Pack does the reverse thing that UnpackInterface does.
It returns Golang interface{} object that current Interface describes.
*/
func (i Interface) Pack() (i2 interface{}) {

	if i.Type == 0 {
		return nil
	}

	i2I := (*Interface)(unsafe.Pointer(&i2))

	i2I.Type = i.Type
	i2I.Word = i.Word

	return
}

/*
UnpackInterface exposes Golang interface{} internal parts
and returns it.
If passed argument is absolutely nil, returns an empty Interface object.
*/
func UnpackInterface(i interface{}) Interface {
	if i == nil {
		return Interface{}
	}
	return *(*Interface)(unsafe.Pointer(&i))
}
