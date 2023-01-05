// Copyright © 2020-2022. All rights reserved.
// Author: Ilya Stroy.
// Contacts: iyuryevich@pm.me, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekaunsafe

import (
	"time"
	"unsafe"

	"github.com/modern-go/reflect2"
)

//goland:noinspection GoVarAndConstTypeMayBeOmitted,GoRedundantConversion,GoBoolExpressions
var (
	rTypeBool               = reflect2.RTypeOf(bool(0 == 0))
	rTypeByte               = reflect2.RTypeOf(byte(0))
	rTypeRune               = reflect2.RTypeOf(rune(0))
	rTypeInt                = reflect2.RTypeOf(int(0))
	rTypeInt8               = reflect2.RTypeOf(int8(0))
	rTypeInt16              = reflect2.RTypeOf(int16(0))
	rTypeInt32              = reflect2.RTypeOf(int32(0))
	rTypeInt64              = reflect2.RTypeOf(int64(0))
	rTypeUint               = reflect2.RTypeOf(uint(0))
	rTypeUint8              = reflect2.RTypeOf(uint8(0))
	rTypeUint16             = reflect2.RTypeOf(uint16(0))
	rTypeUint32             = reflect2.RTypeOf(uint32(0))
	rTypeUint64             = reflect2.RTypeOf(uint64(0))
	rTypeFloat32            = reflect2.RTypeOf(float32(0))
	rTypeFloat64            = reflect2.RTypeOf(float64(0))
	rTypeComplex64          = reflect2.RTypeOf(complex64(0))
	rTypeComplex128         = reflect2.RTypeOf(complex128(0))
	rTypeString             = reflect2.RTypeOf(string(""))
	rTypeStringArray        = reflect2.RTypeOf([]string(nil))
	rTypeBytes              = reflect2.RTypeOf([]byte(nil))
	rTypeBytesArray         = reflect2.RTypeOf([][]byte(nil))
	rTypeMapStringString    = reflect2.RTypeOf(map[string]string(nil))
	rTypeMapStringInterface = reflect2.RTypeOf(map[string]any(nil))
	rTypeUintptr            = reflect2.RTypeOf(uintptr(0))
	rTypeUnsafePointer      = reflect2.RTypeOf(unsafe.Pointer(nil))
	rTypeTimeTime           = reflect2.RTypeOf(time.Time{})
	rTypeTimeDuration       = reflect2.RTypeOf(time.Duration(0))
)

func RTypeBool() uintptr               { return rTypeBool }
func RTypeByte() uintptr               { return rTypeByte }
func RTypeRune() uintptr               { return rTypeRune }
func RTypeInt() uintptr                { return rTypeInt }
func RTypeInt8() uintptr               { return rTypeInt8 }
func RTypeInt16() uintptr              { return rTypeInt16 }
func RTypeInt32() uintptr              { return rTypeInt32 }
func RTypeInt64() uintptr              { return rTypeInt64 }
func RTypeUint() uintptr               { return rTypeUint }
func RTypeUint8() uintptr              { return rTypeUint8 }
func RTypeUint16() uintptr             { return rTypeUint16 }
func RTypeUint32() uintptr             { return rTypeUint32 }
func RTypeUint64() uintptr             { return rTypeUint64 }
func RTypeFloat32() uintptr            { return rTypeFloat32 }
func RTypeFloat64() uintptr            { return rTypeFloat64 }
func RTypeComplex64() uintptr          { return rTypeComplex64 }
func RTypeComplex128() uintptr         { return rTypeComplex128 }
func RTypeString() uintptr             { return rTypeString }
func RTypeStringArray() uintptr        { return rTypeStringArray }
func RTypeBytes() uintptr              { return rTypeBytes }
func RTypeBytesArray() uintptr         { return rTypeBytesArray }
func RTypeMapStringString() uintptr    { return rTypeMapStringString }
func RTypeMapStringInterface() uintptr { return rTypeMapStringInterface }
func RTypeUintptr() uintptr            { return rTypeUintptr }
func RTypeUnsafePointer() uintptr      { return rTypeUnsafePointer }
func RTypeTimeTime() uintptr           { return rTypeTimeTime }
func RTypeTimeDuration() uintptr       { return rTypeTimeDuration }

func RTypeIsAnyNumeric(rtype uintptr) bool {
	return RTypeIsIntAny(rtype) || RTypeIsUintAny(rtype)
}

func RTypeIsAnyReal(rtype uintptr) bool {
	return RTypeIsAnyNumeric(rtype) || RTypeIsFloatAny(rtype)
}

func RTypeIsIntAny(rtype uintptr) bool {
	return rtype == rTypeInt || RTypeIsIntFixed(rtype)
}

func RTypeIsIntFixed(rtype uintptr) bool {
	switch rtype {
	case rTypeInt8, rTypeInt16, rTypeInt32, rTypeInt64:
		return true
	default:
		return false
	}
}

func RTypeIsUintAny(rtype uintptr) bool {
	return rtype == rTypeUint || RTypeIsUintFixed(rtype)
}

func RTypeIsUintFixed(rtype uintptr) bool {
	switch rtype {
	case rTypeUint8, rTypeUint16, rTypeUint32, rTypeUint64:
		return true
	default:
		return false
	}
}

func RTypeIsFloatAny(rtype uintptr) bool {
	switch rtype {
	case rTypeFloat32, rTypeFloat64:
		return true
	default:
		return false
	}
}

func RTypeIsComplexAny(rtype uintptr) bool {
	switch rtype {
	case rTypeComplex64, rTypeComplex128:
		return true
	default:
		return false
	}
}
