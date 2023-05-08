// Copyright © 2020. All rights reserved.
// Author: Ilya Stroy.
// Contacts: iyuryevich@pm.me, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekaunsafe

import (
	"unsafe"
)

// Interface represents what "any" is in internal Golang parts.
type Interface struct {
	Type uintptr        // pointer to the type definition struct
	Word unsafe.Pointer // pointer to the value
}

// Pack does the reverse thing that UnpackInterface does.
// It returns Golang any object that current Interface describes.
func (i Interface) Pack() any {

	if i.Type == 0 {
		return nil
	}

	var ret any
	var iRet = (*Interface)(unsafe.Pointer(&ret))

	iRet.Type = i.Type
	iRet.Word = i.Word

	return ret
}

// Tuple returns Interface arguments as a tuple. It's just for convenient usage
// while chaining or unpacking interface and passing its parts to anywhere else.
func (i Interface) Tuple() (uintptr, unsafe.Pointer) {
	return i.Type, i.Word
}

// PackInterface is the same as Interface.Pack method,
// but it is just a small convenient alias.
func PackInterface(rtype uintptr, word unsafe.Pointer) any {
	return Interface{rtype, word}.Pack()
}

// UnpackInterface exposes Golang any internal parts and returns it.
// If passed argument is absolutely nil, returns an empty Interface object.
func UnpackInterface(i any) Interface {
	if i == nil {
		return Interface{}
	}
	return *(*Interface)(unsafe.Pointer(&i))
}