// Copyright © 2020. All rights reserved.
// Author: Ilya Stroy.
// Contacts: iyuryevich@pm.me, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekaclike

import (
	"time"
	"unsafe"

	"github.com/modern-go/reflect2"
)

//goland:noinspection GoVarAndConstTypeMayBeOmitted,GoRedundantConversion,GoBoolExpressions
var (
	RTypeBool               = reflect2.RTypeOf(bool(0 == 0))
	RTypeByte               = reflect2.RTypeOf(byte(0))
	RTypeRune               = reflect2.RTypeOf(rune(0))
	RTypeInt                = reflect2.RTypeOf(int(0))
	RTypeInt8               = reflect2.RTypeOf(int8(0))
	RTypeInt16              = reflect2.RTypeOf(int16(0))
	RTypeInt32              = reflect2.RTypeOf(int32(0))
	RTypeInt64              = reflect2.RTypeOf(int64(0))
	RTypeUint               = reflect2.RTypeOf(uint(0))
	RTypeUint8              = reflect2.RTypeOf(uint8(0))
	RTypeUint16             = reflect2.RTypeOf(uint16(0))
	RTypeUint32             = reflect2.RTypeOf(uint32(0))
	RTypeUint64             = reflect2.RTypeOf(uint64(0))
	RTypeFloat32            = reflect2.RTypeOf(float32(0))
	RTypeFloat64            = reflect2.RTypeOf(float64(0))
	RTypeComplex64          = reflect2.RTypeOf(complex64(0))
	RTypeComplex128         = reflect2.RTypeOf(complex128(0))
	RTypeString             = reflect2.RTypeOf(string(""))
	RTypeStringArray        = reflect2.RTypeOf([]string(nil))
	RTypeBytes              = reflect2.RTypeOf([]byte(nil))
	RTypeBytesArray         = reflect2.RTypeOf([][]byte(nil))
	RTypeMapStringString    = reflect2.RTypeOf(map[string]string(nil))
	RTypeMapStringInterface = reflect2.RTypeOf(map[string]any(nil))
	RTypeUintptr            = reflect2.RTypeOf(uintptr(0))
	RTypeUnsafePointer      = reflect2.RTypeOf(unsafe.Pointer(nil))
	RTypeTimeTime           = reflect2.RTypeOf(time.Time{})
	RTypeTimeDuration       = reflect2.RTypeOf(time.Duration(0))
)
