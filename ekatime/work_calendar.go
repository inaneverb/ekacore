// Copyright © 2021. All rights reserved.
// Author: Ilya Stroy.
// Contacts: qioalice@gmail.com, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekatime

import (
	"github.com/qioalice/ekago/v3/ekamath"
)

type (
	// WorkCalendar is a RAM friendly data structure that allows you to keep
	// a 365 days of some year with flags whether day is day off or a workday,
	// store a reason of that and also binary/text encoding/decoding.
	//
	// WARNING!
	// You MUST use NewWorkCalendar() constructor to construct this object.
	// If you just instantiate an object it will be considered as invalid,
	// and almost all methods will return you an unexpected, bad result.
	WorkCalendar struct {

		// The year this calendar of.
		year Year

		// Flag whether the current year is leap or not.
		// Less computations, more RAM consumption.
		isLeap bool

		// Bitset of days in calendar.
		// 0 means work day, 1 means day off.
		// The index of bit is a day of year.
		dayOff *ekamath.BitSet

		// A set of Event that is used to overwrite default values of calendar.
		// Nil if `enableCause` is false at the NewWorkCalendar() call.
		cause map[uint]EventID
	}
)

// IsValid reports whether current WorkCalendar is valid and not malformed.
func (wc *WorkCalendar) IsValid() bool {
	return wc != nil && wc.dayOff != nil
}

// Clear clears the current WorkCalendar marking ALL days as a workdays
// and removing all events.
// Does nothing if current WorkCalendar is nil or malformed.
func (wc *WorkCalendar) Clear() {
	if wc.IsValid() {
		if wc.cause != nil {
			wc.cause = make(map[uint]EventID, _CALENDAR2_CAUSE_DEFAULT_CAPACITY)
		}
		wc.dayOff.Clear()
	}
}

// Clone returns a full-copy of the current WorkCalendar.
// Returns nil if current WorkCalendar is nil or malformed.
func (wc *WorkCalendar) Clone() *WorkCalendar {

	if !wc.IsValid() {
		return nil
	}

	cloned := WorkCalendar{
		year:   wc.year,
		dayOff: wc.dayOff.Clone(),
	}

	if wc.cause != nil {
		cloned.cause = make(map[uint]EventID, len(wc.cause))
		for k, v := range wc.cause {
			cloned.cause[k] = v
		}
	}

	return &cloned
}

// Year returns a Year this calendar of.
// Returns 0 if current WorkCalendar is nil or malformed.
func (wc *WorkCalendar) Year() Year {
	if !wc.IsValid() {
		return 0
	}
	return wc.year
}

// OverrideDate allows you to change the default (or previous defined)
// day type (day-off / workday) of the provided Date `dd` in the current WorkCalendar.
func (wc *WorkCalendar) OverrideDate(dd Date, isDayOff bool) {
	wc.overrideDate(dd, 0, isDayOff, false)
}

// AddEvent adds a new Event to the WorkCalendar.
// It's the same as just OverrideDate() but provides an ability to set an EventID
// of such overwriting rule.
//
// Requirements:
//
//  - WorkCalendar is valid and not malformed object,
//  - Causing feature is enabled (`enableCause` being `true` at the NewWorkCalendar() call),
//  - Provided Event (`e`) is valid and belongs to Year, this WorkCalendar of.
//
// Does nothing if any of requirements is failed.
func (wc *WorkCalendar) AddEvent(e Event) {
	if e.IsValid() {
		eventID, dd, isDayOff := e.Split()
		wc.overrideDate(dd, eventID, isDayOff, true)
	}
}

// IsDayOff reports whether required day is day off.
// If you have a Date object just call Date.Days() method.
//
// Requirements:
//
//  - WorkCalendar is valid and not malformed object,
//  - Requested Date is valid and belongs to Year, this WorkCalendar of.
//
// Returns false if requested day is workday or if any of requirements is failed.
func (wc *WorkCalendar) IsDayOff(dd Date) bool {
	return wc.IsValid() &&
		dd.IsValid() &&
		dd.Year() == wc.year &&
		wc.dayOff.IsSetUnsafe(wc.dateToIndex(dd))
}

// NextWorkDay returns a next work day followed by provided day of year.
//
// Requirements:
//
//  - WorkCalendar is valid and not malformed object,
//  - Requested Date is valid and belongs to Year, this WorkCalendar of.
//
// Returns an invalid date if there's no remaining workdays after requested
// or if any of requirements is failed.
func (wc *WorkCalendar) NextWorkDay(dd Date) Date {
	return wc.nextDay(dd, false)
}

// NextDayOff returns a next day off followed by provided day of year.
//
// Requirements:
//
//  - WorkCalendar is valid and not malformed object,
//  - Requested Date is valid and belongs to Year, this WorkCalendar of.
//
// Returns an invalid date if there's no remaining days off after requested
// or if any of requirements is failed.
func (wc *WorkCalendar) NextDayOff(dd Date) Date {
	return wc.nextDay(dd, true)
}

// EventOfDate returns an Event because of which the type of the current date is changed.
//
// Requirements:
//
//  - WorkCalendar is valid and not malformed object,
//  - Causing feature is enabled (`enableCause` being `true` at the NewWorkCalendar() call),
//  - Requested date (`dayOfYear`) is valid and belongs the year, this calendar belongs also to.
//
// Returns an invalid event if there's no registered event with passed date,
// or if any of requirements is failed.
func (wc *WorkCalendar) EventOfDate(dd Date) Event {

	if !(wc.IsValid() && wc.cause != nil && dd.IsValid() && dd.Year() == wc.year) {
		return _EVENT_INVALID
	}

	idx := wc.dateToIndex(dd)
	eventID, ok := wc.cause[idx]
	if !ok {
		return _EVENT_INVALID
	}

	isDayOff := wc.dayOff.IsSetUnsafe(idx)

	return NewEvent(dd, eventID, isDayOff)
}

// WorkDays returns an array of work days of the provided Month.
//
// Requirements:
//
//  - WorkCalendar is valid and not malformed object,
//  - Requested Month (`m`) is valid.
//
// WARNING:
// It takes quite a lot of time to prepare return data, because of the way
// the data is stored internally. So, if you need to access it often,
// consider caching data (generate output data once per time its still valid
// and then use it later).
//
// Returns an empty set if any of requirements is failed.
func (wc *WorkCalendar) WorkDays(m Month) []Day {
	return wc.daysIn(m, false)
}

// DaysOff returns an array of days off of the provided Month.
//
// See requirements, warnings and return section of WorkDays().
// It works the same way here.
func (wc *WorkCalendar) DaysOff(m Month) []Day {
	return wc.daysIn(m, true)
}

// WorkDaysCount returns a number of working days in the provided Month.
//
// Requirements:
//
//  - WorkCalendar is valid and not malformed object,
//  - Requested Month (`m`) is valid.
//
// Returns 0 if any of requirements is failed.
func (wc *WorkCalendar) WorkDaysCount(m Month) Days {
	return wc.daysInCount(m, false)
}

// DaysOffCount returns a number of days off in the provided Month.
//
// See requirements, warnings and return section of WorkDaysCount().
// It works the same way here.
func (wc *WorkCalendar) DaysOffCount(m Month) Days {
	return wc.daysInCount(m, true)
}

// NewWorkCalendar is a WorkCalendar constructor.
// Returns an initialized, ready to use object.
//
// You MUST specify a valid Year, otherwise nil is returned.
//
// It's allowed to pass Year < 1900 or > 4095
// (that Year for which Year.IsValid() method will return false).
//
// If `saturdayAndSunday` is true, these days will be marked as days off.
// Keep in mind, that marking all saturdays and sundays as days off is quite heavy op.
// It takes 425ns for i7-9750H CPU @ 2.60GHz.
// Maybe it will be better for you to generate once a "template" of that
// and then just call WorkCalendar.Clone() if you need many WorkCalendar objects
// for the same year.
// For configuration above the cloning without `enableCause` feature (read later)
// it takes just 95ns. Its faster than filling each object up to x9 times.
//
// If `enableCause` is true, it also pre-allocates 512 bytes to be able to keep
// 64+ reasons of when the default type of specific day is changed.
// If you don't need that (and EventOfDay(), EventOfDate() methods), just pass `false`.
func NewWorkCalendar(y Year, saturdayAndSunday, enableCause bool) *WorkCalendar {

	if !IsValidDate(y, MONTH_JANUARY, 1) {
		return nil
	}

	wc := WorkCalendar{
		year:   y,
		isLeap: y.IsLeap(),
		dayOff: ekamath.NewBitSet(_CALENDAR2_DEFAULT_CAPACITY),
		cause:  nil,
	}

	if enableCause {
		wc.cause = make(map[uint]EventID, _CALENDAR2_CAUSE_DEFAULT_CAPACITY)
	}

	if saturdayAndSunday {
		wc.doSaturdayAndSundayDayOff()
	}

	return &wc
}