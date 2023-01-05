// Copyright © 2020. All rights reserved.
// Author: Ilya Stroy.
// Contacts: iyuryevich@pm.me, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekastr

import (
	"fmt"
	"reflect"
	"strconv"
	"unsafe"

	"github.com/qioalice/ekago/ekaunsafe/v4"

	"github.com/modern-go/reflect2"
)

//goland:noinspection GoSnakeCaseUsage
const (
	ToStringHandleArrays        uint8 = 0x01
	ToStringHandleStructs       uint8 = 0x02
	ToStringHandleMaps          uint8 = 0x04
	ToStringDereferencePointers uint8 = 0x08
)

// ToString is a special function that allows you to get a string representation
// of any type's passed argument. It returns an empty string for nil interface.
func ToString(i any) string {
	var iFace = ekaunsafe.UnpackInterface(i)
	return ToStringUnsafe(iFace.Type, iFace.Word, 0xFF)
}

// ToStringUnsafe is the same as ToString but allows you to reject handling
// arrays, structs, maps, pointers.
// It also awaits to get unpacked Golang interface as both of rtype's
// and data's pointers. Zero safe - returns an empty string then.
func ToStringUnsafe(rtype uintptr, word unsafe.Pointer, mask uint8) string {

	if word == nil && rtype != 0 {
		switch {
		case rtype == ekaunsafe.RTypeBool() || ekaunsafe.RTypeIsStringLike(rtype):
			return ""
		case ekaunsafe.RTypeIsRealAny(rtype) || ekaunsafe.RTypeIsComplexAny(rtype):
			return "0"
		}
		return toStringUnsafeComplex(rtype, word, mask)
	}

	switch {

	case rtype == 0:
		return ""

	case rtype == ekaunsafe.RTypeString():
		return *(*string)(word)

	case rtype == ekaunsafe.RTypeBytes():
		return B2S(*(*[]byte)(word))

	case rtype == ekaunsafe.RTypeBool():
		if *(*bool)(word) {
			return "true"
		} else {
			return "false"
		}

	case rtype == ekaunsafe.RTypeInt():
		return strconv.FormatInt(int64(*(*int)(word)), 10)
	case rtype == ekaunsafe.RTypeInt8():
		return strconv.FormatInt(int64(*(*int8)(word)), 10)
	case rtype == ekaunsafe.RTypeInt16():
		return strconv.FormatInt(int64(*(*int16)(word)), 10)
	case rtype == ekaunsafe.RTypeInt32():
		return strconv.FormatInt(int64(*(*int32)(word)), 10)
	case rtype == ekaunsafe.RTypeInt64():
		return strconv.FormatInt(*(*int64)(word), 10)

	case rtype == ekaunsafe.RTypeUint():
		return strconv.FormatUint(uint64(*(*uint)(word)), 10)
	case rtype == ekaunsafe.RTypeUint8():
		return strconv.FormatUint(uint64(*(*uint8)(word)), 10)
	case rtype == ekaunsafe.RTypeUint16():
		return strconv.FormatUint(uint64(*(*uint16)(word)), 10)
	case rtype == ekaunsafe.RTypeUint32():
		return strconv.FormatUint(uint64(*(*uint32)(word)), 10)
	case rtype == ekaunsafe.RTypeUint64():
		return strconv.FormatUint(*(*uint64)(word), 10)

	case rtype == ekaunsafe.RTypeFloat32():
		return strconv.FormatFloat(float64(*(*float32)(word)), 'f', 2, 32)
	case rtype == ekaunsafe.RTypeFloat64():
		return strconv.FormatFloat(*(*float64)(word), 'f', 2, 64)

	case rtype == ekaunsafe.RTypeComplex64():
		return strconv.FormatComplex(complex128(*(*complex64)(word)), 'f', 2, 64)
	case rtype == ekaunsafe.RTypeComplex128():
		return strconv.FormatComplex(*(*complex128)(word), 'f', 2, 128)
	}

	return toStringUnsafeComplex(rtype, word, mask)
}

// toStringUnsafeComplex is a part of ToStringUnsafe, but handles complex cases.
func toStringUnsafeComplex(rtype uintptr, word unsafe.Pointer, mask uint8) string {

	// All standard types are over.
	// Next types will be complex.

	var iFace = ekaunsafe.PackInterface(rtype, word)
	var kind = reflect2.TypeOf(iFace).Kind()

	switch {
	case kind == reflect.Array && mask&ToStringHandleArrays == 0:
		return ""
	case kind == reflect.Slice && mask&ToStringHandleArrays == 0:
		return ""
	case kind == reflect.Struct && mask&ToStringHandleStructs == 0:
		return ""
	case kind == reflect.Map && mask&ToStringHandleMaps == 0:
		return ""
	case kind == reflect.Ptr && mask&ToStringDereferencePointers == 0:
		return ""
	}

	if word == nil {
		switch kind {
		case reflect.Array, reflect.Slice:
			return "[]"
		case reflect.Struct, reflect.Map:
			return "{}"
		}
	}

	// TODO: Support fmt.Stringer interface.
	// TODO: Separate handle all complex cases more optimised way.

	return fmt.Sprintf("%+v", iFace)
}
