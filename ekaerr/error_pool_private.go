// Copyright © 2020-2021. All rights reserved.
// Author: Ilya Stroy.
// Contacts: iyuryevich@pm.me, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekaerr

import (
	"runtime"
	"sync"
	"sync/atomic"

	"github.com/qioalice/ekago/v3/internal/ekaletter"
)

type (
	// ErrorPoolStat is an internal struct that allows you to inspect
	// how Error's pool utilized and it's current state.
	ErrorPoolStat struct {

		// AllocCalls is how much absolutely new Error objects
		// with included ekaletter.LetterField are created using new RAM slice.
		AllocCalls uint64

		// NewCalls is how much attempts both of to allocate a new Error objects
		// or pop an one from Error's pool were here.
		//
		// Once again.
		// It contains AllocCalls + popping an oldest (and prepared to reuse)
		// Error objects from its pool.
		NewCalls uint64

		// ReleaseCalls is how much Error objects were returned to its pool
		// and prepared for being reused.
		ReleaseCalls uint64
	}
)

// EPS returns an EntryPoolStat object that contains an info about utilizing
// Logger's Entry pool. Using that info you can figure out how often
// a new Logger's Entry objects are created and how often the oldest ones
// are reused.
//
// In 99% cases you don't need to know that stat,
// and you should not to worry about that.
func EPS() (eps ErrorPoolStat) {
	eps.AllocCalls = atomic.LoadUint64(&eps.AllocCalls)
	eps.NewCalls = atomic.LoadUint64(&eps.NewCalls)
	eps.ReleaseCalls = atomic.LoadUint64(&eps.ReleaseCalls)
	return
}

var (
	// errorPool is the pool of Error (with allocated ekaletter.Letter) objects
	// for being reused.
	errorPool sync.Pool

	// eps contains current state of Error's pool utilizing,
	// and its copy is returned by EPS() function.
	eps ErrorPoolStat
)

// allocError creates a new Error object, creates a new ekaletter.Letter object inside,
// performs base initialization and returns it.
func allocError() interface{} {

	e := new(Error)
	e.letter = new(ekaletter.Letter)
	e.letter.Messages = make([]ekaletter.LetterMessage, 0, 8)
	e.letter.Fields = make([]ekaletter.LetterField, 0, 16)

	runtime.SetFinalizer(e, releaseErrorForFinalizer)
	e.needSetFinalizer = false

	// SystemFields is used for saving Error's meta data.

	e.letter.SystemFields = make([]ekaletter.LetterField, 4)

	e.letter.SystemFields[_ERR_SYS_FIELD_IDX_CLASS_ID].Key = "class_id"
	e.letter.SystemFields[_ERR_SYS_FIELD_IDX_CLASS_ID].Kind |=
		ekaletter.KIND_FLAG_SYSTEM | ekaletter.KIND_SYS_TYPE_EKAERR_CLASS_ID

	e.letter.SystemFields[_ERR_SYS_FIELD_IDX_CLASS_NAME].Key = "class_name"
	e.letter.SystemFields[_ERR_SYS_FIELD_IDX_CLASS_NAME].Kind |=
		ekaletter.KIND_FLAG_SYSTEM | ekaletter.KIND_SYS_TYPE_EKAERR_CLASS_NAME

	e.letter.SystemFields[_ERR_SYS_FIELD_IDX_ERROR_ID].Key = "error_id"
	e.letter.SystemFields[_ERR_SYS_FIELD_IDX_ERROR_ID].Kind |=
		ekaletter.KIND_FLAG_SYSTEM | ekaletter.KIND_SYS_TYPE_EKAERR_UUID

	atomic.AddUint64(&eps.AllocCalls, 1)
	return e
}

// acquireError returns a new *Error object from the Error's pool or newly instantiated.
func acquireError() *Error {
	atomic.AddUint64(&eps.NewCalls, 1)
	return errorPool.Get().(*Error).prepare()
}

// releaseError returns Error to the Error's pool for being reused in the future
// and that Error could be obtained later using acquireError().
func releaseError(e *Error) {
	atomic.AddUint64(&eps.ReleaseCalls, 1)
	errorPool.Put(e.cleanup())
}

// releaseErrorForFinalizer is a callback for runtime.SetFinalizer()
// that allows to return an Error to its pool if it's gone out of scope
// without automatic returning to its pool by any ekalog.Logger's finisher.
func releaseErrorForFinalizer(e *Error) {
	e.needSetFinalizer = true
	releaseError(e)
}
