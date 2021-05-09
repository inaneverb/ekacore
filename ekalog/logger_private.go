// Copyright © 2018-2021. All rights reserved.
// Author: Ilya Stroy.
// Contacts: qioalice@gmail.com, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekalog

import (
	"fmt"
	"time"
	"unsafe"

	"github.com/qioalice/ekago/v3/ekadeath"
	"github.com/qioalice/ekago/v3/ekaerr"
	"github.com/qioalice/ekago/v3/internal/ekaclike"
	"github.com/qioalice/ekago/v3/internal/ekaletter"

	"github.com/modern-go/reflect2"
)

var (
	// baseLogger is default package-level Logger, that used by all package-level
	// logger functions.
	baseLogger *Logger

	// nopLogger is a special Logger that is returned to indicate,
	// that next methods must do nothing.
	nopLogger *Logger
)

func (l *Logger) assert() {
	if !l.IsValid() {
		panic("Failed to do something with Logger. Logger is malformed.")
	}
}

// levelEnabled reports whether Entry with provided Level should be handled.
func (l *Logger) levelEnabled(lvl Level) bool {
	return lvl <= l.integrator.MinLevelEnabled()
}

// derive returns a new Logger with cloned Entry based on current Logger.
func (l *Logger) derive() (newLogger *Logger) {
	return new(Logger).setIntegrator(l.integrator).setEntry(l.entry.clone())
}

// setIntegrator changes the Logger's Integrator to the passed.
// It's just assignment nothing more. Useful at the method chaining.
func (l *Logger) setIntegrator(newIntegrator Integrator) (this *Logger) {
	l.integrator = newIntegrator
	return l
}

// setEntry changes the Logger's Entry to the passed. Also changes parent ptr
// in newEntry to being pointed to the l. Useful at the method chaining.
func (l *Logger) setEntry(newEntry *Entry) (this *Logger) {
	l.entry = newEntry
	l.entry.l = l
	return l
}

// addField checks whether Logger is valid, not nop Logger, makes a copy
// if it's requested and adds an ekaletter.LetterField to current Logger's Entry,
// if field is addable.
// Returns modified current Logger or its modified copy.
func (l *Logger) addField(f ekaletter.LetterField) *Logger {
	l.assert()
	if l == nopLogger || f.IsInvalid() || f.RemoveVary() && f.IsZero() {
		return l
	}
	ekaletter.LAddField(l.entry.LogLetter, f)
	return l
}

// addFields is the same as addField() but works with an array of ekaletter.LetterField,
// making a copy only once if it's requested.
func (l *Logger) addFields(fs []ekaletter.LetterField) *Logger {
	l.assert()
	if l == nopLogger || len(fs) == 0 {
		return l
	}
	for i, n := 0, len(fs); i < n; i++ {
		ekaletter.LAddFieldWithCheck(l.entry.LogLetter, fs[i])
	}
	return l
}

// addFieldsParse creates a ekaletter.LetterField objects based on passed values,
// try to treating them as a key-value pairs of that fields.
// It also makes a copy if it's requested and adds generated ekaletter.LetterField
// to the destination Logger's Entry only if those fields are addable.
// Returns modified current Logger or its modified copy.
func (l *Logger) addFieldsParse(fs []interface{}) *Logger {
	l.assert()
	if l == nopLogger || len(fs) == 0 {
		return l
	}
	ekaletter.LParseTo(l.entry.LogLetter, fs, true)
	return l
}

// log starts and manages all processes that should be passed between
// forming log message and writing it.
//
// There are 4 big things this method shall done:
//
// 1. Figure out what will be used as log message's body and construct it
//    if necessary (wasn't done before):
//
//    - Uses 'format' as just string if len(args) == 0;
//    - Uses 'format' as printf-like format string if it's contains printf verbs
//      and len(args) > 0, and then uses args[:N] as printf values;
//    - Uses err's message as log's string if 'format' == "" and len(args) == 0;
//    - Tries to extract string (or something string-like, e.g. fmt.Sprinter)
//      or Golang error from 'args[0]' and then use it in an one of four cases
//      described above. In that case only 'args[1:]' are allowed to be processed then.
//
// 2. Figure out what fields should be attached to the log message,
//    Convert implicit fields (both of named/unnamed) to the explicit
//    (using Entry.parseLogArgs method);
//
//    It guarantees that only one of 'args' or 'explicitFields' provided
//    at the same time.
//    If there's printf message, only 'args[N:]' or 'args[N+1:]' uses as
//    explicit/implicit args, where N - printf verbs.
//    N or N+1 depends by whether 'format' or 'err' was provided
//    as printf-like string, or it was extracted from 'args[0]'.
//
// 3. Adds caller and stacktrace (if it necessary and if it wasn't provided
//    by gext.errors.Error) (using Entry.addStacktrace method);
//
//    Uses gext.errors.Error.Stacktrace as stacktrace avoiding another one
//    stacktrace generation procedure.
//
// 4. Finally write a message and call then death.Die() if it's fatal level.
func (l *Logger) log(

	lvl      Level,
	format   string,
	err      *ekaerr.Error,
	args     []interface{},
	fields   []ekaletter.LetterField,

) *Logger {

	l.assert()
	if l == nopLogger || !l.levelEnabled(lvl) {
		return l
	}

	// empty messages are skipped by default, but who knows?
	if err.IsNil() && format == "" && len(args) == 0 && len(fields) == 0 {
		return l
	}

	// At this code point Entry should be copied anyway even if it's empty message,
	// even w/o body, even w/o fields, because of at least timestamp.
	//
	// But we don't need to copy an internal Logger's entry:
	// There is no necessary to keep (and then return) Logger with modified entry.
	// Because otherwise I guess it will lead to very unpredictable behaviour.
	workTempEntry := l.entry.clone()

	workTempEntry.Level = lvl
	workTempEntry.Time = time.Now()

	var (
		onlyFields = false
		errLetter = ekaletter.BridgeErrorGetLetter(unsafe.Pointer(err))
	)

	// Try to use ekaerr.Error's last message or first arg from args
	// if format is not presented.

	if format == "" {
		if errLetter != nil {
			format = ekaletter.LPopLastMessage(errLetter)
		} else if len(args) > 0 {
			var (
				typ1stArg = reflect2.TypeOf(args[0])
				rtype1stArg = uintptr(0)
			)

			if args[0] != nil {
				rtype1stArg = typ1stArg.RType()
			}

			if rtype1stArg == ekaclike.RTypeString {
				var str string
				typ1stArg.UnsafeSet(unsafe.Pointer(&str), reflect2.PtrOf(args[0]))
				args = args[1:]
				if format = str; format == "" {
					// if user pass empty string as argument,
					// it means he don't want log message, only fields
					onlyFields = true
				}

			} else if typ1stArg.Implements(ekaletter.TypeFmtStringer) {
				stringer := args[0].(fmt.Stringer)
				if stringer != nil {
					format = stringer.String()
					args = args[1:]
					if format == "" {
						onlyFields = true
					}
				}
			}
		}
	}

	workTempEntry.LogLetter.Messages[0].Body = format
	workTempEntry.ErrLetter = errLetter

	if lvl <= l.integrator.MinLevelForStackTrace() {
		workTempEntry.addStacktraceIfNotPresented()
	}

	// Try to extract message from 'args' if 'errLetter' == nil ('onlyFields' == false),
	// but if 'errLetter' is set, it's OK to log w/o message.
	switch {
	case len(args) > 0:
		ekaletter.LParseTo(workTempEntry.LogLetter, args, onlyFields)
	case len(fields) > 0:
		workTempEntry.LogLetter.Fields = fields
	}

	l.integrator.EncodeAndWrite(workTempEntry)

	ekaerr.ReleaseError(err)
	releaseEntry(workTempEntry)

	switch workTempEntry.Level {
	case LEVEL_EMERGENCY:
		ekadeath.Die()
	}

	return l
}
