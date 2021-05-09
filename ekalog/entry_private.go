// Copyright © 2020-2021. All rights reserved.
// Author: Ilya Stroy.
// Contacts: qioalice@gmail.com, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekalog

import (
	"runtime"

	"github.com/qioalice/ekago/v3/ekasys"
	"github.com/qioalice/ekago/v3/internal/ekaletter"
)

// prepare prepares current Entry for being used assuming that Entry has been
// obtained from the Entry's pool. Returns prepared Entry.
func (e *Entry) prepare() *Entry {

	// Because we can't detect when work with Entry is done while user chains
	// Logger's method (with Entry cloning), we must use runtime.SetFinalizer()
	// to return Entry to the its pool.
	if e.needSetFinalizer {
		runtime.SetFinalizer(e, releaseEntryForFinalizer)
		e.needSetFinalizer = false
	}

	return e
}

// cleanup frees all allocated resources (RAM in 99% cases) by Entry, preparing
// it for being returned to the pool and being reused in the future.
func (e *Entry) cleanup() (this *Entry) {

	e.l = nil
	e.LogLetter.StackTrace = nil
	e.ErrLetter = nil

	ekaletter.LReset(e.LogLetter)
	return e
}

// clone clones the current Entry and returns it copy. It takes a new Entry
// object from its pool to avoid unnecessary RAM allocations.
func (e *Entry) clone() *Entry {

	clonedEntry := acquireEntry()

	// Clone Fields using most efficient way.
	// Do not allocate RAM if it's already allocated (but nulled).
	if lFrom := len(e.LogLetter.Fields); lFrom > 0 {
		if cTo := cap(clonedEntry.LogLetter.Fields); cTo < lFrom {
			clonedEntry.LogLetter.Fields = make([]ekaletter.LetterField, lFrom)
		} else {
			// lFrom <= cTo, it's ok to do that
			clonedEntry.LogLetter.Fields =
				clonedEntry.LogLetter.Fields[:lFrom]
		}
		for i := 0; i < lFrom; i++ {
			clonedEntry.LogLetter.Fields[i] = e.LogLetter.Fields[i]
		}
	}

	// There is no need to zero Time, Level, LetterMessage fields
	// because they used only in one place and will be overwritten anyway.

	return clonedEntry
}

// addStacktraceIfNotPresented generates and adds stacktrace
// (if it's not presented by ErrLetter's field).
func (e *Entry) addStacktraceIfNotPresented() (this *Entry) {
	if e.ErrLetter == nil {
		e.LogLetter.StackTrace = ekasys.GetStackTrace(3, -1).ExcludeInternal()
	}
	return e
}
