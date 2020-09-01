// Copyright © 2020. All rights reserved.
// Author: Ilya Stroy.
// Contacts: qioalice@gmail.com, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekatime

import "fmt"

type (
	// Day is a special type that has enough space to store Day's number.
	// Useless just by yourself but is a part of Date object.
	// Valid values: [1..31].
	Day int8

	// Month is a special type that has enough space to store Month's number.
	// Useless just by yourself but is a part of Date object.
	// Valid values: [1..12]. Use predefined constants to make it clear.
	Month int8

	// Year is a special type that has enough space to store Year's number.
	// Useless just by yourself but is a part of Date object.
	// Valid values: [1900..4095].
	Year int16

	// Date is a special object that has enough space to store a some Date
	// (including Day, Month, Year) but does it most RAM efficient way
	// taking only 4 bytes.
	//
	// WARNING!
	// DO NOT COMPARE Date OBJECTS JUST BY EQUAL OPERATOR! PREPARE THEM BEFORE.
	//
	//     Because of the internal parts, the Date objects that represents the same
	//     date (e.g: 01 Jan 1970) may not be equal by just eq comparision, like:
	//
	//         d1 := UnixFrom(NewDate(1970, 1, 1), 0).Date()
	//         d2 := NewDate(1970, 1, 1)
	//         d1 == d2 // false <---
	//
	//     For being able to use eq operator, call ToCmp() before,
	//     or use Equal() method:
	//
	//         d1 := UnixFrom(NewDate(1970, 1, 1), 0).Date()
	//         d2 := NewDate(1970, 1, 1)
	//         d1.ToCmp() == d2.ToCmp() // true <---
	//         d1.Equal(d2) // true <---
	//
	Date uint32
)

//noinspection GoSnakeCaseUsage,GoUnusedConst
const (
	MONTH_JANUARY Month = 1 + iota
	MONTH_FEBRUARY
	MONTH_MARCH
	MONTH_APRIL
	MONTH_MAY
	MONTH_JUNE
	MONTH_JULY
	MONTH_AUGUST
	MONTH_SEPTEMBER
	MONTH_OCTOBER
	MONTH_NOVEMBER
	MONTH_DECEMBER
)

// IsLeap returns true if the current year is leap (e.g. 1992, 1996, 2000, 2004, etc).
func (y Year) IsLeap() bool {
	return IsLeap(y)
}

// ToCmp returns the current Date object ready for being compared using eq operator ==.
// Yes you MUST NOT compare Date objects directly w/o call this method.
// See Date's doc for more details.
func (dd Date) ToCmp() Date {
	return dd &^ (Date(_DATE_MASK_WEEKDAY) << _DATE_OFFSET_WEEKDAY)
}

// Equal returns true if the current Date is the same as 'other'.
func (dd Date) Equal(other Date) bool {
	return dd.ToCmp() == other.ToCmp()
}

// Year returns the year number the current Date includes which.
//
// It guarantees that Year() returns the valid year number Date is of,
// only if Date has not been created manually but using constructors
// like NewDate(), Event.Date(), Timestamp.Date(), Timestamp.Split(), etc.
func (dd Date) Year() Year {
	return Year(dd >> _DATE_OFFSET_YEAR) & _DATE_MASK_YEAR
}

// Month returns the month number the current Date includes which.
// January is 1, not 0.
//
// It guarantees that Month() returns the valid month number Date is of,
// only if Date has not been created manually but using constructors
// like NewDate(), Event.Date(), Timestamp.Date(), Timestamp.Split(), etc.
func (dd Date) Month() Month {
	return Month(dd >> _DATE_OFFSET_MONTH) & _DATE_MASK_MONTH
}

// Day returns the day number the current Date includes which.
//
// It guarantees that Day() returns the valid day number for the month Date is of,
// only if Date has not been created manually but using constructors
// like NewDate(), Event.Date(), Timestamp.Date(), Timestamp.Split(), etc.
func (dd Date) Day() Day {
	return Day(dd >> _DATE_OFFSET_DAY) & _DATE_MASK_DAY
}

// Weekday returns the current Date's day of week.
func (dd Date) Weekday() Weekday {
	if w := Weekday(dd >> _DATE_OFFSET_WEEKDAY) & _DATE_MASK_WEEKDAY; w > 0 {
		return w-1
	} else {
		return UnixFrom(dd, 0).weekday()
	}
}

// Split returns the year number, month number and day number the current Date
// includes which.
// It's just like a separate Year(), Month(), Day() calls.
func (dd Date) Split() (y Year, m Month, d Day) {
	return dd.Year(), dd.Month(), dd.Day()
}

// NewDate creates a new Date object using provided year number, month number,
// day number, normalizing these values and shifting date if it's required.
//
// Totally, it just stores the provided data if values are in their valid ranges,
// like: [1900..4095] for year, [1..12] for month and [1..X] for day,
// where X depends by month and year.
//
// If they are not, the date may be (will be) shifted. E.g:
// 0 May 2016 (y == 2016, m == 5, d == 0) -> 30 April 2016 (y == 2016, m == 4, d == 30).
func NewDate(y Year, m Month, d Day) Date {

	if y > 4095 {
		y = 4095
	}

	y, m, d = normalizeDate(y, m, d)

	// Do not forgot 'bitwise AND' between y, m, d and their bitmasks
	// if you will change the logic of getting valid y, m, d from invalid ones.
	// Now it's unnecessary (redundant).
	//y, m, d = y & _DATE_MASK_YEAR, m & _DATE_MASK_MONTH, d & _DATE_MASK_DAY

	return (Date(y) << _DATE_OFFSET_YEAR) |
		(Date(m) << _DATE_OFFSET_MONTH) |
		(Date(d) << _DATE_OFFSET_DAY)
}

// String returns the current Date's string representation in the following format:
// "YYYY/MM/DD".
func (dd Date) String() string {
	return fmt.Sprintf("%04d/%02d/%02d", dd.Year(), dd.Month(), dd.Day())
}

// IsLeap returns true if 'y' Year is leap (e.g. 1992, 1996, 2000, 2004, etc).
func IsLeap(y Year) bool {
	return y%400 == 0 || (y%4 == 0 && y%100 != 0)
}

// InMonth returns how much seconds month 'm' in year 'y' contains.
//
// Returns it as Timestamp because of easy arithmetic operations,
// but it's NOT A TIMESTAMP!
func InMonth(y Year, m Month) Timestamp {
	y, m, _ = normalizeDate(y, m, 1)

	sec := _Table1[m-1]
	if m == MONTH_FEBRUARY && IsLeap(y) {
		sec += SECONDS_IN_DAY
	}

	return sec
}

// InYear returns how much seconds year 'y' contains.
//
// Returns it as Timestamp because of easy arithmetic operations,
// but it's NOT A TIMESTAMP!
func InYear(y Year) Timestamp {
	if IsLeap(y) {
		return SECONDS_IN_366_YEAR
	} else {
		return SECONDS_IN_365_YEAR
	}
}

// BeginningOfYear returns the Timestamp of the 'y' year beginning:
// 1 January, 00:00:00 (12:00:00 AM).
//
// It's up to 10 times faster for [1970..N+10] years, where N is current year.
func BeginningOfYear(y Year) Timestamp {
	return getBeginningOfYear(y) // there is no need of normalizing
}

// EndOfYear returns the Timestamp of the 'y' year end:
// 31 December, 23:59:59 (11:59:59 PM).
//
// It's up to 4.8 times faster for [1970..N+10] years, where N is current year.
func EndOfYear(y Year) Timestamp {
	return getBeginningOfYear(y) + InYear(y) -1 // there is no need of normalizing
}

// BeginningAndEndOfYear is like BeginningOfYear() and EndOfYear() calls.
func BeginningAndEndOfYear(y Year) TimestampPair {
	by := getBeginningOfYear(y) // there is no need of normalizing
	return TimestampPair{by, by + InYear(y) -1} // there is no need of normalizing
}

// BeginningOfMonth returns the Timestamp of the 'm' month beginning in 'y' year:
// 1st day, 00:00:00 (12:00:00 AM).
//
// It's up to 7 times faster for [1970..N+10] years, where N is current year.
func BeginningOfMonth(y Year, m Month) Timestamp {
	return getBeginningOfMonth(y, m) // there is no need of normalizing
}

// EndOfMonth returns the Timestamp of the 'm' month end in 'y' year:
// 29/30/31 (depends by month) day, 23:59:59 (11:59:59 PM).
//
// It's up to 2.1 times faster for [1970..N+10] years, where N is current year.
func EndOfMonth(y Year, m Month) Timestamp {
	return getBeginningOfMonth(y, m) + InMonth(y, m) -1 // there is no need of normalizing
}

// BeginningAndEndOfMonth is like BeginningOfMonth() and EndOfMonth() calls.
func BeginningAndEndOfMonth(y Year, m Month) TimestampPair {
	bm := getBeginningOfMonth(y, m) // there is no need of normalizing
	return TimestampPair{bm, bm + InMonth(y, m) -1} // there is no need of normalizing
}
