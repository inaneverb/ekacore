// Copyright © 2020. All rights reserved.
// Author: Ilya Stroy.
// Contacts: qioalice@gmail.com, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekatime

import (
	"sync/atomic"
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
	// because of Now() call more expensive than these calls (up to 8x times).
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
//     ekatime.OnceInHour.Call(func(ts Timestamp, _ Date, _ Time){
//         fmt.Println(ts)
//     })
//
// will call provided callback every hour, printing the UNIX timestamp of the time
// when that hour has come.
//
// Does nothing if you pass a nil cb.
// That callback will be rejected with no-op.
//
// WARNING! IMPOSSIBLE TO STOP!
// YOU CAN NOT "CANCEL" A CALLBACK YOU ONCE ADDED TO PLANNER.
// IT WILL BE CALLED EVERY TIME UNTIL THE END.
// If you need to stop, handle it manually!
//
// WARNING! ONE THREAD!
// All callbacks associated with the same onceIn delayer (e.g. "once in hour",
// "once in day", etc) are called consistently one by one, AT THE SAME goroutine!
// So, if there is some "big" work, wrap your callback manually to the closure with
// "go callback(ts, dd, t)" call (spawn a separate goroutine).
func (oiu *onceInUpdater) Call(cb onceInCallback) {

	if cb == nil {
		return
	}

	oiu.cbsMutex.Lock()
	defer oiu.cbsMutex.Unlock()

	oiu.cbs = append(oiu.cbs, cb)
}
