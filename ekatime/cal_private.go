// Copyright © 2020. All rights reserved.
// Author: Ilya Stroy.
// Contacts: qioalice@gmail.com, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekatime

import (
	"sync/atomic"
	"time"
	"unsafe"
)

var (
	calDummyPointer unsafe.Pointer
)

// newDayHasCome is a special function that is called when new day (00:00 AM) has come.
// Called once for each new day. Updates timer for being called one more time (next day),
// and calls registered callback by WhenNewDay() method.
func (c *Calendar) newDayHasCome() {

	newToday := c.updateToday()
	c.newDayTimer.Reset(TillNextMidnight())

	if c.newDayCallback != nil {
		c.newDayCallback(newToday)
	}
}

// updateToday creates a new Today object, fills them, stores them into the
// Calendar's internal fast access cache and returns.
func (c *Calendar) updateToday() *Today {

	newToday := new(Today)

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

	for i, n := 0, len(c.events); i < n; i++ {
		if c.events[i].Year() == newToday.Year && c.events[i].Month() == newToday.Month {

			if c.events[i].Day() == newToday.Day {
				newToday.IsDayOff = c.events[i].IsDayOff()
			}

			if c.events[i].Day() <= newToday.Day &&
				c.events[i].Weekday().IsDayOff() != c.events[i].IsDayOff() {

				if c.events[i].IsDayOff() {
					newToday.WorkDayCurrent--
					newToday.WorkDayTotal--
				} else {
					newToday.WorkDayCurrent++
					newToday.WorkDayTotal++
				}
			}
		}
	}

	newToday.DayOffTotal = newToday.DaysInMonth - newToday.WorkDayTotal

	if c.todayEncoderJson != nil {
		newToday.AsJson = c.todayEncoderJson(newToday)
	}
	if c.todayEncoderUsers1 != nil {
		newToday.AsYourOwn1 = c.todayEncoderUsers1(newToday)
	}

	atomic.StorePointer(&c.today, unsafe.Pointer(newToday))
	return newToday
}

//
func (c *Calendar) getToday() *Today {

	if c == nil {
		return nil
	}

	switch today := atomic.LoadPointer(&c.today); {

	case today == nil &&
		atomic.CompareAndSwapPointer(&c.today, calDummyPointer, unsafe.Pointer(c)):
		// We'll do initialize.
		return c.updateToday()

	case today == nil || today == unsafe.Pointer(c):
		// Another one goroutine does initializing.
		time.Sleep(10 * time.Microsecond)
		return c.getToday()

	default:
		// All is good. Today is ready.
		return (*Today)(today)
	}
}

//
func (c *Calendar) eventIdx(event Event) int {
	eventDate := event.Date().ToCmp()
	for i, n := 0, len(c.events); i < n; i++ {
		if c.events[i].Date().ToCmp() == eventDate {
			return i
		}
	}
	return -1
}

// initCal initializes dummy pointer that is used to check whether Calendar
// is initialized or not.
func initCal() {
	calDummyPointer = unsafe.Pointer(new(int))
}