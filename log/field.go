// Copyright © 2019. All rights reserved.
// Author: Ilya Yuryevich.
// Contacts: qioalice@gmail.com, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package log

import (
	"fmt"
	"math"
	"unsafe"

	"github.com/modern-go/reflect2"
)

// Field is an explicit field's type.
//
// Field stores some data most optimized way providing ability to use it
// as replacing of Golang interface{} but with more clear and optimized type checks.
// Thus you can then write your own Integrator and encode log's Entry fields that way
// you want.
//
// TBH, all implicit fields (both of named/unnamed) are converts to the explicit
// ones by internal private methods and thus all key-value pairs you want to attach
// to log's entry will be stored using this type - Field.
//
// WARNING!
// Do not instantiate Field object manually if you dont want how you have to do it.
// In most cases of bad initializations that Field will be considered invalid
// and do not handled then at the logging.
// Use constructors or make sure you understand what do you do.
type Field struct {

	// Relevant methods:
	// Entry.parseLogArgs, Entry.with, Entry.tryAddExplicitField, Entry.addField.

	// Key is a field's name.
	// Empty if it's unnamed explicit/implicit field.
	Key string

	// Kind represents what kind of field it is.
	// It's uint8 in the following format: XXXYYYYY, where:
	//   XXX - 3 highest bits - kind flags: nil, array, something else
	//   YYYYY - 5 lowest bits - used to store const of base type field's value.
	Kind FieldKind

	IValue int64  // for all ints, uints, floats, bool, complex64, pointers
	SValue string // for string, []byte, fmt.Stringer

	Value interface{} // nil for almost all easy cases
}

// FieldKind is an alias to uint8. Generally it's a way to store field's base type
// predefined const and flags. As described in 'Field.Kind' comments:
//
// It's uint8 in the following format: XXXYYYYY, where:
//   XXX - 3 highest bits - kind flags: nil, array, something else
//   YYYYY - 5 lowest bits - used to store const of base type field's value.
type FieldKind uint8

//
// TODO: Array support
// TODO: Map support
// TODO: Struct (classes) support

const (
	FieldKindMaskBaseType  = 0b_0001_1111
	FieldKindFlagArray     = 0b_0010_0000
	FieldKindFlagNil       = 0b_0100_0000
	fieldKindFlagReserved1 = 0b_1000_0000 // private because reserved

	/**/ // save this comment line to avoid gofmt
	/* 00 */
	FieldKindInvalid = iota // can't be handled in almost all cases
	/* 01 */ FieldKindUnknown // almost all implicit fields have this kind
	/* 02 */ _ // reserved
	/* 03 */ FieldKindBool
	/* 04 */ FieldKindInt
	/* 05 */ FieldKindInt8
	/* 06 */ FieldKindInt16
	/* 07 */ FieldKindInt32
	/* 08 */ FieldKindInt64
	/* 09 */ FieldKindUint
	/* 10 */ FieldKindUint8
	/* 11 */ FieldKindUint16
	/* 12 */ FieldKindUint32
	/* 13 */ FieldKindUint64
	/* 14 */ FieldKindUintptr
	/* 15 */ FieldKindFloat32
	/* 16 */ FieldKindFloat64
	/* 17 */ FieldKindComplex64
	/* 18 */ FieldKindComplex128
	/* 19 */ FieldKindString
	/* 20 */ _ // reserved
	/* 21 */ FieldKindAddr

	// -----
	// WARNING!
	// Keep in mind that max value of FieldKind base type is 31
	// (because of FieldKindMaskBaseType == 0b00011111 == 0x1F == 31).
	// -----
)

var (
	// Used for internal type comparision.
	reflectTypeField       = reflect2.TypeOf(Field{})
	reflectTypeFieldPtr    = reflect2.TypeOf((*Field)(nil))
	reflectTypeFmtStringer = reflect2.TypeOfPtr((*fmt.Stringer)(nil)).Elem()
)

// BaseType extracts only 5 lowest bits from fk and returns it (ignore flags).
//
// Call fk.BaseType() and then you can compare returned value with predefined
// FieldKind... constants. Do not compare directly, because fk can contain flags
// and then regular equal check (==) will fail.
func (fk FieldKind) BaseType() FieldKind {

	return fk & FieldKindMaskBaseType
}

// IsArray reports whether fk represents an array with some base type.
func (fk FieldKind) IsArray() bool {

	return fk&FieldKindFlagArray != 0
}

// IsNil reports whether fk represents a nil value.
//
// Returns true for both cases:
//   - fk is nil with some base type,
//   - fl is absolutely untyped nil.
func (fk FieldKind) IsNil() bool {

	return fk&FieldKindFlagNil != 0
}

// BaseType returns f's kind base type. You can use direct comparision operators
// (==, !=, etc) with returned value and FieldKind... constants.
func (f Field) BaseType() FieldKind {

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

// reset prepares f to be reused.
func (f *Field) reset() {

	f.Key = ""
	f.Kind = 0
	f.IValue, f.SValue, f.Value = 0, "", nil
}

// --------------------------- EASY CASES GENERATORS -------------------------- //

// Bool constructs a field with the given key and value.
func Bool(key string, value bool) Field {

	if value {
		return Field{Key: key, IValue: 1, Kind: FieldKindBool}
	} else {
		return Field{Key: key, IValue: 0, Kind: FieldKindBool}
	}
}

// Int constructs a field with the given key and value.
func Int(key string, value int) Field {
	return Field{Key: key, IValue: int64(value), Kind: FieldKindInt}
}

// Int8 constructs a field with the given key and value.
func Int8(key string, value int8) Field {
	return Field{Key: key, IValue: int64(value), Kind: FieldKindInt8}
}

// Int16 constructs a field with the given key and value.
func Int16(key string, value int16) Field {
	return Field{Key: key, IValue: int64(value), Kind: FieldKindInt16}
}

// Int32 constructs a field with the given key and value.
func Int32(key string, value int32) Field {
	return Field{Key: key, IValue: int64(value), Kind: FieldKindInt32}
}

// Int64 constructs a field with the given key and value.
func Int64(key string, value int64) Field {
	return Field{Key: key, IValue: value, Kind: FieldKindInt64}
}

// Uint constructs a field with the given key and value.
func Uint(key string, value uint) Field {
	return Field{Key: key, IValue: int64(value), Kind: FieldKindUint}
}

// Uint8 constructs a field with the given key and value.
func Uint8(key string, value uint8) Field {
	return Field{Key: key, IValue: int64(value), Kind: FieldKindUint8}
}

// Uint16 constructs a field with the given key and value.
func Uint16(key string, value uint16) Field {
	return Field{Key: key, IValue: int64(value), Kind: FieldKindUint16}
}

// Uint32 constructs a field with the given key and value.
func Uint32(key string, value uint32) Field {
	return Field{Key: key, IValue: int64(value), Kind: FieldKindUint32}
}

// Uint64 constructs a field with the given key and value.
func Uint64(key string, value uint64) Field {
	return Field{Key: key, IValue: int64(value), Kind: FieldKindUint64}
}

// Uintptr constructs a field with the given key and value.
func Uintptr(key string, value uintptr) Field {
	return Field{Key: key, IValue: int64(value), Kind: FieldKindUintptr}
}

// Float32 constructs a field with the given key and value.
func Float32(key string, value float32) Field {
	return Field{Key: key, IValue: int64(math.Float32bits(value)), Kind: FieldKindFloat32}
}

// Float64 constructs a field with the given key and value.
func Float64(key string, value float64) Field {
	return Field{Key: key, IValue: int64(math.Float64bits(value)), Kind: FieldKindFloat64}
}

// Complex64 constructs a field with the given key and value.
func Complex64(key string, value complex64) Field {
	return Field{Key: key, Value: value, Kind: FieldKindComplex64}
}

// Complex128 constructs a field with the given key and value.
func Complex128(key string, value complex128) Field {
	return Field{Key: key, Value: value, Kind: FieldKindComplex128}
}

// String constructs a field with the given key and value.
func String(key string, value string) Field {
	return Field{Key: key, SValue: value, Kind: FieldKindString}
}

// ------------------------- POINTER CASES GENERATORS ------------------------- //

// Boolp constructs a field that carries a *bool. The returned Field will safely
// and explicitly represent `nil` when appropriate.
func Boolp(key string, value *bool) Field {
	if value == nil {
		return fieldNilValue(key, FieldKindBool)
	}
	return Bool(key, *value)
}

// Intp constructs a field that carries a *int. The returned Field will safely
// and explicitly represent `nil` when appropriate.
func Intp(key string, value *int) Field {
	if value == nil {
		return fieldNilValue(key, FieldKindInt)
	}
	return Int(key, *value)
}

// Int8p constructs a field that carries a *int8. The returned Field will safely
// and explicitly represent `nil` when appropriate.
func Int8p(key string, value *int8) Field {
	if value == nil {
		return fieldNilValue(key, FieldKindInt8)
	}
	return Int8(key, *value)
}

// Int16p constructs a field that carries a *int16. The returned Field will safely
// and explicitly represent `nil` when appropriate.
func Int16p(key string, value *int16) Field {
	if value == nil {
		return fieldNilValue(key, FieldKindInt16)
	}
	return Int16(key, *value)
}

// Int32p constructs a field that carries a *int32. The returned Field will safely
// and explicitly represent `nil` when appropriate.
func Int32p(key string, value *int32) Field {
	if value == nil {
		return fieldNilValue(key, FieldKindInt32)
	}
	return Int32(key, *value)
}

// Int64p constructs a field that carries a *int64. The returned Field will safely
// and explicitly represent `nil` when appropriate.
func Int64p(key string, value *int64) Field {
	if value == nil {
		return fieldNilValue(key, FieldKindInt64)
	}
	return Int64(key, *value)
}

// Uintp constructs a field that carries a *uint. The returned Field will safely
// and explicitly represent `nil` when appropriate.
func Uintp(key string, value *uint) Field {
	if value == nil {
		return fieldNilValue(key, FieldKindUint)
	}
	return Uint(key, *value)
}

// Uint8p constructs a field that carries a *uint8. The returned Field will safely
// and explicitly represent `nil` when appropriate.
func Uint8p(key string, value *uint8) Field {
	if value == nil {
		return fieldNilValue(key, FieldKindUint8)
	}
	return Uint8(key, *value)
}

// Uint16p constructs a field that carries a *uint16. The returned Field will safely
// and explicitly represent `nil` when appropriate.
func Uint16p(key string, value *uint16) Field {
	if value == nil {
		return fieldNilValue(key, FieldKindUint16)
	}
	return Uint16(key, *value)
}

// Uint32p constructs a field that carries a *uint32. The returned Field will safely
// and explicitly represent `nil` when appropriate.
func Uint32p(key string, value *uint32) Field {
	if value == nil {
		return fieldNilValue(key, FieldKindUint32)
	}
	return Uint32(key, *value)
}

// Uint64p constructs a field that carries a *uint64. The returned Field will safely
// and explicitly represent `nil` when appropriate.
func Uint64p(key string, value *uint64) Field {
	if value == nil {
		return fieldNilValue(key, FieldKindUint64)
	}
	return Uint64(key, *value)
}

// Float32p constructs a field that carries a *float32. The returned Field will safely
// and explicitly represent `nil` when appropriate.
func Float32p(key string, value *float32) Field {
	if value == nil {
		return fieldNilValue(key, FieldKindFloat32)
	}
	return Float32(key, *value)
}

// Float64p constructs a field that carries a *float64. The returned Field will safely
// and explicitly represent `nil` when appropriate.
func Float64p(key string, value *float64) Field {
	if value == nil {
		return fieldNilValue(key, FieldKindFloat64)
	}
	return Float64(key, *value)
}

// ------------------------ COMPLEX CASES GENERATORS -------------------------- //

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
		return fieldNilValue(key, FieldKindString)
	}
	return Field{Key: key, SValue: value.String(), Kind: FieldKindString}
}

// Addr constructs a field that carries an any addr as is. E.g. If you want to print
// exactly addr of some var instead of its dereferenced value use this generator.
//
// All other generators that takes any pointer finally prints a value,
// addr points to. This func used to print exactly addr. Nil-safe.
func Addr(key string, value interface{}) Field {

	if value != nil {

		type golangInterface struct {
			typ  uintptr
			word unsafe.Pointer
		}

		v := (*golangInterface)(unsafe.Pointer(&value))
		return Field{Key: key, IValue: int64(uintptr(v.word)), Kind: FieldKindAddr}

	} else {
		return Field{Key: key, Kind: FieldKindAddr}
	}
}

// ---------------------- INTERNAL AUXILIARY FUNCTIONS ------------------------ //

// fieldNilValue creates a special field that indicates its store a nil value
// (nil pointer) to some baseType.
func fieldNilValue(key string, baseType FieldKind) Field {

	return Field{Key: key, Kind: baseType | FieldKindFlagNil}
}
