// Copyright © 2020. All rights reserved.
// Author: Ilya Stroy.
// Contacts: qioalice@gmail.com, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekaexp

import (
	"github.com/qioalice/ekago/v2/internal/field"
)

// To see docs and comments,
// navigate to the origin package.

type (
	Field = field.Field
	FieldKind = field.Kind
)

//noinspection GoSnakeCaseUsage,GoUnusedConst
const (
	FIELD_KIND_MASK_BASE_TYPE = field.KIND_MASK_BASE_TYPE
	FIELD_KIND_FLAG_ARRAY     = field.KIND_FLAG_ARRAY
	FIELD_KIND_FLAG_NULL      = field.KIND_FLAG_NULL
	FIELD_KIND_FLAG_SYSTEM    = field.KIND_FLAG_SYSTEM

	FIELD_KIND_TYPE_INVALID = field.KIND_TYPE_INVALID

	FIELD_KIND_SYS_TYPE_EKAERR_UUID           = field.KIND_SYS_TYPE_EKAERR_UUID
	FIELD_KIND_SYS_TYPE_EKAERR_CLASS_ID       = field.KIND_SYS_TYPE_EKAERR_CLASS_ID
	FIELD_KIND_SYS_TYPE_EKAERR_CLASS_NAME     = field.KIND_SYS_TYPE_EKAERR_CLASS_NAME
	FIELD_KIND_SYS_TYPE_EKAERR_PUBLIC_MESSAGE = field.KIND_SYS_TYPE_EKAERR_PUBLIC_MESSAGE

	FIELD_KIND_TYPE_BOOL        = field.KIND_TYPE_BOOL
	FIELD_KIND_TYPE_INT         = field.KIND_TYPE_INT
	FIELD_KIND_TYPE_INT_8       = field.KIND_TYPE_INT_8
	FIELD_KIND_TYPE_INT_16      = field.KIND_TYPE_INT_16
	FIELD_KIND_TYPE_INT_32      = field.KIND_TYPE_INT_32
	FIELD_KIND_TYPE_INT_64      = field.KIND_TYPE_INT_64
	FIELD_KIND_TYPE_UINT        = field.KIND_TYPE_UINT
	FIELD_KIND_TYPE_UINT_8      = field.KIND_TYPE_UINT_8
	FIELD_KIND_TYPE_UINT_16     = field.KIND_TYPE_UINT_16
	FIELD_KIND_TYPE_UINT_32     = field.KIND_TYPE_UINT_32
	FIELD_KIND_TYPE_UINT_64     = field.KIND_TYPE_UINT_64
	FIELD_KIND_TYPE_UINTPTR     = field.KIND_TYPE_UINTPTR
	FIELD_KIND_TYPE_FLOAT_32    = field.KIND_TYPE_FLOAT_32
	FIELD_KIND_TYPE_FLOAT_64    = field.KIND_TYPE_FLOAT_64
	FIELD_KIND_TYPE_COMPLEX_64  = field.KIND_TYPE_COMPLEX_64
	FIELD_KIND_TYPE_COMPLEX_128 = field.KIND_TYPE_COMPLEX_128
	FIELD_KIND_TYPE_STRING      = field.KIND_TYPE_STRING
	FIELD_KIND_TYPE_ADDR        = field.KIND_TYPE_ADDR
)

//noinspection GoUnusedGlobalVariable
var (
	FieldReset = field.Reset
	FieldBool = field.Bool
	FieldInt = field.Int
	FieldInt8 = field.Int8
	FieldInt16 = field.Int16
	FieldInt32 = field.Int32
	FieldInt64 = field.Int64
	FieldUint = field.Uint
	FieldUint8 = field.Uint8
	FieldUint16 = field.Uint16
	FieldUint32 = field.Uint32
	FieldUint64 = field.Uint64
	FieldUintptr = field.Uintptr
	FieldFloat32 = field.Float32
	FieldFloat64 = field.Float64
	FieldComplex64 = field.Complex64
	FieldComplex128 = field.Complex128
	FieldString = field.String
	FieldBoolp = field.Boolp
	FieldIntp = field.Intp
	FieldInt8p = field.Int8p
	FieldInt16p = field.Int16p
	FieldInt32p = field.Int32p
	FieldInt64p = field.Int64p
	FieldUintp = field.Uintp
	FieldUint8p = field.Uint8p
	FieldUint16p = field.Uint16p
	FieldUint32p = field.Uint32p
	FieldUint64p = field.Uint64p
	FieldFloat32p = field.Float32p
	FieldFloat64p = field.Float64p
	FieldType = field.Type
	FieldStringer = field.Stringer
	FieldTime = field.Time
	FieldDuration = field.Duration
	FieldNilValue = field.NilValue
)

