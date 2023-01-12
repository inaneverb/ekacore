// Copyright © 2020. All rights reserved.
// Author: Ilya Stroy.
// Contacts: iyuryevich@pm.me, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekaerr

import (
	"github.com/qioalice/ekago/v4/internal/ekaletter"
)

//import (
//	"unsafe"
//
//	"github.com/qioalice/ekago/v3/ekasys"
//	"github.com/qioalice/ekago/v3/internal/ekaletter"
//)
//
//type (
//	// CbErrorUpdateStacktrace is a special func alias that is an argument
//	// of ErrorUpdateStacktrace function.
//	CbErrorUpdateStacktrace func(oldStacktrace ekasys.StackTrace) (newStacktrace ekasys.StackTrace)
//)
//

// ErrorGetLetter returns an underlying ekaletter.Letter from provided ekaerr.Error object.
// Returns nil if err is not valid.
func GetLetter(err *Error) *ekaletter.Letter {
	return err.letter
}

//
//// ErrorUpdateStacktrace calls provided callback to change ekaerr.Error's stacktrace.
//// Your callback must not be nil, error must be valid
//// and your callback should return a new (modified) stacktrace that will be used.
//func ErrorUpdateStacktrace(err *Error, cb CbErrorUpdateStacktrace) {
//
//	if cb == nil || err.IsNil() {
//		return
//	}
//
//	l := ErrorGetLetter(err)
//	oldStacktraceLen := int16(len(l.StackTrace))
//
//	l.StackTrace = cb(l.StackTrace)
//	newStacktraceLen := int16(len(l.StackTrace))
//
//	if newStacktraceLen >= oldStacktraceLen {
//		return
//	}
//
//	errStackIdx := ekaletter.BridgeErrorGetStackIdx(unsafe.Pointer(err))
//	if errStackIdx < newStacktraceLen {
//		return
//	}
//
//	ekaletter.BridgeErrorSetStackIdx(unsafe.Pointer(err), newStacktraceLen-1)
//}
