// Copyright © 2019. All rights reserved.
// Author: Ilya Yuryevich.
// Contacts: qioalice@gmail.com, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package field

import (
	"fmt"
	"math"

	"github.com/qioalice/ekago/ekadanger"

	"github.com/modern-go/reflect2"
)

// Inspired by: https://github.com/uber-go/zap/blob/master/zapcore/field.go

// TODO: Array support
// TODO: Map support
// TODO: Struct (classes) support

type (
	// Field is an explicit logger or error field's type.
	//
	// Field stores some data most optimized way providing ability to use it
	// as replacing of Golang interface{} but with more clear and optimized type
	// checks. Thus you can then write your own Integrator and encode log Entry's
	// or Error's fields the way you want.
	//
	// TBH, all implicit fields (both of named/unnamed) are converts
	// to the explicit ones by internal private methods and thus all key-value
	// pairs you want to attach to log's entry will be stored using this type.
	//
	// WARNING!
	// DO NOT INSTANTIATE FIELD OBJECT MANUALLY IF YOU DONT KNOW HOW TO USE IT.
	// In most cases of bad initializations that Field will be considered invalid
	// and do not handled then at the logging.
	// USE CONSTRUCTORS OR MAKE SURE YOU UNDERSTAND WHAT DO YOU DO.
	Field struct {

		// Key is a field's name.
		// Empty if it's unnamed explicit/implicit field.
		Key string

		// Kind represents what kind of field it is.
		//
		// WARNING!
		// If you want to compare Kind with any of Kind... constants, use
		// Field.Kind.BaseType() or Field.BaseType() method before!
		// For more info see these methods docs.
		Kind Kind

		IValue int64  // for all ints, uints, floats, bool, complex64, pointers
		SValue string // for string, []byte, fmt.Stringer

		Value interface{} // for all not easy cases
	}

	// Kind is an alias to uint8. Generally it's a way to store field's base type
	// predefined const and flags. As described in 'Field.Kind' comments:
	//
	// It's uint8 in the following format: XXXYYYYY, where:
	//   XXX - 3 highest bits - kind flags: nil, array, something else
	//   YYYYY - 5 lowest bits - used to store const of base type field's value.
	Kind uint8
)

//noinspection GoSnakeCaseUsage
const (
	KIND_MASK_BASE_TYPE = 0b_0001_1111
	KIND_FLAG_ARRAY     = 0b_0010_0000
	KIND_FLAG_NULL      = 0b_0100_0000
	KIND_FLAG_SYSTEM    = 0b_1000_0000

	KIND_TYPE_INVALID = 0 // can't be handled in almost all cases

	// field.Kind & KIND_MASK_BASE_TYPE could be any of listed below,
	// only if field.Kind KIND_FLAG_INTERNAL_SYS != 0 (system letter's field)

	KIND_SYS_TYPE_EKAERR_UUID           = 1
	KIND_SYS_TYPE_EKAERR_CLASS_ID       = 2
	KIND_SYS_TYPE_EKAERR_CLASS_NAME     = 3
	KIND_SYS_TYPE_EKAERR_PUBLIC_MESSAGE = 4

	// field.Kind & KIND_MASK_BASE_TYPE could be any of listed below,
	// only if field.Kind & KIND_FLAG_INTERNAL_SYS == 0 (user's field)

	_                     = 1  // reserved
	_                     = 2  // reserved
	KIND_TYPE_BOOL        = 3  // uses IValue to store bool
	KIND_TYPE_INT         = 4  // uses IValue to store int
	KIND_TYPE_INT_8       = 5  // uses IValue to store int8
	KIND_TYPE_INT_16      = 6  // uses IValue to store int16
	KIND_TYPE_INT_32      = 7  // uses IValue to store int32
	KIND_TYPE_INT_64      = 8  // uses IValue to store int64
	KIND_TYPE_UINT        = 9  // uses IValue to store uint
	KIND_TYPE_UINT_8      = 10 // uses IValue to store uint8
	KIND_TYPE_UINT_16     = 11 // uses IValue to store uint16
	KIND_TYPE_UINT_32     = 12 // uses IValue to store uint32
	KIND_TYPE_UINT_64     = 13 // uses IValue to store uint64
	KIND_TYPE_UINTPTR     = 14 // uses IValue to store uintptr
	KIND_TYPE_FLOAT_32    = 15 // uses IValue to store float32 (bits)
	KIND_TYPE_FLOAT64     = 16 // uses IValue to store float64 (bits)
	KIND_TYPE_COMPLEX_64  = 17 // uses IValue to store complex64
	KIND_TYPE_COMPLEX_128 = 18 // uses Value (interface{}) to store complex128
	KIND_TYPE_STRING      = 19 // uses SValue to store string
	_                     = 20 // reserved
	KIND_TYPE_ADDR        = 21 // uses IValue to store some addr (like uintptr)
	_                     = 31 // reserved, max, range [21..31] is free now

	// --------------------------------------------------------------------- //
	//                                WARNING                                //
	// Keep in mind that max value of Kind base type is 31              //
	// (because of KIND_MASK_BASE_TYPE == 0b00011111 == 0x1F == 31).       //
	// DO NOT OVERFLOW THIS VALUE WHEN YOU WILL ADD A NEW CONSTANTS          //
	// --------------------------------------------------------------------- //
)

var (
	// Used for type comparision.
	ReflectedType            = reflect2.TypeOf(Field{})
	ReflectedTypePtr         = reflect2.TypeOf((*Field)(nil))
	ReflectedTypeFmtStringer = reflect2.TypeOfPtr((*fmt.Stringer)(nil)).Elem()
)

// BaseType extracts only 5 lowest bits from fk and returns it (ignore flags).
//
// Call fk.BaseType() and then you can compare returned value with predefined
// Kind... constants. DO NOT COMPARE DIRECTLY, because fk can contain flags
// and then regular equal check (==) will fail.
func (fk Kind) BaseType() Kind {
	return fk & KIND_MASK_BASE_TYPE
}

// IsArray reports whether fk represents an array with some base type.
func (fk Kind) IsArray() bool {
	return fk&KIND_FLAG_ARRAY != 0
}

// IsNil reports whether fk represents a nil value.
//
// Returns true for both cases:
//   - fk is nil with some base type,
//   - fk is absolutely untyped nil.
func (fk Kind) IsNil() bool {
	return fk&KIND_FLAG_NULL != 0
}

// IsSystem reports whether fk represents a *Letter system field.
// See https://github.com/qioalice/ekago/internal/letter/letter.go .
func (fk Kind) IsSystem() bool {
	return fk&KIND_FLAG_SYSTEM != 0
}

// BaseType returns f's kind base type. You can use direct comparision operators
// (==, !=, etc) with returned value and Kind... constants.
func (f Field) BaseType() Kind {
	return f.Kind.BaseType()
}

// IsArray reports whether f represents an array with some base type.
func (f Field) IsArray() bool {
	return f.Kind.IsArray()
}

// IsNil reports whether f represents a nil value.
//
// Returns true for both cases:
//   - f stores nil as value of some base type,
//   - f stores nil and its absolutely untyped nil.
func (f Field) IsNil() bool {
	return f.Kind.IsNil()
}

// IsSystem reports whether f represents a *Letter system field.
// See https://github.com/qioalice/ekago/internal/letter/letter.go .
func (f Field) IsSystem() bool {
	return f.Kind.IsSystem()
}

// Reset frees all allocated resources (RAM in 99% cases) by Field f, preparing
// it for being reused in the future.
func Reset(f *Field) {
	f.Key = ""
	f.Kind = 0
	f.IValue, f.SValue, f.Value = 0, "", nil
}

// --------------------------- EASY CASES GENERATORS -------------------------- //
// ---------------------------------------------------------------------------- //

// Bool constructs a field with the given key and value.
func Bool(key string, value bool) Field {
	if value {
		return Field{Key: key, IValue: 1, Kind: KIND_TYPE_BOOL}
	} else {
		return Field{Key: key, IValue: 0, Kind: KIND_TYPE_BOOL}
	}
}

// Int constructs a field with the given key and value.
func Int(key string, value int) Field {
	return Field{Key: key, IValue: int64(value), Kind: KIND_TYPE_INT}
}

// Int8 constructs a field with the given key and value.
func Int8(key string, value int8) Field {
	return Field{Key: key, IValue: int64(value), Kind: KIND_TYPE_INT_8}
}

// Int16 constructs a field with the given key and value.
func Int16(key string, value int16) Field {
	return Field{Key: key, IValue: int64(value), Kind: KIND_TYPE_INT_16}
}

// Int32 constructs a field with the given key and value.
func Int32(key string, value int32) Field {
	return Field{Key: key, IValue: int64(value), Kind: KIND_TYPE_INT_32}
}

// Int64 constructs a field with the given key and value.
func Int64(key string, value int64) Field {
	return Field{Key: key, IValue: value, Kind: KIND_TYPE_INT_64}
}

// Uint constructs a field with the given key and value.
func Uint(key string, value uint) Field {
	return Field{Key: key, IValue: int64(value), Kind: KIND_TYPE_UINT}
}

// Uint8 constructs a field with the given key and value.
func Uint8(key string, value uint8) Field {
	return Field{Key: key, IValue: int64(value), Kind: KIND_TYPE_UINT_8}
}

// Uint16 constructs a field with the given key and value.
func Uint16(key string, value uint16) Field {
	return Field{Key: key, IValue: int64(value), Kind: KIND_TYPE_UINT_16}
}

// Uint32 constructs a field with the given key and value.
func Uint32(key string, value uint32) Field {
	return Field{Key: key, IValue: int64(value), Kind: KIND_TYPE_UINT_32}
}

// Uint64 constructs a field with the given key and value.
func Uint64(key string, value uint64) Field {
	return Field{Key: key, IValue: int64(value), Kind: KIND_TYPE_UINT_64}
}

// Uintptr constructs a field with the given key and value.
func Uintptr(key string, value uintptr) Field {
	return Field{Key: key, IValue: int64(value), Kind: KIND_TYPE_UINTPTR}
}

// Float32 constructs a field with the given key and value.
func Float32(key string, value float32) Field {
	return Field{Key: key, IValue: int64(math.Float32bits(value)), Kind: KIND_TYPE_FLOAT_32}
}

// Float64 constructs a field with the given key and value.
func Float64(key string, value float64) Field {
	return Field{Key: key, IValue: int64(math.Float64bits(value)), Kind: KIND_TYPE_FLOAT64}
}

// Complex64 constructs a field with the given key and value.
func Complex64(key string, value complex64) Field {
	return Field{Key: key, Value: value, Kind: KIND_TYPE_COMPLEX_64}
}

// Complex128 constructs a field with the given key and value.
func Complex128(key string, value complex128) Field {
	return Field{Key: key, Value: value, Kind: KIND_TYPE_COMPLEX_128}
}

// String constructs a field with the given key and value.
func String(key string, value string) Field {
	return Field{Key: key, SValue: value, Kind: KIND_TYPE_STRING}
}

// ------------------------- POINTER CASES GENERATORS ------------------------- //
// ---------------------------------------------------------------------------- //

// Boolp constructs a field that carries a *bool. The returned Field will safely
// and explicitly represent `nil` when appropriate.
func Boolp(key string, value *bool) Field {
	if value == nil {
		return NilValue(key, KIND_TYPE_BOOL)
	}
	return Bool(key, *value)
}

// Intp constructs a field that carries a *int. The returned Field will safely
// and explicitly represent `nil` when appropriate.
func Intp(key string, value *int) Field {
	if value == nil {
		return NilValue(key, KIND_TYPE_INT)
	}
	return Int(key, *value)
}

// Int8p constructs a field that carries a *int8. The returned Field will safely
// and explicitly represent `nil` when appropriate.
func Int8p(key string, value *int8) Field {
	if value == nil {
		return NilValue(key, KIND_TYPE_INT_8)
	}
	return Int8(key, *value)
}

// Int16p constructs a field that carries a *int16. The returned Field will safely
// and explicitly represent `nil` when appropriate.
func Int16p(key string, value *int16) Field {
	if value == nil {
		return NilValue(key, KIND_TYPE_INT_16)
	}
	return Int16(key, *value)
}

// Int32p constructs a field that carries a *int32. The returned Field will safely
// and explicitly represent `nil` when appropriate.
func Int32p(key string, value *int32) Field {
	if value == nil {
		return NilValue(key, KIND_TYPE_INT_32)
	}
	return Int32(key, *value)
}

// Int64p constructs a field that carries a *int64. The returned Field will safely
// and explicitly represent `nil` when appropriate.
func Int64p(key string, value *int64) Field {
	if value == nil {
		return NilValue(key, KIND_TYPE_INT_64)
	}
	return Int64(key, *value)
}

// Uintp constructs a field that carries a *uint. The returned Field will safely
// and explicitly represent `nil` when appropriate.
func Uintp(key string, value *uint) Field {
	if value == nil {
		return NilValue(key, KIND_TYPE_UINT)
	}
	return Uint(key, *value)
}

// Uint8p constructs a field that carries a *uint8. The returned Field will safely
// and explicitly represent `nil` when appropriate.
func Uint8p(key string, value *uint8) Field {
	if value == nil {
		return NilValue(key, KIND_TYPE_UINT_8)
	}
	return Uint8(key, *value)
}

// Uint16p constructs a field that carries a *uint16. The returned Field will safely
// and explicitly represent `nil` when appropriate.
func Uint16p(key string, value *uint16) Field {
	if value == nil {
		return NilValue(key, KIND_TYPE_UINT_16)
	}
	return Uint16(key, *value)
}

// Uint32p constructs a field that carries a *uint32. The returned Field will safely
// and explicitly represent `nil` when appropriate.
func Uint32p(key string, value *uint32) Field {
	if value == nil {
		return NilValue(key, KIND_TYPE_UINT_32)
	}
	return Uint32(key, *value)
}

// Uint64p constructs a field that carries a *uint64. The returned Field will safely
// and explicitly represent `nil` when appropriate.
func Uint64p(key string, value *uint64) Field {
	if value == nil {
		return NilValue(key, KIND_TYPE_UINT_64)
	}
	return Uint64(key, *value)
}

// Float32p constructs a field that carries a *float32. The returned Field will safely
// and explicitly represent `nil` when appropriate.
func Float32p(key string, value *float32) Field {
	if value == nil {
		return NilValue(key, KIND_TYPE_FLOAT_32)
	}
	return Float32(key, *value)
}

// Float64p constructs a field that carries a *float64. The returned Field will safely
// and explicitly represent `nil` when appropriate.
func Float64p(key string, value *float64) Field {
	if value == nil {
		return NilValue(key, KIND_TYPE_FLOAT64)
	}
	return Float64(key, *value)
}

// ------------------------ COMPLEX CASES GENERATORS -------------------------- //
// ---------------------------------------------------------------------------- //

// Type constructs a field that holds on value's type as string.
func Type(key string, value interface{}) Field {
	if value == nil {
		return String(key, "<nil>")
	}
	return String(key, reflect2.TypeOf(value).String()) // SIGSEGV if value == nil
}

// Stringer constructs a field that holds on string generated by fmt.Stringer.String().
// The returned Field will safely and explicitly represent `nil` when appropriate.
func Stringer(key string, value fmt.Stringer) Field {
	if value == nil {
		return NilValue(key, KIND_TYPE_STRING)
	}
	return Field{Key: key, SValue: value.String(), Kind: KIND_TYPE_STRING}
}

// Addr constructs a field that carries an any addr as is. E.g. If you want to print
// exactly addr of some var instead of its dereferenced value use this generator.
//
// All other generators that takes any pointer finally prints a value,
// addr points to. This func used to print exactly addr. Nil-safe.
func Addr(key string, value interface{}) Field {
	if value != nil {

		addr := ekadanger.TakeRealAddr(value)
		return Field{Key: key, IValue: int64(uintptr(addr)), Kind: KIND_TYPE_ADDR}

	} else {
		return Field{Key: key, Kind: KIND_TYPE_ADDR}
	}
}

// ---------------------- INTERNAL AUXILIARY FUNCTIONS ------------------------ //
// ---------------------------------------------------------------------------- //

// NilValue creates a special field that indicates its store a nil value
// (nil pointer) to some baseType.
func NilValue(key string, baseType Kind) Field {
	return Field{Key: key, Kind: baseType | KIND_FLAG_NULL}
}
