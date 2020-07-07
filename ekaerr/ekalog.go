// Copyright © 2020. All rights reserved.
// Author: Ilya Stroy.
// Contacts: qioalice@gmail.com, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekaerr

import (
	"unsafe"

	"github.com/qioalice/ekago/ekalog"
	"github.com/qioalice/ekago/internal/letter"
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
func (e *Error) LogAs(level ekalog.Level) {
	if e.IsValid() {
		if level < ekalog.LEVEL_WARNING {
			level = ekalog.LEVEL_ERROR
		}
		errLetterCopy := e.letter
		e.letter = nil
		// e saved inside e.letter.something thus Error won't be GC'ed while
		// its *Letter is alive.
		letter.GLogErrThroughDefaultLogger(uint8(level), errLetterCopy)
	}
}

// LogAsThrough logs current Error if it's valid using provided 'logger' as Logger
// through which Error will be logged (only if it's valid too) with the log level 'level'.
// YOU CAN NOT USE ERROR OBJECT AFTER THAT CALL. IT WILL BE BROKEN!
//
// Requirements:
// 'level' >= ekalog.LEVEL_WARNING. Otherwise overwritten by ekalog.LEVER_ERROR.
func (e *Error) LogAsThrough(level ekalog.Level, logger *ekalog.Logger) {
	if e.IsValid() && logger.IsValid() {
		if level < ekalog.LEVEL_WARNING {
			level = ekalog.LEVEL_ERROR
		}
		letter.GLogErr(unsafe.Pointer(logger), uint8(level), e.letter)
	}
}

// LogAsWarning logs current Error if it's valid using standard ekalog package's logger
// with 'LEVEL_WARNING' log level.
// YOU CAN NOT USE ERROR OBJECT AFTER THAT CALL. IT WILL BE BROKEN!
func (e *Error) LogAsWarning() {
	e.LogAs(ekalog.LEVEL_WARNING)
}

// LogAsError logs current Error if it's valid using standard ekalog package's logger
// with 'LEVEL_ERROR' log level.
// YOU CAN NOT USE ERROR OBJECT AFTER THAT CALL. IT WILL BE BROKEN!
func (e *Error) LogAsError() {
	e.LogAs(ekalog.LEVEL_ERROR)
}

// LogAsFatal logs current Error if it's valid using standard ekalog package's logger
// with 'LEVEL_FATAL' log level.
// YOU CAN NOT USE ERROR OBJECT AFTER THAT CALL. IT WILL BE BROKEN!
//
// Warning:
// LogAsFatal calls ekadeath.Die(1) (os.Exit(1) synonym) then.
// Make sure this is what you want.
func (e *Error) LogAsFatal() {
	e.LogAs(ekalog.LEVEL_FATAL)
}

// LogAsWarningThrough logs current Error if it's valid using provided 'logger' as Logger
// through which Error will be logged (only if it's valid too) with the 'LEVEL_WARNING'
// log level.
// YOU CAN NOT USE ERROR OBJECT AFTER THAT CALL. IT WILL BE BROKEN!
func (e *Error) LogAsWarningThrough(logger *ekalog.Logger) {
	e.LogAsThrough(ekalog.LEVEL_WARNING, logger)
}

// LogAsErrorThrough logs current Error if it's valid using provided 'logger' as Logger
// through which Error will be logged (only if it's valid too) with the 'LEVEL_ERROR'
// log level.
// YOU CAN NOT USE ERROR OBJECT AFTER THAT CALL. IT WILL BE BROKEN!
func (e *Error) LogAsErrorThrough(logger *ekalog.Logger) {
	e.LogAsThrough(ekalog.LEVEL_ERROR, logger)
}

// LogAsFatalThrough logs current Error if it's valid using provided 'logger' as Logger
// through which Error will be logged (only if it's valid too) with the 'LEVEL_FATAL'
// log level.
// YOU CAN NOT USE ERROR OBJECT AFTER THAT CALL. IT WILL BE BROKEN!
//
// Warning:
// LogAsFatalThrough calls ekadeath.Die(1) (os.Exit(1) synonym) then.
// Make sure this is what you want.
func (e *Error) LogAsFatalThrough(logger *ekalog.Logger) {
	e.LogAsThrough(ekalog.LEVEL_FATAL, logger)
}
