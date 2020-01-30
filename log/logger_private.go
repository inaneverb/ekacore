// Copyright © 2019. All rights reserved.
// Author: Ilya Yuryevich.
// Contacts: qioalice@gmail.com, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package log

import (
	"fmt"
	"time"

	"github.com/qioalice/gext/death"
	"github.com/qioalice/gext/ec"
	"github.com/qioalice/gext/errors"
)

// canContinue reports whether l is valid (not nil and has valid core and entry).
func (l *Logger) canContinue() bool {

	// core and Entry (if they're not nil) can not be invalid (internally)
	// because they are created only by internal functions and they're private.
	// So 3 nil checks are enough here and of course check ptr equal.

	return l != nil &&
		l.integrator != nil &&
		l.entry != nil &&
		l == l.entry.l
}

// levelEnabled reports whether log's entry with level 'lvl' should be handled.
func (l *Logger) levelEnabled(lvl Level) bool {

	return lvl >= l.integrator.MinLevelEnabled()
}

// checkErr returns l if err != nil, otherwise nil.
// Make sure you have canContinue call before next ops over returned object.
//
// It's very useful at the method chaining.
// If err == nil, nil will be returned. Any Logger's public does nothing
// (w/o panicking) if Logger's receiver is nil.
func (l *Logger) checkErr(err error) (this *Logger) {

	if err == nil {
		return nil
	}
	return l
}

// apply changes l's behaviour according with new rules described by passed args.
func (l *Logger) apply(args []interface{}) (copy *Logger) {

	newIntegrator, options := parseOptions(args)
	return l.derive(newIntegrator).entry.apply(options).l
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
func (l *Logger) log(lvl Level, format string, err error, args []interface{}, explicitFields []Field) *Logger {

	if !(l.canContinue() && l.levelEnabled(lvl)) {
		return l
	}

	// empty messages are skipped by default, but who knows?
	if format == "" && err == nil && len(args) == 0 &&
		!l.entry.testFlag(bEntryFlagDontSkipEmptyMessages) {
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

	// if both of format and explicit error is empty, maybe first passed arg is err?
	// then handle it later as explicit error
	if format == "" && err == nil && len(args) > 0 {

		if err_, ok := args[0].(error); ok {
			err = err_ // will be handled later
			args = args[1:]
		}
	}

	switch {
	// use explicit error as message's body and moreover use its stacktrace
	// if it's led's package error
	case format == "" && err != nil:

		if gextErr := errors.Unwrap(err); gextErr != nil {
			workTempEntry.StackTrace = gextErr.StackTrace
		}
		format = err.Error()

	// maybe first arg is something like string (ducktypes)?
	// if it so, use it as message's body
	case format == "" && len(args) > 0:

		if str, ok := args[0].(string); ok && str != "" {
			format = str
			args = args[1:]

		} else if stringer, ok := args[0].(fmt.Stringer); ok && stringer != nil {
			format = stringer.String()
			args = args[1:]
		}
	}

	workTempEntry.
		parseLogArgs(format, args, explicitFields).
		addStacktrace()

	if workTempEntry.beforeWrite != nil {
		returned := workTempEntry.beforeWrite(workTempEntry)
		if returned == nil || returned != workTempEntry {
			goto exit
		}
	}

	l.integrator.Write(workTempEntry)

exit:
	switch {
	case workTempEntry.Level.mustDie():
		death.Die()
	}

	reuseEntry(workTempEntry)
	return l
}

// logec is the same as log but generates unique error's UUID,
// adds it to the log's entry, calls log and returns ec.ECXT based on passed 'errorCode'
// instead just returning Logger object.
func (l *Logger) logec(lvl Level, err error, errorCode ec.EC, args []interface{}, explicitFields []Field) (ret ec.ECXT) {

	if !(l.canContinue() && l.levelEnabled(lvl)) {
		return ec.EOK.ECXT()
	}

	ret = errorCode.ECXTForce()
	l.entry.ssfp++

	switch lA, lEF := len(args), len(explicitFields); {

	case lA == 0: // for both cases: 'lEF' ==/!= 0
		l.log(lvl, "", err, nil,
			append(explicitFields, String("error_id", ret.UUID.String())),
		)

	case lEF == 0:
		l.log(lvl, "", err,
			append(args, String("error_id", ret.UUID.String())),
			nil)
	}

	l.entry.ssfp--
	return
}

// derive returns a new Logger with at least clone Entry based on l's one
// and probably with a new integrator if 'newIntegrator' is not nil.
func (l *Logger) derive(newIntegrator Integrator) (newLogger *Logger) {

	// Entry clones any way because of that fact that entry stores ptr to parent.
	// if we won't clone we'll have a two different loggers with the same entry's ptrs
	// that are points to the only one (first or second) logger.

	if newIntegrator == nil {
		newIntegrator = l.integrator
	}
	return new(Logger).setIntegrator(newIntegrator).setEntry(l.entry.clone())
}

// setIntegrator changes the l's integrator to the passed. It's just assignment no more.
// Useful in method chaining.
func (l *Logger) setIntegrator(newIntegrator Integrator) (this *Logger) {

	// It guarantees that newCore != nil

	l.integrator = newIntegrator
	return l
}

// setEntry changes the l's base entry to the passed. Also changes parent ptr
// in newEntry to being pointed to the l. Useful in method chaining.
func (l *Logger) setEntry(newEntry *Entry) (this *Logger) {

	// It guarantees that newEntry != nil

	l.entry = newEntry
	l.entry.l = l

	return l
}
