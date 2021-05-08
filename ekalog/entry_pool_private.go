// Copyright © 2020-2021. All rights reserved.
// Author: Ilya Stroy.
// Contacts: qioalice@gmail.com, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekalog

import (
	"runtime"
	"sync"
	"sync/atomic"

	"github.com/qioalice/ekago/v2/internal/ekaletter"
)

type (
	// EntryPoolStat is an internal struct that allows you to inspect
	// how Logger's Entry pool utilized and it's current state.
	EntryPoolStat struct {

		// AllocCalls is how much absolutely new Entry objects
		// with included ekaletter.LetterField are created using new RAM slice.
		AllocCalls uint64

		// NewCalls is how much attempts both of to allocate a new Entry objects
		// or pop an one from Entry's pool were here.
		//
		// Once again.
		// It contains AllocCalls + popping an oldest (and prepared to reuse)
		// Entry objects from its pool.
		NewCalls uint64

		// ReleaseCalls is how much Entry objects were returned to its pool
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
func EPS() (eps EntryPoolStat) {
	eps.AllocCalls = atomic.LoadUint64(&eps.AllocCalls)
	eps.NewCalls = atomic.LoadUint64(&eps.NewCalls)
	eps.ReleaseCalls = atomic.LoadUint64(&eps.ReleaseCalls)
	return
}

var (
	// entryPool is the pool of Entry (with allocated ekaletter.Letter) objects
	// for being reused.
	entryPool sync.Pool

	// eps contains current state of Entry's pool utilizing,
	// and its copy is returned by EPS() function.
	eps EntryPoolStat
)

// allocEntry creates a new Entry object, creates a new ekaletter.Letter object inside,
// performs base initialization and returns it.
func allocEntry() interface{} {

	e := new(Entry)
	e.LogLetter = new(ekaletter.Letter)
	e.LogLetter.Messages = make([]ekaletter.LetterMessage, 1)

	runtime.SetFinalizer(e, releaseEntryForFinalizer)
	e.needSetFinalizer = false

	// SystemFields is used for saving Entry's meta data.
	// https://github.com/qioalice/ekago/internal/letter/letter.go

	// TODO: Is there any meta info we have to save to the SystemFields?
	// TODO: Is there something we need to save to e.LogLetter.something?

	atomic.AddUint64(&eps.AllocCalls, 1)
	return e
}

// acquireEntry returns a new *Entry object from the Entry's pool or newly instantiated.
func acquireEntry() *Entry {
	atomic.AddUint64(&eps.NewCalls, 1)
	return entryPool.Get().(*Entry).prepare()
}

// releaseEntry returns Entry to the Entry's pool for being reused in the future
// and that Entry could be obtained later using acquireEntry().
func releaseEntry(e *Entry) {
	atomic.AddUint64(&eps.ReleaseCalls, 1)
	entryPool.Put(e.cleanup())
}

// releaseEntryForFinalizer is a callback for runtime.SetFinalizer()
// that allows to return an Entry to its pool if it's gone out of scope
// without automatic returning to its pool by any Logger's finisher.
func releaseEntryForFinalizer(e *Entry) {
	e.needSetFinalizer = true
	releaseEntry(e)
}
