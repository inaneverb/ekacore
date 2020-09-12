// Copyright © 2020. All rights reserved.
// Author: Ilya Stroy.
// Contacts: qioalice@gmail.com, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekaunsafe

import (
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

	FIELD_KIND_TYPE_INVALID = ekafield.KIND_TYPE_INVALID

	FIELD_KIND_SYS_TYPE_EKAERR_UUID           = ekafield.KIND_SYS_TYPE_EKAERR_UUID
	FIELD_KIND_SYS_TYPE_EKAERR_CLASS_ID       = ekafield.KIND_SYS_TYPE_EKAERR_CLASS_ID
	FIELD_KIND_SYS_TYPE_EKAERR_CLASS_NAME     = ekafield.KIND_SYS_TYPE_EKAERR_CLASS_NAME
	FIELD_KIND_SYS_TYPE_EKAERR_PUBLIC_MESSAGE = ekafield.KIND_SYS_TYPE_EKAERR_PUBLIC_MESSAGE

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

//noinspection GoUnusedGlobalVariable
var (
	FieldReset = ekafield.Reset
	FieldBool = ekafield.Bool
	FieldInt = ekafield.Int
	FieldInt8 = ekafield.Int8
	FieldInt16 = ekafield.Int16
	FieldInt32 = ekafield.Int32
	FieldInt64 = ekafield.Int64
	FieldUint = ekafield.Uint
	FieldUint8 = ekafield.Uint8
	FieldUint16 = ekafield.Uint16
	FieldUint32 = ekafield.Uint32
	FieldUint64 = ekafield.Uint64
	FieldUintptr = ekafield.Uintptr
	FieldFloat32 = ekafield.Float32
	FieldFloat64 = ekafield.Float64
	FieldComplex64 = ekafield.Complex64
	FieldComplex128 = ekafield.Complex128
	FieldString = ekafield.String
	FieldBoolp = ekafield.Boolp
	FieldIntp = ekafield.Intp
	FieldInt8p = ekafield.Int8p
	FieldInt16p = ekafield.Int16p
	FieldInt32p = ekafield.Int32p
	FieldInt64p = ekafield.Int64p
	FieldUintp = ekafield.Uintp
	FieldUint8p = ekafield.Uint8p
	FieldUint16p = ekafield.Uint16p
	FieldUint32p = ekafield.Uint32p
	FieldUint64p = ekafield.Uint64p
	FieldFloat32p = ekafield.Float32p
	FieldFloat64p = ekafield.Float64p
	FieldType = ekafield.Type
	FieldStringer = ekafield.Stringer
	FieldTime = ekafield.Time
	FieldDuration = ekafield.Duration
	FieldNilValue = ekafield.NilValue
)

