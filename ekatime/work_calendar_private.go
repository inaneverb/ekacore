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

func (wc *WorkCalendar) nextDay(dayOfYear Days, isDayOff bool) Date {

	if !(wc.IsValid() && dayOfYear.BelongsToYear(wc.year)) {
		return _DATE_INVALID
	}

	var (
		nextDay uint
		exist   bool
	)

	if isDayOff {
		nextDay, exist = wc.dayOff.NextUp(uint(dayOfYear))
	} else {
		nextDay, exist = wc.dayOff.NextDown(uint(dayOfYear))
	}

	if !exist {
		return _DATE_INVALID
	}

	return NewDateFromDays(wc.year, Days(nextDay))
}

func (wc *WorkCalendar) daysIn(m Month, isDayOff bool) []Day {

	if !(wc.IsValid() && m.IsValid()) {
		return nil
	}

	// We don't use Month.DaysInForYear() method here,
	// because it treats years not in range [1900..4095] as invalid years.

	d1 := _Table2[m-1]
	d2 := d1 + Days(m.DaysInIgnoreYear()) // d2 is the 1st day of the next month
	if m == MONTH_FEBRUARY && wc.isLeap {
		d2++
	}

	ret := make([]Day, 0, 31)

	if isDayOff {
		for v, e := wc.dayOff.NextUp(uint(d1 - 1)); e && v < uint(d2); v, e = wc.dayOff.NextUp(v) {
			ret = append(ret, Day(v))
		}
	} else {
		for v, e := wc.dayOff.NextDown(uint(d1 - 1)); e && v < uint(d2); v, e = wc.dayOff.NextDown(v) {
			ret = append(ret, Day(v))
		}
	}

	return ret
}

func (wc *WorkCalendar) daysInCount(m Month, isDayOff bool) Days {

	if !(wc.IsValid() && m.IsValid()) {
		return 0
	}

	// We don't use Month.DaysInForYear() method here,
	// because it treats years not in range [1900..4095] as invalid years.

	d := Days(m.DaysInIgnoreYear())
	if m == MONTH_FEBRUARY && wc.isLeap {
		d++
	}

	d1 := _Table2[m-1]
	d2 := d1 + d - 1 // -1 for d2 be the last day of month

	c := Days(wc.dayOff.CountBetween(uint(d1), uint(d2)))
	if !isDayOff {
		c = d - c
	}

	return c
}

func (wc *WorkCalendar) DoSaturdayAndSundayDayOff() {

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

	for ; idx <= 366; idx += 7 {
		wc.dayOff.Up(idx)
		wc.dayOff.Up(idx + 1)
	}
}
