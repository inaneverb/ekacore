// Copyright © 2020. All rights reserved.
// Author: Ilya Stroy.
// Contacts: qioalice@gmail.com, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekatime

import (
	"sync/atomic"
	"time"
	"unsafe"

	"github.com/qioalice/ekago/v2/ekalog"
)

//noinspection GoSnakeCaseUsage
const (
	// Delay in the ns that must be passed since the Calendar's setter has been called
	// at the running Calendar for logging a current status of called event setters.
	_CAL_PENDING_LOG_STAT_DELAY = 30 * time.Second
)

// confirmPending flushes all Calendar's pending changes into confirmed.
// Do nothing if nothing changed.
func (c *Calendar) confirmPending() *Calendar {

	addTotalCounter := c.pendingEventsAddTotalCounter
	removeTotalCounter := c.pendingEventsRemoveTotalCounter
	if addTotalCounter != 0 || removeTotalCounter != 0 {
		// https://github.com/go101/go101/wiki/How-to-perfectly-clone-a-slice
		c.confirmedEvents = append(c.pendingEvents[:0:0], c.pendingEvents...)
	}

	c.confirmedNewDayCallback = c.pendingNewDayCallback
	c.confirmedTodayEncoderJson = c.pendingTodayEncoderJson
	c.confirmedTodayEncoderUsers1 = c.pendingTodayEncoderUsers1

	return c
}

// deferredLoggingOfPendingStat is logging function that writes logs if there are
// new pending added/removed events.
func (c *Calendar) deferredLoggingOfPendingStat() {

	c.mu.Lock()
	addCounter := c.pendingEventsAddCounter
	removeCounter := c.pendingEventsRemoveCounter
	c.pendingEventsAddCounter = 0
	c.pendingEventsRemoveCounter = 0
	c.mu.Unlock()

	if addCounter == 0 && removeCounter == 0 {
		return
	}

	// There is no need to check c.disableLogging, because timer won't be created
	// if logging is disabled and this function will never be called.

	ekalog.Debug("ekatime.Calendar has a new pending events will be processed later.",
		"c_pending_added_new_events_count?", addCounter,
		"c_pending_removed_old_events_count?", removeCounter)
}

// newDayHasCome is a special function that is called when new day (00:00 AM) has come.
// Called once for each new day. Updates timer for being called one more time (next day),
// and calls registered callback by WhenNewDay() method.
func (c *Calendar) newDayHasCome() {

	c.mu.Lock()
	c.confirmPending()
	addTotalCounter := c.pendingEventsAddTotalCounter
	removeTotalCounter := c.pendingEventsRemoveTotalCounter
	c.pendingEventsAddTotalCounter = 0
	c.pendingEventsRemoveTotalCounter = 0
	disableLogging := c.disableLogging
	c.mu.Unlock()

	newToday := c.updateToday()
	c.newDayTimer.Reset(TillNextMidnight())

	if !disableLogging && (addTotalCounter != 0 || removeTotalCounter != 0) {
		ekalog.Debug("ekatime.Calendar has been updated, and event list has been changed.",
			"c_added_new_events_count?", addTotalCounter,
			"c_removed_old_events_count?", removeTotalCounter)
	}

	if c.confirmedNewDayCallback != nil {
		c.wg.Add(1)
		defer c.wg.Done()
		c.confirmedNewDayCallback(newToday)
	}
}

// updateToday creates a new Today object, fills them, stores them into the
// Calendar's internal fast access cache and returns.
func (c *Calendar) updateToday() *Today {

	newToday := new(Today)
	newToday.c = c

	newToday.Timestamp = Now()

	newToday.Date, newToday.Time = newToday.Timestamp.Split()
	newToday.Year, newToday.Month, newToday.Day = newToday.Date.Split()
	newToday.Hour, newToday.Minute, newToday.Second = newToday.Time.Split()

	newToday.Weekday = newToday.Date.Weekday()
	newToday.IsDayOff = newToday.Weekday.IsDayOff() // may be overwritten later

	// Calculate work days counters

	_1stDayWeekday := newToday.Weekday

	newToday.DaysInMonth = _Table0[newToday.Month-1]
	if newToday.Month == MONTH_FEBRUARY && newToday.Year.IsLeap() {
		newToday.DaysInMonth++
	}

	if newToday.Day > 1 {
		_1stDayWeekday = NewDate(newToday.Year, newToday.Month, 1).Weekday()
	}

	for day := Day(1); day <= newToday.DaysInMonth; day++ {

		if !_1stDayWeekday.IsDayOff() {
			newToday.WorkDayTotal++

			if day <= newToday.Day {
				newToday.WorkDayCurrent++
			}
		}

		_1stDayWeekday = _1stDayWeekday.Next()
	}

	// Apply all events.
	// Also correct counters of workdays.

	for i, n := 0, len(c.confirmedEvents); i < n; i++ {
		if c.confirmedEvents[i].Year() == newToday.Year &&
			c.confirmedEvents[i].Month() == newToday.Month {

			if c.confirmedEvents[i].Day() == newToday.Day {
				newToday.IsDayOff = c.confirmedEvents[i].IsDayOff()
			}

			// If by default this day is the same as at the event, there is no need
			// to change something.
			if c.confirmedEvents[i].Weekday().IsDayOff() != c.confirmedEvents[i].IsDayOff() {
				if c.confirmedEvents[i].IsDayOff() {
					newToday.WorkDayTotal--
					if c.confirmedEvents[i].Day() <= newToday.Day {
						newToday.WorkDayCurrent--
					}
				} else {
					newToday.WorkDayTotal++
					if c.confirmedEvents[i].Day() <= newToday.Day {
						newToday.WorkDayCurrent++
					}
				}
			}
		}
	}

	newToday.DayOffTotal = newToday.DaysInMonth - newToday.WorkDayTotal

	if c.confirmedTodayEncoderJson != nil {
		newToday.AsJson = c.confirmedTodayEncoderJson(newToday)
	}
	if c.confirmedTodayEncoderUsers1 != nil {
		newToday.AsYourOwn1 = c.confirmedTodayEncoderUsers1(newToday)
	}

	old := atomic.SwapPointer(&c.today, unsafe.Pointer(newToday))
	if old != nil {
		(*Today)(old).c = nil // allow to being GC'ed
	}

	return newToday
}

// destructor provides you a graceful shutdown for user defined new day callback
// (by WhenNewDay() method) allowing to pause shutting down until callback is done their work.
func (c *Calendar) destructor() {
	c.wg.Wait()
}

// pendingEventIdx returns the idx of 'event' in the pending event's slice
// in the current Calendar's object or -1 if not found.
func (c *Calendar) pendingEventIdx(event Event) int {

	eventDate := event.Date().ToCmp()
	for i, n := 0, len(c.pendingEvents); i < n; i++ {
		if c.pendingEvents[i].Date().ToCmp() == eventDate {
			return i
		}
	}
	return -1
}
