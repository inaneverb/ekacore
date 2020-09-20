// Copyright © 2020. All rights reserved.
// Author: Ilya Stroy.
// Contacts: qioalice@gmail.com, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekatime

import (
	"github.com/qioalice/ekago/v2/ekaerr"
)

type (
	// Today is a very useful part of Calendar.
	// Honestly, it's a main feature and power of Calendar.
	//
	// When you call Calendar.Run(), the Calendar allows you to use Calendar.Today()
	// method that returns cached Today object. This object.
	// The Today object updates only once in day (at the midnight), and then
	// you will always get a cached pointer of that object.
	//
	// Thus, it's blazing fast (and ofc thread-safety, but only RO).
	// You may get an info about current day, about current workday according
	// with all Event s loaded into Calendar.
	//
	// Moreover, at the Today object creating, its also is encoded to at least JSON
	// (and you may set an encoder using Calendar.RegJsonEncoder()).
	// It's also is encoded once in day and available at the AsJson field.
	// (if encoded failed, check AsJsonErr field).
	//
	// Moreover, you may use your own custom encoder.
	// Check also Calendar.RegYourOwnEncoder() and corresponding Today's fields
	// E.g: AsYourOwn1, AsYourOwn1Err, etc.
	//
	// WARNING!
	// THE TIME PRESENTED BY THIS OBJECT IS FROZEN AT THE OBJECT CREATION TIME,
	// SO IT ALWAYS BE ABOUT 00:00:XX AM.
	// IT IS AVAILABLE ONLY FOR HISTORICAL PURPOSES AND MAY BE REMOVED!
	//
	// WARNING!
	// DO NOT UPDATE THIS OBJECT FIELDS MANUALLY!
	// DATA RACE CAUTION! THREAD SAFETY ISSUES CAUTION!
	// REGARDLESS OF THE FACT THAT YOU GETTING THIS OBJECT BY A POINTER
	// (INSTEAD OF BY VALUE TO DECREASE UNNECESSARY COPYING),
	// YOU SHALL NOT MODIFYING IT!
	// POSSIBLE DATA CORRUPTION AND UNDEFINED BEHAVIOUR INSTEAD.
	//
	// If you want to take a copy of the Today object, take your look at the
	// Today.Copy() or Today.CopyWithEncodedData() methods. That's what you need,
	// if you wanna modify Today object.
	Today struct {

		// TODO: Implement new methods:
		//  - ApplyEvents(events): Recalculates Today object
		//    (make a copy if necessary - c.today addr check),
		//    applying only passed events instead of using Calendar's ones.
		//    Allow to work with nil as Calendar.
		//  - ApplyDefaultEventsAnd(events): The same as ApplyEvents(),
		//    but does not overwrites default Calendar events,
		//    but adds passed events to them.
		//    If Calendar is not presented, does the same thing as ApplyEvents().

		// Parent Calendar object. Could be nil if Today is instantiated manually.
		c                *Calendar     `json:"-"`

		// Deprecated: Does not represent an actual Time.
		Timestamp        Timestamp     `json:"ts"`

		// Deprecated: Does not represent an actual Time.
		Time             Time          `json:"-"`

		Date             Date          `json:"-"`

		Year             Year          `json:"year"`
		Month            Month         `json:"month"`
		Day              Day           `json:"day"`

		Weekday          Weekday       `json:"weekday"`

		// Deprecated: Does not represent an actual Hour
		Hour             Hour          `json:"-"`

		// Deprecated: Does not represent an actual Minute
		Minute           Minute        `json:"-"`

		// Deprecated: Does not represent an actual Second
		Second           Second        `json:"-"`

		WorkDayCurrent   Day           `json:"work_day_current"`
		IsDayOff         bool          `json:"is_dayoff"`

		DaysInMonth      Day           `json:"days_in_month"`
		WorkDayTotal     Day           `json:"work_day_total"` // in month
		DayOffTotal      Day           `json:"dayoff_total"` // in month

		WorkDays         []Day         `json:"work_days"` // days in month
		DayOffs          []Day         `json:"dayoffs"` // days in month

		AsJson           []byte        `json:"-"`
		AsJsonErr        *ekaerr.Error `json:"-"`

		AsYourOwn1       []byte        `json:"-"`
		AsYourOwn1Err    *ekaerr.Error `json:"-"`
	}

	// TodayEncoder is a function alias that represents a function that
	// encodes Today object to the some data.
	// Used at the Calendar.RegJsonEncoder(), Calendar.RegYourOwnEncoder()
	// to register a custom encoders of Today object, the output data of will be stored
	// into Today's corresponding fields.
	TodayEncoder func(today *Today) ([]byte, *ekaerr.Error)
)

// Copy returns the current Today's object copy but with the same encoded data
// (only pointers are copied, not internal data).
// If you want to copy them too, use CopyWithEncodedData() instead.
func (t *Today) Copy() *Today {
	cp := *t
	return &cp
}

// CopyWithEncodedData returns the current Today's object copy even with encoded data.
// Yes, a new buffers are allocated and the encoded data explicitly copied.
func (t *Today) CopyWithEncodedData() *Today {
	cp := t.Copy()
	// https://github.com/go101/go101/wiki/How-to-perfectly-clone-a-slice
	cp.AsJson = append(t.AsJson[:0:0], t.AsJson...)
	cp.AsYourOwn1 = append(t.AsYourOwn1[:0:0], t.AsYourOwn1...)
	return cp
}

// Cal returns the Calendar, the current Today is generated from.
func (t *Today) Cal() *Calendar {
	return t.c
}
