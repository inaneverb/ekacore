// Copyright © 2021. All rights reserved.
// Author: Ilya Stroy.
// Contacts: qioalice@gmail.com, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekatime

//goland:noinspection GoSnakeCaseUsage
const (
	_CALENDAR2_DEFAULT_CAPACITY       = 366
	_CALENDAR2_CAUSE_DEFAULT_CAPACITY = 64
)

func (wc *Calendar) dateToIndex(dd Date) uint {

	// WARNING!
	// Method assumes that wc.year == dd.Year().

	// Calendar always works as in leap year.
	// So, we need to increase doy +1 if it's not leap year and march+ month.

	doy := dd.DayOfYear()
	if dd.Month() >= MONTH_MARCH && !wc.isLeap {
		doy++
	}

	return uint(doy)
}

func (wc *Calendar) rangeOfMonth(m Month) (uint, uint) {

	d := Days(_Table0[m-1])
	if m == MONTH_FEBRUARY && wc.isLeap {
		d++
	}

	d1 := _Table2[m-1]
	d2 := d1 + d - 1 // -1 for d2 be the last day of month

	return uint(d1), uint(d2)
}

func (wc *Calendar) overrideDate(dd Date, eventID EventID, isDayOff, useEventID bool) {

	// 3rd bool argument:
	//
	// A: Use provided EventID (`useEventID`),
	// B: Calendar causing is enabled (`wc.cause != nil`).
	//
	// Truth table:
	// A B Res
	// 0 0 1
	// 0 1 1
	// 1 0 0
	// 1 1 1
	//
	// Thus, it's a material conditional (implication).
	// Read more: https://en.wikipedia.org/wiki/Material_conditional

	if !(wc.IsValid() && dd.Year() == wc.year && (!useEventID || (wc.cause != nil))) {
		return
	}

	idx := wc.dateToIndex(dd)
	wc.dayOff.Set(idx, isDayOff)

	// SAFETY:
	// Condition above guarantees if `useEventID` is true,
	// `wc.cause` is also not nil.
	if useEventID {
		wc.cause[idx] = eventID
	}
}

func (wc *Calendar) nextDay(dd Date, isDayOff bool) Date {

	if !(wc.IsValid() && dd.IsValid() && dd.Year() == wc.year) {
		return _DATE_INVALID
	}

	var (
		nextDay uint
		exist   bool
	)

	if isDayOff {
		nextDay, exist = wc.dayOff.NextUp(wc.dateToIndex(dd))
	} else {
		nextDay, exist = wc.dayOff.NextDown(wc.dateToIndex(dd))
	}

	if !exist {
		return _DATE_INVALID
	}

	return NewDateFromDayOfYear(wc.year, Days(nextDay))
}

func (wc *Calendar) daysIn(m Month, isDayOff bool) []Day {

	if !(wc.IsValid() && m.IsValid()) {
		return nil
	}

	// We don't use Month.DaysInForYear() method here,
	// because it treats years not in range [1900..4095] as invalid years.

	d1, d2 := wc.rangeOfMonth(m)
	ret := make([]Day, 0, 31)

	if isDayOff {
		for v, e := wc.dayOff.NextUp(d1 - 1); e && v <= d2; v, e = wc.dayOff.NextUp(v) {
			ret = append(ret, Day(v-d1)+1)
		}
	} else {
		for v, e := wc.dayOff.NextDown(d1 - 1); e && v <= d2; v, e = wc.dayOff.NextDown(v) {
			ret = append(ret, Day(v-d1)+1)
		}
	}

	return ret
}

func (wc *Calendar) daysInCount(m Month, isDayOff bool) Days {

	if !(wc.IsValid() && m.IsValid()) {
		return 0
	}

	// We don't use Month.DaysInForYear() method here,
	// because it treats years not in range [1900..4095] as invalid years.

	d1, d2 := wc.rangeOfMonth(m)
	d := Days(d2 - d1 + 1)

	c := Days(wc.dayOff.CountBetween(d1, d2))
	if !isDayOff {
		c = d - c
	}

	return c
}

func (wc *Calendar) doSaturdayAndSundayDayOff() {

	if !wc.IsValid() {
		return
	}

	w := NewDate(wc.year, MONTH_JANUARY, 1).Weekday()
	var idx uint = 1

	if w == WEEKDAY_SATURDAY {
		// do nothing, idx already is 1

	} else if w == WEEKDAY_SUNDAY {
		wc.dayOff.Up(idx)
		idx += 6

	} else {
		idx += uint(WEEKDAY_SATURDAY.To06() - w.To06())
	}

	for ; idx <= uint(_Table2[MONTH_MARCH-1]); idx += 7 {
		wc.dayOff.Up(idx)
		wc.dayOff.Up(idx + 1)
	}

	if !wc.isLeap {
		// Calendar's year is always leap (29Feb).
		// But sometimes the real year might be not leap and weekdays are:
		// 28 feb sat and 1 mar sun.
		// 29 feb will be marked as weekday (in loop above).
		// Index increasing is placed below (after this condition).
		// And the last thing we need to do is mark 1 mar sun as weekday.
		// This is exactly what this code do.
		if idx-6 == uint(_Table2[MONTH_MARCH-1]-1) {
			wc.dayOff.Up(idx - 5)
		}
		idx++
	}

	for ; idx <= 366; idx += 7 {
		wc.dayOff.Up(idx)
		wc.dayOff.Up(idx + 1)
	}
}
