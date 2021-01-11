// Copyright © 2020. All rights reserved.
// Author: Ilya Stroy.
// Contacts: qioalice@gmail.com, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekaunsafe

import (
	"fmt"
	"time"

	"github.com/qioalice/ekago/v2/internal/ekafield"
)

// To see docs and comments,
// navigate to the origin package.

type (
	Field = ekafield.Field
	FieldKind = ekafield.Kind
)

//noinspection GoSnakeCaseUsage,GoUnusedConst
const (
	FIELD_KIND_MASK_BASE_TYPE = ekafield.KIND_MASK_BASE_TYPE
	FIELD_KIND_FLAG_ARRAY     = ekafield.KIND_FLAG_ARRAY
	FIELD_KIND_FLAG_NULL      = ekafield.KIND_FLAG_NULL
	FIELD_KIND_FLAG_SYSTEM    = ekafield.KIND_FLAG_SYSTEM
)

//noinspection GoSnakeCaseUsage,GoUnusedConst
const (
	FIELD_KIND_TYPE_INVALID = ekafield.KIND_TYPE_INVALID
)

//noinspection GoSnakeCaseUsage,GoUnusedConst
const (
	FIELD_KIND_SYS_TYPE_EKAERR_UUID           = ekafield.KIND_SYS_TYPE_EKAERR_UUID
	FIELD_KIND_SYS_TYPE_EKAERR_CLASS_ID       = ekafield.KIND_SYS_TYPE_EKAERR_CLASS_ID
	FIELD_KIND_SYS_TYPE_EKAERR_CLASS_NAME     = ekafield.KIND_SYS_TYPE_EKAERR_CLASS_NAME
	FIELD_KIND_SYS_TYPE_EKAERR_PUBLIC_MESSAGE = ekafield.KIND_SYS_TYPE_EKAERR_PUBLIC_MESSAGE
)

//noinspection GoSnakeCaseUsage,GoUnusedConst
const (
	FIELD_KIND_TYPE_BOOL        = ekafield.KIND_TYPE_BOOL
	FIELD_KIND_TYPE_INT         = ekafield.KIND_TYPE_INT
	FIELD_KIND_TYPE_INT_8       = ekafield.KIND_TYPE_INT_8
	FIELD_KIND_TYPE_INT_16      = ekafield.KIND_TYPE_INT_16
	FIELD_KIND_TYPE_INT_32      = ekafield.KIND_TYPE_INT_32
	FIELD_KIND_TYPE_INT_64      = ekafield.KIND_TYPE_INT_64
	FIELD_KIND_TYPE_UINT        = ekafield.KIND_TYPE_UINT
	FIELD_KIND_TYPE_UINT_8      = ekafield.KIND_TYPE_UINT_8
	FIELD_KIND_TYPE_UINT_16     = ekafield.KIND_TYPE_UINT_16
	FIELD_KIND_TYPE_UINT_32     = ekafield.KIND_TYPE_UINT_32
	FIELD_KIND_TYPE_UINT_64     = ekafield.KIND_TYPE_UINT_64
	FIELD_KIND_TYPE_UINTPTR     = ekafield.KIND_TYPE_UINTPTR
	FIELD_KIND_TYPE_FLOAT_32    = ekafield.KIND_TYPE_FLOAT_32
	FIELD_KIND_TYPE_FLOAT_64    = ekafield.KIND_TYPE_FLOAT_64
	FIELD_KIND_TYPE_COMPLEX_64  = ekafield.KIND_TYPE_COMPLEX_64
	FIELD_KIND_TYPE_COMPLEX_128 = ekafield.KIND_TYPE_COMPLEX_128
	FIELD_KIND_TYPE_STRING      = ekafield.KIND_TYPE_STRING
	FIELD_KIND_TYPE_ADDR        = ekafield.KIND_TYPE_ADDR
)

//noinspection GoUnusedGlobalVariable
var (
	ErrFieldUnsupportedKind = ekafield.ErrUnsupportedKind
)

func FieldReset(f *Field) {
	ekafield.Reset(f)
}

func FieldBool(key string, value bool) Field {
	return ekafield.Bool(key, value)
}

func FieldInt(key string, value int) Field {
	return ekafield.Int(key, value)
}

func FieldInt8(key string, value int8) Field {
	return ekafield.Int8(key, value)
}

func FieldInt16(key string, value int16) Field {
	return ekafield.Int16(key, value)
}

func FieldInt32(key string, value int32) Field {
	return ekafield.Int32(key, value)
}

func FieldInt64(key string, value int64) Field {
	return ekafield.Int64(key, value)
}

func FieldUint(key string, value uint) Field {
	return ekafield.Uint(key, value)
}

func FieldUint8(key string, value uint8) Field {
	return ekafield.Uint8(key, value)
}

func FieldUint16(key string, value uint16) Field {
	return ekafield.Uint16(key, value)
}

func FieldUint32(key string, value uint32) Field {
	return ekafield.Uint32(key, value)
}

func FieldUint64(key string, value uint64) Field {
	return ekafield.Uint64(key, value)
}

func FieldUintptr(key string, value uintptr) Field {
	return ekafield.Uintptr(key, value)
}

func FieldFloat32(key string, value float32) Field {
	return ekafield.Float32(key, value)
}

func FieldFloat64(key string, value float64) Field {
	return ekafield.Float64(key, value)
}

func FieldComplex64(key string, value complex64) Field {
	return ekafield.Complex64(key, value)
}

func FieldComplex128(key string, value complex128) Field {
	return ekafield.Complex128(key, value)
}

func FieldString(key string, value string) Field {
	return ekafield.String(key, value)
}

func FieldBoolp(key string, value *bool) Field {
	return ekafield.Boolp(key, value)
}

func FieldIntp(key string, value *int) Field {
	return ekafield.Intp(key, value)
}

func FieldInt8p(key string, value *int8) Field {
	return ekafield.Int8p(key, value)
}

func FieldInt16p(key string, value *int16) Field {
	return ekafield.Int16p(key, value)
}

func FieldInt32p(key string, value *int32) Field {
	return ekafield.Int32p(key, value)
}

func FieldInt64p(key string, value *int64) Field {
	return ekafield.Int64p(key, value)
}

func FieldUintp(key string, value *uint) Field {
	return ekafield.Uintp(key, value)
}

func FieldUint8p(key string, value *uint8) Field {
	return ekafield.Uint8p(key, value)
}

func FieldUint16p(key string, value *uint16) Field {
	return ekafield.Uint16p(key, value)
}

func FieldUint32p(key string, value *uint32) Field {
	return ekafield.Uint32p(key, value)
}

func FieldUint64p(key string, value *uint64) Field {
	return ekafield.Uint64p(key, value)
}

func FieldFloat32p(key string, value *float32) Field {
	return ekafield.Float32p(key, value)
}

func FieldFloat64p(key string, value *float64) Field {
	return ekafield.Float64p(key, value)
}

func FieldType(key string, value interface{}) Field {
	return ekafield.Type(key, value)
}

func FieldStringer(key string, value fmt.Stringer) Field {
	return ekafield.Stringer(key, value)
}

func FieldAddr(key string, value interface{}) Field {
	return ekafield.Addr(key, value)
}

func FieldTime(key string, value time.Time) Field {
	return ekafield.Time(key, value)
}

func FieldDuration(key string, value time.Duration) Field {
	return ekafield.Duration(key, value)
}

func FieldNilValue(key string, baseType FieldKind) Field {
	return ekafield.NilValue(key, baseType)
}
