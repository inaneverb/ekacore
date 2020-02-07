// Copyright © 2019. All rights reserved.
// Author: Ilya Yuryevich.
// Contacts: qioalice@gmail.com, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package log

// Logger is used to generate and write log messages.
//
// You can instantiate as many your own loggers with different behaviour,
// different contexts, as you want. But also you can just use package level logger,
// modernize and configure it the same way as any instantiated Logger object.
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
// 		log := Package("main")
// 		log.With("key", "value").Warn("It's dangerous!")
//
// But in that case, after finishing execution of second line,
// 'log' variable won't contain add field "key".
// But package name "main" will contain.
// Want a different behaviour and want to have Logger with these fields?
// No problem, save generated Logger object:
//
// 		log := Package("main")
// 		log = log.With("key", "value")
// 		log.Warn("It's dangerous!")
//
// Because of all finishers (methods that actually writes a log message, e.g:
// Debug, Debugf, Debugw, Warn, Warnf, Warnw, etc...) also returns a Logger
// object they uses to generate log Entry, you can save it too, and finally
// it's the same as in the example above:
//
// 		log := Package("main")
// 		log = log.With("key", "value").Warn("It's dangerous!")
//
// but it's strongly not recommended to do so, because it made code less clear.
//
// There are 5 Logger constructors:
//
// 		Package(packageName, options...)
// 		Func(funcName, options...)
// 		Class(className, options...)
// 		Method(className, methodName, options...)
// 		New(options...)
//
// You can instantiate Logger object using any of constructor listed above.
// First four are used to create Logger object that binds to some Golang entity,
// and their output will contain field with 'sys.func' key and your passed value.
//
// You can combine them then to explicitly create an exactly Logger you want:
// E.g.: Method(className, methodName) == Class(className).Method(methodName).
// In that case 'sys.func' will have this value: 'className.methodName'.
//
// And the fifth creates a common regular Logger object, but it will contain
// 'sys.func' field also. Because there is auto-generation reflect information
// by default (based by stacktrace).
// You can disable this behaviour applying 'Options.EnableCallerInfo(false)'
// option to Logger's constructor or using 'Apply' method.
type Logger struct {

	// integrator is the log's entry writing destination and it's formatting way.
	// integrator determines how entry will be written.
	integrator Integrator

	// entry is what log message is.
	// entry is it's stacktrace, caller info, timestamp, level, message, group,
	// flags, etc.
	entry *Entry
}

// Sync forces to flush all l's integrator's buffer and makes sure all pending
// log's entries are written.
func (l *Logger) Sync() error {

	if !l.canContinue() {
		return errLoggerObjectInvalid
	}
	return l.integrator.Sync()
}

// Apply overwrites current Logger's copy behaviour by provided reasons.
// Read more: Options.
func (l *Logger) Apply(options ...interface{}) (copy *Logger) {

	if len(options) == 0 || !l.canContinue() {
		return l
	}
	return l.apply(options)
}

// With adds fields to the current Logger's copy.
//
// You can pass both of explicit or implicit fields. Even both of named/unnamed
// implicit fields, but names (keys) should be only string.
// Neither string-like (fmt.Stringer) nor string-cast ([]byte). Only strings.
func (l *Logger) With(fields ...interface{}) (copy *Logger) {

	if len(fields) == 0 || !l.canContinue() {
		return l
	}
	return l.derive(nil).entry.with(fields, nil).l
}

// WithStrict adds an explicit fields to the current Logger's copy.
func (l *Logger) WithStrict(fields ...Field) (copy *Logger) {

	if len(fields) == 0 || !l.canContinue() {
		return l
	}
	return l.derive(nil).entry.with(nil, fields).l
}

// SkipStackFrames specified how much stack frames shall be skipped
// at the stacktrace generation. Forces stacktrace generation if it's not so.
func (l *Logger) SkipStackFrames(n int) (copy *Logger) {

	if !l.canContinue() {
		return l
	}
	return l.derive(nil).entry.forceStacktrace(n).l
}

// If returns current logger if 'cond' == 'true', otherwise nil. Thus it's useful
// to chaining methods - next methods in chaining will be done only if 'cond' == true.
func (l *Logger) If(cond bool) *Logger {

	if cond {
		return l
	} else {
		return nil
	}
}
