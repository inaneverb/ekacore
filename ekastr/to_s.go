// Copyright © 2020. All rights reserved.
// Author: Ilya Stroy.
// Contacts: qioalice@gmail.com, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekastr

import (
	"fmt"
	"reflect"
	"strconv"
	"unsafe"

	"github.com/qioalice/ekago/v3/internal/ekaclike"

	"github.com/modern-go/reflect2"
)

//goland:noinspection GoSnakeCaseUsage
const (
	TO_S_HANDLE_ARRAYS   uint8 = 0x01
	TO_S_HANDLE_STRUCTS  uint8 = 0x02
	TO_S_HANDLE_MAPS     uint8 = 0x04
	TO_S_DEREFERENCE_PTR uint8 = 0x08
)

/*
ToString is a special type that allows you to get a string representation
of any type's passed argument. It returns an empty string for nil interface.
*/
func ToString(i interface{}) string {
	iface := ekaclike.UnpackInterface(i)
	return ToStringUnsafe(iface.Type, iface.Word, 0xFF)
}

/*
ToStringUnsafe is the same as ToString but allows you to reject handling
arrays, structs, maps, pointers. It also awaits to get unpacked Golang interface
as both of rtype's and data's pointers. Zero safe - returns an empty string then.
*/
func ToStringUnsafe(rtype uintptr, word unsafe.Pointer, mask uint8) string {

	if word == nil && rtype != 0 {
		switch rtype {
		case ekaclike.RTypeBool, ekaclike.RTypeString, ekaclike.RTypeBytes:
			return ""
		case ekaclike.RTypeInt,
		ekaclike.RTypeInt8, ekaclike.RTypeInt16,
		ekaclike.RTypeInt32, ekaclike.RTypeInt64,
		ekaclike.RTypeUint,
		ekaclike.RTypeUint8, ekaclike.RTypeUint16,
		ekaclike.RTypeUint32, ekaclike.RTypeUint64,
		ekaclike.RTypeFloat32, ekaclike.RTypeFloat64,
		ekaclike.RTypeComplex64, ekaclike.RTypeComplex128:
			return "0"
		}
		switch reflect2.TypeOf(ekaclike.Interface{Type: rtype}.Pack()).Kind() {
		case reflect.Array, reflect.Slice:
			return "[]"
		case reflect.Struct, reflect.Map:
			return "{}"
		default:
			return ""
		}
	}

	switch {

	case rtype == 0:
		return ""

	case rtype == ekaclike.RTypeString:
		return *(*string)(word)

	case rtype == ekaclike.RTypeBytes:
		return B2S(*(*[]byte)(word))

	case rtype == ekaclike.RTypeBool:
		if *(*bool)(word) {
			return "true"
		} else {
			return "false"
		}

	case rtype == ekaclike.RTypeInt:
		return strconv.FormatInt(int64(*(*int)(word)), 10)
	case rtype == ekaclike.RTypeInt8:
		return strconv.FormatInt(int64(*(*int8)(word)), 10)
	case rtype == ekaclike.RTypeInt16:
		return strconv.FormatInt(int64(*(*int16)(word)), 10)
	case rtype == ekaclike.RTypeInt32:
		return strconv.FormatInt(int64(*(*int32)(word)), 10)
	case rtype == ekaclike.RTypeInt64:
		return strconv.FormatInt(*(*int64)(word), 10)

	case rtype == ekaclike.RTypeUint:
		return strconv.FormatUint(uint64(*(*uint)(word)), 10)
	case rtype == ekaclike.RTypeUint8:
		return strconv.FormatUint(uint64(*(*uint8)(word)), 10)
	case rtype == ekaclike.RTypeUint16:
		return strconv.FormatUint(uint64(*(*uint16)(word)), 10)
	case rtype == ekaclike.RTypeUint32:
		return strconv.FormatUint(uint64(*(*uint32)(word)), 10)
	case rtype == ekaclike.RTypeUint64:
		return strconv.FormatUint(*(*uint64)(word), 10)

	case rtype == ekaclike.RTypeFloat32:
		return strconv.FormatFloat(float64(*(*float32)(word)), 'f', 2, 32)
	case rtype == ekaclike.RTypeFloat64:
		return strconv.FormatFloat(*(*float64)(word), 'f', 2, 64)

	case rtype == ekaclike.RTypeComplex64:
		return strconv.FormatComplex(complex128(*(*complex64)(word)), 'f', 2, 64)
	case rtype == ekaclike.RTypeComplex128:
		return strconv.FormatComplex(*(*complex128)(word), 'f', 2, 128)
	}

	// All standard types are over.
	// Next types will be complex.

	var (
		eface = ekaclike.Interface{Type: rtype, Word: word}.Pack()
		typ   = reflect2.TypeOf(eface)
		kind  = typ.Kind()
	)

	switch {
	case kind == reflect.Array && mask & TO_S_HANDLE_ARRAYS == 0:
		return ""
	case kind == reflect.Slice && mask & TO_S_HANDLE_ARRAYS == 0:
		return ""
	case kind == reflect.Struct && mask & TO_S_HANDLE_STRUCTS == 0:
		return ""
	case kind == reflect.Map && mask & TO_S_HANDLE_MAPS == 0:
		return ""
	case kind == reflect.Ptr && mask & TO_S_DEREFERENCE_PTR == 0:
		return ""
	}

	// Well, it's complex case and it must be handled.

	// TODO: Support fmt.Stringer interface.
	// TODO: Separate handle all complex cases more optimised way.
	return fmt.Sprintf("%+v", eface)
}
