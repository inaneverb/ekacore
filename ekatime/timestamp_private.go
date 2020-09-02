// Copyright © 2020. All rights reserved.
// Author: Ilya Stroy.
// Contacts: qioalice@gmail.com, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekatime

import (
	"time"
)

// beginningAndEndOf returns the [A,B] pair, and A <= 'ts' <= B, and A-B == 'range_'.
// There is a special formula that allows to get start and ending of 'ts'
// related: day, month, year.
func (ts Timestamp) beginningAndEndOf(range_ Timestamp) TimestampPair {
	x := ts + (range_- ts % range_)
	return TimestampPair{x - range_, x - 1}
}

// tillNext returns how much ns (as time.Duration) must be passwd until next time
// 'range_' will end for the current Timestamp 'ts'.
func (ts Timestamp) tillNext(range_ Timestamp) time.Duration {
	return time.Duration(ts + (range_- ts % range_) - ts) * time.Second
}

// weekday returns the current Timestamp 'ts' the number of day in week.
func (ts Timestamp) weekday() Weekday {
	return Weekday(((ts + SECONDS_IN_DAY) % SECONDS_IN_WEEK) / SECONDS_IN_DAY)
}
