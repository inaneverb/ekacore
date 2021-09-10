package ekatime

import (
	"github.com/bits-and-blooms/bitset"
)

type (
	// WorkCalendar is a RAM friendly data structure that allows you to keep
	// a 365 days of some year with flags whether day is day off or a workday,
	// store a reason of that and also binary/text encoding/decoding.
	WorkCalendar struct {

		// The year this calendar of.
		year Year

		// Bitset of days in calendar.
		// 0 means work day, 1 means day off.
		// The index of bit is a day of year.
		dayOff *bitset.BitSet

		// A set of Event that is used to overwrite default values of calendar.
		// Nil if `enableCause` is false at the NewWorkCalendar() call.
		cause map[Days]EventID
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
			wc.cause = make(map[Days]EventID, _CALENDAR2_CAUSE_DEFAULT_CAPACITY)
		}
		wc.dayOff.ClearAll()
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
		cloned.cause = make(map[Days]EventID, len(wc.cause))
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

// Reports whether required day is day off.
// If you have a Date object just call Date.Days() method.
//
// Requirements:
//
//  - WorkCalendar is valid and not malformed object,
//  - Requested day of year (`dayOfYear`) belongs the range [1..365/366].
//
// Returns false if requested day is workday or if any of requirements is failed.
func (wc *WorkCalendar) IsDayOffDay(dayOfYear Days) bool {
	return wc.IsValid() && dayOfYear.BelongsToYear(wc.year) && wc.dayOff.Test(uint(dayOfYear))
}

// Returns a next work day followed by provided day of year.
// If you have a Date object just call Date.Days() method.
//
// Requirements:
//
//  - WorkCalendar is valid and not malformed object,
//  - Requested day of year (`dayOfYear`) belongs the range [1..365/366].
//
// Returns an invalid date if there's no remaining workdays after requested
// or if any of requirements is failed.
func (wc *WorkCalendar) NextWorkDay(dayOfYear Days) Date {
	return wc.nextDay(dayOfYear, false)
}

// Returns a next day off followed by provided day of year.
// If you have a Date object just call Date.Days() method.
//
// Requirements:
//
//  - WorkCalendar is valid and not malformed object,
//  - Requested day of year (`dayOfYear`) belongs the range [1..365/366].
//
// Returns an invalid date if there's no remaining days off after requested
// or if any of requirements is failed.
func (wc *WorkCalendar) NextDayOff(dayOfYear Days) Date {
	return wc.nextDay(dayOfYear, true)
}

// EventOfDay returns an Event because of which the type of the current day is changed.
//
// Requirements:
//
//  - WorkCalendar is valid and not malformed object,
//  - Causing feature is enabled (`enableCause` being `true` at the NewWorkCalendar() call),
//  - Requested day of year (`dayOfYear`) belongs the range [1..365/366].
//
// Returns an invalid event if there's no registered event with passed day of year,
// or if any of requirements is failed.
func (wc *WorkCalendar) EventOfDay(dayOfYear Days) Event {

	if !(wc.IsValid() && wc.cause != nil) {
		return _EVENT_INVALID
	}

	eventID, ok := wc.cause[dayOfYear]
	if !ok {
		return _EVENT_INVALID
	}

	dd := NewDateFromDays(wc.year, dayOfYear)
	isDayOff := wc.dayOff.Test(uint(dayOfYear))

	return NewEvent(dd, eventID, isDayOff)
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

	dayOfYear := dd.Days()

	eventID, ok := wc.cause[dayOfYear]
	if !ok {
		return _EVENT_INVALID
	}

	isDayOff := wc.dayOff.Test(uint(dayOfYear))

	return NewEvent(dd, eventID, isDayOff)
}

// NewWorkCalendar is a WorkCalendar constructor.
// Returns an initialized, ready to use object.
//
// You MUST specify a valid Year, otherwise nil is returned.
//
// It's allowed to pass Year < 1900 or > 4095
// (that Year for which IsValid() method will return false).
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

	var cause map[Days]EventID = nil
	if enableCause {
		cause = make(map[Days]EventID, _CALENDAR2_CAUSE_DEFAULT_CAPACITY)
	}

	wc := WorkCalendar{
		year:   y,
		dayOff: bitset.New(_CALENDAR2_DEFAULT_CAPACITY),
		cause:  cause,
	}

	if saturdayAndSunday {
		wc.doSaturdayAndSundayDayOff()
	}

	return &wc
}
