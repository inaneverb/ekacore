// Copyright © 2020. All rights reserved.
// Author: Ilya Stroy.
// Contacts: qioalice@gmail.com, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekatime

import (
	"sync"
	"sync/atomic"
	"time"
	"unsafe"

	"github.com/qioalice/ekago/v2/ekadeath"

)

type (
	// Calendar is a special type that represents a gregorian calendar
	// extending it with business calendar features, like:
	// - Supporting workdays and day offs,
	// - Supporting holidays or unscheduled work days,
	// - Caching the "today" state,
	// - Encoding the "today" day as JSON or user's defined encoder (+caching),
	// - Callback when new day has come (ctrl+c protected, graceful shutdown)
	// - Thread-safety updating when running,
	// - Updating logging,
	// etc.
	//
	// So, in most cases you need ONLY ONE INSTANCE of Calendar for the whole
	// your service. Trust me, it's enough, even in some cases when you need
	// different events.
	//
	// Create a Calendar object, configure it (add events, encoders), and call Run().
	// IT IS IMPORTANT! THE CALENDAR WON'T WORK CORRECTLY IF Run() WOULD NOT BE CALLED!
	//
	// -----
	//
	// When Calendar will be running you may access to the current's day state by:
	//
	// 1. Your own registered callback by WhenNewDay() method, that will be called
	//    when new day has come (at midnight, 00:00:00 (12:00:00AM)),
	//    and *Today object will be passed.
	//
	// 2. Use Today() method.
	//    It will return the cached *Today object (very fast accessing),
	//    ready to any usage.
	//
	// -----
	//
	// Using Today object you may get an info about current day and current month
	// like:
	// - What today the day is,
	// - Whether it's day off or not,
	// - Workdays number in month,
	// - How much workdays has been passed (or what is working day today),
	// etc.
	//
	// See https://github.com/qioalice/ekago/ekatime/today.go for more info.
	//
	Calendar struct {

		today unsafe.Pointer // *Today object, protected by atomic operations

		newDayTimer *time.Timer // timer that fires when new day's midnight has come
		// timer that fires after N seconds since new event has been added or
		// and old one has been removed
		logPendingEventsChangedTimer *time.Timer

		mu sync.Mutex // all fields named "pending..." protector
		wg sync.WaitGroup // graceful shutdown for user's callback and ekadeath

		isStarted bool // flag: Run() must be called only once
		disableLogging bool // flag: whether logging must be disabled

		// All next fields (except counters) has 2 variants: pending and confirmed.
		// By default, when any Calendar's setter is called, the pending related
		// field is changed.
		//
		// Later, when *Today object must be updated, the pending fields
		// will be flashed into the confirmed ones.
		//
		// It allows us to avoid protecting 'today' field using mutex
		// keeping thread safety even at the *Today reconstructing!

		confirmedNewDayCallback func(*Today) // user defined callback when new day has come
		pendingNewDayCallback func(*Today)

		confirmedEvents []Event // user defined events that may change month counters
		pendingEvents []Event

		confirmedTodayEncoderJson func(*Today) []byte
		pendingTodayEncoderJson func(*Today) []byte

		confirmedTodayEncoderUsers1 func(*Today) []byte
		pendingTodayEncoderUsers1 func(*Today) []byte

		// -----

		// The next two counters is how much new events has been added or/and
		// old one has been removed since the last 30s.
		// The counters will be used at the logging allowing us to get an actual state
		// of Calendar through logs and delay allows us not to spam logs
		// (timer resets when another one Calendar's setter has been called).

		pendingEventsAddCounter int
		pendingEventsRemoveCounter int

		// The next two counters is how much new events has been added or/and
		// old one has been removed since the last *Today object updating.
		// The counters will be used at the logging allowing us to get
		// how CONFIRMED events changes since the last day.

		pendingEventsAddTotalCounter int
		pendingEventsRemoveTotalCounter int
	}
)

// DisableLogging disables logging the internal state of Calendar when it under updating.
// If logging once disabled it cannot be enabled again.
// Nil safe. Thread-safety.
//
// DOES NOTHING IF CALENDAR ALREADY RUNNING.
// CALL THIS METHOD BEFORE Run() IS CALLED!
//
// By default after Run() is called, if you call any Calendar's setter,
// sometime the info messages will be logged about Calendar's update processing.
// You may disable it using this method.
func (c *Calendar) DisableLogging() *Calendar {

	if c == nil {
		return nil
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.isStarted {
		c.disableLogging = true
	}

	return c
}

// EventAdd adds a new event to the Calendar that will be applied when new day will come.
// If the same day event already exists, does nothing (remove it before).
// Nil safe. Thread-safety.
func (c *Calendar) EventAdd(event Event) *Calendar {

	if c == nil {
		return nil
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	if c.pendingEventIdx(event) == -1 {
		c.pendingEvents = append(c.pendingEvents, event)

		c.pendingEventsAddCounter++
		c.pendingEventsAddTotalCounter++

		if c.logPendingEventsChangedTimer != nil && !c.disableLogging {
			c.logPendingEventsChangedTimer.Reset(_CAL_PENDING_LOG_STAT_DELAY)
		}
	}

	return c
}

// EventRemove removes the provided event from the Calendar when new day will come.
// If there is no the same day event, does nothing.
// Nil safe. Thread-safety.
func (c *Calendar) EventRemove(event Event) *Calendar {

	if c == nil {
		return nil
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	if idx := c.pendingEventIdx(event); idx != -1 {
		// https://github.com/golang/go/wiki/SliceTricks#delete-without-preserving-order
		c.pendingEvents[idx] = c.pendingEvents[len(c.pendingEvents)-1]
		c.pendingEvents = c.pendingEvents[:len(c.pendingEvents)-1]

		c.pendingEventsRemoveCounter++
		c.pendingEventsRemoveTotalCounter++

		if c.logPendingEventsChangedTimer != nil && !c.disableLogging {
			c.logPendingEventsChangedTimer.Reset(_CAL_PENDING_LOG_STAT_DELAY)
		}
	}

	return c
}

// WhenNewDay registers the 'cb' as callback that will be called when new day has come.
// Keep in mind, you must call Run() method to start the internal goroutine of Calendar.
// Does nothing if Run() has been called before.
// Nil safe. Thread-safety.
//
// It's guaranteed that the service won't be stopped from the ekadeath's package
// calls until the 'cb' finished their work, but if there is a lot of work,
// IT IS STRONGLY RECOMMEND TO SPAWN YOUR OWN WORK GOROUTINE AND HANDLE IT MANUALLY!
func (c *Calendar) WhenNewDay(cb func(*Today)) *Calendar {

	if c == nil || cb == nil {
		return c
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	c.pendingNewDayCallback = cb // just overwrite
	return nil
}

// RegJsonEncoder registers 'encoder' as a new JSON *Today's object encoder,
// allowing to encoded it once and then get the same encoded data in the same day.
// The encoder will be used when the new day will come.
// You may get an access to the encoded data using 'Today.AsJson' field.
// Nil safe. Thread-safety.
//
// YOU MUST NOT HOLD THE Today's POINTER INSIDE ENCODER!
// DATA RACE AND MEMORY LEAKS OTHERWISE!
func (c *Calendar) RegJsonEncoder(encoder func(*Today) []byte) *Calendar {

	if c == nil || encoder == nil {
		return c
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	c.pendingTodayEncoderJson = encoder // just overwrite
	return c
}

// RegYourOwnEncoder registers 'encoder' as a new *Today's object encoder
// (the way you want), allowing to encode it once and then get the same encoded data
// in the same day. The encoder will be used when the new day will come.
// You may get an access to the encoded data using 'Today.AsUser<number>' field,
// where <number> is 'number' arg.
// Nil safe. Thread-safety.
//
// YOU MUST NOT HOLD THE Today's POINTER INSIDE ENCODER!
// DATA RACE AND MEMORY LEAKS OTHERWISE!
func (c *Calendar) RegYourOwnEncoder(number int, encoder func(*Today) []byte) *Calendar {

	if c == nil || encoder == nil {
		return c
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	c.pendingTodayEncoderUsers1 = encoder // just overwrite
	return c
}

// Today returns the cached Today's pointer. It's blazing fast access.
// You must call Run() before, otherwise nil is returned.
// Nil safe. Thread-safety.
func (c *Calendar) Today() *Today {

	if c == nil {
		return nil
	}

	return (*Today)(atomic.LoadPointer(&c.today))
}

// Run starts the Calendar's internal timers, goroutines, etc, allowing you
// to use Today() method and getting an actual Today object in your registered
// callback (by WhenNewDay()) when new day has come.
// Nil safe. Thread-safety.
func (c *Calendar) Run() {

	if c == nil {
		return
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	if c.isStarted {
		return
	}

	if !c.disableLogging {
		c.logPendingEventsChangedTimer =
			time.AfterFunc(_CAL_PENDING_LOG_STAT_DELAY, c.deferredLoggingOfPendingStat)
		c.logPendingEventsChangedTimer.Stop()
	}

	// Call confirmPending() before zeroing the counters.
	c.confirmPending().updateToday()

	// Null the counters. It's first run.
	c.pendingEventsAddCounter = 0
	c.pendingEventsRemoveCounter = 0
	c.pendingEventsAddTotalCounter = 0
	c.pendingEventsRemoveTotalCounter = 0

	ekadeath.Reg(c.destructor)
	c.newDayTimer = time.AfterFunc(TillNextMidnight(), c.newDayHasCome)
}
