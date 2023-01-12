// Copyright © 2020. All rights reserved.
// Author: Ilya Stroy.
// Contacts: iyuryevich@pm.me, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekaunsafe

import (
	"fmt"
	"reflect"
	"strconv"
	"unsafe"

	"github.com/modern-go/reflect2"

	"github.com/qioalice/ekago/v4/ekaext"
)

const (
	ToStringHandleArrays        uint8 = 0x01 // Print arrays
	ToStringHandleStructs       uint8 = 0x02 // Print structs
	ToStringHandleMaps          uint8 = 0x04 // Print maps
	ToStringDereferencePointers uint8 = 0x08 // Dereference pointers before
)

// ToInt64 casts (or converts) 'v' to int64 and returns it.
// If that operation is impossible or meaningless, (0, false) is returned.
//
// (impossible or meaningless operations are those, that are not listed
// in the list of rules below).
//
// Rules:
//
//   - Any int, int8, int16, int32, int64, uint, uint8, uint16, uint32
//     are best case scenario. They will be converted to int64 and that's all.
//
//   - uint64 is "on the edge" scenario. It will be converted using standard
//     Golang's op. So, you should expect to get a negative int64 values,
//     for those uint64, that are overflows int64 max. But it's likely UB.
//
//   - float32, float64 are also converted using standard Golang's op.
//     You should expect a truncating.
//
//   - 1 for boolean true, 0 for boolean false. Like C.
//
//   - Equivalent of len() call for non-empty string, array, slice, map.
//     0 for either nil or empty ones.
//
//   - 1 for non-nil func, chan, pointer. 0 for nil ones.
//
//   - Structs are not convertible. So, (0, false) is returned.
//     (But pointer to struct is allowed, see above).
//
// WARNING!
// It's C-style. Use it on your own risk, and make sure
// you're passing the type that is covered by any rule. UB otherwise.
func ToInt64(v any) (int64, bool) {

	var rtype, word = UnpackInterface(v).Tuple()
	if n, ok := toInt64FastUnsafe(rtype, word); ok || rtype == 0 {
		return n, ok
	} else {
		return toInt64Complex(v, word)
	}
}

// ToInt64Fast is a simple version of ToInt64() that does the same job,
// but some rules are discarded because of rejecting of reflect or reflect2
// packages usage. It improves an execution time comparable to ToInt64().
//
// Difference with ToInt64():
//
//   - Supporting only base types.
//     int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64,
//     bool, byte, rune, float32, float64, string.
//
//   - For string also equivalent of len(string) call.
//
//   - Supports nil pointers and types, that can be treated as nil pointer.
//     Like []T(nil), map[T]V(nil). 0 is returned for them.
//
//   - For all other types (0, false) is returned.
//
// WARNING!
// Keep in mind, that ToInt64() fully supports arrays, slices, maps.
// For example, ToInt64() of non-empty slice will return len of that slice,
// and ToInt64Fast() return (0, false). But for nil (not just empty),
// both of ToInt64(), ToInt64Fast() returns (0, true).
func ToInt64Fast(v any) (int64, bool) {
	return toInt64FastUnsafe(UnpackInterface(v).Tuple())
}

// ToBool is the same as ToInt64() but with checking 1st arg is not zero,
// and 2nd one is true. Return that.
func ToBool(v any) bool {
	var n, ok = ToInt64(v)
	return ok && n != 0
}

// ToString is a function that allows you to get a string representation
// of any type's passed argument. It returns an empty string for nil.
func ToString(v any) string {
	var rtype, word = UnpackInterface(v).Tuple()
	return ToStringUnsafe(rtype, word, 0xFF)
}

// ToStringUnsafe is the same as ToString() but allows you to reject handling
// arrays, structs, maps, pointers.
//
// It also expects to get unpacked Golang interface as both of rtype's
// and data's pointers. Zero safe. Returns an empty string then.
func ToStringUnsafe(rtype uintptr, word unsafe.Pointer, mask uint8) string {

	if word == nil && rtype != 0 {
		switch {
		case rtype == rTypeBool || RTypeIsStringLike(rtype):
			return ""
		case RTypeIsRealAny(rtype) || RTypeIsComplexAny(rtype):
			return "0"
		}
		return toStringUnsafeComplex(rtype, word, mask)
	}

	switch {

	case rtype == 0:
		return ""

	case rtype == rTypeString:
		return *(*string)(word)

	case rtype == rTypeBytes:
		return BytesToString(*(*[]byte)(word))

	case rtype == rTypeBool:
		return ekaext.If(*(*bool)(word), "true", "false")

	case rtype == rTypeInt:
		return strconv.FormatInt(int64(*(*int)(word)), 10)
	case rtype == rTypeInt8:
		return strconv.FormatInt(int64(*(*int8)(word)), 10)
	case rtype == rTypeInt16:
		return strconv.FormatInt(int64(*(*int16)(word)), 10)
	case rtype == rTypeInt32:
		return strconv.FormatInt(int64(*(*int32)(word)), 10)
	case rtype == rTypeInt64:
		return strconv.FormatInt(*(*int64)(word), 10)

	case rtype == rTypeUint:
		return strconv.FormatUint(uint64(*(*uint)(word)), 10)
	case rtype == rTypeUint8:
		return strconv.FormatUint(uint64(*(*uint8)(word)), 10)
	case rtype == rTypeUint16:
		return strconv.FormatUint(uint64(*(*uint16)(word)), 10)
	case rtype == rTypeUint32:
		return strconv.FormatUint(uint64(*(*uint32)(word)), 10)
	case rtype == rTypeUint64:
		return strconv.FormatUint(*(*uint64)(word), 10)

	case rtype == rTypeUintptr:

	case rtype == rTypeFloat32:
		return strconv.FormatFloat(float64(*(*float32)(word)), 'f', 2, 32)
	case rtype == rTypeFloat64:
		return strconv.FormatFloat(*(*float64)(word), 'f', 2, 64)

	case rtype == rTypeComplex64:
		return strconv.FormatComplex(complex128(*(*complex64)(word)), 'f', 2, 64)
	case rtype == rTypeComplex128:
		return strconv.FormatComplex(*(*complex128)(word), 'f', 2, 128)
	}

	return toStringUnsafeComplex(rtype, word, mask)
}

// toInt64FastUnsafe does the job, ToInt64Fast() describes of.
func toInt64FastUnsafe(rtype uintptr, word unsafe.Pointer) (int64, bool) {

	if rtype == 0 || word == nil {
		// If 'word' is nil, it's "complex" type - a func, maybe pointer, map.
		// If 'rtype' is 0, it's a very strange case and definitely bug.
		return 0, rtype != 0
	}

	switch {

	case rtype == rTypeString:
		return int64(len(*(*string)(word))), true

	case rtype == rTypeBool:
		return ekaext.If(*(*bool)(word), int64(1), 0), true

	case rtype == rTypeRune:
		return int64(*(*rune)(word)), true

	case rtype == rTypeByte:
		return int64(*(*byte)(word)), true

	case rtype == rTypeInt:
		return int64(*(*int)(word)), true
	case rtype == rTypeInt8:
		return int64(*(*int8)(word)), true
	case rtype == rTypeInt16:
		return int64(*(*int16)(word)), true
	case rtype == rTypeInt32:
		return int64(*(*int32)(word)), true
	case rtype == rTypeInt64:
		return *(*int64)(word), true

	case rtype == rTypeUint:
		return int64(*(*uint)(word)), true
	case rtype == rTypeUint8:
		return int64(*(*uint8)(word)), true
	case rtype == rTypeUint16:
		return int64(*(*uint16)(word)), true
	case rtype == rTypeUint32:
		return int64(*(*uint32)(word)), true
	case rtype == rTypeUint64:
		return int64(*(*uint64)(word)), true

	case rtype == rTypeUintptr:
		return int64(*(*uintptr)(word)), true

	case rtype == rTypeFloat32:
		return int64(*(*float32)(word)), true
	case rtype == rTypeFloat64:
		return int64(*(*float64)(word)), true

	default:
		return 0, false
	}
}

// toInt64Complex does the job, that ToInt64Fast() does not do,
// but ToInt64() should.
func toInt64Complex(v any, word unsafe.Pointer) (int64, bool) {

	var rv = reflect.ValueOf(v)

	switch rv.Type().Kind() {
	case reflect.Array, reflect.Slice, reflect.Map:
		return int64(rv.Len()), true

	case reflect.Func, reflect.Chan, reflect.Pointer, reflect.UnsafePointer:
		return ekaext.If(word != nil && !rv.IsNil(), int64(1), 0), true

	default:
		return 0, false
	}
}

// toStringUnsafeComplex is a part of ToStringUnsafe, but handles complex cases.
func toStringUnsafeComplex(rtype uintptr, word unsafe.Pointer, mask uint8) string {

	// All standard types are over.
	// Next types will be complex.

	var v = PackInterface(rtype, word)
	var kind = reflect2.TypeOf(v).Kind()

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

	return fmt.Sprintf("%+v", v)
}
