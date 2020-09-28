// Copyright © 2020. All rights reserved.
// Author: Ilya Stroy.
// Contacts: qioalice@gmail.com, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekatime

import (
	"sync"
	"sync/atomic"
	"time"
)

type (
	// onceInCallback is an alias to the function that user must define.
	//
	// That function user could pass to the OnceIn<interval>.Call() call.
	// Then it will be called each time when time has come.
	//
	// An arguments represents a moment when a callback chain has been started
	// for being called.
	onceInCallback func(ts Timestamp, dd Date, t Time)

	// onceInUpdater is a special internal struct that gets the current timestamp
	// once in some period and caching it allowing to get the cached data by getters.
	onceInUpdater struct {
		cbs []onceInCallback // callbacks that must be called when time has come
		cbsMutex sync.Mutex
		ts Timestamp // cached current Timestamp
		d Date // cached current Date
		t Time // cached current Time
		tm *time.Timer // timer that allows to update the cached data
		updateDelayInSec Timestamp // timer delay
	}
)

// update updates the cached data inside the current onceInUpdater to the actual one
// using atomic operations. Returns the new actual Timestamp.
func (oiu *onceInUpdater) update() Timestamp {
	ts := Now()
	d, t := ts.Split()
	atomic.StoreInt64((*int64)(&oiu.ts), int64(ts))
	atomic.StoreUint32((*uint32)(&oiu.d), uint32(d))
	atomic.StoreUint32((*uint32)(&oiu.t), uint32(t))
	oiu.callbacksCaller()
	return ts
}

// callbacksCaller calls a registered delayed callbacks at the one goroutine,
// passing the current Timestamp, Date, Time to them.
func (oiu *onceInUpdater) callbacksCaller() {

	// Make a copy.
	// We don't know how long these callbacks will be running.
	//
	// There's also no need to use atomic package, cause there will no data race.
	// These values just were updated.
	ts, dd, t := oiu.ts, oiu.d, oiu.t

	oiu.cbsMutex.Lock()
	defer oiu.cbsMutex.Unlock()

	if len(oiu.cbs) == 0 {
		return
	}

	callbacksCopy := make([]onceInCallback, len(oiu.cbs))
	copy(callbacksCopy, oiu.cbs)

	f := func(cbs []onceInCallback, ts Timestamp, dd Date, t Time) {
		for i, n := 0, len(cbs); i < n; i++ {
			cbs[i](ts, dd, t)
		}
	}

	go f(callbacksCopy, ts, dd, t)
}

// tick is a special method that is called when onceInUpdater's timer is triggered.
// Updates the onceInUpdater's cached data, plans and runs timer at least one more time.
func (oiu *onceInUpdater) tick() {
	oiu.tm.Reset(oiu.update().tillNext(oiu.updateDelayInSec))
}

// run starts the onceInUpdater internal timer, fills the cached data by initial values.
func (oiu *onceInUpdater) run(delayInSec Timestamp) {
	oiu.update()
	oiu.updateDelayInSec = delayInSec
	oiu.tm = time.AfterFunc(Now().tillNext(delayInSec), oiu.tick)
}

// initOnceIn initializes all package level onceInUpdater global variables.
func initOnceIn() {
	OnceInMinute.run(SECONDS_IN_MINUTE)
	OnceIn10Minutes.run(SECONDS_IN_MINUTE * 10)
	OnceIn15Minutes.run(SECONDS_IN_MINUTE * 15)
	OnceIn30Minutes.run(SECONDS_IN_MINUTE * 30)
	OnceInHour.run(SECONDS_IN_HOUR)
	OnceIn12Hours.run(SECONDS_IN_HOUR * 12)
	OnceInDay.run(SECONDS_IN_DAY)
}
