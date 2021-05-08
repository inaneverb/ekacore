// Copyright © 2020. All rights reserved.
// Author: Ilya Stroy.
// Contacts: qioalice@gmail.com, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekalog

import (
	"fmt"
	"time"

	"github.com/qioalice/ekago/v2/internal/ekaclike"
	"github.com/qioalice/ekago/v2/internal/ekaletter"
)

type (
	// Logger is used to generate and write log messages.
	//
	// You can instantiate as many your own loggers with different behaviour,
	// different Integrator, as you want.
	// But also you can just use package level logger, modify it and configure it
	// the same way as any instantiated Logger object.
	//
	// Inheritance.
	//
	// Remember!
	// No one func or method does not change the current object.
	// They always creates and returns a copy of the current object with applied
	// your changes except in cases in which the opposite is not explicitly indicated.
	//
	// Of course you can chain all setters and then call log message generator, like this:
	//
	// 		WithString("key", "value").Warn("It's dangerous!")
	//
	// But in that case, after finishing execution of second line,
	// 'log' variable won't contain added field "key".
	// Want a different behaviour and want to have Logger with these fields?
	// No problem, save generated Logger object:
	//
	// 		log := WithString("key", "value")
	// 		log.Warn("It's dangerous!")
	//
	// Because of all finishers (methods that actually writes a log message, e.g:
	// Debug, Debugf, Debugw, Warn, Warnf, Warnw, etc...) also returns a Logger
	// object that is used to generate log Entry, you can save it too, and finally
	// it's the same as in the example above:
	//
	// 		log := log.WithString("key", "value").Warn("It's dangerous!")
	//
	// but it's strongly not recommended to do so, because it made code less clear.
	Logger struct {

		// integrator is the way how Entry will be encoded and what encoded Entry
		// will be written to.
		integrator Integrator

		// entry is _WHAT_ log message is.
		// entry is it's stacktrace, caller info, timestamp, level, message, group,
		// flags, etc.
		entry *Entry
	}
)


// ------------------------------ COMMON METHODS ------------------------------ //
// ---------------------------------------------------------------------------- //


// IsValid reports whether current Logger is valid.
// Nil safe.
//
// It returns false if Logger is nil or has not been initialized properly
// (instantiated manually instead of Logger's constructors calling).
func (l *Logger) IsValid() bool {

	// Integrator and Entry (if they're not nil) can not be invalid (internally)
	// because they are created only by internal functions and they're private.
	// So 3 nil checks are enough here and of course check ptr equal.

	return l != nil && l.integrator != nil && l.entry != nil && l == l.entry.l
}

// Copy returns a copy of the current Logger. Does nothing for 'nopLogger'.
//
// Copy is useful when you need to build your Entry step-by-step,
// adding fields, messages, etc.
func (l *Logger) Copy() (copy *Logger) {
	l.assert()
	if l == nopLogger {
		return nopLogger
	}
	return l.derive()
}

// Sync forces to flush all Integrator buffers of current Logger
// and makes sure all pending Entry are written.
// Nil safe.
//
// Requirements:
//  - Logger != nil, panic otherwise;
//  - Logger is initialized properly and has registered Integrator, panic otherwise.
func (l *Logger) Sync() error {
	l.assert()
	if l == nopLogger {
		return nil
	}
	return l.integrator.Sync()
}


// --------------------------- FIELDS ADDING METHODS -------------------------- //
// ---------------------------------------------------------------------------- //

// Methods below are code-generated.

// With adds presented ekaletter.LetterField to the current Logger if it's addable.
// With DO NOT makes a copy of current Logger and adds field in-place.
// If you need to construct an logger with your own fields that will be used later,
// you need to call Copy() manually before.
func (l *Logger) With(f ekaletter.LetterField) *Logger { return l.addField(f) }

func (l *Logger) WithBool(key string, value bool) *Logger { return l.addField(ekaletter.FBool(key, value)) }
func (l *Logger) WithInt(key string, value int) *Logger { return l.addField(ekaletter.FInt(key, value)) }
func (l *Logger) WithInt8(key string, value int8) *Logger { return l.addField(ekaletter.FInt8(key, value)) }
func (l *Logger) WithInt16(key string, value int16) *Logger { return l.addField(ekaletter.FInt16(key, value)) }
func (l *Logger) WithInt32(key string, value int32) *Logger { return l.addField(ekaletter.FInt32(key, value)) }
func (l *Logger) WithInt64(key string, value int64) *Logger { return l.addField(ekaletter.FInt64(key, value)) }
func (l *Logger) WithUint(key string, value uint) *Logger { return l.addField(ekaletter.FUint(key, value)) }
func (l *Logger) WithUint8(key string, value uint8) *Logger { return l.addField(ekaletter.FUint8(key, value)) }
func (l *Logger) WithUint16(key string, value uint16) *Logger { return l.addField(ekaletter.FUint16(key, value)) }
func (l *Logger) WithUint32(key string, value uint32) *Logger { return l.addField(ekaletter.FUint32(key, value)) }
func (l *Logger) WithUint64(key string, value uint64) *Logger { return l.addField(ekaletter.FUint64(key, value)) }
func (l *Logger) WithUintptr(key string, value uintptr) *Logger { return l.addField(ekaletter.FUintptr(key, value)) }
func (l *Logger) WithFloat32(key string, value float32) *Logger { return l.addField(ekaletter.FFloat32(key, value)) }
func (l *Logger) WithFloat64(key string, value float64) *Logger { return l.addField(ekaletter.FFloat64(key, value)) }
func (l *Logger) WithComplex64(key string, value complex64) *Logger { return l.addField(ekaletter.FComplex64(key, value)) }
func (l *Logger) WithComplex128(key string, value complex128) *Logger { return l.addField(ekaletter.FComplex128(key, value)) }
func (l *Logger) WithString(key string, value string) *Logger { return l.addField(ekaletter.FString(key, value)) }
func (l *Logger) WithBoolp(key string, value *bool) *Logger { return l.addField(ekaletter.FBoolp(key, value)) }
func (l *Logger) WithIntp(key string, value *int) *Logger { return l.addField(ekaletter.FIntp(key, value)) }
func (l *Logger) WithInt8p(key string, value *int8) *Logger { return l.addField(ekaletter.FInt8p(key, value)) }
func (l *Logger) WithInt16p(key string, value *int16) *Logger { return l.addField(ekaletter.FInt16p(key, value)) }
func (l *Logger) WithInt32p(key string, value *int32) *Logger { return l.addField(ekaletter.FInt32p(key, value)) }
func (l *Logger) WithInt64p(key string, value *int64) *Logger { return l.addField(ekaletter.FInt64p(key, value)) }
func (l *Logger) WithUintp(key string, value *uint) *Logger { return l.addField(ekaletter.FUintp(key, value)) }
func (l *Logger) WithUint8p(key string, value *uint8) *Logger { return l.addField(ekaletter.FUint8p(key, value)) }
func (l *Logger) WithUint16p(key string, value *uint16) *Logger { return l.addField(ekaletter.FUint16p(key, value)) }
func (l *Logger) WithUint32p(key string, value *uint32) *Logger { return l.addField(ekaletter.FUint32p(key, value)) }
func (l *Logger) WithUint64p(key string, value *uint64) *Logger { return l.addField(ekaletter.FUint64p(key, value)) }
func (l *Logger) WithFloat32p(key string, value *float32) *Logger { return l.addField(ekaletter.FFloat32p(key, value)) }
func (l *Logger) WithFloat64p(key string, value *float64) *Logger { return l.addField(ekaletter.FFloat64p(key, value)) }
func (l *Logger) WithType(key string, value interface{}) *Logger { return l.addField(ekaletter.FType(key, value)) }
func (l *Logger) WithStringer(key string, value fmt.Stringer) *Logger { return l.addField(ekaletter.FStringer(key, value)) }
func (l *Logger) WithAddr(key string, value interface{}) *Logger { return l.addField(ekaletter.FAddr(key, value)) }
func (l *Logger) WithUnixFromStd(key string, value time.Time) *Logger { return l.addField(ekaletter.FUnixFromStd(key, value)) }
func (l *Logger) WithUnixNanoFromStd(key string, value time.Time) *Logger { return l.addField(ekaletter.FUnixNanoFromStd(key, value)) }
func (l *Logger) WithUnix(key string, value int64) *Logger { return l.addField(ekaletter.FUnix(key, value)) }
func (l *Logger) WithUnixNano(key string, value int64) *Logger { return l.addField(ekaletter.FUnixNano(key, value)) }
func (l *Logger) WithDuration(key string, value time.Duration) *Logger { return l.addField(ekaletter.FDuration(key, value)) }
func (l *Logger) WithArray(key string, value interface{}) *Logger { return l.addField(ekaletter.FArray(key, value)) }
func (l *Logger) WithObject(key string, value interface{}) *Logger { return l.addField(ekaletter.FObject(key, value)) }
func (l *Logger) WithMap(key string, value interface{}) *Logger { return l.addField(ekaletter.FMap(key, value)) }
func (l *Logger) WithExtractedMap(key string, value map[string]interface{}) *Logger { return l.addField(ekaletter.FExtractedMap(key, value)) }
func (l *Logger) WithAny(key string, value interface{}) *Logger { return l.addField(ekaletter.FAny(key, value)) }
func (l *Logger) WithMany(fields ...ekaletter.LetterField) *Logger { return l.addFields(fields) }
func (l *Logger) WithManyAny(fields ...interface{}) *Logger { return l.addFieldsParse(fields) }


// ------------------------ CONDITIONAL LOGGING METHODS ----------------------- //
// ---------------------------------------------------------------------------- //


// If returns current Logger if cond is true, otherwise nop Logger is returned.
// Thus it's useful to chaining methods - next methods in chaining will be done
// only if cond is true.
//
// Requirements:
//  - Logger is not nil, panic otherwise.
func (l *Logger) If(cond bool) *Logger {
	l.assert()
	if cond {
		return l
	} else {
		return nopLogger
	}
}


// ------------------------------ UTILITY METHODS ----------------------------- //
// ---------------------------------------------------------------------------- //


// ReplaceIntegrator replaces Integrator for the current Logger object
// to the passed one.
//
// Requirements:
//  - Logger is not nil, panic otherwise;
//  - Integrator is not nil (even typed nil), panic otherwise;
//  - If Integrator is CommonIntegrator
//    it must not be registered with some Logger before, panic otherwise;
//  - If Integrator is CommonIntegrator
//    it must have at least 1 registered io.Writer, panic otherwise.
//
// WARNING.
// Replacing Integrator will drop all pre-encoded ekaletter.LetterField fields
// that are might be added already to the current Integrator.
func (l *Logger) ReplaceIntegrator(newIntegrator Integrator) {
	l.assert()
	if l == nopLogger {
		return
	}
	if ekaclike.TakeRealAddr(newIntegrator) == nil {
		panic("Failed to change Integrator. New Integrator is nil.")
	}
	if ci, ok := newIntegrator.(*CommonIntegrator); ok {
		ci.build()
	}
	baseLogger.setIntegrator(newIntegrator)
}
