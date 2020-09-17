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
	"github.com/qioalice/ekago/v2/ekamath"
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

	newToday.Date, newToday.Time                    = newToday.Timestamp.Split()
	newToday.Year, newToday.Month, newToday.Day     = newToday.Date.Split()
	newToday.Hour, newToday.Minute, newToday.Second = newToday.Time.Split()
	newToday.Weekday                                = newToday.Date.Weekday()
	newToday.DaysInMonth                            = newToday.Date.DaysInMonth()

	newToday.WorkDays = make([]Day, 0, 31)
	newToday.DayOffs = make([]Day, 0, 31)

	newToday.WorkDayCurrent, newToday.WorkDayTotal, newToday.IsDayOff =
		workdaysFor(
			NewDate(newToday.Year, newToday.Month, 1),
			newToday.Day,
			c.confirmedEvents,
			&newToday.WorkDays,
		)

	// newToday.WorkDays already filled. Fill newToday.DayOffs.
	var daysBitSet ekamath.Flags32
	for _, d := range newToday.WorkDays {
		daysBitSet |= 1 << (d-1)
	}
	// Wanna day offs? No problem. XOR is your friend.
	daysBitSet ^= 0xFF_FF_FF_FF
	for i := Day(0); i < newToday.DaysInMonth; i++ {
		if daysBitSet & (1 << i) != 0 {
			newToday.DayOffs = append(newToday.DayOffs, i+1)
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

// See Calendar.WorkdaysFor docs.
// isDayOff reports whether d1 is day off.
// workDaysOut will contain working days if it's not nil.
func workdaysFor(dd Date, d1 Day, events []Event, workDaysOut *[]Day) (current, total Day, isDayOff bool) {

	y, m, d := normalizeDate(dd.Split())
	dd = NewDate(y, m, d)

	daysInMonth := dd.DaysInMonth()

	switch {
	case d1 < d:
		d1 = d
	case d1 > daysInMonth:
		d1 = daysInMonth
	}

	var workDaysBitSet ekamath.Flags32

	for d, wd := d, dd.Weekday(); d <= daysInMonth; d++ {
		if !wd.IsDayOff() {
			total++
			workDaysBitSet |= 1 << (d-1)
			if d <= d1 {
				current++
			}
		}
		wd = wd.Next()
	}

	if d1 != d {
		isDayOff = NewDate(y, m, d1).Weekday().IsDayOff()
	} else {
		isDayOff = dd.Weekday().IsDayOff()
	}

	for _, ce := range events {
		if ce.Year() != y || ce.Month() != m || ce.Day() < d {
			continue
		}
		if ce.Day() == d1 {
			isDayOff = ce.IsDayOff()
		}
		if ce.Weekday().IsDayOff() == ce.IsDayOff() {
			continue
		}
		if ce.IsDayOff() {
			total--
			workDaysBitSet &^= 1 << (d-1)
			if ce.Day() <= d1 {
				current--
			}
		} else {
			total++
			workDaysBitSet |= 1 << (d-1)
			if ce.Day() <= d1 {
				current++
			}
		}
	}

	if workDaysOut != nil {
		// Represent bit set as array
		for i := Day(0); i < 31; i++ {
			if workDaysBitSet & (1 << i) != 0 {
				*workDaysOut = append(*workDaysOut, i+1)
			}
		}
	}

	return current, total, isDayOff
}