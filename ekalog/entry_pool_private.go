// Copyright © 2020. All rights reserved.
// Author: Ilya Stroy.
// Contacts: qioalice@gmail.com, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekalog

import (
	"runtime"
	"sync"
	"sync/atomic"

	"github.com/qioalice/ekago/internal/letter"
)

type (
	entryPoolStat struct {
		AllocCalls   uint64
		NewCalls     uint64
		ReleaseCalls uint64
	}
)

//noinspection GoSnakeCaseUsage
const (
	// _ENTRY_POOL_INIT_COUNT is how much letterPoolForLogEntries pool
	// will contain *Letter items at the start.
	_ENTRY_POOL_INIT_COUNT = 128
)

var (
	// entryPool is the pool of *Entry (with *Letter) objects for being reused.
	entryPool sync.Pool

	eps_ entryPoolStat
)

//
//noinspection GoExportedFuncWithUnexportedType
func EPS() (eps entryPoolStat) {
	eps.AllocCalls = atomic.LoadUint64(&eps_.AllocCalls)
	eps.NewCalls = atomic.LoadUint64(&eps_.NewCalls)
	eps.ReleaseCalls = atomic.LoadUint64(&eps_.ReleaseCalls)
	return
}

// allocEntry creates a new *Entry object, creates a new *Letter object inside Entry,
// performs base initialization and returns it.
func allocEntry() interface{} {

	e := new(Entry)
	e.LogLetter = new(letter.Letter)

	runtime.SetFinalizer(e, releaseEntryForFinalizer)
	e.needSetFinalizer = false

	// Alloc exactly one *LetterItem. We don't need more.
	e.LogLetter.Items = new(letter.LetterItem)
	letter.L_SetLastItem(e.LogLetter, e.LogLetter.Items)

	// SystemFields is used for saving Entry's meta data.
	// https://github.com/qioalice/ekago/internal/letter/letter.go

	// TODO: Is there any meta info we have to save to the SystemFields?
	// TODO: Is there something we need to save to e.LogLetter.something?

	atomic.AddUint64(&eps_.AllocCalls, 1)
	return e
}

// initEntryPool initializes entryPool creating and storing
// exactly _ENTRY_POOL_INIT_COUNT *Entry objects to that pool.
func initEntryPool() {
	entryPool.New = allocEntry
	for i := 0; i < _ENTRY_POOL_INIT_COUNT; i++ {
		entryPool.Put(allocEntry())
	}
}

// acquireEntry returns a new *Entry object from the Entry's pool or newly instantiated.
func acquireEntry() *Entry {
	atomic.AddUint64(&eps_.NewCalls, 1)
	return entryPool.Get().(*Entry).prepare()
}

// releaseEntry returns 'e' to the Entry's pool for being reused in the future
// and that Entry could be obtained later using acquireEntry().
func releaseEntry(e *Entry) {
	atomic.AddUint64(&eps_.ReleaseCalls, 1)
	entryPool.Put(e.cleanup())
}

//
func releaseEntryForFinalizer(e *Entry) {
	e.needSetFinalizer = true
	releaseEntry(e)
}
