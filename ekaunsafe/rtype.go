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

// RTypeBool returns rtype of bool.
func RTypeBool() uintptr { return rTypeBool }

// RTypeByte returns rtype of byte.
func RTypeByte() uintptr { return rTypeByte }

// RTypeRune returns rtype of rune.
func RTypeRune() uintptr { return rTypeRune }

// RTypeInt returns rtype of int.
func RTypeInt() uintptr { return rTypeInt }

// RTypeInt8 returns rtype of int8.
func RTypeInt8() uintptr { return rTypeInt8 }

// RTypeInt16 returns rtype of int16.
func RTypeInt16() uintptr { return rTypeInt16 }

// RTypeInt32 returns rtype of int32.
func RTypeInt32() uintptr { return rTypeInt32 }

// RTypeInt64 returns rtype of int64.
func RTypeInt64() uintptr { return rTypeInt64 }

// RTypeUint returns rtype of uint.
func RTypeUint() uintptr { return rTypeUint }

// RTypeUint8 returns rtype of uint8.
func RTypeUint8() uintptr { return rTypeUint8 }

// RTypeUint16 returns rtype of uint16.
func RTypeUint16() uintptr { return rTypeUint16 }

// RTypeUint32 returns rtype of uint32.
func RTypeUint32() uintptr { return rTypeUint32 }

// RTypeUint64 returns rtype of uint64.
func RTypeUint64() uintptr { return rTypeUint64 }

// RTypeFloat32 returns rtype of float32.
func RTypeFloat32() uintptr { return rTypeFloat32 }

// RTypeFloat64 returns rtype of float64.
func RTypeFloat64() uintptr { return rTypeFloat64 }

// RTypeComplex64 returns rtype of complex64.
func RTypeComplex64() uintptr { return rTypeComplex64 }

// RTypeComplex128 returns rtype of complex128.
func RTypeComplex128() uintptr { return rTypeComplex128 }

// RTypeString returns rtype of string.
func RTypeString() uintptr { return rTypeString }

// RTypeStringArray returns rtype of []string.
func RTypeStringArray() uintptr { return rTypeStringArray }

// RTypeBytes returns rtype of []byte.
func RTypeBytes() uintptr { return rTypeBytes }

// RTypeBytesArray returns rtype of [][]byte.
func RTypeBytesArray() uintptr { return rTypeBytesArray }

// RTypeMapStringString returns rtype of map[string]string.
func RTypeMapStringString() uintptr { return rTypeMapStringString }

// RTypeMapStringAny returns rtype of map[string]any.
func RTypeMapStringAny() uintptr { return rTypeMapStringInterface }

// RTypeUintptr returns rtype of uintptr.
func RTypeUintptr() uintptr { return rTypeUintptr }

// RTypeUnsafePointer returns rtype of unsafe.Pointer.
func RTypeUnsafePointer() uintptr { return rTypeUnsafePointer }

// RTypeTimeTime returns rtype of time.Time.
func RTypeTimeTime() uintptr { return rTypeTimeTime }

// RTypeTimeDuration returns rtype of time.Duration.
func RTypeTimeDuration() uintptr { return rTypeTimeDuration }

// RTypeIsNumericAny returns true if provided rtype is of any int or uint type.
// Covers: int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64.
func RTypeIsNumericAny(rtype uintptr) bool {
	return RTypeIsIntAny(rtype) || RTypeIsUintAny(rtype)
}

// RTypeIsRealAny returns true if provided rtype is any numeric or float type.
// Covers: int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64,
// float32, float64.
func RTypeIsRealAny(rtype uintptr) bool {
	return RTypeIsNumericAny(rtype) || RTypeIsFloatAny(rtype)
}

// RTypeIsIntAny returns true if provided rtype is any int type.
// Covers: int, int8, int16, int32, int64.
func RTypeIsIntAny(rtype uintptr) bool {
	return rtype == rTypeInt || RTypeIsIntFixed(rtype)
}

// RTypeIsIntFixed returns true if provided rtype is fixed int type.
// Covers: int8, int16, int32, int64.
func RTypeIsIntFixed(rtype uintptr) bool {
	return rtype == rTypeInt8 || rtype == rTypeInt16 || rtype == rTypeInt32 ||
		rtype == rTypeInt64
}

// RTypeIsUintAny returns true if provided rtype is any uint type.
// Covers: uint, uint8, uint16, uint32, uint64.
func RTypeIsUintAny(rtype uintptr) bool {
	return rtype == rTypeUint || RTypeIsUintFixed(rtype)
}

// RTypeIsUintFixed returns true if provided rtype is fixed uint type.
// Covers: uint8, uint16, uint32, uint64.
func RTypeIsUintFixed(rtype uintptr) bool {
	return rtype == rTypeUint8 || rtype == rTypeUint16 ||
		rtype == rTypeUint32 || rtype == rTypeUint64
}

// RTypeIsFloatAny returns true if provided rtype is any float type.
// Covers: float32, float64.
func RTypeIsFloatAny(rtype uintptr) bool {
	return rtype == rTypeFloat32 || rtype == rTypeFloat64
}

// RTypeIsComplexAny returns true if provided rtype is any complex type.
// Covers: complex64, complex128.
func RTypeIsComplexAny(rtype uintptr) bool {
	return rtype == rTypeComplex64 || rtype == rTypeComplex128
}

// RTypeIsStringLike returns true if provided rtype is anything string like
// (or can be cast to string). Covers: string, []byte.
func RTypeIsStringLike(rtype uintptr) bool {
	return rtype == rTypeString || rtype == rTypeBytes
}
