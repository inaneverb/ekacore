// Copyright © 2020. All rights reserved.
// Author: Ilya Stroy.
// Contacts: qioalice@gmail.com, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekatime

import (
	"fmt"
	"time"

	"github.com/modern-go/reflect2"
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

//noinspection GoSnakeCaseUsage
var (
	_TIME_PART_AS_NUM_STR [60][]byte
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

// initTimeNumStr initializes _TIME_PART_AS_NUM_STR,
// that are used at the Time.String(), Time.AppendValue()
// and other methods to get a string (or []byte) representation of Time.
func initTimeNumStr() {
	for i := 0; i <= 59; i++ {
		_TIME_PART_AS_NUM_STR[i] =
			reflect2.UnsafeCastString(fmt.Sprintf("%02d", i))
	}
}
