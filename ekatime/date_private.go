// Copyright © 2020. All rights reserved.
// Author: Ilya Stroy.
// Contacts: qioalice@gmail.com, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekatime

import (
	"time"
)

//noinspection GoSnakeCaseUsage
const (
	_DATE_OFFSET_DAY     uint8   = 0
	_DATE_OFFSET_MONTH   uint8   = _DATE_OFFSET_DAY + 5
	_DATE_OFFSET_YEAR    uint8   = _DATE_OFFSET_MONTH + 4
	_DATE_OFFSET_WEEKDAY uint8   = _DATE_OFFSET_YEAR + 12
	_DATE_OFFSET_UNUSED  uint8   = _DATE_OFFSET_WEEKDAY + 3

	_DATE_MASK_DAY       Day     = 0x1F
	_DATE_MASK_MONTH     Month   = 0x0F
	_DATE_MASK_YEAR      Year    = 0x0FFF
	_DATE_MASK_WEEKDAY   Weekday = 0x07

	_DATE_MASK_DATE      Date    = (Date(1) << _DATE_OFFSET_UNUSED) - 1
)

// normalizeDate shifts Date, 'y', 'm', and 'd' represents which if they are not
// in their valid ranges. Returns the fixed values (if they has been).
func normalizeDate(y Year, m Month, d Day) (Year, Month, Day) {
	if !IsValidDate(y, m, d) {
		t := time.Date(int(y), time.Month(m), int(d), 0, 0, 0, 0, time.UTC)
		ty, tm, td := t.Date()
		y, m, d = Year(ty), Month(tm), Day(td)
	}
	return y, m, d
}

// ensureWeekdayExist returns the current Date w/o modifications if it already
// has a weekday or adds it and to and returns a copy.
func (dd Date) ensureWeekdayExist() Date {
	return dd | dd.Weekday().asPartOfDate()
}
