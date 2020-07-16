// Copyright © 2020. All rights reserved.
// Author: Ilya Stroy.
// Contacts: qioalice@gmail.com, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekalog

import (
	"fmt"
	"os"
	"time"
	"unsafe"

	"github.com/qioalice/ekago/v2/ekadeath"
	"github.com/qioalice/ekago/v2/internal/field"
	"github.com/qioalice/ekago/v2/internal/letter"
)

var (
	// baseLogger is default package-level logger, that used by all package-level
	// logger functions such a Debug, Debugf, Debugw, With, Group
	baseLogger *Logger
)

// levelEnabled reports whether log's entry with level 'lvl' should be handled.
func (l *Logger) levelEnabled(lvl Level) bool {
	return lvl >= l.integrator.MinLevelEnabled()
}

// derive returns a new *Logger with at least cloned *Entry based on l's one
// and probably with a new Integrator if 'newIntegrator' is not nil.
func (l *Logger) derive(newIntegrator Integrator) (newLogger *Logger) {

	// Entry clones any way because of that fact that entry stores ptr to parent.
	// if we won't clone we'll have a two different loggers with the same entry's ptrs
	// that are points to the only one (first or second) logger.

	if newIntegrator == nil {
		newIntegrator = l.integrator
	}
	return new(Logger).setIntegrator(newIntegrator).setEntry(l.entry.clone())
}

// setIntegrator changes the l's integrator to the passed.
// It's just assignment nothing more. Useful at the method chaining.
func (l *Logger) setIntegrator(newIntegrator Integrator) (this *Logger) {
	l.integrator = newIntegrator
	return l
}

// setEntry changes the l's entry to the passed. Also changes parent ptr
// in newEntry to being pointed to the l. Useful at the method chaining.
func (l *Logger) setEntry(newEntry *Entry) (this *Logger) {
	l.entry = newEntry
	l.entry.l = l
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

	lvl Level,
	format string,
	errLetter *letter.Letter,
	args []interface{},
	explicitFields []field.Field,

) *Logger {

	if !(l.IsValid() && l.levelEnabled(lvl)) {
		return l
	}

	// empty messages are skipped by default, but who knows?
	if errLetter == nil && format == "" &&
		len(args) == 0 && len(explicitFields) == 0 &&
		!l.entry.LogLetter.Items.Flags.TestAll(FLAG_ALLOW_EMPTY_MESSAGES) {

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

	// maybe first arg is something like string (ducktypes)?
	// if it so, use it as message's body
	if errLetter == nil && format == "" && len(args) > 0 {

		if str, ok := args[0].(string); ok && str != "" {
			format = str
			args = args[1:]

		} else if stringer, ok := args[0].(fmt.Stringer); ok && stringer != nil {
			format = stringer.String()
			args = args[1:]
		}
	}

	workTempEntry.LogLetter.Items.Message = format
	workTempEntry.ErrLetter = errLetter
	workTempEntry.addStacktrace()

	// Try to extract message from 'args' if 'errLetter' == nil ('onlyFields' == false),
	// but if 'errLetter' is set, it's OK to log w/o message.
	if len(args) > 0 || len(explicitFields) > 0 {
		letter.ParseTo(workTempEntry.LogLetter.Items, args, explicitFields, errLetter != nil)
	}

	l.integrator.Write(workTempEntry)

	if !l.integrator.IsAsync() {
		if errLetter != nil {
			letter.GErrRelease(errLetter)
		}
		releaseEntry(workTempEntry)
	}

	switch {
	case workTempEntry.Level.mustDie():
		ekadeath.Die()
	}

	return l
}

// logErr is just the same as 'logger'.log() but only for *Error's *Letter 'errLetter'.
// Assumes that 'logger' is valid logger.
//
// Requirements:
// 'logger'.IsValid() == true. Otherwise UB (may panic).
func logErr(logger unsafe.Pointer, level uint8, errLetter *letter.Letter) {
	(*Logger)(logger).log(Level(level), "", errLetter, nil, nil)
}

// logErrThroughDefaultLogger is just the same as baseLogger.log() but only for
// *Error's *Letter 'errLetter'.
func logErrThroughDefaultLogger(level uint8, errLetter *letter.Letter) {
	baseLogger.log(Level(level), "", errLetter, nil, nil)
}

// initBaseLogger performs a baseLogger initialization.
func initBaseLogger() {

	defaultConsoleEncoder = new(CI_ConsoleEncoder).FreezeAndGetEncoder()
	defaultJSONEncoder = new(CI_JSONEncoder).FreezeAndGetEncoder()

	integrator := new(CommonIntegrator).
		WithEncoder(defaultConsoleEncoder).
		WithMinLevel(LEVEL_DEBUG).
		WriteTo(os.Stdout)

	entry := acquireEntry()
	baseLogger = new(Logger).setIntegrator(integrator).setEntry(entry)
}
