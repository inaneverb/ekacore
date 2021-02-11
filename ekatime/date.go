// Copyright © 2020. All rights reserved.
// Author: Ilya Stroy.
// Contacts: qioalice@gmail.com, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekatime

type (
	// Day is a special type that has enough space to store Day's number.
	// Useless just by yourself but is a part of Date object.
	// Valid values: [1..31].
	Day int8

	// Days is like time.Duration from std Golang lib,
	// but indicates time's range in passed days.
	Days int32

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
	//     date (e.g: 01 Jan 1970) may not be equal by just eq comparison, like:
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

// DaysInForYear is an alias for DaysInMonth(y, m), where m is the current Month.
func (m Month) DaysInForYear(y Year) Day {
	return DaysInMonth(y, m)
}

// DaysInIgnoreYear is an alias for DaysInMonthIgnoreYear(m),
// where m is the current Month.
func (m Month) DaysInIgnoreYear() Day {
	return DaysInMonthIgnoreYear(m)
}

// IsLeap returns true if the current year is leap (e.g. 1992, 1996, 2000, 2004, etc).
func (y Year) IsLeap() bool {
	return IsLeap(y)
}

// IsValidDate reports whether 'y', 'm' and 'd' in their valid ranges.
// Also checks leap year, february 28,29, correct day number for month.
func IsValidDate(y Year, m Month, d Day) bool {
	return !(y < 0 || y > 4095 || m < 1 || m > 12 || d < 1 ||
		(IsLeap(y) && m == 2 && d > 29) || d > _Table0[m-1])
}

// IsValid is an alias for IsValidDate().
func (dd Date) IsValid() bool {
	return IsValidDate(dd.Split())
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

// DaysInMonth returns how much days the month contains the current Date includes which.
func (dd Date) DaysInMonth() Day {
	y, m, d := normalizeDate(dd.Split())
	d = _Table0[m-1]
	if m == MONTH_FEBRUARY && y.IsLeap() {
		d++
	}
	return d
}

// DaysInMonth returns how much days requested Month contains in passed Year.
// Returns -1, if Year and Month not in their allowed ranges.
func DaysInMonth(y Year, m Month) Day {
	if _YEAR_MIN <= y && y <= _YEAR_MAX &&
			MONTH_JANUARY <= m && m <= MONTH_DECEMBER {
		d := _Table0[m-1]
		if m == MONTH_FEBRUARY && y.IsLeap() {
			d++
		}
		return d
	}
	return -1
}

// DaysInMonthIgnoreYear does the same thing as DaysInMonth,
// but does not require the Year, meaning that it will be always 28 for MONTH_FEBRUARY,
// even for leap years like 2000, 2004, etc.
func DaysInMonthIgnoreYear(m Month) Day {
	if MONTH_JANUARY <= m && m <= MONTH_DECEMBER {
		return _Table0[m-1]
	} else {
		return -1
	}
}

// Weekday returns the current Date's day of week.
func (dd Date) Weekday() Weekday {
	if w := Weekday(dd >> _DATE_OFFSET_WEEKDAY) & _DATE_MASK_WEEKDAY; w > 0 {
		return w-1
	} else {
		return dd.WithTime(0, 0, 0).weekday()
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

// Replace returns a new Date based on the current.
// It returns the current Date with changed Year, Month, Day to those passed values,
// which are in their allowed ranges. Does not doing date addition. Only replacing.
// For date addition, use Add() method.
//
// Examples:
//  NewDate(2021, MONTH_FEBRUARY, 10)    // -> 10 Feb 2021
//    .Replace(2013, MONTH_JANUARY, 2)   // -> 02 Jan 2013
//    .Replace(2020, -2, 4)              // -> 04 Jan 2020
//    .Replace(0, 0, -50)                // -> 04 Jan 2020
//    .Replace(1899, 61, 31)             // -> 31 Jan 2020 (1900 is min year)
//    .Replace(2013, MONTH_FEBRUARY, -2) // -> 31 Jan 2013 (31 Feb is not allowed).
//    .Replace(2014, MONTH_FEBRUARY, 30) // -> 30 Jan 2014
//    .Replace(2000, MONTH_FEBRUARY, 29) // -> 29 Feb 2000
//    .Replace(2001, -1, -1)             // -> 29 Feb 2000 (29 Feb in 2001 is not allowed)
func (dd Date) Replace(y Year, m Month, d Day) Date {
	y_, m_, d_ := dd.Split()

	// There are some edge cases, that must be handled:
	// 1. dd is 29 Feb Leap's year, changing year is requested, check it.
	// 2. dd is 31 of some month, changing to 30-days month is requested, check it.

	var (
		ndOK = 1 <= d && d <= 31 && d != d_
		nmOK = MONTH_JANUARY <= m && m <= MONTH_DECEMBER && m != m_
		nyOK = _YEAR_MIN <= y && y <= _YEAR_MAX && y != y_
	)

	//goland:noinspection GoSnakeCaseUsage
	var (
		Zym   = DaysInMonth(y, m)
		Zy_m  = DaysInMonth(y_, m)
		Zy_m_ = DaysInMonth(y_, m_)
		Zym_  = DaysInMonth(y, m_)
	)

	//switch {
	//case nmOK && nyOK && ndOK && d <= Zym:
	//	y_, m_, d_ = y, m, d
	//
	//case nmOK && nyOK && ndOK && d <= Zy_m:
	//	m_, d_ = m, d
	//
	//case nmOK && nyOK && ndOK && d <= Zym_:
	//	y_, d_ = y, d
	//
	//case nmOK && nyOK && ndOK && d <= Zy_m_:
	//	d_ = d
	//
	//case nmOK && nyOK && !ndOK && d_ <= Zym:
	//	y_, m_ = y, m
	//
	//case nmOK && nyOK && !ndOK && d_ <= Zy_m:
	//	m_ = m
	//
	//case nmOK && nyOK && !ndOK && d_ <= Zym_:
	//	y_ = y
	//
	//case !nmOK && nyOK && ndOK && d <= Zym_:
	//	y_, d_ = y, d
	//
	//case !nmOK && nyOK && ndOK && d_ <= Zym_:
	//	y_ = y
	//
	//case nmOK && !nyOK && ndOK && d <= Zy_m:
	//	m_, d_ = m, d
	//
	//case nmOK && !nyOK && ndOK && d_ <= Zy_m:
	//	m_ = m
	//
	//case nmOK && !nyOK && !ndOK && d_ <= Zy_m:
	//	m_ = m
	//
	//case !nmOK && nyOK && !ndOK && d_ <= Zym_:
	//	y_ = y
	//
	//case !nmOK && !nyOK && ndOK && d <= Zy_m_:
	//	d_ = d
	//}

	// Switch above is simplified, using CNDF:
	// https://en.wikipedia.org/wiki/Canonical_normal_form

	var (
		A = nmOK
		B = nyOK
		C = ndOK
		D = d <= Zym
		E = d <= Zy_m
		F = d <= Zy_m_
		G = d <= Zym_
		H = d_ <= Zym
		I = d_ <= Zy_m
		K = d_ <= Zym_
	)

	if A && B && (C && (D || G) || H || K) || B && (C && (G || K) || K) {
		y_ = y
	}
	if A && B && (C && (D || E) || H || I) || A && (C && (E || I) || I) {
		m_ = m
	}
	if A && B && C && (D || E || G || F) || C && (B && G || A && E || F) {
		d_ = d
	}

	return NewDate(y_, m_, d_)
}

// Add returns a new Date based on the current.
// It returns the current Date with changed Year, Month, Day using passed values
// as their addition's deltas.
//
// Examples:
//  NewDate(2021, MONTH_FEBRUARY, 10) // 10 Feb 2021
//    .Add(1, 2, 3)                   // 13 Apr 2022
//    .Add(0, -1, -13)                // 28 Feb 2022 (0 Mar 2022 -> 28 Feb 2022)
//    .Add(0, 1, 3)                   // 31 Mar 2022
//    .Add(0, 1, 0)                   // 01 May 2022 (31 Apr 2022 (not exist) -> 01 May 2022)
//    .Add(0, 127, 0)                 // 01 Dec 2032 (works OK with potential integer overflow)
func (dd Date) Add(y Year, m Month, d Day) Date {
	y_, m_, d_ := dd.Split()

	if y > _YEAR_MAX {
		y = _YEAR_MAX
	}
	if y < -(y_ - _YEAR_MIN) {
		y = -(y_ - _YEAR_MIN)
	}

	y_ += y
	y_ += Year(m / 12)

	m_ += m % 12

	if y_ > _YEAR_MAX {
		y_ = _YEAR_MAX
	}
	if y_ < _YEAR_MIN {
		y_ = _YEAR_MIN
	}

	return NewDate(y_, m_, d_).AddDays(Days(d))
}

// Days returns an accumulated number of days that has been passed since 1 Jan.
// Returns 0 if current Date is not valid.
func (dd Date) Days() Days {
	if !dd.IsValid() {
		return 0
	}
	d := Days(dd.Day())
	for i := MONTH_JANUARY; i < dd.Month(); i++ {
		d += Days(_Table0[i-1])
	}
	if dd.Month() > MONTH_FEBRUARY && dd.Year().IsLeap() {
		d++
	}
	return d
}

// AddDays returns a new Date based on the current Date.
// It additions exactly passed days to the current Date and returns a result.
func (dd Date) AddDays(days Days) Date {
	return (dd.WithTime(0, 0, 0) + Timestamp(days) * SECONDS_IN_DAY).Date()
}

// WithTime returns the current Date with the presented Time's hour, minute, second
// as a new Timestamp object.
func (dd Date) WithTime(hh Hour, mm Minute, ss Second) Timestamp {
	y, m, d := dd.Split()
	return UnixFrom(y, m, d, hh, mm,ss)
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
