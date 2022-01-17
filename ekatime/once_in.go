// Copyright © 2020. All rights reserved.
// Author: Ilya Stroy.
// Contacts: iyuryevich@pm.me, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekatime

import (
	"sync/atomic"
)

type (
	// OnceInCallback is an alias to the function that user must define.
	//
	// That function user could pass to the OnceIn<interval>.Call() call.
	// Then it will be called each time when time has come.
	//
	// An arguments represents a moment when a callback is called.
	OnceInCallback func(ts Timestamp)

	// OnceInPanicCallback is an alias to the function that user may define.
	//
	// That function user could pass along with OnceInCallback
	// to the OnceIn<interval>.Call() call.
	// This function will be called only if original OnceInCallback panics.
	//
	// An argument represents a recovered panic object.
	OnceInPanicCallback func(panicObj interface{})
)

//noinspection GoUnusedGlobalVariable
var (
	// -----
	// OnceIn<period> are a special objects that allows you to get the current time
	// that is not really actual but under updating each time when specified
	// period has come.
	//
	// For example, the OnceInMinute provides a new time for you each minute,
	// meaning that inside one minute you will get the same unix timestamp.
	//
	// It's useful when you do not need an exact time or time with highest precision
	// because of NewTimestampNow() call more expensive than these calls (up to 8x times).
	// -----

	// OnceInMinute allows you to get an actual time once in 60 seconds (1 minute).
	OnceInMinute onceInUpdater

	// OnceIn10Minutes allows you to get an actual time once in 10 minutes.
	OnceIn10Minutes onceInUpdater

	// OnceIn15Minutes allows you to get an actual time once in 15 minutes.
	OnceIn15Minutes onceInUpdater

	// OnceIn30Minutes allows you to get an actual time once in 30 minutes.
	OnceIn30Minutes onceInUpdater

	// OnceInHour allows you to get an actual time once in 60 minutes (1 hour).
	OnceInHour onceInUpdater

	// OnceIn2Hour allows you to get an actual time once in 120 minutes (2 hours).
	OnceIn2Hour onceInUpdater

	// OnceIn3Hour allows you to get an actual time once in 180 minutes (3 hours).
	OnceIn3Hour onceInUpdater

	// OnceIn6Hour allows you to get an actual time once in 1440 minutes (6 hours).
	OnceIn6Hour onceInUpdater

	// OnceIn12Hours allows you to get an actual time once in 12 hours.
	OnceIn12Hours onceInUpdater

	// OnceInDay allows you to get an actual time once in 24 hours (1 day).
	OnceInDay onceInUpdater
)

// Now returns the cached unix Timestamp from the current onceInUpdater that caches
// the current Timestamp once in the specified period.
func (oiu *onceInUpdater) Now() Timestamp {
	return Timestamp(atomic.LoadInt64((*int64)(&oiu.ts)))
}

// Date returns the cached Date from the current onceInUpdater that caches
// the current Date once in the specified period.
func (oiu *onceInUpdater) Date() Date {
	return Date(atomic.LoadUint32((*uint32)(&oiu.d)))
}

// Time returns the cached Time from the current onceInUpdater that caches
// the current Time once in the specified period.
func (oiu *onceInUpdater) Time() Time {
	return Time(atomic.LoadUint32((*uint32)(&oiu.t)))
}

// Call calls cb every time when associated onceIn updater's time has come.
// So, it means that
//
//     ekatime.OnceInHour.Call(func(ts Timestamp){
//         fmt.Println(ts)
//     })
//
// will call provided callback every hour, printing the UNIX timestamp of the time
// when that hour has come.
//
// Does nothing if you pass a nil cb.
// That callback will be rejected with no-op.
//
// WARNING!
// IMPOSSIBLE TO STOP! YOU CANNOT "CANCEL" A CALLBACK YOU ONCE ADDED TO PLANNER.
// It will be called every time until the end. If you need to stop, handle it manually!
//
// WARNING!
// ONE THREAD! ALL CALLBACKS SHARES THE SAME WORKER GOROUTINE.
// So, if there is some "big" work, wrap your callback manually to the closure with
// "go callback(ts, dd, t)" call (spawn a separate goroutine).
func (oiu *onceInUpdater) Call(invokeNow bool, cb OnceInCallback, panicCb ...OnceInPanicCallback) {
	if cb != nil {
		cbNum := atomic.AddUint32(&oiu.cbNum, 1) - 1
		onceInRegister(cb, panicCb, oiu.repeatDelay, 1, invokeNow, true, cbNum)
	}
}

// After is the same as Call() but calls cb after dur is passed.
//
// Does nothing if dur >= OnceIn's delayer time.
// It means, you cannot delay up to 1h function execution using OnceInHour.
func (oiu *onceInUpdater) After(delayInSec Timestamp, invokeNow bool, cb OnceInCallback, panicCb ...OnceInPanicCallback) {

	if delayInSec < 1 {
		delayInSec = 1
	}

	if cb != nil && delayInSec < oiu.repeatDelay {
		cbNum := atomic.AddUint32(&oiu.cbNum, 1) - 1
		onceInRegister(cb, panicCb, oiu.repeatDelay, delayInSec, invokeNow, true, cbNum)
	}
}
