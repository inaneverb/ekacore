// Copyright © 2020. All rights reserved.
// Author: Ilya Stroy.
// Contacts: qioalice@gmail.com, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekaunsafe

import (
	"github.com/qioalice/ekago/v3/internal/ekaclike"
)

/*
RTypeBool is a "constant" function.
Always returns the same value.

Returns an integer representation of pointer to the type that describes
builtin Golang "bool" type.

Useful along with reflect2.RTypeOf() function.
*/
//go:inline
func RTypeBool() uintptr { return ekaclike.RTypeBool }

/*
RTypeByte is a "constant" function.
Always returns the same value.

Returns an integer representation of pointer to the type that describes
builtin Golang "byte" type.

Useful along with reflect2.RTypeOf() function.
*/
//go:inline
func RTypeByte() uintptr { return ekaclike.RTypeByte }

/*
RTypeRune is a "constant" function.
Always returns the same value.

Returns an integer representation of pointer to the type that describes
builtin Golang "rune" type.

Useful along with reflect2.RTypeOf() function.
*/
//go:inline
func RTypeRune() uintptr { return ekaclike.RTypeRune }

/*
RTypeInt is a "constant" function.
Always returns the same value.

Returns an integer representation of pointer to the type that describes
builtin Golang "int" type (not "int8", "int16", "int32", "int64"!)

Useful along with reflect2.RTypeOf() function.
*/
//go:inline
func RTypeInt() uintptr { return ekaclike.RTypeInt }

/*
RTypeInt8 is a "constant" function.
Always returns the same value.

Returns an integer representation of pointer to the type that describes
builtin Golang "int8" type.

Useful along with reflect2.RTypeOf() function.
*/
//go:inline
func RTypeInt8() uintptr { return ekaclike.RTypeInt8 }

/*
RTypeInt16 is a "constant" function.
Always returns the same value.

Returns an integer representation of pointer to the type that describes
builtin Golang "int16" type.

Useful along with reflect2.RTypeOf() function.
*/
//go:inline
func RTypeInt16() uintptr { return ekaclike.RTypeInt16 }

/*
RTypeInt32 is a "constant" function.
Always returns the same value.

Returns an integer representation of pointer to the type that describes
builtin Golang "int32" type.

Useful along with reflect2.RTypeOf() function.
*/
//go:inline
func RTypeInt32() uintptr { return ekaclike.RTypeInt32 }

/*
RTypeInt64 is a "constant" function.
Always returns the same value.

Returns an integer representation of pointer to the type that describes
builtin Golang "int64" type.

Useful along with reflect2.RTypeOf() function.
*/
//go:inline
func RTypeInt64() uintptr { return ekaclike.RTypeInt64 }

/*
RTypeUint is a "constant" function.
Always returns the same value.

Returns an integer representation of pointer to the type that describes
builtin Golang "uint" type (not "uint8", "uint16", "uint32", "uint64"!)

Useful along with reflect2.RTypeOf() function.
*/
//go:inline
func RTypeUint() uintptr { return ekaclike.RTypeUint }

/*
RTypeUint8 is a "constant" function.
Always returns the same value.

Returns an integer representation of pointer to the type that describes
builtin Golang "uint8" type.

Useful along with reflect2.RTypeOf() function.
*/
//go:inline
func RTypeUint8() uintptr { return ekaclike.RTypeUint8 }

/*
RTypeUint16 is a "constant" function.
Always returns the same value.

Returns an integer representation of pointer to the type that describes
builtin Golang "uint16" type.

Useful along with reflect2.RTypeOf() function.
*/
//go:inline
func RTypeUint16() uintptr { return ekaclike.RTypeUint16 }

/*
RTypeUint32 is a "constant" function.
Always returns the same value.

Returns an integer representation of pointer to the type that describes
builtin Golang "uint32" type.

Useful along with reflect2.RTypeOf() function.
*/
//go:inline
func RTypeUint32() uintptr { return ekaclike.RTypeUint32 }

/*
RTypeUint64 is a "constant" function.
Always returns the same value.

Returns an integer representation of pointer to the type that describes
builtin Golang "uint64" type.

Useful along with reflect2.RTypeOf() function.
*/
//go:inline
func RTypeUint64() uintptr { return ekaclike.RTypeUint64 }

/*
RTypeFloat32 is a "constant" function.
Always returns the same value.

Returns an integer representation of pointer to the type that describes
builtin Golang "float32" type.

Useful along with reflect2.RTypeOf() function.
*/
//go:inline
func RTypeFloat32() uintptr { return ekaclike.RTypeFloat32 }

/*
RTypeFloat64 is a "constant" function.
Always returns the same value.

Returns an integer representation of pointer to the type that describes
builtin Golang "float64" type.

Useful along with reflect2.RTypeOf() function.
*/
//go:inline
func RTypeFloat64() uintptr { return ekaclike.RTypeFloat64 }

/*
RTypeComplex64 is a "constant" function.
Always returns the same value.

Returns an integer representation of pointer to the type that describes
builtin Golang "complex64" type.

Useful along with reflect2.RTypeOf() function.
*/
//go:inline
func RTypeComplex64() uintptr { return ekaclike.RTypeComplex64 }

/*
RTypeComplex128 is a "constant" function.
Always returns the same value.

Returns an integer representation of pointer to the type that describes
builtin Golang "complex128" type.

Useful along with reflect2.RTypeOf() function.
*/
//go:inline
func RTypeComplex128() uintptr { return ekaclike.RTypeComplex128 }

/*
RTypeString is a "constant" function.
Always returns the same value.

Returns an integer representation of pointer to the type that describes
builtin Golang "string" type.

Useful along with reflect2.RTypeOf() function.
*/
//go:inline
func RTypeString() uintptr { return ekaclike.RTypeString }

/*
RTypeStringArray is a "constant" function.
Always returns the same value.

Returns an integer representation of pointer to the type that describes
complex Golang "[]string" type.

Useful along with reflect2.RTypeOf() function.
*/
//go:inline
func RTypeStringArray() uintptr { return ekaclike.RTypeStringArray }

/*
RTypeBytes is a "constant" function.
Always returns the same value.

Returns an integer representation of pointer to the type that describes
complex Golang "[]byte" type.

Useful along with reflect2.RTypeOf() function.
*/
//go:inline
func RTypeBytes() uintptr { return ekaclike.RTypeBytes }

/*
RTypeBytesArray is a "constant" function.
Always returns the same value.

Returns an integer representation of pointer to the type that describes
complex Golang "[][]byte" type.

Useful along with reflect2.RTypeOf() function.
*/
//go:inline
func RTypeBytesArray() uintptr { return ekaclike.RTypeBytesArray }

/*
RTypeMapStringString is a "constant" function.
Always returns the same value.

Returns an integer representation of pointer to the type that describes
complex Golang "map[string]string" type.

Useful along with reflect2.RTypeOf() function.
*/
//go:inline
func RTypeMapStringString() uintptr { return ekaclike.RTypeMapStringString }

/*
RTypeMapStringInterface is a "constant" function.
Always returns the same value.

Returns an integer representation of pointer to the type that describes
complex Golang "map[string]interface{}" type.

Useful along with reflect2.RTypeOf() function.
*/
//go:inline
func RTypeMapStringInterface() uintptr { return ekaclike.RTypeMapStringInterface }

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
	return rtype == ekaclike.RTypeInt || RTypeIsIntFixed(rtype)
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
	case ekaclike.RTypeInt8, ekaclike.RTypeInt16, ekaclike.RTypeInt32, ekaclike.RTypeInt64:
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
	return rtype == ekaclike.RTypeUint || RTypeIsUintFixed(rtype)
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
	case ekaclike.RTypeUint8, ekaclike.RTypeUint16, ekaclike.RTypeUint32, ekaclike.RTypeUint64:
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
	case ekaclike.RTypeFloat32, ekaclike.RTypeFloat64:
		return true
	default:
		return false
	}
}

/*
RTypeIsComplexAny returns true if passed rtype is any of complexes only.
It means that it's true for complex64, complex128.
*/
func RTypeIsComplexAny(rtype uintptr) bool {
	switch rtype {
	case ekaclike.RTypeComplex64, ekaclike.RTypeComplex128:
		return true
	default:
		return false
	}
}
