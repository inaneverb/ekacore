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
	_TIME_OFFSET_HOUR   uint8 = 0
	_TIME_OFFSET_MINUTE uint8 = _TIME_OFFSET_HOUR + 5
	_TIME_OFFSET_SECOND uint8 = _TIME_OFFSET_MINUTE + 6
	_TIME_OFFSET_UNUSED uint8 = _TIME_OFFSET_SECOND + 6

	_TIME_MASK_HOUR     Hour   = 0x1F
	_TIME_MASK_MINUTE   Minute = 0x3F
	_TIME_MASK_SECOND   Second = 0x3F

	_TIME_MASK_TIME     Time   = (Time(1) << _TIME_OFFSET_UNUSED) - 1
)

// normalizeTime shifts Time, 'h', 'm' and 's' represents which if they are not
// in their valid ranges. Returns the fixed values (if they has been).
func normalizeTime(h Hour, m Minute, s Second) (Hour, Minute, Second) {
	if !IsValidTime(h, m, s) {
		t := time.Date(0, 0, 0, int(h), int(m), int(s), 0, time.UTC)
		th, tm, ts := t.Clock()
		h, m, s = Hour(th), Minute(tm), Second(ts)
	}
	return h, m, s
}
