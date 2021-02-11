// Copyright © 2020. All rights reserved.
// Author: Ilya Stroy.
// Contacts: qioalice@gmail.com, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekatime

import (
	"fmt"
	"strconv"
	"time"

	"github.com/modern-go/reflect2"
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

//noinspection GoSnakeCaseUsage
const (
	// Min and max values to use _YEAR_AS_NUM_STR instead of strconv.Itoa()
	_YEAR_AS_NUM_STR_MIN      = _YEAR_MIN
	_YEAR_AS_NUM_STR_MAX Year = 2100

	_YEAR_MIN Year = 1900
	_YEAR_MAX Year = 4095
)

//noinspection GoSnakeCaseUsage
var (
	_YEAR_AS_NUM_STR [_YEAR_AS_NUM_STR_MAX-_YEAR_AS_NUM_STR_MIN+1][]byte
	_DAY_AS_NUM_STR [31][]byte
)

// normalizeDate shifts Date, 'y', 'm', and 'd' represents which if they are not
// in their valid ranges. Returns the fixed values (if they has been).
func normalizeDate(y Year, m Month, d Day) (Year, Month, Day) {
	if !IsValidDate(y, m, d) {
		t := time.Date(int(y), time.Month(m), int(d), 0, 0, 0, 0, time.UTC)
		ty, tm, td := t.Date()
		if y > _YEAR_MAX {
			y = _YEAR_MAX
		}
		y, m, d = Year(ty), Month(tm), Day(td)
	}
	return y, m, d
}

// ensureWeekdayExist returns the current Date w/o modifications if it already
// has a weekday or adds it and to and returns a copy.
func (dd Date) ensureWeekdayExist() Date {
	return dd | dd.Weekday().asPartOfDate()
}

// initDateNumStr initializes _YEAR_AS_NUM_STR and _DAY_AS_NUM_STR
// that are used at the Date.String(), Date.AppendTo()
// and other methods to get a string (or []byte) representation of Date.
func initDateNumStr() {
	for y := _YEAR_AS_NUM_STR_MIN; y <= _YEAR_AS_NUM_STR_MAX; y++ {
		_YEAR_AS_NUM_STR[y-_YEAR_AS_NUM_STR_MIN] =
			reflect2.UnsafeCastString(strconv.Itoa(int(y)))
	}
	for d := Day(1); d <= 31; d++ {
		_DAY_AS_NUM_STR[d-1] =
			reflect2.UnsafeCastString(fmt.Sprintf("%02d", d))
	}
}
