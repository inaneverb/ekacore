// Copyright © 2021. All rights reserved.
// Author: Ilya Stroy.
// Contacts: iyuryevich@pm.me, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekatime

import (
	"encoding/base64"
	"encoding/binary"
	"errors"

	"github.com/qioalice/ekago/v4/ekamath"
)

// Calendar is a RAM friendly data structure that allows you to keep
// 365 days of some year with flags whether day is day off or a workday,
// store a reason of that and also binary/text encoding/decoding.
//
// WARNING!
// You MUST use NewCalendar() constructor to construct this object.
// If you just instantiate an object it will be considered as invalid,
// and almost all methods will return you an unexpected, bad result.
//
// WARNING!
// Encode/decode operations DOES NOT SUPPORT causing feature (for now).
// It will be fixed in the future.
type Calendar struct {

	// -------------------------------------------------------------
	//   Binary encoding/decoding protocol.
	//   Version: 1.0
	//
	//   [0..3] bytes:  Reserved.
	//   [4..5] bytes:  Field `year`, big endian.
	//   [6] byte:      Field `isLeap`.
	//   [7] byte:      Reserved.
	//   [8..] bytes:   Field `dayOff`. BitSet as binary encoded.
	// -------------------------------------------------------------

	// TODO: Add support of causing feature for encode/decode operations.

	// The year this calendar of.
	year Year

	// Flag whether the current year is leap or not.
	// Fewer computations, more RAM consumption.
	isLeap bool

	// Bitset of days in calendar.
	// 0 means work day, 1 means day off.
	// The index of bit is a day of year.
	dayOff *ekamath.BitSet

	// A set of Event that is used to overwrite default values of calendar.
	// Nil if `enableCause` is false at the NewCalendar() call.
	cause map[uint]EventID

	// A map of EventID's descriptions.
	// Nil if `enableCause` is false at the NewCalendar() call.
	eventDescriptions map[EventID]string
}

var (
	ErrCalendarInvalid             = errors.New("invalid Calendar")
	ErrCalendarInvalidDataToDecode = errors.New("invalid data to decode to Calendar")
)

// ---------------------------------------------------------------------------- //

// IsValid reports whether current Calendar is valid and not malformed.
func (wc *Calendar) IsValid() bool {
	return wc != nil && wc.dayOff != nil
}

// Clear clears the current Calendar marking ALL days as a workdays
// and removing all events.
// Does nothing if current Calendar is nil or malformed.
func (wc *Calendar) Clear() {
	if wc.IsValid() {
		if wc.cause != nil {
			wc.cause = make(map[uint]EventID, _CALENDAR2_CAUSE_DEFAULT_CAPACITY)
			wc.eventDescriptions = make(map[EventID]string, _CALENDAR2_EVENT_DESCRIPTIONS_DEFAULT_CAPACITY)
		}
		wc.dayOff.Clear()
	}
}

// Clone returns a full-copy of the current Calendar.
// Returns nil if current Calendar is nil or malformed.
func (wc *Calendar) Clone() *Calendar {

	if !wc.IsValid() {
		return nil
	}

	cloned := Calendar{
		year:   wc.year,
		dayOff: wc.dayOff.Clone(),
	}

	if wc.cause != nil {
		cloned.cause = make(map[uint]EventID, len(wc.cause))
		for k, v := range wc.cause {
			cloned.cause[k] = v
		}
		cloned.eventDescriptions = make(map[EventID]string, len(wc.eventDescriptions))
		for k, v := range wc.eventDescriptions {
			cloned.eventDescriptions[k] = v
		}
	}

	return &cloned
}

// Year returns a Year this calendar of.
// Returns 0 if current Calendar is nil or malformed.
func (wc *Calendar) Year() Year {
	if !wc.IsValid() {
		return 0
	}
	return wc.year
}

// OverrideDate allows you to change the default (or previous defined)
// day type (day-off / workday) of the provided Date `dd` in the current Calendar.
func (wc *Calendar) OverrideDate(dd Date, isDayOff bool) {
	wc.overrideDate(dd, 0, isDayOff, false)
}

// AddEvent adds a new Event to the Calendar.
// It's the same as just OverrideDate() but provides an ability to set an EventID
// of such overwriting rule.
//
// Requirements:
//   - Calendar is valid and not malformed object,
//   - Causing feature is enabled (`enableCause` being `true` at the NewCalendar() call),
//   - Provided Event (`e`) is valid and belongs to Year, this Calendar of.
//
// Does nothing if any of requirements is failed.
func (wc *Calendar) AddEvent(e Event) {
	if e.IsValid() {
		eventID, dd, isDayOff := e.Split()
		wc.overrideDate(dd, eventID, isDayOff, true) // contains all checks
	}
}

// AddEventDescription adds a new EventID's description to the Calendar.
// Using that you can describe your EventID and figure out event's name.
//
// Requirements:
//   - Calendar is valid and not malformed object,
//   - Causing feature is enabled (`enableCause` being `true` at the NewCalendar() call),
//   - Provided EventID's name (`desc`) is not empty
//
// Does nothing if any of requirements is failed.
func (wc *Calendar) AddEventDescription(eid EventID, desc string) {
	if wc.IsValid() && wc.cause != nil && desc != "" {
		wc.eventDescriptions[eid] = desc
	}
}

// IsDayOff reports whether required day is day off.
// If you have a Date object just call Date.Days() method.
//
// Requirements:
//   - Calendar is valid and not malformed object,
//   - Requested Date is valid and belongs to Year, this Calendar of.
//
// Returns false if requested day is workday or if any of requirements is failed.
func (wc *Calendar) IsDayOff(dd Date) bool {
	return wc.IsValid() &&
		dd.IsValid() &&
		dd.Year() == wc.year &&
		wc.dayOff.IsSetUnsafe(wc.dateToIndex(dd))
}

// NextWorkDay returns a next work day followed by provided day of year.
//
// Requirements:
//   - Calendar is valid and not malformed object,
//   - Requested Date is valid and belongs to Year, this Calendar of.
//
// Returns an invalid date if there's no remaining workdays after requested
// or if any of requirements is failed.
func (wc *Calendar) NextWorkDay(dd Date) Date {
	return wc.nextDay(dd, false)
}

// NextDayOff returns a next day off followed by provided day of year.
//
// Requirements:
//   - Calendar is valid and not malformed object,
//   - Requested Date is valid and belongs to Year, this Calendar of.
//
// Returns an invalid date if there's no remaining days off after requested
// or if any of requirements is failed.
func (wc *Calendar) NextDayOff(dd Date) Date {
	return wc.nextDay(dd, true)
}

// EventOfDate returns an Event because of which the type of the current date is changed.
//
// Requirements:
//   - Calendar is valid and not malformed object,
//   - Causing feature is enabled (`enableCause` being `true` at the NewCalendar() call),
//   - Requested date (`dayOfYear`) is valid and belongs the year, this calendar belongs also to.
//
// Returns an invalid event if there's no registered event with passed date,
// or if any of requirements is failed.
func (wc *Calendar) EventOfDate(dd Date) Event {

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

// DescriptionOfEvent returns an EventID's description (name).
//
// Requirements:
//   - Calendar is valid and not malformed object,
//   - Causing feature is enabled (`enableCause` being `true` at the NewCalendar() call),
//   - Calendar has at least one Event with requested EventID,
//
// Returns an empty string if there's no registered such EventID,
// or if any of requirements is failed.
func (wc *Calendar) DescriptionOfEvent(eid EventID) string {

	if !(wc.IsValid() && wc.cause != nil) {
		return ""
	}

	return wc.eventDescriptions[eid]
}

// WorkDays returns an array of work days of the provided Month.
//
// Requirements:
//   - Calendar is valid and not malformed object,
//   - Requested Month (`m`) is valid.
//
// WARNING:
// It takes quite a lot of time to prepare return data, because of the way
// the data is stored internally. So, if you need to access it often,
// consider caching data (generate output data once per time its still valid
// and then use it later).
//
// Returns an empty set if any of requirements is failed.
func (wc *Calendar) WorkDays(m Month) []Day {
	return wc.daysIn(m, false)
}

// DaysOff returns an array of days off of the provided Month.
//
// See requirements, warnings and return section of WorkDays().
// It works the same way here.
func (wc *Calendar) DaysOff(m Month) []Day {
	return wc.daysIn(m, true)
}

// WorkDaysCount returns a number of working days in the provided Month.
//
// Requirements:
//   - Calendar is valid and not malformed object,
//   - Requested Month (`m`) is valid.
//
// Returns 0 if any of requirements is failed.
func (wc *Calendar) WorkDaysCount(m Month) Days {
	return wc.daysInCount(m, false)
}

// DaysOffCount returns a number of days off in the provided Month.
//
// See requirements, warnings and return section of WorkDaysCount().
// It works the same way here.
func (wc *Calendar) DaysOffCount(m Month) Days {
	return wc.daysInCount(m, true)
}

// ---------------------------------------------------------------------------- //

// MarshalBinary implements BinaryMarshaler interface encoding current Calendar
// in binary form.
//
// It guarantees that if Calendar is valid, the MarshalBinary() cannot fail.
// There's no guarantees about algorithm that will be used to encode/decode.
//
// Requirements:
//   - Calendar is valid and not malformed object,
//   - User MUST NOT modify returned data. If you need it, clone it firstly.
//
// Limitations:
// - Encode/decode operations DOES NOT SUPPORT causing feature.
func (wc *Calendar) MarshalBinary() ([]byte, error) {

	if !wc.IsValid() {
		return nil, ErrCalendarInvalid
	}

	dayOffEncoded, err := wc.dayOff.MarshalBinary()
	if err != nil {
		return nil, ErrCalendarInvalidDataToDecode
	}

	// For more info about binary encode/decode protocol,
	// see Calendar's internal docs (at the Calendar struct declaration).

	buf := make([]byte, len(dayOffEncoded)+8)
	binary.BigEndian.PutUint16(buf[4:], uint16(wc.year))
	if wc.isLeap {
		buf[6] = 1
	}
	copy(buf[8:], dayOffEncoded)

	return buf, nil
}

// UnmarshalBinary implements BinaryUnmarshaler interface decoding provided `data`
// from binary form.
//
// The current Calendar's data will be overwritten by the decoded one
// if decoding operation has been completed successfully.
//
// There's no guarantees about algorithm that will be used to encode/decode.
// Does nothing (and returns nil) if provided `data` is empty.
//
// Requirements:
//   - Provided `data` MUST BE obtained by calling Calendar.MarshalBinary() method.
//   - Provided `data` MUST BE valid, ErrCalendarInvalidDataToDecode returned otherwise.
//   - User MUST NOT use provided `data` after passing to this method. UB otherwise.
//
// Limitations:
// - Encode/decode operations DOES NOT SUPPORT causing feature.
func (wc *Calendar) UnmarshalBinary(data []byte) error {

	// It's ok for Calendar to be invalid - it will be overwritten anyway.
	// But it must be not nil.

	switch {
	case len(data) == 0:
		return nil

	case len(data) < 9:
		return ErrCalendarInvalidDataToDecode

	case wc == nil:
		return ErrCalendarInvalid
	}

	// For more info about binary encode/decode protocol,
	// see Calendar's internal docs (at the Calendar struct declaration).

	wc.cause = nil
	wc.eventDescriptions = nil
	wc.year = Year(binary.BigEndian.Uint16(data[4:]))
	wc.isLeap = data[6] == 1

	if wc.dayOff == nil {
		wc.dayOff = new(ekamath.BitSet)
	}

	return wc.dayOff.UnmarshalBinary(data[8:])
}

// MarshalText implements TextMarshaler interface encoding current Calendar
// in text form.
//
// It guarantees that if Calendar is valid, the MarshalText() cannot fail.
// MarshalText guarantees that output data will be base64 encoded.
//
// Requirements:
//   - Calendar is valid and not malformed object,
//   - User MUST NOT modify returned data. If you need it, clone it firstly.
//
// Limitations:
// - Encode/decode operations DOES NOT SUPPORT causing feature.
// - Provided base64 data is NO URL FRIENDLY!
func (wc *Calendar) MarshalText() ([]byte, error) {

	binaryEncodedData, err := wc.MarshalBinary()
	if err != nil {
		return nil, err
	}

	buf := make([]byte, base64.StdEncoding.EncodedLen(len(binaryEncodedData)))
	base64.StdEncoding.Encode(buf, binaryEncodedData)

	return buf, nil
}

// UnmarshalText implements TextUnmarshaler interface decoding provided `data`
// from text form.
//
// The current Calendar's data will be overwritten by the decoded one
// if decoding operation has been completed successfully.
//
// Does nothing (and returns nil) if provided `data` is empty.
//
// Requirements:
//   - Provided `data` MUST BE obtained by calling Calendar.MarshalBinary() method.
//   - Provided `data` MUST BE valid, ErrCalendarInvalidDataToDecode returned otherwise.
//   - User MUST NOT use provided `data` after passing to this method. UB otherwise.
//
// Limitations:
// - Encode/decode operations DOES NOT SUPPORT causing feature.
func (wc *Calendar) UnmarshalText(data []byte) error {

	switch {
	case len(data) == 0:
		return nil

	case wc == nil:
		return ErrCalendarInvalid
	}

	buf := make([]byte, base64.StdEncoding.DecodedLen(len(data)))
	n, err := base64.StdEncoding.Decode(buf, data)
	if err != nil {
		return err
	}

	return wc.UnmarshalBinary(buf[:n])
}

// ---------------------------------------------------------------------------- //

// NewCalendar is a Calendar constructor.
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
// and then just call Calendar.Clone() if you need many Calendar objects
// for the same year.
// For configuration above the cloning without `enableCause` feature (read later)
// it takes just 95ns. Its faster than filling each object up to x9 times.
//
// If `enableCause` is true, it also pre-allocates ~512 bytes to be able to keep
// 64+ reasons of when the default type of specific day is changed.
// If you don't need that (and EventOfDay(), EventOfDate() methods), just pass `false`.
//
// WARNING!
// Encode/decode operations DOES NOT SUPPORT causing feature.
func NewCalendar(y Year, saturdayAndSunday, enableCause bool) *Calendar {

	if !IsValidDate(y, MONTH_JANUARY, 1) {
		return nil
	}

	wc := Calendar{
		year:   y,
		isLeap: y.IsLeap(),
		dayOff: ekamath.NewBitSet(_CALENDAR2_DEFAULT_CAPACITY),
	}

	if enableCause {
		wc.cause = make(map[uint]EventID, _CALENDAR2_CAUSE_DEFAULT_CAPACITY)
		wc.eventDescriptions = make(map[EventID]string, _CALENDAR2_EVENT_DESCRIPTIONS_DEFAULT_CAPACITY)
	}

	if saturdayAndSunday {
		wc.doSaturdayAndSundayDayOff()
	}

	return &wc
}
