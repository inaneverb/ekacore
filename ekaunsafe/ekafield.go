// Copyright © 2020. All rights reserved.
// Author: Ilya Stroy.
// Contacts: iyuryevich@pm.me, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekaunsafe

import (
	"fmt"
	"time"

	"github.com/qioalice/ekago/v3/internal/ekaletter"
)

// To see docs and comments,
// navigate to the origin package.

type (
	LetterField     = ekaletter.LetterField
	LetterFieldKind = ekaletter.LetterFieldKind
)

// noinspection GoSnakeCaseUsage,GoUnusedConst
const (
	FIELD_KIND_MASK_BASE_TYPE    = ekaletter.KIND_MASK_BASE_TYPE
	FIELD_KIND_FLAG_USER_DEFINED = ekaletter.KIND_FLAG_USER_DEFINED
	FIELD_KIND_FLAG_NULL         = ekaletter.KIND_FLAG_NULL
	FIELD_KIND_FLAG_SYSTEM       = ekaletter.KIND_FLAG_SYSTEM
)

// noinspection GoSnakeCaseUsage,GoUnusedConst
const (
	FIELD_KIND_TYPE_INVALID = ekaletter.KIND_TYPE_INVALID
)

// noinspection GoSnakeCaseUsage,GoUnusedConst
const (
	FIELD_KIND_SYS_TYPE_EKAERR_UUID       = ekaletter.KIND_SYS_TYPE_EKAERR_UUID
	FIELD_KIND_SYS_TYPE_EKAERR_CLASS_ID   = ekaletter.KIND_SYS_TYPE_EKAERR_CLASS_ID
	FIELD_KIND_SYS_TYPE_EKAERR_CLASS_NAME = ekaletter.KIND_SYS_TYPE_EKAERR_CLASS_NAME
)

// noinspection GoSnakeCaseUsage,GoUnusedConst
const (
	FIELD_KIND_TYPE_BOOL        = ekaletter.KIND_TYPE_BOOL
	FIELD_KIND_TYPE_INT         = ekaletter.KIND_TYPE_INT
	FIELD_KIND_TYPE_INT_8       = ekaletter.KIND_TYPE_INT_8
	FIELD_KIND_TYPE_INT_16      = ekaletter.KIND_TYPE_INT_16
	FIELD_KIND_TYPE_INT_32      = ekaletter.KIND_TYPE_INT_32
	FIELD_KIND_TYPE_INT_64      = ekaletter.KIND_TYPE_INT_64
	FIELD_KIND_TYPE_UINT        = ekaletter.KIND_TYPE_UINT
	FIELD_KIND_TYPE_UINT_8      = ekaletter.KIND_TYPE_UINT_8
	FIELD_KIND_TYPE_UINT_16     = ekaletter.KIND_TYPE_UINT_16
	FIELD_KIND_TYPE_UINT_32     = ekaletter.KIND_TYPE_UINT_32
	FIELD_KIND_TYPE_UINT_64     = ekaletter.KIND_TYPE_UINT_64
	FIELD_KIND_TYPE_UINTPTR     = ekaletter.KIND_TYPE_UINTPTR
	FIELD_KIND_TYPE_FLOAT_32    = ekaletter.KIND_TYPE_FLOAT_32
	FIELD_KIND_TYPE_FLOAT_64    = ekaletter.KIND_TYPE_FLOAT_64
	FIELD_KIND_TYPE_COMPLEX_64  = ekaletter.KIND_TYPE_COMPLEX_64
	FIELD_KIND_TYPE_COMPLEX_128 = ekaletter.KIND_TYPE_COMPLEX_128
	FIELD_KIND_TYPE_STRING      = ekaletter.KIND_TYPE_STRING
	FIELD_KIND_TYPE_ADDR        = ekaletter.KIND_TYPE_ADDR
	FIELD_KIND_TYPE_UNIX        = ekaletter.KIND_TYPE_UNIX
	FIELD_KIND_TYPE_UNIX_NANO   = ekaletter.KIND_TYPE_UNIX_NANO
	FIELD_KIND_TYPE_DURATION    = ekaletter.KIND_TYPE_DURATION
	FIELD_KIND_TYPE_ARRAY       = ekaletter.KIND_TYPE_ARRAY
	FIELD_KIND_TYPE_MAP         = ekaletter.KIND_TYPE_MAP
	FIELD_KIND_TYPE_EXTMAP      = ekaletter.KIND_TYPE_EXTMAP
	FIELD_KIND_TYPE_STRUCT      = ekaletter.KIND_TYPE_STRUCT
)

// noinspection GoUnusedGlobalVariable
var (
	ErrFieldUnsupportedKind = ekaletter.ErrUnsupportedKind
)

func FieldReset(f *LetterField) {
	ekaletter.FieldReset(f)
}

func FBool(key string, value bool) LetterField              { return ekaletter.FBool(key, value) }
func FInt(key string, value int) LetterField                { return ekaletter.FInt(key, value) }
func FInt8(key string, value int8) LetterField              { return ekaletter.FInt8(key, value) }
func FInt16(key string, value int16) LetterField            { return ekaletter.FInt16(key, value) }
func FInt32(key string, value int32) LetterField            { return ekaletter.FInt32(key, value) }
func FInt64(key string, value int64) LetterField            { return ekaletter.FInt64(key, value) }
func FUint(key string, value uint) LetterField              { return ekaletter.FUint(key, value) }
func FUint8(key string, value uint8) LetterField            { return ekaletter.FUint8(key, value) }
func FUint16(key string, value uint16) LetterField          { return ekaletter.FUint16(key, value) }
func FUint32(key string, value uint32) LetterField          { return ekaletter.FUint32(key, value) }
func FUint64(key string, value uint64) LetterField          { return ekaletter.FUint64(key, value) }
func FUintptr(key string, value uintptr) LetterField        { return ekaletter.FUintptr(key, value) }
func FFloat32(key string, value float32) LetterField        { return ekaletter.FFloat32(key, value) }
func FFloat64(key string, value float64) LetterField        { return ekaletter.FFloat64(key, value) }
func FComplex64(key string, value complex64) LetterField    { return ekaletter.FComplex64(key, value) }
func FComplex128(key string, value complex128) LetterField  { return ekaletter.FComplex128(key, value) }
func FString(key string, value string) LetterField          { return ekaletter.FString(key, value) }
func FBoolp(key string, value *bool) LetterField            { return ekaletter.FBoolp(key, value) }
func FIntp(key string, value *int) LetterField              { return ekaletter.FIntp(key, value) }
func FInt8p(key string, value *int8) LetterField            { return ekaletter.FInt8p(key, value) }
func FInt16p(key string, value *int16) LetterField          { return ekaletter.FInt16p(key, value) }
func FInt32p(key string, value *int32) LetterField          { return ekaletter.FInt32p(key, value) }
func FInt64p(key string, value *int64) LetterField          { return ekaletter.FInt64p(key, value) }
func FUintp(key string, value *uint) LetterField            { return ekaletter.FUintp(key, value) }
func FUint8p(key string, value *uint8) LetterField          { return ekaletter.FUint8p(key, value) }
func FUint16p(key string, value *uint16) LetterField        { return ekaletter.FUint16p(key, value) }
func FUint32p(key string, value *uint32) LetterField        { return ekaletter.FUint32p(key, value) }
func FUint64p(key string, value *uint64) LetterField        { return ekaletter.FUint64p(key, value) }
func FFloat32p(key string, value *float32) LetterField      { return ekaletter.FFloat32p(key, value) }
func FFloat64p(key string, value *float64) LetterField      { return ekaletter.FFloat64p(key, value) }
func FType(key string, value any) LetterField               { return ekaletter.FType(key, value) }
func FStringer(key string, value fmt.Stringer) LetterField  { return ekaletter.FStringer(key, value) }
func FAddr(key string, value any) LetterField               { return ekaletter.FAddr(key, value) }
func FUnixFromStd(key string, t time.Time) LetterField      { return ekaletter.FUnixFromStd(key, t) }
func FUnixNanoFromStd(key string, t time.Time) LetterField  { return ekaletter.FUnixNanoFromStd(key, t) }
func FUnix(key string, unix int64) LetterField              { return ekaletter.FUnix(key, unix) }
func FUnixNano(key string, unixNano int64) LetterField      { return ekaletter.FUnixNano(key, unixNano) }
func FDuration(key string, value time.Duration) LetterField { return ekaletter.FDuration(key, value) }
func FArray(key string, value any) LetterField              { return ekaletter.FArray(key, value) }
func FObject(key string, value any) LetterField             { return ekaletter.FObject(key, value) }
func FMap(key string, value any) LetterField                { return ekaletter.FMap(key, value) }
func FExtractedMap(key string, value map[string]any) LetterField {
	return ekaletter.FExtractedMap(key, value)
}
func FAny(key string, value any) LetterField                { return ekaletter.FAny(key, value) }
func FNil(key string, baseType LetterFieldKind) LetterField { return ekaletter.FNil(key, baseType) }
func FInvalid(key string) LetterField                       { return ekaletter.FInvalid(key) }
