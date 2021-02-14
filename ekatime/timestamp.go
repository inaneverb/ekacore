// Copyright © 2020. All rights reserved.
// Author: Ilya Stroy.
// Contacts: qioalice@gmail.com, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekatime

import (
	"time"
)

type (
	// Timestamp is just unix timestamp type (stores a seconds that has been passed
	// since 01 January 1970 00:00:00 (12:00:00 AM)).
	Timestamp int64

	// TimestampPair is just two some unix timestamps.
	// Used in those methods that returns 2 timestamps.
	TimestampPair [2]Timestamp
)

// I64 returns int64 representation of the current Timestamp 'ts.
func (ts Timestamp) I64() int64 {
	return int64(ts)
}

// Std returns standard Golang's time.Time object with the same values
// as current Timestamp have.
func (ts Timestamp) Std() time.Time {
	return time.Unix(ts.I64(), 0).UTC()
}

// Split splits the current TimestampPair 'tsp' into two separate Timestamps.
func (tsp TimestampPair) Split() (Timestamp, Timestamp) {
	return tsp[0], tsp[1]
}

// I64 returns two int64 values which are int64 representation of each value
// in the current TimestampPair 'tsp'.
func (tsp TimestampPair) I64() (int64, int64) {
	return int64(tsp[0]), int64(tsp[1])
}

// Date returns the Date object, the current Timestamp includes which.
// Returns 0/0/0 Date if Timestamp is 0 (01 Jan 1970 00:00:00).
func (ts Timestamp) Date() Date {
	if ts == 0 {
		return 0
	}
	return NewDate(dateFromUnix(ts)) | ts.Weekday().asPartOfDate()
}

// Time returns the Time object, the current Timestamp includes which.
func (ts Timestamp) Time() Time {
	if ts == 0 {
		return 0
	}
	return NewTime(timeFromUnix(ts))
}

// Year returns the year number, the current Timestamp includes which.
//
// If you need also at least one more Date's parameter (Month or Day),
// avoid to call many methods explicitly. Call Date() instead (and Split() then).
func (ts Timestamp) Year() Year {
	y, _, _ := dateFromUnix(ts)
	return y
}

// Month returns the month number, the current Timestamp includes which.
//
// If you need also at least one more Date's parameter (Year or Day),
// avoid to call many methods explicitly. Call Date() instead (and Split() then).
func (ts Timestamp) Month() Month {
	_, m, _ := dateFromUnix(ts)
	return m
}

// Day returns the day number, the current Timestamp includes which.
//
// If you need also at least one more Date's parameter (Month or Year),
// avoid to call many methods explicitly. Call Date() instead (and Split() then).
func (ts Timestamp) Day() Day {
	_, _, d := dateFromUnix(ts)
	return d
}

// Hour returns the hour number, the current Timestamp includes which.
//
// If you need also at least one more Time's parameter (Minute or Second),
// avoid to call many methods explicitly. Call Time() instead (and Split() then).
func (ts Timestamp) Hour() Hour {
	h, _, _ := timeFromUnix(ts)
	return h
}

// Minute returns the minute number, the current Timestamp includes which.
//
// If you need also at least one more Time's parameter (Hour or Second),
// avoid to call many methods explicitly. Call Time() instead (and Split() then).
func (ts Timestamp) Minute() Minute {
	_, m, _ := timeFromUnix(ts)
	return m
}

// Second returns the second number, the current Timestamp includes which.
//
// If you need also at least one more Time's parameter (Minute or Hour),
// avoid to call many methods explicitly. Call Time() instead (and Split() then).
func (ts Timestamp) Second() Second {
	_, _, s := timeFromUnix(ts)
	return s
}

// Split returns the Date and Time, the current Timestamp includes which.
// It's just like a separate Date(), Time() calls.
func (ts Timestamp) Split() (d Date, t Time) {
	return ts.Date(), ts.Time()
}

// Weekday returns the current Timestamp ts the number of day in week.
func (ts Timestamp) Weekday() Weekday {
	return Weekday(((ts + SECONDS_IN_DAY) % SECONDS_IN_WEEK) / SECONDS_IN_DAY)
}

// Now is just the same as time.Now().
func Now() Timestamp {
	return UnixFromStd(time.Now())
}

// UnixFrom creates and returns Timestamp object from the presented Date 'd'
// and Time 't'.
//
// WARNING!
// It's 0 value Timestamp if you use 01 Jan 1970 00:00:00,
// and Timestamp.MarshalJSON() will return JSON null for that value,
// and Timestamp.Date() will return 0/0/0 Date, NOT 01 Jan 1970!
// To avoid it, use 00:00:01 as Time.
func UnixFrom(y Year, m Month, d Day, hh Hour, mm Minute, ss Second) Timestamp {
	if y > 4095 {
		y = 4095
	}
	tt := time.Date(int(y), time.Month(m), int(d), int(hh), int(mm), int(ss), 0, time.UTC)
	return Timestamp(tt.Unix())
}

// UnixFromStd creates and returns Timestamp object from the standard Golang's
// time.Time object (UTC time).
func UnixFromStd(t time.Time) Timestamp {
	return Timestamp(t.Unix())
}

// BeginningOfDay returns the day beginning of the current timestamp 'ts'.
// E.g: 12/11/2019, 15:46:40 (3:46:40 PM) -> 12/11/2019 00:00:00 (12:00:00 AM).
func (ts Timestamp) BeginningOfDay() Timestamp {
	return ts.beginningAndEndOf(SECONDS_IN_DAY)[0]
}

// EndOfDay returns the day ending of the current timestamp 'ts'.
// E.g: 12/11/2019, 15:46:40 (3:46:40 PM) -> 12/11/2019 23:59:59 (11:59:59 PM).
func (ts Timestamp) EndOfDay() Timestamp {
	return ts.beginningAndEndOf(SECONDS_IN_DAY)[1]
}

// BeginningAndEndOfDay is the same as BeginningOfDay() and EndOfDay() calls.
func (ts Timestamp) BeginningAndEndOfDay() TimestampPair {
	return ts.beginningAndEndOf(SECONDS_IN_DAY)
}

// BeginningOfMonth returns the month beginning of the current timestamp 'ts'.
// E.g: 12/11/2019, 15:46:40 (3:46:40 PM) -> 1/11/2019 00:00:00 (12:00:00 AM).
func (ts Timestamp) BeginningOfMonth() Timestamp {
	y, m, _ := dateFromUnix(ts)
	return ts.beginningAndEndOf(InMonth(y, m))[0]
}

// EndOfMonth returns the month ending of the current timestamp 'ts'.
// E.g: 12/11/2019, 15:46:40 (3:46:40 PM) -> 30/11/2019 23:59:59 (11:59:59 PM).
func (ts Timestamp) EndOfMonth() Timestamp {
	y, m, _ := dateFromUnix(ts)
	return ts.beginningAndEndOf(InMonth(y, m))[1]
}

// BeginningAndEndOfMonth is the same as BeginningOfMonth() and EndOfMonth() calls.
func (ts Timestamp) BeginningAndEndOfMonth() TimestampPair {
	y, m, _ := dateFromUnix(ts)
	return ts.beginningAndEndOf(InMonth(y, m))
}

// BeginningOfYear returns the year beginning of the current timestamp 'ts'.
// E.g: 12/11/2019, 15:46:40 (3:46:40 PM) -> 1/1/2019 00:00:00 (12:00:00 AM).
func (ts Timestamp) BeginningOfYear() Timestamp {
	return ts.beginningAndEndOf(InYear(ts.Year()))[0]
}

// EndOfYear returns the year ending of the current timestamp 'ts'.
// E.g: 12/11/2019, 15:46:40 (3:46:40 PM) -> 31/12/2019 23:59:59 (11:59:59 PM).
func (ts Timestamp) EndOfYear() Timestamp {
	return ts.beginningAndEndOf(InYear(ts.Year()))[1]
}

// BeginningAndEndOfYear is the same as BeginningOfYear() and EndOfYear() calls.
func (ts Timestamp) BeginningAndEndOfYear() TimestampPair {
	return ts.beginningAndEndOf(InYear(ts.Year()))
}
