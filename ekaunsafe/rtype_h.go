// Copyright © 2020-2023. All rights reserved.
// Author: Ilya Stroy.
// Contacts: iyuryevich@pm.me, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekaunsafe

import (
	"reflect"
	"unsafe"
)

var (
	sizeOfRType       uintptr // how many bytes reflect.rtype type requires
	rtypeRTypePtr     uintptr // addr of *reflect.rtype type in the Go types
	rtypeRTypeItabPtr uintptr // addr of *rtype itab struct in reflect.Type
)

////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////

func init() {
	type Ptr = unsafe.Pointer
	var rt = reflect.TypeOf(reflect.TypeOf(uint64(0x6B617479616C7675)))
	sizeOfRType = rt.Elem().Size()
	rtypeRTypePtr = UnpackInterface(rt).Type
	rtypeRTypeItabPtr = uintptr(*(*Ptr)(Ptr(&rt)))
	_ = rtypeRTypePtr
}

////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////

// RTypeOfDeref returns a unique type's addr of T,
// assuming that given 'x' is a *T (points to T).
// WARNING! UB, MAY PANIC IF 'x' IS EITHER NIL OR NOT A POINTER!
func RTypeOfDeref(x any) uintptr {
	return RTypeOfDerefRType(UnpackInterface(x).Type)
}

// RTypeOfDerefRType returns a unique type's addr of T,
// assuming that given 'rtype' is a *T rtype (points to T).
// WARNING! UB, MAY PANIC IF 'rtype' IS EITHER 0 OR NOT A POINTER TYPE!
func RTypeOfDerefRType(rtype uintptr) uintptr {

	// All actions here are presented to step-by-step for better understanding.

	// Golang's rtype struct:
	// | rtype base fields | *rtype to base type if this is the pointer type
	//   ^                   ^
	//   |                   This is where we want to have pointer to
	//   This is where we have rtype pointer points to from unpacked interface

	var _1 = rtype + sizeOfRType // Shift pointer to (*rtype) of base type

	// The statement below is valid. Because Go tells us, that convert
	// from uintptr to unsafe.Pointer is invalid, referencing the cases,
	// when uintptr may point to data, that may be GC'ed already.
	// But that's not our case. uintptr points to internal Golang runtime
	// type table, so it lives forever.

	//goland:noinspection GoVetUnsafePointer
	var _2 = unsafe.Pointer(_1)    // For further casting
	var _3 = (*unsafe.Pointer)(_2) // That's a pointer to pointer, right?

	// Deref one layer.
	// So, now it's a pointer to rtype of base type.

	return uintptr(*_3)
}

// RTypeOfReflectType returns a unique type's addr of given reflect.Type.
// WARNING! UB, MAY PANIC IF 't' IS ZERO.
func RTypeOfReflectType(t reflect.Type) uintptr {
	return uintptr(UnpackInterface(t).Word)
}

// ReflectTypeOfRType returns reflect.Type, constructing it from given 'rtype'.
// Provided addr must be obtained from rtype getters, or unpacking interface.
// WARNING! UB IF EITHER RTYPE IS 0 OR ANY RANDOM DATA (SET OR MODIFIED).
func ReflectTypeOfRType(rtype uintptr) reflect.Type {

	// This implementation is much faster, than just packing an interface
	// using PackInterface and then cast it manually to reflect.Type,
	// since during interface casting there's embedded Go's reflections.

	//return PackInterface(rtypeRTypePtr, unsafe.Add(nil, rtype)).(reflect.Type)

	var r reflect.Type
	var ri = (*Interface)(unsafe.Pointer(&r))

	ri.Type = rtypeRTypeItabPtr
	ri.Word = unsafe.Add(nil, rtype)

	return r
}
