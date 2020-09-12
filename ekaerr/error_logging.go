// Copyright © 2020. All rights reserved.
// Author: Ilya Stroy.
// Contacts: qioalice@gmail.com, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekaerr

import (
	"fmt"
	"unsafe"

	"github.com/qioalice/ekago/v2/ekalog"
	"github.com/qioalice/ekago/v2/internal/ekafield"
	"github.com/qioalice/ekago/v2/internal/ekaletter"
)

// ---------------------------------------------------------------------------- //
//                               IMPORTANT NOTICE                               //
//                                                                              //
// I specially DO NOT allow you (do not provide methods) log errors fast way    //
// with "Debug" or "Info" level, because it's an error!                         //
// Error couldn't be debug or info. Error is error.                             //
//                                                                              //
// You even CAN NOT avoid this limitation by using LogAs() method,              //
// because a minimum log level you can use with this method is "Warning".       //
// Thus, you can use this method to log error with your own log level but only  //
// if your log level has value > "Warning" level's one (70).                    //
//                                                                              //
// But it's an error and we can't just ignore it if you trying to log it with   //
// "Debug" level, for example. Then we'll log it but with overwritten "Error"   //
// log level.                                                                   //
//                                                                              //
// ---------------------------------------------------------------------------- //

// LogAs logs current Error if it's valid using standard ekalog package's logger
// with the provided log level 'level'.
// YOU CAN NOT USE ERROR OBJECT AFTER THAT CALL. IT WILL BE BROKEN!
//
// Requirements:
// 'level' >= ekalog.LEVEL_WARNING. Otherwise overwritten by ekalog.LEVER_ERROR.


//
func (e *Error) Log(level ekalog.Level, args ...interface{}) {
	e.log(nil, level, args)
}

//
func (e *Error) Logf(level ekalog.Level, format string, args ...interface{}) {
	e.logf(nil, level, format, args)
}

//
func (e *Error) Logw(level ekalog.Level, message string, fields ...ekafield.Field) {
	e.logw(nil, level, message, fields)
}

//
func (e *Error) Logww(level ekalog.Level, message string, fields []ekafield.Field) {
	e.logw(nil, level, message, fields)
}

//
func (e *Error) LogAsWarn(args ...interface{}) {
	e.log(nil, ekalog.LEVEL_WARNING, args)
}

//
func (e *Error) LogAsWarnf(format string, args ...interface{}) {
	e.logf(nil, ekalog.LEVEL_WARNING, format, args)
}

//
func (e *Error) LogAsWarnw(message string, fields ...ekafield.Field) {
	e.logw(nil, ekalog.LEVEL_WARNING, message, fields)
}

//
func (e *Error) LogAsWarnww(message string, fields []ekafield.Field) {
	e.logw(nil, ekalog.LEVEL_WARNING, message, fields)
}

//
func (e *Error) LogAsError(args ...interface{}) {
	e.log(nil, ekalog.LEVEL_ERROR, args)
}

//
func (e *Error) LogAsErrorf(format string, args ...interface{}) {
	e.logf(nil, ekalog.LEVEL_ERROR, format, args)
}

//
func (e *Error) LogAsErrorw(message string, fields ...ekafield.Field) {
	e.logw(nil, ekalog.LEVEL_ERROR, message, fields)
}

//
func (e *Error) LogAsErrorww(message string, fields []ekafield.Field) {
	e.logw(nil, ekalog.LEVEL_ERROR, message, fields)
}

//
func (e *Error) LogAsFatal(args ...interface{}) {
	e.log(nil, ekalog.LEVEL_FATAL, args)
}

//
func (e *Error) LogAsFatalf(format string, args ...interface{}) {
	e.logf(nil, ekalog.LEVEL_FATAL, format, args)
}

//
func (e *Error) LogAsFatalw(message string, fields ...ekafield.Field) {
	e.logw(nil, ekalog.LEVEL_FATAL, message, fields)
}

//
func (e *Error) LogAsFatalww(message string, fields []ekafield.Field) {
	e.logw(nil, ekalog.LEVEL_FATAL, message, fields)
}

//
func (e *Error) LogUsing(logger *ekalog.Logger, level ekalog.Level, args ...interface{}) {
	e.log(logger, level, args)
}

//
func (e *Error) LogfUsing(logger *ekalog.Logger, level ekalog.Level, format string, args ...interface{}) {
	e.logf(logger, level, format, args)
}

//
func (e *Error) LogwUsing(logger *ekalog.Logger, level ekalog.Level, message string, fields ...ekafield.Field) {
	e.logw(logger, level, message, fields)
}

//
func (e *Error) LogwwUsing(logger *ekalog.Logger, level ekalog.Level, message string, fields []ekafield.Field) {
	e.logw(logger, level, message, fields)
}

//
func (e *Error) LogAsWarnUsing(logger *ekalog.Logger, args ...interface{}) {
	e.log(logger, ekalog.LEVEL_WARNING, args)
}

//
func (e *Error) LogAsWarnfUsing(logger *ekalog.Logger, format string, args ...interface{}) {
	e.logf(logger, ekalog.LEVEL_WARNING, format, args)
}

//
func (e *Error) LogAsWarnwUsing(logger *ekalog.Logger, message string, fields ...ekafield.Field) {
	e.logw(logger, ekalog.LEVEL_WARNING, message, fields)
}

//
func (e *Error) LogAsWarnwwUsing(logger *ekalog.Logger, message string, fields []ekafield.Field) {
	e.logw(logger, ekalog.LEVEL_WARNING, message, fields)
}

//
func (e *Error) LogAsErrorUsing(logger *ekalog.Logger, args ...interface{}) {
	e.log(logger, ekalog.LEVEL_ERROR, args)
}

//
func (e *Error) LogAsErrorfUsing(logger *ekalog.Logger, format string, args ...interface{}) {
	e.logf(logger, ekalog.LEVEL_ERROR, format, args)
}

//
func (e *Error) LogAsErrorwUsing(logger *ekalog.Logger, message string, fields ...ekafield.Field) {
	e.logw(logger, ekalog.LEVEL_ERROR, message, fields)
}

//
func (e *Error) LogAsErrorwwUsing(logger *ekalog.Logger, message string, fields []ekafield.Field) {
	e.logw(logger, ekalog.LEVEL_ERROR, message, fields)
}

//
func (e *Error) LogAsFatalUsing(logger *ekalog.Logger, args ...interface{}) {
	e.log(logger, ekalog.LEVEL_FATAL, args)
}

//
func (e *Error) LogAsFatalfUsing(logger *ekalog.Logger, format string, args ...interface{}) {
	e.logf(logger, ekalog.LEVEL_FATAL, format, args)
}

//
func (e *Error) LogAsFatalwUsing(logger *ekalog.Logger, message string, fields ...ekafield.Field) {
	e.logw(logger, ekalog.LEVEL_FATAL, message, fields)
}

//
func (e *Error) LogAsFatalwwUsing(logger *ekalog.Logger, message string, fields []ekafield.Field) {
	e.logw(logger, ekalog.LEVEL_FATAL, message, fields)
}

// LogAsThrough logs current Error if it's valid using provided 'logger' as Logger
// through which Error will be logged (only if it's valid too) with the log level 'level'.
// YOU CAN NOT USE ERROR OBJECT AFTER THAT CALL. IT WILL BE BROKEN!
//
// Requirements:
// 'level' >= ekalog.LEVEL_WARNING. Otherwise overwritten by ekalog.LEVER_ERROR.


// LogAsWarning logs current Error if it's valid using standard ekalog package's logger
// with 'LEVEL_WARNING' log level.
// YOU CAN NOT USE ERROR OBJECT AFTER THAT CALL. IT WILL BE BROKEN!

// LogAsError logs current Error if it's valid using standard ekalog package's logger
// with 'LEVEL_ERROR' log level.
// YOU CAN NOT USE ERROR OBJECT AFTER THAT CALL. IT WILL BE BROKEN!


// LogAsFatal logs current Error if it's valid using standard ekalog package's logger
// with 'LEVEL_FATAL' log level.
// YOU CAN NOT USE ERROR OBJECT AFTER THAT CALL. IT WILL BE BROKEN!
//
// Warning:
// LogAsFatal calls ekadeath.Die(1) (os.Exit(1) synonym) then.
// Make sure this is what you want.

// LogAsWarningThrough logs current Error if it's valid using provided 'logger' as Logger
// through which Error will be logged (only if it's valid too) with the 'LEVEL_WARNING'
// log level.
// YOU CAN NOT USE ERROR OBJECT AFTER THAT CALL. IT WILL BE BROKEN!


// LogAsErrorThrough logs current Error if it's valid using provided 'logger' as Logger
// through which Error will be logged (only if it's valid too) with the 'LEVEL_ERROR'
// log level.
// YOU CAN NOT USE ERROR OBJECT AFTER THAT CALL. IT WILL BE BROKEN!


// LogAsFatalThrough logs current Error if it's valid using provided 'logger' as Logger
// through which Error will be logged (only if it's valid too) with the 'LEVEL_FATAL'
// log level.
// YOU CAN NOT USE ERROR OBJECT AFTER THAT CALL. IT WILL BE BROKEN!
//
// Warning:
// LogAsFatalThrough calls ekadeath.Die(1) (os.Exit(1) synonym) then.
// Make sure this is what you want.

//
func (e *Error) log(logger *ekalog.Logger, level ekalog.Level, args []interface{}) {
	if e.IsNotNil() && (logger == nil || logger.IsValid()) {
		level, errLetter := e.logPreparations(level)
		ekaletter.BridgeLogErr2(unsafe.Pointer(logger), uint8(level), errLetter, args)
	}
}

//
func (e *Error) logf(logger *ekalog.Logger, level ekalog.Level, format string, args []interface{}) {
	if e.IsNotNil() && (logger == nil || logger.IsValid()) {
		level, errLetter := e.logPreparations(level)
		ekaletter.BridgeLogwErr2(unsafe.Pointer(logger), uint8(level), errLetter, fmt.Sprintf(format, args...), nil)
	}
}

//
func (e *Error) logw(logger *ekalog.Logger, level ekalog.Level, message string, fields []ekafield.Field) {
	if e.IsNotNil() && (logger == nil || logger.IsValid()) {
		level, errLetter := e.logPreparations(level)
		ekaletter.BridgeLogwErr2(unsafe.Pointer(logger), uint8(level), errLetter, message, fields)
	}
}

//
func (e *Error) logPreparations(level ekalog.Level) (newLevel ekalog.Level, errLetter *ekaletter.Letter) {
	
	if level < ekalog.LEVEL_WARNING {
		level = ekalog.LEVEL_ERROR
	}

	errLetterCopy := e.letter
	e.letter = nil

	// e saved inside e.letter.something thus Error won't be GC'ed while
	// its *Letter is alive.
	
	return level, errLetterCopy
}