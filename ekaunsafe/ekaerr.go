// Copyright © 2020. All rights reserved.
// Author: Ilya Stroy.
// Contacts: qioalice@gmail.com, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekaunsafe

import (
	"unsafe"

	"github.com/qioalice/ekago/v2/ekaerr"
	"github.com/qioalice/ekago/v2/ekasys"

	"github.com/qioalice/ekago/v2/internal/letter"
)

type (
	//
	CbErrorUpdateStacktrace func(oldStacktrace ekasys.StackTrace) (newStacktrace ekasys.StackTrace)
)

//
func ErrorGetLetter(err *ekaerr.Error) *letter.Letter {
	return letter.BridgeErrorGetLetter(unsafe.Pointer(err))
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

	errStackIdx := letter.BridgeErrorGetStackIdx(unsafe.Pointer(err))
	if errStackIdx < newStacktraceLen {
		return
	}

	letter.BridgeErrorSetStackIdx(unsafe.Pointer(err), newStacktraceLen-1)
}