// Copyright © 2020. All rights reserved.
// Author: Ilya Stroy.
// Contacts: qioalice@gmail.com, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekaunsafe

import (
	"unsafe"

	"github.com/qioalice/ekago/v2/ekaerr"
	"github.com/qioalice/ekago/v2/ekasys"

	"github.com/qioalice/ekago/v2/internal/ekaletter"
)

type (
	//
	CbErrorUpdateStacktrace func(oldStacktrace ekasys.StackTrace) (newStacktrace ekasys.StackTrace)
)

//
func ErrorGetLetter(err *ekaerr.Error) *ekaletter.Letter {
	return ekaletter.BridgeErrorGetLetter(unsafe.Pointer(err))
}

//
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