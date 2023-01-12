// Copyright © 2017-2021. All rights reserved.
// Author: Ilya Stroy.
// Contacts: iyuryevich@pm.me, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekalog

import (
	"fmt"
	"io"
	"time"

	"github.com/qioalice/ekago/v4/internal/ekaletter"
)

// ------------------------------ COMMON METHODS ------------------------------ //
// ---------------------------------------------------------------------------- //

// Copy returns a copy of the package-level Logger. Does nothing for 'nopLogger'.
//
// Copy is useful when you need to build your Entry step-by-step,
// adding fields, messages, etc.
func Copy() (copy *Logger) {
	return baseLogger.derive()
}

// Sync forces to flush all Integrator buffers of package Logger
// and makes sure all pending Entry are written.
func Sync() error {
	return baseLogger.Sync()
}

// --------------------------- FIELDS ADDING METHODS -------------------------- //
// ---------------------------------------------------------------------------- //

// Methods below are code-generated.

func With(f ekaletter.LetterField) *Logger {
	return baseLogger.addField(f)
}
func WithBool(key string, value bool) *Logger {
	return baseLogger.addField(ekaletter.FBool(key, value))
}
func WithInt(key string, value int) *Logger {
	return baseLogger.addField(ekaletter.FInt(key, value))
}
func WithInt8(key string, value int8) *Logger {
	return baseLogger.addField(ekaletter.FInt8(key, value))
}
func WithInt16(key string, value int16) *Logger {
	return baseLogger.addField(ekaletter.FInt16(key, value))
}
func WithInt32(key string, value int32) *Logger {
	return baseLogger.addField(ekaletter.FInt32(key, value))
}
func WithInt64(key string, value int64) *Logger {
	return baseLogger.addField(ekaletter.FInt64(key, value))
}
func WithUint(key string, value uint) *Logger {
	return baseLogger.addField(ekaletter.FUint(key, value))
}
func WithUint8(key string, value uint8) *Logger {
	return baseLogger.addField(ekaletter.FUint8(key, value))
}
func WithUint16(key string, value uint16) *Logger {
	return baseLogger.addField(ekaletter.FUint16(key, value))
}
func WithUint32(key string, value uint32) *Logger {
	return baseLogger.addField(ekaletter.FUint32(key, value))
}
func WithUint64(key string, value uint64) *Logger {
	return baseLogger.addField(ekaletter.FUint64(key, value))
}
func WithUintptr(key string, value uintptr) *Logger {
	return baseLogger.addField(ekaletter.FUintptr(key, value))
}
func WithFloat32(key string, value float32) *Logger {
	return baseLogger.addField(ekaletter.FFloat32(key, value))
}
func WithFloat64(key string, value float64) *Logger {
	return baseLogger.addField(ekaletter.FFloat64(key, value))
}
func WithComplex64(key string, value complex64) *Logger {
	return baseLogger.addField(ekaletter.FComplex64(key, value))
}
func WithComplex128(key string, value complex128) *Logger {
	return baseLogger.addField(ekaletter.FComplex128(key, value))
}
func WithString(key string, value string) *Logger {
	return baseLogger.addField(ekaletter.FString(key, value))
}
func WithStringFromBytes(key string, value []byte) *Logger {
	return baseLogger.addField(ekaletter.FStringFromBytes(key, value))
}
func WithBoolp(key string, value *bool) *Logger {
	return baseLogger.addField(ekaletter.FBoolp(key, value))
}
func WithIntp(key string, value *int) *Logger {
	return baseLogger.addField(ekaletter.FIntp(key, value))
}
func WithInt8p(key string, value *int8) *Logger {
	return baseLogger.addField(ekaletter.FInt8p(key, value))
}
func WithInt16p(key string, value *int16) *Logger {
	return baseLogger.addField(ekaletter.FInt16p(key, value))
}
func WithInt32p(key string, value *int32) *Logger {
	return baseLogger.addField(ekaletter.FInt32p(key, value))
}
func WithInt64p(key string, value *int64) *Logger {
	return baseLogger.addField(ekaletter.FInt64p(key, value))
}
func WithUintp(key string, value *uint) *Logger {
	return baseLogger.addField(ekaletter.FUintp(key, value))
}
func WithUint8p(key string, value *uint8) *Logger {
	return baseLogger.addField(ekaletter.FUint8p(key, value))
}
func WithUint16p(key string, value *uint16) *Logger {
	return baseLogger.addField(ekaletter.FUint16p(key, value))
}
func WithUint32p(key string, value *uint32) *Logger {
	return baseLogger.addField(ekaletter.FUint32p(key, value))
}
func WithUint64p(key string, value *uint64) *Logger {
	return baseLogger.addField(ekaletter.FUint64p(key, value))
}
func WithFloat32p(key string, value *float32) *Logger {
	return baseLogger.addField(ekaletter.FFloat32p(key, value))
}
func WithFloat64p(key string, value *float64) *Logger {
	return baseLogger.addField(ekaletter.FFloat64p(key, value))
}
func WithType(key string, value any) *Logger {
	return baseLogger.addField(ekaletter.FType(key, value))
}
func WithStringer(key string, value fmt.Stringer) *Logger {
	return baseLogger.addField(ekaletter.FStringer(key, value))
}
func WithAddr(key string, value any) *Logger {
	return baseLogger.addField(ekaletter.FAddr(key, value))
}
func WithUnixFromStd(key string, value time.Time) *Logger {
	return baseLogger.addField(ekaletter.FUnixFromStd(key, value))
}
func WithUnixNanoFromStd(key string, value time.Time) *Logger {
	return baseLogger.addField(ekaletter.FUnixNanoFromStd(key, value))
}
func WithUnix(key string, value int64) *Logger {
	return baseLogger.addField(ekaletter.FUnix(key, value))
}
func WithUnixNano(key string, value int64) *Logger {
	return baseLogger.addField(ekaletter.FUnixNano(key, value))
}
func WithDuration(key string, value time.Duration) *Logger {
	return baseLogger.addField(ekaletter.FDuration(key, value))
}
func WithArray(key string, value any) *Logger {
	return baseLogger.addField(ekaletter.FArray(key, value))
}
func WithObject(key string, value any) *Logger {
	return baseLogger.addField(ekaletter.FObject(key, value))
}
func WithMap(key string, value any) *Logger {
	return baseLogger.addField(ekaletter.FMap(key, value))
}
func WithExtractedMap(key string, value map[string]any) *Logger {
	return baseLogger.addField(ekaletter.FExtractedMap(key, value))
}
func WithAny(key string, value any) *Logger {
	return baseLogger.addField(ekaletter.FAny(key, value))
}
func WithMany(fields ...ekaletter.LetterField) *Logger {
	return baseLogger.addFields(fields)
}
func WithManyAny(fields ...any) *Logger {
	return baseLogger.addFieldsParse(fields)
}

// ------------------------ CONDITIONAL LOGGING METHODS ----------------------- //
// ---------------------------------------------------------------------------- //

// If returns package Logger if cond is true, otherwise nop Logger is returned.
// Thus it's useful to chaining methods - next methods in chaining will be done
// only if cond is true.
func If(cond bool) (defaultLogger *Logger) {
	return baseLogger.If(cond)
}

// ------------------------------ UTILITY METHODS ----------------------------- //
// ---------------------------------------------------------------------------- //

// ReplaceIntegrator replaces Integrator for the package Logger to the passed one.
//
// Requirements:
//   - Integrator is not nil (even typed nil), panic otherwise;
//   - If Integrator is CommonIntegrator
//     it must not be registered with some Logger before, panic otherwise;
//   - If Integrator is CommonIntegrator
//     it must have at least 1 registered io.Writer, panic otherwise.
//
// WARNING.
// Replacing Integrator will drop all pre-encoded ekaletter.LetterField fields
// that are might be added already to the current Integrator.
func ReplaceIntegrator(newIntegrator Integrator) {
	baseLogger.ReplaceIntegrator(newIntegrator)
}

// ReplaceEncoder is an alias for creating a new CommonIntegrator,
// setting provided CI_Encoder for them and register it as a new integrator.
// The synced stdout is used as writer if writer's set is empty.
func ReplaceEncoder(cie CI_Encoder, writers ...io.Writer) {
	ci := new(CommonIntegrator).WithEncoder(cie).WriteTo(writers...)
	baseLogger.ReplaceIntegrator(ci)
}
