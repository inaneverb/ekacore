// Copyright © 2020. All rights reserved.
// Author: Ilya Stroy.
// Contacts: qioalice@gmail.com, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekaunsafe

import (
	"github.com/modern-go/reflect2"
)

//goland:noinspection GoVarAndConstTypeMayBeOmitted,GoRedundantConversion,GoBoolExpressions
var (
	rtypeBool        uintptr = reflect2.RTypeOf(bool(0 == 0))
	rtypeByte        uintptr = reflect2.RTypeOf(byte(0))
	rtypeRune        uintptr = reflect2.RTypeOf(rune(0))
	rtypeInt         uintptr = reflect2.RTypeOf(int(0))
	rtypeInt8        uintptr = reflect2.RTypeOf(int8(0))
	rtypeInt16       uintptr = reflect2.RTypeOf(int16(0))
	rtypeInt32       uintptr = reflect2.RTypeOf(int32(0))
	rtypeInt64       uintptr = reflect2.RTypeOf(int64(0))
	rtypeUint        uintptr = reflect2.RTypeOf(uint(0))
	rtypeUint8       uintptr = reflect2.RTypeOf(uint8(0))
	rtypeUint16      uintptr = reflect2.RTypeOf(uint16(0))
	rtypeUint32      uintptr = reflect2.RTypeOf(uint32(0))
	rtypeUint64      uintptr = reflect2.RTypeOf(uint64(0))
	rtypeFloat32     uintptr = reflect2.RTypeOf(float32(0))
	rtypeFloat64     uintptr = reflect2.RTypeOf(float64(0))
	rtypeString      uintptr = reflect2.RTypeOf(string(""))
	rtypeStringArray uintptr = reflect2.RTypeOf([]string(nil))
	rtypeBytes       uintptr = reflect2.RTypeOf([]byte(nil))
	rtypeBytesArray  uintptr = reflect2.RTypeOf([][]byte(nil))
)

/*
RTypeBool is a "constant" function.
Always returns the same value.

Returns an integer representation of pointer to the type that describes
builtin Golang "bool" type.

Useful along with reflect2.RTypeOf() function.
*/
//go:inline
func RTypeBool() uintptr { return rtypeBool }

/*
RTypeByte is a "constant" function.
Always returns the same value.

Returns an integer representation of pointer to the type that describes
builtin Golang "byte" type.

Useful along with reflect2.RTypeOf() function.
*/
//go:inline
func RTypeByte() uintptr { return rtypeByte }

/*
RTypeRune is a "constant" function.
Always returns the same value.

Returns an integer representation of pointer to the type that describes
builtin Golang "rune" type.

Useful along with reflect2.RTypeOf() function.
*/
//go:inline
func RTypeRune() uintptr { return rtypeRune }

/*
RTypeInt is a "constant" function.
Always returns the same value.

Returns an integer representation of pointer to the type that describes
builtin Golang "int" type (not "int8", "int16", "int32", "int64"!)

Useful along with reflect2.RTypeOf() function.
*/
//go:inline
func RTypeInt() uintptr { return rtypeInt }

/*
RTypeInt8 is a "constant" function.
Always returns the same value.

Returns an integer representation of pointer to the type that describes
builtin Golang "int8" type.

Useful along with reflect2.RTypeOf() function.
*/
//go:inline
func RTypeInt8() uintptr { return rtypeInt8 }

/*
RTypeInt16 is a "constant" function.
Always returns the same value.

Returns an integer representation of pointer to the type that describes
builtin Golang "int16" type.

Useful along with reflect2.RTypeOf() function.
*/
//go:inline
func RTypeInt16() uintptr { return rtypeInt16 }

/*
RTypeInt32 is a "constant" function.
Always returns the same value.

Returns an integer representation of pointer to the type that describes
builtin Golang "int32" type.

Useful along with reflect2.RTypeOf() function.
*/
//go:inline
func RTypeInt32() uintptr { return rtypeInt32 }

/*
RTypeInt64 is a "constant" function.
Always returns the same value.

Returns an integer representation of pointer to the type that describes
builtin Golang "int64" type.

Useful along with reflect2.RTypeOf() function.
*/
//go:inline
func RTypeInt64() uintptr { return rtypeInt64 }

/*
RTypeUint is a "constant" function.
Always returns the same value.

Returns an integer representation of pointer to the type that describes
builtin Golang "uint" type (not "uint8", "uint16", "uint32", "uint64"!)

Useful along with reflect2.RTypeOf() function.
*/
//go:inline
func RTypeUint() uintptr { return rtypeUint }

/*
RTypeUint8 is a "constant" function.
Always returns the same value.

Returns an integer representation of pointer to the type that describes
builtin Golang "uint8" type.

Useful along with reflect2.RTypeOf() function.
*/
//go:inline
func RTypeUint8() uintptr { return rtypeUint8 }

/*
RTypeUint16 is a "constant" function.
Always returns the same value.

Returns an integer representation of pointer to the type that describes
builtin Golang "uint16" type.

Useful along with reflect2.RTypeOf() function.
*/
//go:inline
func RTypeUint16() uintptr { return rtypeUint16 }

/*
RTypeUint32 is a "constant" function.
Always returns the same value.

Returns an integer representation of pointer to the type that describes
builtin Golang "uint32" type.

Useful along with reflect2.RTypeOf() function.
*/
//go:inline
func RTypeUint32() uintptr { return rtypeUint32 }

/*
RTypeUint64 is a "constant" function.
Always returns the same value.

Returns an integer representation of pointer to the type that describes
builtin Golang "uint64" type.

Useful along with reflect2.RTypeOf() function.
*/
//go:inline
func RTypeUint64() uintptr { return rtypeUint64 }

/*
RTypeFloat32 is a "constant" function.
Always returns the same value.

Returns an integer representation of pointer to the type that describes
builtin Golang "float32" type.

Useful along with reflect2.RTypeOf() function.
*/
//go:inline
func RTypeFloat32() uintptr { return rtypeFloat32 }

/*
RTypeFloat64 is a "constant" function.
Always returns the same value.

Returns an integer representation of pointer to the type that describes
builtin Golang "float64" type.

Useful along with reflect2.RTypeOf() function.
*/
//go:inline
func RTypeFloat64() uintptr { return rtypeFloat64 }

/*
RTypeString is a "constant" function.
Always returns the same value.

Returns an integer representation of pointer to the type that describes
builtin Golang "string" type.

Useful along with reflect2.RTypeOf() function.
*/
//go:inline
func RTypeString() uintptr { return rtypeString }

/*
RTypeStringArray is a "constant" function.
Always returns the same value.

Returns an integer representation of pointer to the type that describes
builtin Golang "[]string" type.

Useful along with reflect2.RTypeOf() function.
*/
//go:inline
func RTypeStringArray() uintptr { return rtypeStringArray }

/*
RTypeBytes is a "constant" function.
Always returns the same value.

Returns an integer representation of pointer to the type that describes
builtin Golang "[]byte" type.

Useful along with reflect2.RTypeOf() function.
*/
//go:inline
func RTypeBytes() uintptr { return rtypeBytes }

/*
RTypeBytesArray is a "constant" function.
Always returns the same value.

Returns an integer representation of pointer to the type that describes
builtin Golang "[][]byte" type.

Useful along with reflect2.RTypeOf() function.
*/
//go:inline
func RTypeBytesArray() uintptr { return rtypeBytesArray }

/*
RTypeIsAnyNumeric returns true if passed rtype is any of signed or unsigned integers.
It means, that it's true for
 - int, int8, int16, int32, int64,
 - uint, uint8, uint16, uint32, uint64
but it also is true for
 - byte (because it's Golang alias to uint8),
 - rune (because it's Golang alias to int32).
*/
//go:inline
func RTypeIsAnyNumeric(rtype uintptr) bool {
	return RTypeIsIntAny(rtype) || RTypeIsUintAny(rtype)
}
/*
RTypeIsAnyReal returns true if passed rtype is any of signed or unsigned integers
or floats.
It means that it's true for
 - int, int8, int16, int32, int64,
 - uint, uint8, uint16, uint32, uint64,
 - float32, float64
but it also is true for
 - byte (because it's Golang alias to uint8),
 - rune (because it's Golang alias to int32).
*/
//go:inline
func RTypeIsAnyReal(rtype uintptr) bool {
	return RTypeIsAnyNumeric(rtype) || RTypeIsFloatAny(rtype)
}

/*
RTypeIsIntAny returns true if passed rtype is any of signed integers only.
It means that it's true for int, int8, int16, int32, int64,
but it also true for rune, because it's Golang alias for int32.
*/
//go:inline
func RTypeIsIntAny(rtype uintptr) bool {
	return rtype == rtypeInt || RTypeIsIntFixed(rtype)
}

/*
RTypeIsIntFixed returns true if passed rtype is any of signed fixed length integers only.
It means that it's true for int8, int16, int32, int64,
but it also true for rune, because it's Golang alias for int32.

NOTE!
Again. It false for just "int"!
*/
//go:inline
func RTypeIsIntFixed(rtype uintptr) bool {
	switch rtype {
	case rtypeInt8, rtypeInt16, rtypeInt32, rtypeInt64:
		return true
	default:
		return false
	}
}

/*
RTypeIsUintAny returns true if passed rtype is any of unsigned integers only.
It means that it's true for uint, uint8, uint16, uint32, uint64,
but it also true for byte, because it's Golang alias for uint8.
*/
//go:inline
func RTypeIsUintAny(rtype uintptr) bool {
	return rtype == rtypeUint || RTypeIsUintFixed(rtype)
}

/*
RTypeIsUintFixed returns true if passed rtype is any of unsigned fixed length integers only.
It means that it's true for uint8, uint16, uint32, uint64,
but it also true for byte, because it's Golang alias for uint8.

NOTE!
Again. It false for just "uint"!
*/
//go:inline
func RTypeIsUintFixed(rtype uintptr) bool {
	switch rtype {
	case rtypeUint8, rtypeUint16, rtypeUint32, rtypeUint64:
		return true
	default:
		return false
	}
}

/*
RTypeIsFloatAny returns true if passed rtype is any of floats only.
It means that it's true for float32, float64.
*/
//go:inline
func RTypeIsFloatAny(rtype uintptr) bool {
	switch rtype {
	case rtypeFloat32, rtypeFloat64:
		return true
	default:
		return false
	}
}
