// Copyright © 2019. All rights reserved.
// Author: Ilya Yuryevich.
// Contacts: qioalice@gmail.com, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package log

import (
	"os"
)

// -----
// In the process of initiating the idea of this package, improving it,
// developing it, this package was a separate entity, not a part of LED tool
// And this package was called 'logintegro'.
//
// Now this is just echo of the past,
// but the idea of this package, its modernization has not left my head since 2017.
// For this reason, this file, which is the package entry point,
// contains constructors and other package-level functions, has this name.
// -----

// ------------------------------- ATTENTION ------------------------------- //
// This is the disclaimer of package author.
// Please, let it be as is, and be sure to read it completely,
// because it's important to understand why some code is written this way.
//
// 1. REPEATABLE CODE - AVOIDING UNNECESSARY COPYING
//
// 1.1. Functions Package, Func, Class and Methods which are constructors,
// don't use the similar methods on created Logger object because Logger
// methods will cause unnecessary copying of entry object in that case.
//
// 1.2. Functions Debug, Debugf, Debugw, Info, Infof, Infow, etc don't use
// the similar methods of baseLogger object (default Logger instance) because
// it also will cause unnecessary copying but not a entry object but
// various length arguments' slice of printf-style functions arguments or
// arguments which represents log message fields.
//
// 1.3. Functions Debugw, Infow, Warnw, Errorw, Fatalw as well as
// Logger class methods with the same name don't use With (or with)
// baseLogger object's method or With (or with) Logger method for avoiding
// unnecessary copying at least various length arguments' slice of
// log message's fields and as most - the whole Entry object.
//
// 1.4. Functions Apply, With, Group which applies something to the default
// package-level baseLogger object, also don't use Logger methods to avoiding
// unnecessary copying and calls internal parts directly.
// ------------------------------------------------------------------------- //

// baseLogger is default package-level logger, that used by all package-level
// logger functions such a Debug, Debugf, Debugw, With, Group
var baseLogger *Logger

// init performs a baseLogger initialization.
func init() {

	var encoder commonIntegratorEncoderGenerator = &ConsoleEncoder{
		format: "{{l}} {{t}}\n{{w}}\n{{m}}\n{{f}}\n{{s}}\n\n",
	}

	integrator := new(CommonIntegrator).
		WithEncoder(encoder.FreezeAndGetEncoder()).
		WithMinLevel(lvlDebug).
		WriteTo(os.Stdout)

	entry := getEntry().reset()

	baseLogger = new(Logger).setIntegrator(integrator).setEntry(entry)
}

// SyncThis forces to flush all default package logger's integrator's buffer
// and makes sure all pending log's entries are written.
func SyncThis() error {

	// baseLogger.canContinue() always returns true
	// there is no need the same check as Logger.Sync has.
	return baseLogger.Sync()
}

// Apply overwrites the behaviour of default package logger by provided reasons.
//
// This function works the same as any Logger constructor (New, Package, Func,
// Class, Method) but all these things will be made on default baseLogger.
// And it will be returned.
func ApplyThis(options ...interface{}) (defaultLogger *Logger) {

	defaultLogger = baseLogger

	if len(options) > 0 {
		defaultLogger = defaultLogger.apply(options)
	}
	return
}

// Constructor.
func New(options ...interface{}) *Logger {

	// apply is private instead of public, because there is only 2 diff:
	// 1. public's has canContinue check (baseLogger always passes which)
	// 2. public's has empty options check (but it does not matter,
	// cause we shall clone Logger anyway).
	return baseLogger.apply(options) // has derive() call
}

// With adds the fields to the default package logger's copy.
//
// You can pass both of explicit or implicit fields. Even both of named/unnamed
// implicit fields, but names (keys) should be only string.
// Neither string-like (fmt.Stringer) nor string-cast ([]byte). Only strings.
func With(fields ...interface{}) (copy *Logger) {

	if len(fields) == 0 {
		return baseLogger // avoid unnecessary copy
	}
	return baseLogger.derive(nil).entry.with(fields, nil).l
}

// WithThis is the same as With but doesn't create a copy of default package
// logger. Modifies it in-place and returns then.
func WithThis(fields ...interface{}) (defaultLogger *Logger) {

	defaultLogger = baseLogger

	if len(fields) > 0 {
		defaultLogger.entry.with(fields, nil)
	}
	return
}

// WithStrict adds an explicit fields to the default package logger's copy.
func WithStrict(fields ...Field) (copy *Logger) {

	if len(fields) == 0 {
		return baseLogger // avoid unnecessary copy
	}
	return baseLogger.derive(nil).entry.with(nil, fields).l
}

// WithStrictThis is the same as WithStrict but doesn't create a copy of default
// package logger. Modifies it in-place and returns then.
func WithStrictThis(fields ...Field) (defaultLogger *Logger) {

	defaultLogger = baseLogger

	if len(fields) > 0 {
		defaultLogger.entry.with(nil, fields)
	}
	return
}

// SkipStackFrames specified how much stack frames shall be skipped
// at the stacktrace generation. Forces stacktrace generation if it's not so.
func SkipStackFrames(n int) (copy *Logger) {

	return baseLogger.derive(nil).entry.forceStacktrace(n).l
}

// SkipStackFramesThis is the same as SkipStackFrames but doesn't create a copy
// of default package logger. Modifies it in-place and returns then.
func SkipStackFramesThis(n int) (defaultLogger *Logger) {

	defaultLogger = baseLogger
	defaultLogger.entry.forceStacktrace(n)

	return
}

// If returns package logger if 'cond' == 'true', otherwise nil. Thus it's useful
// to chaining methods - next methods in chaining will be done only if 'cond' == true.
func If(cond bool) (defaultLogger *Logger) {

	if cond {
		return baseLogger
	} else {
		return nil
	}
}

// allocLogger just allocates the memory to the new Logger object and then
// corrects the all internal pointers.
//func allocLogger() *Logger {
//	l := new(Logger)
//	// Don't allocate memory for core object, because in the default baseLogger
//	// it will be overwritten and for other objects, if core isn't specified,
//	// the baseLogger's core will be used.
//	l.entry = new(Entry)
//	l.entry.createdBy = l
//	return l
//}
