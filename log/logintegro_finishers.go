// Copyright © 2019. All rights reserved.
// Author: Ilya Yuryevich.
// Contacts: qioalice@gmail.com, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package log

import (
	"fmt"

	"github.com/qioalice/gext/ec"
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

// -----
// This file contains only package baseLogger's finishers:
// Those functions that really generates log message's body and starts
// a process of writing'em.
//
// They're moved to a separately file because of bleeding my eyes,
// while I trying to work with another Logger's methods.
//
// Separated by 3 spaces into log level categories.
//
// Generally, there are common, regular functions, that are really documented
// and could be used to write log message with user defined log level.
//
// All rest methods are just "aliases" to most used (and useful) log levels.
//
// TODO: Supporting Options by finishers
// -----

// Log writes log message with desired 'level',
// analyzing 'args' in the most powerful and smart way:
//
// - args[0] could be printf-like format string, then next N args will be
//   its printf values (N - num of format's printf verbs),
//   and M-N (where M is total count of args) will be treated as
//   explicit or implicit fields (depends on what kind of fields are allowed);
//
// - args[0] could be an error => error's message will be used as log's one,
//   and wherein, error's message could be printf-like format string,
//   then next N args will be handled as in the case above,
//   and M-N (where M is total count of args) will be treated as
//   explicit or implicit fields (depends on what kind of fields are allowed);
//
// - If only explicit fields are enabled (Options.OnlyExplicitFields(true)),
//   all args except explicit fields will be form log's message using fmt.Sprint
//   (the same way as LogStrict does), but you still can pass explicit fields also;
//
// - If both of explicit and implicit fields are enabled (by default),
//   only first non-explicit field arg forms log's message using "%+v" printf verb.
//   Others ones are treated as explicit and implicit (including unnamed) fields.
//
// Not bad, huh?
func Log(level Level, args ...interface{}) *Logger {

	return baseLogger.log(level, "", nil, args, nil)
}

// LogStrict writes log message with desired 'level',
// generating log's message using fmt.Sprint(args...).
//
// No explicit/implicit fields supporting, no printf style supporting,
// no error supporting, no group supporting, no options supporting, ...
//
// ... and, as a conclusion, no reflections (usage of Golang 'reflect' package).
func LogStrict(level Level, args ...interface{}) *Logger {

	return baseLogger.log(level, fmt.Sprint(args...), nil, nil, nil)
}

// Logf writes log message with desired 'level',
// generating log's message using fmt.Sprintf(format, args[:N]...),
// where N is a number of format's printf verbs, and uses args[N:] as explicit
// or implicit fields (depends on what kind of fields are allowed).
func Logf(level Level, format string, args ...interface{}) *Logger {

	return baseLogger.log(level, format, nil, args, nil)
}

// LogfStrict is the same as just Logf but
//
// no explicit/implicit fields supporting, no group supporting,
// no options supporting, ...
//
// ... and, as a conclusion, no reflections (usage of Golang 'reflect' package).
func LogfStrict(level Level, format string, args ...interface{}) *Logger {

	return baseLogger.log(level, fmt.Sprintf(format, args...), nil, nil, nil)
}

// Logw writes log's message 'msg' with desired 'level', and passed implicit fields.
func Logw(level Level, msg string, fields ...Field) *Logger {

	return baseLogger.log(level, msg, nil, nil, fields)
}

// Loge writes log's message with desired 'level'
// using err's message as log's one but does anything only if err != nil.
//
// Uses args as printf args if err's message is printf-like format string
// and if args are left or it err's message is just a message but there are args,
// args are used as explicit or implicit fields (depends on what kind of fields are allowed).
func Loge(level Level, err error, args ...interface{}) *Logger {

	return baseLogger.checkErr(err).log(level, "", err, args, nil)
}

// LogeStrict is the same as just Loge but
//
// no implicit fields supporting, no err's printf style supporting,
// no group supporting, no options supporting, ...
//
// ... and, as a conclusion, no reflections (usage of Golang 'reflect' package).
func LogeStrict(level Level, err error, fields ...Field) *Logger {

	return baseLogger.checkErr(err).log(level, "", err, nil, fields)
}

// Logec generates ECXT based on passed 'errorCode', adds ECXT's UUID as error's one
// to the log's entry then writes log's message using err's message as log's one
// with desired 'level' but does that all only if err != nil.
// Do not pass 'ec.EOK' as 'errorCode'! It won't lead to UB but breaks logic.
func Logec(level Level, err error, errorCode ec.EC, args ...interface{}) ec.ECXT {

	return baseLogger.checkErr(err).logec(level, err, errorCode, args, nil)
}

// LogecStrict is the same as just Logec but
//
// no implicit fields supporting, no err's printf style supporting,
// no group supporting, no options supporting, ...
//
// ... and, as a conclusion, no reflections (usage of Golang 'reflect' package).
func LogecStrict(level Level, err error, errorCode ec.EC, fields ...Field) ec.ECXT {

	return baseLogger.checkErr(err).logec(level, err, errorCode, nil, fields)
}

// Debug is the same as Log(Level.Debug, args...).
// Read more: Entry.Log.
func Debug(args ...interface{}) *Logger {

	return baseLogger.log(lvlDebug, "", nil, args, nil)
}

// DebugStrict is the same as LogStrict(Level.Debug, args...).
// Read more: Entry.LogStrict.
func DebugStrict(args ...interface{}) *Logger {

	return baseLogger.log(lvlDebug, fmt.Sprint(args...), nil, nil, nil)
}

// Debugf is the same as Logf(Level.Debug, format, args...).
// Read more: Entry.Logf.
func Debugf(format string, args ...interface{}) *Logger {

	return baseLogger.log(lvlDebug, format, nil, args, nil)
}

// DebugfStrict is the same as LogfStrict(Level.Debug, format, args...).
// Read more: Entry.LogfStrict.
func DebugfStrict(format string, args ...interface{}) *Logger {

	return baseLogger.log(lvlDebug, fmt.Sprintf(format, args...), nil, nil, nil)
}

// Debugw is the same as Logw(Level.Debug, msg, fields...).
// Read more: Entry.Logw.
func Debugw(msg string, fields ...Field) *Logger {

	return baseLogger.log(lvlDebug, msg, nil, nil, fields)
}

// Debuge is the same as Loge(Level.Debug, err, args...).
// Read more: Entry.Loge.
func Debuge(err error, args ...interface{}) *Logger {

	return baseLogger.checkErr(err).log(lvlDebug, "", err, args, nil)
}

// DebugeStrict is the same as LogeStrict(Level.Debug, err, fields...).
// Read more: Entry.LogeStrict.
func DebugeStrict(err error, fields ...Field) *Logger {

	return baseLogger.checkErr(err).log(lvlDebug, "", err, nil, fields)
}

// Debugec is the same as Logec(Level.Debug, err, errorCode, args...).
// Read more: Entry.Logec.
func Debugec(err error, errorCode ec.EC, args ...interface{}) ec.ECXT {

	return baseLogger.checkErr(err).logec(lvlDebug, err, errorCode, args, nil)
}

// DebugecStrict is the same as LogecStrict(Level.Debug, err, errorCode, fields...).
// Read more: Entry.LogecStrict.
func DebugecStrict(err error, errorCode ec.EC, fields ...Field) ec.ECXT {

	return baseLogger.checkErr(err).logec(lvlDebug, err, errorCode, nil, fields)
}

// Info is the same as Log(Level.Info, args...).
// Read more: Entry.Log.
func Info(args ...interface{}) *Logger {

	return baseLogger.log(lvlInfo, "", nil, args, nil)
}

// InfoStrict is the same as LogStrict(Level.Info, args...).
// Read more: Entry.LogStrict.
func InfoStrict(args ...interface{}) *Logger {

	return baseLogger.log(lvlInfo, fmt.Sprint(args...), nil, nil, nil)
}

// Infof is the same as Logf(Level.Info, format, args...).
// Read more: Entry.Logf.
func Infof(format string, args ...interface{}) *Logger {

	return baseLogger.log(lvlInfo, format, nil, args, nil)
}

// InfofStrict is the same as LogfStrict(Level.Info, format, args...).
// Read more: Entry.LogfStrict.
func InfofStrict(format string, args ...interface{}) *Logger {

	return baseLogger.log(lvlInfo, fmt.Sprintf(format, args...), nil, nil, nil)
}

// Infow is the same as Logw(Level.Info, msg, fields...).
// Read more: Entry.Logw.
func Infow(msg string, fields ...Field) *Logger {

	return baseLogger.log(lvlInfo, msg, nil, nil, fields)
}

// Infoe is the same as Loge(Level.Info, err, args...).
// Read more: Entry.Loge.
func Infoe(err error, args ...interface{}) *Logger {

	return baseLogger.checkErr(err).log(lvlInfo, "", err, args, nil)
}

// InfoeStrict is the same as LogeStrict(Level.Info, err, fields...).
// Read more: Entry.LogeStrict.
func InfoeStrict(err error, fields ...Field) *Logger {

	return baseLogger.checkErr(err).log(lvlInfo, "", err, nil, fields)
}

// Infoec is the same as Logec(Level.Info, err, errorCode, args...).
// Read more: Entry.Logec.
func Infoec(err error, errorCode ec.EC, args ...interface{}) ec.ECXT {

	return baseLogger.checkErr(err).logec(lvlInfo, err, errorCode, args, nil)
}

// InfoecStrict is the same as LogecStrict(Level.Debug, err, errorCode, fields...).
// Read more: Entry.LogecStrict.
func InfoecStrict(err error, errorCode ec.EC, fields ...Field) ec.ECXT {

	return baseLogger.checkErr(err).logec(lvlInfo, err, errorCode, nil, fields)
}

// Warn is the same as Log(Level.Warn, args...).
// Read more: Entry.Log.
func Warn(args ...interface{}) *Logger {

	return baseLogger.log(lvlWarning, "", nil, args, nil)
}

// WarnStrict is the same as LogStrict(Level.Warn, args...).
// Read more: Entry.LogStrict.
func WarnStrict(args ...interface{}) *Logger {

	return baseLogger.log(lvlWarning, fmt.Sprint(args...), nil, nil, nil)
}

// Warnf is the same as Logf(Level.Warn, format, args...).
// Read more: Entry.Logf.
func Warnf(format string, args ...interface{}) *Logger {

	return baseLogger.log(lvlWarning, format, nil, args, nil)
}

// WarnfStrict is the same as LogfStrict(Level.Warn, format, args...).
// Read more: Entry.LogfStrict.
func WarnfStrict(format string, args ...interface{}) *Logger {

	return baseLogger.log(lvlWarning, fmt.Sprintf(format, args...), nil, nil, nil)
}

// Warnw is the same as Logw(Level.Warn, msg, fields...).
// Read more: Entry.Logw.
func Warnw(msg string, fields ...Field) *Logger {

	return baseLogger.log(lvlWarning, msg, nil, nil, fields)
}

// Warne is the same as Loge(Level.Warn, err, args...).
// Read more: Entry.Loge.
func Warne(err error, args ...interface{}) *Logger {

	return baseLogger.checkErr(err).log(lvlWarning, "", err, args, nil)
}

// WarneStrict is the same as LogeStrict(Level.Warn, err, fields...).
// Read more: Entry.LogeStrict.
func WarneStrict(err error, fields ...Field) *Logger {

	return baseLogger.checkErr(err).log(lvlWarning, "", err, nil, fields)
}

// Warnec is the same as Logec(Level.Warn, err, errorCode, args...).
// Read more: Entry.Logec.
func Warnec(err error, errorCode ec.EC, args ...interface{}) ec.ECXT {

	return baseLogger.checkErr(err).logec(lvlWarning, err, errorCode, args, nil)
}

// WarnecStrict is the same as LogecStrict(Level.Warn, err, errorCode, fields...).
// Read more: Entry.LogecStrict.
func WarnecStrict(err error, errorCode ec.EC, fields ...Field) ec.ECXT {

	return baseLogger.checkErr(err).logec(lvlWarning, err, errorCode, nil, fields)
}

// Error is the same as Log(Level.Error, args...).
// Read more: Entry.Log.
func Error(args ...interface{}) *Logger {

	return baseLogger.log(lvlError, "", nil, args, nil)
}

// ErrorStrict is the same as LogStrict(Level.Error, args...).
// Read more: Entry.LogStrict.
func ErrorStrict(args ...interface{}) *Logger {

	return baseLogger.log(lvlError, fmt.Sprint(args...), nil, nil, nil)
}

// Errorf is the same as Logf(Level.Error, format, args...).
// Read more: Entry.Logf.
func Errorf(format string, args ...interface{}) *Logger {

	return baseLogger.log(lvlError, format, nil, args, nil)
}

// ErrorfStrict is the same as LogfStrict(Level.Error, format, args...).
// Read more: Entry.LogfStrict.
func ErrorfStrict(format string, args ...interface{}) *Logger {

	return baseLogger.log(lvlError, fmt.Sprintf(format, args...), nil, nil, nil)
}

// Errorw is the same as Logw(Level.Error, msg, fields...).
// Read more: Entry.Logw.
func Errorw(msg string, fields ...Field) *Logger {

	return baseLogger.log(lvlError, msg, nil, nil, fields)
}

// Errore is the same as Loge(Level.Error, err, args...).
// Read more: Entry.Loge.
func Errore(err error, args ...interface{}) *Logger {

	return baseLogger.checkErr(err).log(lvlError, "", err, args, nil)
}

// ErroreStrict is the same as LogeStrict(Level.Error, err, fields...).
// Read more: Entry.LogeStrict.
func ErroreStrict(err error, fields ...Field) *Logger {

	return baseLogger.checkErr(err).log(lvlError, "", err, nil, fields)
}

// Errorec is the same as Logec(Level.Error, err, errorCode, args...).
// Read more: Entry.Logec.
func Errorec(err error, errorCode ec.EC, args ...interface{}) ec.ECXT {

	return baseLogger.checkErr(err).logec(lvlError, err, errorCode, args, nil)
}

// ErrorecStrict is the same as LogecStrict(Level.Error, err, errorCode, fields...).
// Read more: Entry.LogecStrict.
func ErrorecStrict(err error, errorCode ec.EC, fields ...Field) ec.ECXT {

	return baseLogger.checkErr(err).logec(lvlError, err, errorCode, nil, fields)
}

// Fatal is the same as Log(Level.Fatal, args...),
// but also then calls death.Die(1).
// Read more: Entry.Log.
func Fatal(args ...interface{}) *Logger {

	return baseLogger.log(lvlFatal, "", nil, args, nil)
}

// FatalStrict is the same as LogStrict(Level.Fatal, args...),
// but also then calls death.Die(1).
// Read more: Entry.LogStrict.
func FatalStrict(args ...interface{}) *Logger {

	return baseLogger.log(lvlFatal, fmt.Sprint(args...), nil, nil, nil)
}

// Fatalf is the same as Logf(Level.Fatal, format, args...),
// but also then calls death.Die(1).
// Read more: Entry.Logf.
func Fatalf(format string, args ...interface{}) *Logger {

	return baseLogger.log(lvlFatal, format, nil, args, nil)
}

// FatalfStrict is the same as LogfStrict(Level.Fatal, format, args...),
// but also then calls death.Die(1).
// Read more: Entry.LogfStrict.
func FatalfStrict(format string, args ...interface{}) *Logger {

	return baseLogger.log(lvlFatal, fmt.Sprintf(format, args...), nil, nil, nil)
}

// Fatalw is the same as Logw(Level.Fatal, msg, fields...),
// but also then calls death.Die(1).
// Read more: Entry.Logw.
func Fatalw(msg string, fields ...Field) *Logger {

	return baseLogger.log(lvlFatal, msg, nil, nil, fields)
}

// Fatale is the same as Loge(Level.Fatal, err, args...),
// but also then calls death.Die(1).
// Read more: Entry.Loge.
func Fatale(err error, args ...interface{}) *Logger {

	return baseLogger.checkErr(err).log(lvlFatal, "", err, args, nil)
}

// FataleStrict is the same as LogeStrict(Level.Fatal, err, fields...),
// but also then calls death.Die(1).
// Read more: Entry.LogeStrict.
func FataleStrict(err error, fields ...Field) *Logger {

	return baseLogger.checkErr(err).log(lvlFatal, "", err, nil, fields)
}

// Fatalec is the same as Logec(Level.Fatal, err, errorCode, args...).
// Read more: Entry.Logec.
func Fatalec(err error, errorCode ec.EC, args ...interface{}) ec.ECXT {

	return baseLogger.checkErr(err).logec(lvlFatal, err, errorCode, args, nil)
}

// FatalecStrict is the same as LogecStrict(Level.Fatal, err, errorCode, fields...).
// Read more: Entry.LogecStrict.
func FatalecStrict(err error, errorCode ec.EC, fields ...Field) ec.ECXT {

	return baseLogger.checkErr(err).logec(lvlFatal, err, errorCode, nil, fields)
}
