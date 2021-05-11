// Copyright © 2020. All rights reserved.
// Author: Ilya Stroy.
// Contacts: qioalice@gmail.com, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekatime

import (
	"time"
)

// TillNext returns how much ns (as time.Duration) must be passed until next time
// 'range_' will end for the current Timestamp 'ts'.
//
// In most cases you don't need to use this method, but use any of predefined
// instead: TillNextMinute(), TillNextHour(), etc.
//
// Using this function you may get how much time.Duration must be passed,
// until next requested time is passed since now. Range must be passed in seconds.
//
// Examples:
//   d := NewDate(2012, MONTH_JANUARY, 12).WithTime(13, 14, 15) // -> 12 Jan 2012 13:14:15
//   d.TillNext(2 * SECONDS_IN_HOUR) // -> 45m45s (till 14:00:00, range of 2h)
//   d.TillNext(3 * SECONDS_IN_HOUR) // -> 1h45m45s (till 15:00:00, range of 3h)
//   d.TillNext(30 * SECONDS_IN_MINUTE) // -> 15m45s (till 13:30:00, range of 30m).
func (ts Timestamp) TillNext(range_ Timestamp) time.Duration {
	return time.Duration(ts.tillNext(range_)) * time.Second
}

// TillNextMinute returns how much ns (as time.Duration) must be passed until
// next minute (for the current Timestamp 'ts') will came.
func (ts Timestamp) TillNextMinute() time.Duration {
	return ts.TillNext(SECONDS_IN_MINUTE)
}

// TillNextHour returns how much ns (as time.Duration) must be passed until
// next hour (for the current Timestamp 'ts') will came.
func (ts Timestamp) TillNextHour() time.Duration {
	return ts.TillNext(SECONDS_IN_HOUR)
}

// TillNext12h returns how much ns (as time.Duration) must be passed until
// next half day (12h) (for the current Timestamp 'ts') will came.
func (ts Timestamp) TillNext12h() time.Duration {
	return ts.TillNext(SECONDS_IN_12H)
}

// TillNextNoon returns how much ns (as time.Duration) must be passed until
// next noon (12.00 PM) (for the current Timestamp 'ts') will came.
func (ts Timestamp) TillNextNoon() time.Duration {
	d := ts.TillNext(SECONDS_IN_DAY) + 12 * time.Hour
	if d >= 24 * time.Hour {
		d -= 24 * time.Hour
	}
	return d
}

// TillNextMidnight returns how much ns (as time.Duration) must be passed until
// next midnight (12.00 AM) (for the current Timestamp 'ts') will came.
func (ts Timestamp) TillNextMidnight() time.Duration {
	return ts.TillNextDay()
}

// TillNextDay returns how much ns (as time.Duration) must be passed until
// next day (for the current Timestamp 'ts') will came.
func (ts Timestamp) TillNextDay() time.Duration {
	return ts.TillNext(SECONDS_IN_DAY)
}

// TillNextMonth returns how much ns (as time.Duration) must be passed until
// next month (for the current Timestamp 'ts') will came.
func (ts Timestamp) TillNextMonth() time.Duration {
	y, m, _ := dateFromUnix(ts)
	return ts.TillNext(InMonth(y, m))
}

// TillNextYear returns how much ns (as time.Duration) must be passed until
// next year (for the current Timestamp 'ts') will came.
func (ts Timestamp) TillNextYear() time.Duration {
	return ts.TillNext(InYear(ts.Year()))
}

// TillNextMinute is the same as Timestamp.TillNextMinute() but for current time.
func TillNextMinute() time.Duration {
	return Now().TillNextMinute()
}

// TillNextHour is the same as Timestamp.TillNextHour() but for current time.
func TillNextHour() time.Duration {
	return Now().TillNextHour()
}

// TillNext12h is the same as Timestamp.TillNext12h() but for current time.
func TillNext12h() time.Duration {
	return Now().TillNext12h()
}

// TillNextNoon is the same as Timestamp.TillNextNoon() but for current time.
func TillNextNoon() time.Duration {
	return Now().TillNextNoon()
}

// TillNextMidnight is the same as Timestamp.TillNextMidnight() but for current time.
func TillNextMidnight() time.Duration {
	return Now().TillNextMidnight()
}

// TillNextDay is the same as Timestamp.TillNextDay() but for current time.
func TillNextDay() time.Duration {
	return Now().TillNextDay()
}

// TillNextMonth is the same as Timestamp.TillNextMonth() but for current time.
func TillNextMonth() time.Duration {
	return Now().TillNextMonth()
}

// TillNextYear is the same as Timestamp.TillNextYear() but for current time.
func TillNextYear() time.Duration {
	return Now().TillNextYear()
}

func(ts Timestamp) tillNext(range_ Timestamp) Timestamp {
	return ts + (range_- ts % range_) - ts
}