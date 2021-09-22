// Copyright © 2021. All rights reserved.
// Author: Ilya Stroy.
// Contacts: iyuryevich@pm.me, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekatime

import (
	"sync"
	"sync/atomic"
	"time"

	"github.com/qioalice/ekago/v3/ekadeath"

	heap "github.com/theodesp/go-heaps"
	fibheap "github.com/theodesp/go-heaps/fibonacci"
)

/*
"OnceIn" is a concept of repeatable delayed call of some functions when a time is come.
It means you can say "execute function f each 1h" and it will be but starting with next hour.
Your function won't be executed right now (unless otherwise specified).

So, it like Golang's time.Ticker but you don't need to worry about channels
for each time interval, about stopping/GC'ing timers/tickers, etc.

All functions that you register are executed in one "worker" goroutine
that is spawned when you calling some "OnceIn" function first time.
If your function is heavy wrap it by your own runner that will start your function
in a separate goroutine.

It guarantees that ekadeath.Die() or ekadeath.Exit() calls won't shutdown your app
when some your function is under executing, but next won't be executed.
Even if they has the same time of firing.

A minimum value of time you can use is a second. Time tolerance between
"time has come" and "call your function in that time" is about 1 sec.
*/

type (
	// onceInUpdater is a special internal struct that gets the current timestamp
	// once in some period and caching it allowing to get the cached data by getters.
	onceInUpdater struct {

		// WARNING!
		// DO NOT CHANGE THE ORDER OF FIELDS!
		// https://golang.org/pkg/sync/atomic/#pkg-note-BUG :
		//
		//   > On ARM, x86-32, and 32-bit MIPS,
		//   > it is the caller's responsibility to arrange for 64-bit alignment
		//   > of 64-bit words accessed atomically.
		//   > The first word in a variable or in an allocated struct, array,
		//   > or slice can be relied upon to be 64-bit aligned.
		//
		// Also:
		// https://stackoverflow.com/questions/28670232/atomic-addint64-causes-invalid-memory-address-or-nil-pointer-dereference/51012703#51012703

		/* 8b */ ts Timestamp // cached current Timestamp
		/* 4b */ d Date // cached current Date
		/* 4b */ t Time // cached current Time
		/* -- */ repeatDelay Timestamp
	}

	// onceInExeElem represents an execution element of "OnceIn" concept.
	// It contains a function that should be called and a time as unix timestamp
	// of when that function should be called.
	onceInExeElem struct {
		when, repeatDelay, afterDelay Timestamp
		cb                            OnceInCallback
		cbPanic                       OnceInPanicCallback
	}
)

//goland:noinspection GoSnakeCaseUsage
const (
	_ONCE_IN_SLEEP_TIME = 1 * time.Second
)

var (
	// onceInFibHeap is a Fibonacci Heap of onceInExeElem's sorted to the nearest
	// element that must be executed soon.
	// Read more: https://en.wikipedia.org/wiki/Fibonacci_heap .
	onceInFibHeap *fibheap.FibonacciHeap

	// onceInFibHeapMu is a sync.Mutex that provides thread-safety for RW access
	// to onceInFibHeap.
	onceInFibHeapMu sync.Mutex

	// onceInShutdownRequested is a "bool" atomic variable,
	// that is set to 1 when shutdown is requested by ekadeath package.
	onceInShutdownRequested int32

	// onceInShutdownConfirmed is a channel, an ekadeath's destructor will wait a value from
	// as a signal that it's safe to shutdown an app and no user's function
	// is under execution right now.
	// A worker guarantees that this channel will receive a value after
	// onceInShutdownRequested is set to 1.
	onceInShutdownConfirmed chan struct{}
)

// updateAll updates the cached data inside the current onceInUpdater
// to the provided actual ones using atomic operations.
func (oiu *onceInUpdater) updateAll(ts Timestamp, dd Date, t Time) {
	atomic.StoreInt64((*int64)(&oiu.ts), int64(ts))
	atomic.StoreUint32((*uint32)(&oiu.d), uint32(dd))
	atomic.StoreUint32((*uint32)(&oiu.t), uint32(t))
}

func (oiu *onceInUpdater) update(ts Timestamp) {
	dd, t := ts.Split()
	oiu.updateAll(ts, dd, t)
}

// init calls update() and then register this using onceInFibHeap to be updated.
func (oiu *onceInUpdater) init(now, repeatDelay Timestamp) {
	oiu.update(now)
	oiu.repeatDelay = repeatDelay
	onceInRegister(oiu.update, nil, repeatDelay, -1, false, false)
}

// Compare implements `go-heaps.Item` interface.
// It reports the comparing time difference of the current onceInExeElem and provided one.
//
// Returns:
// -1 if current onceInExeElem's time < anotherElem's time,
// 0 if they are the same,
// 1 if current onceInExeElem's time > anotherElem's time.
func (oie onceInExeElem) Compare(anotherOie heap.Item) int {
	return oie.when.Cmp(anotherOie.(onceInExeElem).when)
}

// invoke invokes onceInExeElem's callback passing provided Timestamp,
// checks whether it panics and if it so, calls onPanic callback.
func (oie onceInExeElem) invoke(ts Timestamp) {
	panicProtector := func(cb OnceInPanicCallback) {
		if panicObj := recover(); panicObj != nil {
			cb(panicObj)
		}
	}
	if oie.cbPanic != nil {
		defer panicProtector(oie.cbPanic)
	}
	oie.cb(ts)
}

// onceInWorker is a special worker that is running in a background goroutine,
// pulls nearest (by time) onceInExeElem from onceInFibHeap pool,
// checks whether it time has come and if it so, executes a function.
// Otherwise sleeps goroutine for _ONCE_IN_SLEEP_TIME duration.
func onceInWorker() {

	for atomic.LoadInt32(&onceInShutdownRequested) == 0 {
		ts := NewTimestampNow()
		onceInFibHeapMu.Lock()

		nearestOie := onceInFibHeap.FindMin()
		if nearestOie == nil || nearestOie.(onceInExeElem).when > ts {
			// Pool of onceInExeElems is empty or nearest item's time not come yet.
			// Abort current iteration, sleep, go next.
			onceInFibHeapMu.Unlock()
			time.Sleep(_ONCE_IN_SLEEP_TIME)
			continue
		}

		// If we're here, nearestOie is not nil and its time has come.
		_ = onceInFibHeap.DeleteMin() // the same as nearestOie

		// Register next call.
		nearestOieCopy := nearestOie.(onceInExeElem)
		nearestOieCopy.when = ts + ts.tillNext(nearestOieCopy.repeatDelay) + nearestOieCopy.afterDelay
		onceInFibHeap.Insert(nearestOieCopy)

		onceInFibHeapMu.Unlock()

		nearestOie.(onceInExeElem).invoke(ts)
	}

	// The loop above could be over only if shutdown is requested.
	// So, if we're here, we need to confirm shutdown.
	close(onceInShutdownConfirmed)
}

// onceInRegister registers a new OnceInCallback that must be called each `repeatDelay` time
// waiting for `afterDelay` when time has come, calling `cb` and if it panics,
// call then `panicCb[0]`. Calls right now if `invokeNow` is true.
// Protects access to onceInFibHeap using associated sync.Mutex if `doLock` is true.
func onceInRegister(cb OnceInCallback, panicCb []OnceInPanicCallback, repeatDelay, afterDelay Timestamp, invokeNow, doLock bool) {
	if doLock {
		onceInFibHeapMu.Lock()
	}

	panicCb_ := OnceInPanicCallback(nil)
	if len(panicCb) > 0 && panicCb[0] != nil {
		panicCb_ = panicCb[0]
	}

	if afterDelay < 0 {
		afterDelay = 0
	}

	oie := onceInExeElem{NewTimestampNow(), repeatDelay, afterDelay, cb, panicCb_}
	if !invokeNow {
		oie.when += oie.when.tillNext(repeatDelay) + afterDelay
	}

	onceInFibHeap.Insert(oie)

	if doLock {
		onceInFibHeapMu.Unlock()
	}
}

// initOnceIn initializes all package level onceInUpdater global variables,
// registers onceInFibHeap destructor and starts its worker.
func initOnceIn() {
	onceInFibHeap = fibheap.New()
	onceInShutdownConfirmed = make(chan struct{})

	ekadeath.Reg(func() {
		atomic.StoreInt32(&onceInShutdownRequested, 1)
		<-onceInShutdownConfirmed
	})

	now := NewTimestampNow()

	OnceInMinute.init(now, SECONDS_IN_MINUTE)
	OnceIn10Minutes.init(now, SECONDS_IN_MINUTE*10)
	OnceIn15Minutes.init(now, SECONDS_IN_MINUTE*15)
	OnceIn30Minutes.init(now, SECONDS_IN_MINUTE*30)
	OnceInHour.init(now, SECONDS_IN_HOUR)
	OnceIn2Hour.init(now, SECONDS_IN_HOUR*2)
	OnceIn3Hour.init(now, SECONDS_IN_HOUR*3)
	OnceIn6Hour.init(now, SECONDS_IN_HOUR*6)
	OnceIn12Hours.init(now, SECONDS_IN_HOUR*12)
	OnceInDay.init(now, SECONDS_IN_DAY)

	go onceInWorker()
}
