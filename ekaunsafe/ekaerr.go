// Copyright © 2020. All rights reserved.
// Author: Ilya Stroy.
// Contacts: qioalice@gmail.com, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekaunsafe

import (
	"unsafe"

	"github.com/qioalice/ekago/v3/ekaerr"
	"github.com/qioalice/ekago/v3/ekasys"
	"github.com/qioalice/ekago/v3/internal/ekaletter"
)

type (
	// CbErrorUpdateStacktrace is a special func alias that is an argument
	// of ErrorUpdateStacktrace function.
	CbErrorUpdateStacktrace func(oldStacktrace ekasys.StackTrace) (newStacktrace ekasys.StackTrace)
)

// ErrorGetLetter returns an underlying ekaletter.Letter from provided ekaerr.Error object.
// Returns nil if err is not valid.
func ErrorGetLetter(err *ekaerr.Error) *ekaletter.Letter {
	return ekaletter.BridgeErrorGetLetter(unsafe.Pointer(err))
}

// ErrorUpdateStacktrace calls provided callback to change ekaerr.Error's stacktrace.
// Your callback must not be nil, error must be valid
// and your callback should return a new (modified) stacktrace that will be used.
func ErrorUpdateStacktrace(err *ekaerr.Error, cb CbErrorUpdateStacktrace) {

	if cb == nil || err.IsNil() {
		return
	}

	l := ErrorGetLetter(err)
	oldStacktraceLen := int16(len(l.StackTrace))

	l.StackTrace = cb(l.StackTrace)
	newStacktraceLen := int16(len(l.StackTrace))

	if newStacktraceLen >= oldStacktraceLen {
		return
	}

	errStackIdx := ekaletter.BridgeErrorGetStackIdx(unsafe.Pointer(err))
	if errStackIdx < newStacktraceLen {
		return
	}

	ekaletter.BridgeErrorSetStackIdx(unsafe.Pointer(err), newStacktraceLen-1)
}