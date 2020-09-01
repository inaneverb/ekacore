// Copyright © 2020. All rights reserved.
// Author: Ilya Stroy.
// Contacts: qioalice@gmail.com, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekatime

import (
	"time"
)

// getBeginningOfYear returns a unix timestamp of the
// 00:00:01 (12:00:00 AM) time of the January 1st in the 'y' year.
func getBeginningOfYear(y Year) Timestamp {
	if y >= 1970 && y <= _Table5UpperBound {
		return _Table5[y-1970][0]
	}
	return Timestamp(time.Date(int(y), time.January, 1, 0, 0, 0, 0, time.UTC).Unix())
}

// getBeginningOfMonth returns a unix timestamp of the
// 00:00:01 (12:00:00 AM) time for the 1st day of 'm' month in the 'y' year.
func getBeginningOfMonth(y Year, m Month) Timestamp {
	if y >= 1970 && y <= _Table5UpperBound && m >= MONTH_JANUARY && m <= MONTH_DECEMBER {
		return _Table5[y-1970][m-1]
	}
	return Timestamp(time.Date(int(y), time.Month(m), 1, 0, 0, 0, 0, time.UTC).Unix())
}

// dateFromUnix returns an Year, Month and Day from the current timestamp 'ts'.
func dateFromUnix(t Timestamp) (y Year, m Month, d Day) {
	ty, tm, td := time.Unix(t.I64(), 0).UTC().Date()
	return Year(ty), Month(tm), Day(td)
}

// timeFromUnix returns an Hour, Minute and Second from the current timestamp 'ts'.
func timeFromUnix(t Timestamp) (h Hour, m Minute, s Second) {
	th, tm, ts := time.Unix(t.I64(), 0).UTC().Clock()
	return Hour(th), Minute(tm), Second(ts)
}
