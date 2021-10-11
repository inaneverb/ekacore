// Copyright © 2021. All rights reserved.
// Author: Ilya Yuryevich.
// Contacts: iyuryevich@pm.me, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekaletter

import (
	"fmt"
	"math"
	"reflect"
	"time"
	"unsafe"

	"github.com/qioalice/ekago/v3/ekastr"
	"github.com/qioalice/ekago/v3/internal/ekaclike"

	"github.com/modern-go/reflect2"
)

// Inspired by: https://github.com/uber-go/zap/blob/master/zapcore/field.go

type (
	// LetterField is an abstract type that holds some value and type of that value.
	// It generally uses to represent fields of ekaerr.Error or ekalog.Entry
	// but can be used in your own code.
	//
	// LetterField stores some data most optimized way providing ability to use it
	// as replacing of Golang interface{} but with more clear and optimized type
	// checks. Thus you can then write your own ekalog.Integrator
	// and encode log ekalog.Entry's or ekaerr.Error's fields the way you want.
	//
	// WARNING!
	// DO NOT INSTANTIATE FIELD OBJECT MANUALLY IF YOU DONT KNOW HOW TO USE IT.
	// In most cases of bad initializations that LetterField will be considered invalid
	// and do not handled then at the logging.
	// USE CONSTRUCTORS OR MAKE SURE YOU UNDERSTAND WHAT DO YOU DO.
	LetterField struct {

		// Key is a field's name.
		// Empty if it's unnamed explicit/implicit field.
		Key string

		// Kind represents what the fuck this LetterField is.
		//
		// WARNING!
		// If you want to compare Kind with any of Kind... constants, use
		// LetterField.LetterFieldKind.BaseType() or LetterField.BaseType() method before!
		// For more info see these methods docs.
		Kind LetterFieldKind

		IValue int64  // for all ints, uints, floats, bool, complex64, pointers
		SValue string // for string, []byte, fmt.Stringer (called)

		Value interface{} // for all not easy cases

		// StackFrameIdx contains a number of stack frame, this LetterField
		StackFrameIdx int16
	}

	// LetterFieldKind is an alias to uint8.
	// Generally it's a way to store field's base type predefined const and flags.
	// As described in LetterField.Kind comments:
	//
	// It's uint8 in the following format: XXXYYYYY, where:
	//   XXX - 3 highest bits - kind flags: nil, array, something else
	//   YYYYY - 5 lowest bits - used to store const of base type field's value.
	LetterFieldKind uint8
)

//noinspection GoSnakeCaseUsage
const (
	KIND_MASK_BASE_TYPE = 0b_0001_1111
	KIND_MASK_FLAGS     = 0b_1110_0000

	KIND_FLAG_USER_DEFINED = 0b_0010_0000 // reserved for user's purposes
	KIND_FLAG_NULL         = 0b_0100_0000
	KIND_FLAG_SYSTEM       = 0b_1000_0000

	KIND_TYPE_INVALID = ^LetterFieldKind(0) // can't be handled in almost all cases

	// field.LetterFieldKind & KIND_MASK_BASE_TYPE could be any of listed below,
	// only if field.LetterFieldKind KIND_FLAG_INTERNAL_SYS != 0 (system letter's field)

	KIND_SYS_TYPE_EKAERR_UUID       = 1
	KIND_SYS_TYPE_EKAERR_CLASS_ID   = 2
	KIND_SYS_TYPE_EKAERR_CLASS_NAME = 3

	// field.LetterFieldKind & KIND_MASK_BASE_TYPE could be any of listed below,
	// only if field.LetterFieldKind & KIND_FLAG_INTERNAL_SYS == 0 (user's field)

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
	KIND_TYPE_FLOAT_64    = 16 // uses IValue to store float64 (bits)
	KIND_TYPE_COMPLEX_64  = 17 // uses IValue to store complex64
	KIND_TYPE_COMPLEX_128 = 18 // uses Value (interface{}) to store complex128
	KIND_TYPE_STRING      = 19 // uses SValue to store string
	_                     = 20 // reserved
	KIND_TYPE_ADDR        = 21 // uses IValue to store some addr (like uintptr)
	_                     = 22 // reserved
	KIND_TYPE_UNIX        = 23 // uses IValue to store int64 unixtime sec
	KIND_TYPE_UNIX_NANO   = 24 // uses IValue to store int64 unixtime nanosec
	KIND_TYPE_DURATION    = 25 // uses IValue to store int64 duration in nanosec
	_                     = 26 // reserved
	KIND_TYPE_ARRAY       = 27 // uses Value (interface{}) to store []T or [N]T
	KIND_TYPE_MAP         = 28 // uses Value (interface{}) to store map[T1]T2
	KIND_TYPE_EXTMAP      = 29 // uses Value (interface{}) to store map[T1]T2
	KIND_TYPE_STRUCT      = 30 // uses Value (interface{}) to store struct{<...>}
	_                     = 31 // reserved

	// --------------------------------------------------------------------- //
	//                                WARNING                                //
	// Keep in mind that max value of LetterFieldKind base type is 31        //
	// (because of KIND_MASK_BASE_TYPE == 0b00011111 == 0x1F == 31).         //
	//                                                                       //
	// DO NOT OVERFLOW THIS VALUE WHEN YOU WILL ADD A NEW CONSTANTS          //
	// DO NOT FORGOT TO UPDATE IsValidBaseType() METHOD                      //
	// --------------------------------------------------------------------- //
)

var (
	// Used for type comparison.

	RTypeLetterField    = reflect2.RTypeOf(LetterField{})
	RTypeLetterFieldPtr = reflect2.RTypeOf((*LetterField)(nil))
	TypeFmtStringer     = reflect2.TypeOfPtr((*fmt.Stringer)(nil)).Elem()
)

//noinspection GoErrorStringFormat
var (
	ErrUnsupportedKind = fmt.Errorf("LetterField: Unsupported kind of LetterField.")
)

// BaseType extracts only 5 lowest bits from LetterFieldKind and returns it (ignore flags).
//
// Call BaseType() and then you can compare returned value with predefined
// KIND_<...> constants. DO NOT COMPARE DIRECTLY, because LetterFieldKind can contain flags
// and then regular equal check (==) will fail.
func (fk LetterFieldKind) BaseType() LetterFieldKind {
	return fk & KIND_MASK_BASE_TYPE
}

// IsValidBaseType reports whether LetterFieldKind's BaseType is valid
// (uses any of predefined KIND_TYPE_<...> BaseType's constant).
func (fk LetterFieldKind) IsValidBaseType() bool {
	bt := fk.BaseType()
	return (bt >= KIND_TYPE_BOOL && bt <= KIND_TYPE_STRING) ||
		(bt >= KIND_TYPE_UNIX && bt <= KIND_TYPE_DURATION) ||
		(bt >= KIND_TYPE_ARRAY && bt <= KIND_TYPE_EXTMAP) ||
		bt == KIND_TYPE_ADDR
}

// IsNil reports whether LetterFieldKind represents a nil value.
//
// Returns true for both cases:
//   - LetterFieldKind is nil with some base type (e.g: nil *int, nil []int, etc),
//   - LetterFieldKind is absolutely untyped nil (Golang's nil).
func (fk LetterFieldKind) IsNil() bool {
	return fk&KIND_FLAG_NULL != 0
}

// IsSystem reports whether LetterFieldKind represents a Letter's system field.
func (fk LetterFieldKind) IsSystem() bool {
	return fk&KIND_FLAG_SYSTEM != 0
}

// IsInvalid reports whether LetterFieldKind represents an invalid LetterField.
func (fk LetterFieldKind) IsInvalid() bool {
	return fk == KIND_TYPE_INVALID
}

// BaseType returns LetterField's LetterFieldKind base type.
// It's the same as LetterField.LetterFieldKind.BaseType().
//
// You can use direct comparison operators like EQ, NEQ with returned value
// and KIND_<...> constants.
func (f LetterField) BaseType() LetterFieldKind {
	return f.Kind.BaseType()
}

// IsNil reports whether LetterField represents a nil value.
//
// Returns true for both cases:
//   - LetterField stores nil as value of some base type (e.g: nil *int, nil []int, etc),
//   - LetterField stores nil and its absolutely untyped nil (Golang's nil).
func (f LetterField) IsNil() bool {
	return f.Kind.IsNil()
}

// IsSystem reports whether LetterField represents a Letter's system field.
func (f LetterField) IsSystem() bool {
	return f.Kind.IsSystem()
}

// IsInvalid reports whether LetterField is invalid LetterField.
func (f LetterField) IsInvalid() bool {
	return f.Kind.IsInvalid()
}

// IsZero reports whether LetterField contains zero value of its type (based on kind).
func (f LetterField) IsZero() bool {
	return f.IValue == 0 &&
		(f.SValue == "" || f.SValue == "00000000-0000-0000-0000-000000000000") &&
		f.Value == nil
}

// FieldReset frees all allocated resources (RAM in 99% cases) by LetterField, preparing
// it for being reused in the future.
func FieldReset(f *LetterField) {
	f.Key = ""
	f.Kind = KIND_TYPE_INVALID
	f.IValue, f.SValue, f.Value = 0, "", nil
}

// --------------------------- EASY CASES GENERATORS -------------------------- //
// ---------------------------------------------------------------------------- //

// FBool constructs a field with the given key and value.
func FBool(key string, value bool) LetterField {
	if value {
		return LetterField{Key: key, IValue: 1, Kind: KIND_TYPE_BOOL}
	} else {
		return LetterField{Key: key, IValue: 0, Kind: KIND_TYPE_BOOL}
	}
}

// FInt constructs a field with the given key and value.
func FInt(key string, value int) LetterField {
	return LetterField{Key: key, IValue: int64(value), Kind: KIND_TYPE_INT}
}

// FInt8 constructs a field with the given key and value.
func FInt8(key string, value int8) LetterField {
	return LetterField{Key: key, IValue: int64(value), Kind: KIND_TYPE_INT_8}
}

// FInt16 constructs a field with the given key and value.
func FInt16(key string, value int16) LetterField {
	return LetterField{Key: key, IValue: int64(value), Kind: KIND_TYPE_INT_16}
}

// FInt32 constructs a field with the given key and value.
func FInt32(key string, value int32) LetterField {
	return LetterField{Key: key, IValue: int64(value), Kind: KIND_TYPE_INT_32}
}

// FInt64 constructs a field with the given key and value.
func FInt64(key string, value int64) LetterField {
	return LetterField{Key: key, IValue: value, Kind: KIND_TYPE_INT_64}
}

// FUint constructs a field with the given key and value.
func FUint(key string, value uint) LetterField {
	return LetterField{Key: key, IValue: int64(value), Kind: KIND_TYPE_UINT}
}

// FUint8 constructs a field with the given key and value.
func FUint8(key string, value uint8) LetterField {
	return LetterField{Key: key, IValue: int64(value), Kind: KIND_TYPE_UINT_8}
}

// FUint16 constructs a field with the given key and value.
func FUint16(key string, value uint16) LetterField {
	return LetterField{Key: key, IValue: int64(value), Kind: KIND_TYPE_UINT_16}
}

// FUint32 constructs a field with the given key and value.
func FUint32(key string, value uint32) LetterField {
	return LetterField{Key: key, IValue: int64(value), Kind: KIND_TYPE_UINT_32}
}

// FUint64 constructs a field with the given key and value.
func FUint64(key string, value uint64) LetterField {
	return LetterField{Key: key, IValue: int64(value), Kind: KIND_TYPE_UINT_64}
}

// FUintptr constructs a field with the given key and value.
func FUintptr(key string, value uintptr) LetterField {
	return LetterField{Key: key, IValue: int64(value), Kind: KIND_TYPE_UINTPTR}
}

// FFloat32 constructs a field with the given key and value.
func FFloat32(key string, value float32) LetterField {
	return LetterField{Key: key, IValue: int64(math.Float32bits(value)), Kind: KIND_TYPE_FLOAT_32}
}

// FFloat64 constructs a field with the given key and value.
func FFloat64(key string, value float64) LetterField {
	return LetterField{Key: key, IValue: int64(math.Float64bits(value)), Kind: KIND_TYPE_FLOAT_64}
}

// FComplex64 constructs a field with the given key and value.
func FComplex64(key string, value complex64) LetterField {
	r, i := math.Float32bits(real(value)), math.Float32bits(imag(value))
	return LetterField{Key: key, IValue: (int64(r) << 32) | int64(i), Kind: KIND_TYPE_COMPLEX_64}
}

// FComplex128 constructs a field with the given key and value.
func FComplex128(key string, value complex128) LetterField {
	return LetterField{Key: key, Value: value, Kind: KIND_TYPE_COMPLEX_128}
}

// FString constructs a field with the given key and value.
func FString(key string, value string) LetterField {
	return LetterField{Key: key, SValue: value, Kind: KIND_TYPE_STRING}
}

// FStringFromBytes constructs a field with the given key and value.
func FStringFromBytes(key string, value []byte) LetterField {
	return FString(key, ekastr.B2S(value))
}

// ------------------------- POINTER CASES GENERATORS ------------------------- //
// ---------------------------------------------------------------------------- //

// FBoolp constructs a field that carries a *bool. The returned LetterField will safely
// and explicitly represent nil when appropriate.
func FBoolp(key string, value *bool) LetterField {
	if value == nil {
		return FNil(key, KIND_TYPE_BOOL)
	}
	return FBool(key, *value)
}

// FIntp constructs a field that carries a *int. The returned LetterField will safely
// and explicitly represent nil when appropriate.
func FIntp(key string, value *int) LetterField {
	if value == nil {
		return FNil(key, KIND_TYPE_INT)
	}
	return FInt(key, *value)
}

// FInt8p constructs a field that carries a *int8. The returned LetterField will safely
// and explicitly represent nil when appropriate.
func FInt8p(key string, value *int8) LetterField {
	if value == nil {
		return FNil(key, KIND_TYPE_INT_8)
	}
	return FInt8(key, *value)
}

// FInt16p constructs a field that carries a *int16. The returned LetterField will safely
// and explicitly represent nil when appropriate.
func FInt16p(key string, value *int16) LetterField {
	if value == nil {
		return FNil(key, KIND_TYPE_INT_16)
	}
	return FInt16(key, *value)
}

// FInt32p constructs a field that carries a *int32. The returned LetterField will safely
// and explicitly represent nil when appropriate.
func FInt32p(key string, value *int32) LetterField {
	if value == nil {
		return FNil(key, KIND_TYPE_INT_32)
	}
	return FInt32(key, *value)
}

// FInt64p constructs a field that carries a *int64. The returned LetterField will safely
// and explicitly represent nil when appropriate.
func FInt64p(key string, value *int64) LetterField {
	if value == nil {
		return FNil(key, KIND_TYPE_INT_64)
	}
	return FInt64(key, *value)
}

// FUintp constructs a field that carries a *uint. The returned LetterField will safely
// and explicitly represent nil when appropriate.
func FUintp(key string, value *uint) LetterField {
	if value == nil {
		return FNil(key, KIND_TYPE_UINT)
	}
	return FUint(key, *value)
}

// FUint8p constructs a field that carries a *uint8. The returned LetterField will safely
// and explicitly represent nil when appropriate.
func FUint8p(key string, value *uint8) LetterField {
	if value == nil {
		return FNil(key, KIND_TYPE_UINT_8)
	}
	return FUint8(key, *value)
}

// FUint16p constructs a field that carries a *uint16. The returned LetterField will safely
// and explicitly represent nil when appropriate.
func FUint16p(key string, value *uint16) LetterField {
	if value == nil {
		return FNil(key, KIND_TYPE_UINT_16)
	}
	return FUint16(key, *value)
}

// FUint32p constructs a field that carries a *uint32. The returned LetterField will safely
// and explicitly represent nil when appropriate.
func FUint32p(key string, value *uint32) LetterField {
	if value == nil {
		return FNil(key, KIND_TYPE_UINT_32)
	}
	return FUint32(key, *value)
}

// FUint64p constructs a field that carries a *uint64. The returned LetterField will safely
// and explicitly represent nil when appropriate.
func FUint64p(key string, value *uint64) LetterField {
	if value == nil {
		return FNil(key, KIND_TYPE_UINT_64)
	}
	return FUint64(key, *value)
}

// FFloat32p constructs a field that carries a *float32. The returned LetterField will safely
// and explicitly represent nil when appropriate.
func FFloat32p(key string, value *float32) LetterField {
	if value == nil {
		return FNil(key, KIND_TYPE_FLOAT_32)
	}
	return FFloat32(key, *value)
}

// FFloat64p constructs a field that carries a *float64. The returned LetterField will safely
// and explicitly represent nil when appropriate.
func FFloat64p(key string, value *float64) LetterField {
	if value == nil {
		return FNil(key, KIND_TYPE_FLOAT_64)
	}
	return FFloat64(key, *value)
}

// ------------------------ COMPLEX CASES GENERATORS -------------------------- //
// ---------------------------------------------------------------------------- //

// FType constructs a field that holds on value's type as string.
func FType(key string, value interface{}) LetterField {
	if value == nil {
		return FString(key, "<nil>")
	}
	return FString(key, reflect2.TypeOf(value).String()) // SIGSEGV if value == nil
}

// FStringer constructs a field that holds on string generated by fmt.Stringer.String().
// The returned LetterField will safely and explicitly represent nil when appropriate.
func FStringer(key string, value fmt.Stringer) LetterField {
	if value == nil {
		return FNil(key, KIND_TYPE_STRING)
	}
	return FString(key, value.String())
}

// FAddr constructs a field that carries an any addr as is. E.g. If you want to print
// exactly addr of some var instead of its dereferenced value use this generator.
//
// All other generators that takes any pointer finally prints a value,
// addr points to. This func used to print exactly addr. Nil-safe.
//
// WARNING!
// The resource, value's ptr points to may be GC'ed and unavailable.
// So it's unsafe to try to do something with that addr but comparing.
func FAddr(key string, value interface{}) LetterField {
	if value != nil {
		addr := ekaclike.TakeRealAddr(value)
		return LetterField{Key: key, IValue: int64(uintptr(addr)), Kind: KIND_TYPE_ADDR}
	}
	return FNil(key, KIND_TYPE_ADDR)
}

// FUnixFromStd constructs a field with the given time.Time and its key.
// It treats time.Time as unixtime in sec, so ms, us, ns are unavailable.
func FUnixFromStd(key string, t time.Time) LetterField {
	return LetterField{Key: key, IValue: t.Unix(), Kind: KIND_TYPE_UNIX}
}

// FUnixNanoFromStd constructs a field with the given time.Time and its key.
// It treats time.Time as unixtime in nanosec, so it's more precision then UnixFromStd(),
// but can represent only a [1678..2262] years.
func FUnixNanoFromStd(key string, t time.Time) LetterField {
	return LetterField{Key: key, IValue: t.UnixNano(), Kind: KIND_TYPE_UNIX_NANO}
}

// FUnix constructs a LetterField that represents passed int64 as unixtime in sec.
func FUnix(key string, unix int64) LetterField {
	return LetterField{Key: key, IValue: unix, Kind: KIND_TYPE_UNIX}
}

// FUnixNano constructs a LetterField that represents passed int64 as unixtime in nanosec.
func FUnixNano(key string, unixNano int64) LetterField {
	return LetterField{Key: key, IValue: unixNano, Kind: KIND_TYPE_UNIX_NANO}
}

// FDuration constructs a field with given time.Duration and its key.
func FDuration(key string, d time.Duration) LetterField {
	return LetterField{Key: key, IValue: d.Nanoseconds(), Kind: KIND_TYPE_DURATION}
}

// ----------------------- DIFFICULT CASES GENERATORS ------------------------- //
// ---------------------------------------------------------------------------- //

// FArray returns a LetterField that will represent value only if it's a Golang slice.
// Despite of the name, Golang ARRAYS DOES NOT SUPPORT.
// You can covert Golang's array to Golang's slice using [:] slice operations.
// If value is nil, returns FNil(key, KIND_FLAG_ARRAY).
// If value is not a slice, invalid LetterField is returned.
func FArray(key string, value interface{}) LetterField {
	if value == nil {
		return FNil(key, KIND_TYPE_ARRAY)
	}
	if k := reflect2.TypeOf(value).Kind(); k != reflect.Slice && k != reflect.Array {
		return FInvalid(key)
	}
	return LetterField{Key: key, Value: value, Kind: KIND_TYPE_ARRAY}
}

// FObject returns a LetterField that will represent value only if it's struct{<...>}.
// If value is nil, returns FNil(key, KIND_TYPE_STRUCT).
// If value is not a struct{}, invalid LetterField is returned.
func FObject(key string, value interface{}) LetterField {
	if value == nil {
		return FNil(key, KIND_TYPE_STRUCT)
	}
	if reflect2.TypeOf(value).Kind() != reflect.Struct {
		return FInvalid(key)
	}
	return LetterField{Key: key, Value: value, Kind: KIND_TYPE_STRUCT}
}

// FMap returns a LetterField that will represent value only if it's map[T1]T2.
// If value is nil, returns FNil(key, KIND_TYPE_MAP).
// If value is not a map[T1]T2, invalid LetterField is returned.
func FMap(key string, value interface{}) LetterField {
	if value == nil {
		return FNil(key, KIND_TYPE_MAP)
	}
	if reflect2.TypeOf(value).Kind() != reflect.Map {
		return FInvalid(key)
	}
	return LetterField{Key: key, Value: value, Kind: KIND_TYPE_MAP}
}

// FExtractedMap is the same as FMap(), but a map will be added as it would be
// a LetterField's array to the result storage of LetterField.
//
// WARNING!
// Make sure, your LetterField's worker supports extracted maps.
func FExtractedMap(key string, value map[string]interface{}) LetterField {
	if value == nil {
		return FNil(key, KIND_TYPE_MAP)
	}
	return LetterField{Key: key, Value: value, Kind: KIND_TYPE_EXTMAP}
}

// FAny tries to recognize type of passed value and then call specific
// LetterField generator if it's possible.
//
// If no specific generator can be used, an invalid LetterField is returned.
//
// Value must not be Golang's nil. Otherwise an invalid LetterField will be returned
// and such field is skipped at the parsing.
func FAny(key string, value interface{}) LetterField {
	eface := ekaclike.UnpackInterface(value)

	if eface.Type == 0 && eface.Word == nil {
		return FNil(key, 0)
	}

	typ := reflect2.TypeOf(value)

	switch eface.Type {
	case ekaclike.RTypeBool:
		var boolVal bool
		typ.UnsafeSet(unsafe.Pointer(&boolVal), eface.Word)
		return FBool(key, boolVal)

	case ekaclike.RTypeInt:
		var intVal int
		typ.UnsafeSet(unsafe.Pointer(&intVal), eface.Word)
		return FInt(key, intVal)

	case ekaclike.RTypeInt8:
		var int8Val int8
		typ.UnsafeSet(unsafe.Pointer(&int8Val), eface.Word)
		return FInt8(key, int8Val)

	case ekaclike.RTypeInt16:
		var int16Val int16
		typ.UnsafeSet(unsafe.Pointer(&int16Val), eface.Word)
		return FInt16(key, int16Val)

	case ekaclike.RTypeInt32:
		var int32Val int32
		typ.UnsafeSet(unsafe.Pointer(&int32Val), eface.Word)
		return FInt32(key, int32Val)

	case ekaclike.RTypeInt64:
		var int64Val int64
		typ.UnsafeSet(unsafe.Pointer(&int64Val), eface.Word)
		return FInt64(key, int64Val)

	case ekaclike.RTypeUint:
		var uintVal uint64
		typ.UnsafeSet(unsafe.Pointer(&uintVal), eface.Word)
		return FUint64(key, uintVal)

	case ekaclike.RTypeUint8:
		var uint8Val uint8
		typ.UnsafeSet(unsafe.Pointer(&uint8Val), eface.Word)
		return FUint8(key, uint8Val)

	case ekaclike.RTypeUint16:
		var uint16Val uint16
		typ.UnsafeSet(unsafe.Pointer(&uint16Val), eface.Word)
		return FUint16(key, uint16Val)

	case ekaclike.RTypeUint32:
		var uint32Val uint32
		typ.UnsafeSet(unsafe.Pointer(&uint32Val), eface.Word)
		return FUint32(key, uint32Val)

	case ekaclike.RTypeUint64:
		var uint64Val uint64
		typ.UnsafeSet(unsafe.Pointer(&uint64Val), eface.Word)
		return FUint64(key, uint64Val)

	case ekaclike.RTypeFloat32:
		var float32Val float32
		typ.UnsafeSet(unsafe.Pointer(&float32Val), eface.Word)
		return FFloat32(key, float32Val)

	case ekaclike.RTypeFloat64:
		var float64Val float64
		typ.UnsafeSet(unsafe.Pointer(&float64Val), eface.Word)
		return FFloat64(key, float64Val)

	case ekaclike.RTypeComplex64:
		var complex64Val complex64
		typ.UnsafeSet(unsafe.Pointer(&complex64Val), eface.Word)
		return FComplex64(key, complex64Val)

	case ekaclike.RTypeComplex128:
		var complex128Val complex128
		typ.UnsafeSet(unsafe.Pointer(&complex128Val), eface.Word)
		return FComplex128(key, complex128Val)

	case ekaclike.RTypeString:
		var stringVal string
		typ.UnsafeSet(unsafe.Pointer(&stringVal), eface.Word)
		return FString(key, stringVal)

	case ekaclike.RTypeTimeTime:
		var timeVal time.Time
		typ.UnsafeSet(unsafe.Pointer(&timeVal), eface.Word)
		return FUnixFromStd(key, timeVal)

	case ekaclike.RTypeTimeDuration:
		var durationVal time.Duration
		typ.UnsafeSet(unsafe.Pointer(&durationVal), eface.Word)
		return FDuration(key, durationVal)

	case ekaclike.RTypeUintptr, ekaclike.RTypeUnsafePointer:
		return FAddr(key, value)
	}

	if typ.Implements(TypeFmtStringer) {
		return FStringer(key, value.(fmt.Stringer))
	}

	switch typ.Kind() {

	case reflect.Array, reflect.Slice:
		return FArray(key, value)

	case reflect.Struct:
		return FObject(key, value)

	case reflect.Map:
		return FMap(key, value)

	case reflect.Ptr:
		if eface.Word != nil {
			return FAny(key, typ.Indirect(value))
		} else {
			return FNil(key, KIND_TYPE_ADDR)
		}
	}

	return FInvalid(key)
}

// ---------------------- INTERNAL AUXILIARY FUNCTIONS ------------------------ //
// ---------------------------------------------------------------------------- //

// FNil creates a special field that indicates its store a nil value
// (nil pointer) to some baseType.
func FNil(key string, baseType LetterFieldKind) LetterField {
	return LetterField{Key: key, Kind: baseType | KIND_FLAG_NULL}
}

// FInvalid returns an invalid LetterField with the given key.
func FInvalid(key string) LetterField {
	return LetterField{Key: key, Kind: KIND_TYPE_INVALID}
}
