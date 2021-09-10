package ekatime

import (

)

//goland:noinspection GoSnakeCaseUsage
const (
	_CALENDAR2_DEFAULT_CAPACITY       = 370
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
		nextDay, exist = wc.dayOff.NextSet(uint(dayOfYear))
	} else {
		nextDay, exist = wc.dayOff.NextClear(uint(dayOfYear))
	}

	if !exist {
		return _DATE_INVALID
	}

	return NewDateFromDays(wc.year, Days(nextDay))
}

func (wc *WorkCalendar) doSaturdayAndSundayDayOff() {

	if !wc.IsValid() {
		return
	}

	w := NewDate(wc.year, MONTH_JANUARY, 1).Weekday()
	var idx uint = 1

	if w == WEEKDAY_SATURDAY {
		// do nothing, idx already is 1

	} else if w == WEEKDAY_SUNDAY {
		wc.dayOff.Set(idx)
		idx += 6

	} else {
		idx += uint(WEEKDAY_SATURDAY.To06() - w.To06())
	}

	for ; idx <= 366; idx += 7 {
		wc.dayOff.Set(idx)
		wc.dayOff.Set(idx + 1)
	}
}