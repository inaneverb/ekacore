// Copyright © 2020. All rights reserved.
// Author: Ilya Stroy.
// Contacts: qioalice@gmail.com, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekatime

import (
	"fmt"
)

type (
	// EventID is an alias to Golang unit type
	// that is used to represents an Event's ID.
	EventID uint16

	// Event represents some unusual days - events.
	// Like holidays, personal vacation days or work day off, etc.
	//
	// Contains:
	// - Day (use Day() method, values: 1-31),
	// - Month (use Month() method, values: 1-12),
	// - Year (use Year() method, values: 0-4095),
	// - Working? (use IsWorkday() or IsDayOff() method),
	// - Event type's ID (use ID() method, values: 0-65535).
	//
	// Supports up to 127 types you can bind an Event to. Hope it's enough.
	// Takes just 8 byte. At all! Thanks to bitwise operations.
	Event uint64
)

func (e Event) IsValid() bool {
	return e != _EVENT_INVALID && e.Date().IsValid()
}

//
func (e Event) Year() Year {
	return e.Date().Year()
}

//
func (e Event) Month() Month {
	return e.Date().Month()
}

//
func (e Event) Day() Day {
	return e.Date().Day()
}

//
func (e Event) Weekday() Weekday {
	return e.Date().Weekday()
}

//
func (e Event) Date() Date {
	return Date(e>>_EVENT_OFFSET_DATE) & _DATE_MASK_DATE
}

//
func (e Event) IsWorkday() bool {
	return !e.IsDayOff()
}

//
func (e Event) IsDayOff() bool {
	return uint8((e>>_EVENT_OFFSET_IS_WORKDAY)&_EVENT_MASK_IS_WORKDAY) > 0
}

//
func (e Event) ID() EventID {
	return EventID(e >> _EVENT_OFFSET_ID & _EVENT_MASK_ID)
}

//
func (e Event) Split() (id EventID, dd Date, isDayOff bool) {
	return e.ID(), e.Date(), e.IsDayOff()
}

//
func NewEvent(d Date, id EventID, isDayOff bool) Event {

	isDayOffBit := Event(0)
	if isDayOff {
		isDayOffBit = Event(1) << _EVENT_OFFSET_IS_WORKDAY
	}

	d = d.ensureWeekdayExist()

	return Event(d&_DATE_MASK_DATE)<<_EVENT_OFFSET_DATE | isDayOffBit |
		(Event(id)&_EVENT_MASK_ID)<<_EVENT_OFFSET_ID
}

// String returns the current Event's string representation in the following format:
// "YYYY/MM/DD [<isWorkDay>] ID: <id>", where:
//     <id> is Event's ID as is in the range [1..127],
//     <isWorkDay> is either "Workday" or "Dayoff"
func (e Event) String() string {
	s := "Workday"
	if e.IsDayOff() {
		s = "Dayoff"
	}
	return fmt.Sprintf("%04d/%02d/%02d [%s] ID: %d", e.Year(), e.Month(), e.Day(), s, e.ID())
}
