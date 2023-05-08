// Copyright © 2020. All rights reserved.
// Author: Ilya Stroy.
// Contacts: iyuryevich@pm.me, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekatime

// TillNext returns how much seconds (as Timestamp) must be passed
// until next time 'range_' will end for the current Timestamp 'ts'.
//
// In most cases you don't need to use this method, but use any of predefined
// instead: TillNextMinute(), TillNextHour(), etc.
//
// Using this function you may get how much time.Duration must be passed,
// until next requested time is passed since now. Range must be passed in seconds.
//
// Examples:
//
//	d := NewDate(2012, MONTH_JANUARY, 12).WithTime(13, 14, 15) // -> 12 Jan 2012 13:14:15
//	d.TillNext(2 * SECONDS_IN_HOUR) // -> 45m45s (till 14:00:00, range of 2h)
//	d.TillNext(3 * SECONDS_IN_HOUR) // -> 1h45m45s (till 15:00:00, range of 3h)
//	d.TillNext(30 * SECONDS_IN_MINUTE) // -> 15m45s (till 13:30:00, range of 30m).
func (ts Timestamp) TillNext(range_ Timestamp) Timestamp {
	return ts + (range_ - ts%range_) - ts
}

// TillNextMinute returns how much seconds (as Timestamp) must be passed
// until next minute (for the current Timestamp 'ts') will come.
func (ts Timestamp) TillNextMinute() Timestamp {
	return ts.TillNext(SECONDS_IN_MINUTE)
}

// TillNextHour returns how much seconds (as Timestamp) must be passed
// until next hour (for the current Timestamp 'ts') will come.
func (ts Timestamp) TillNextHour() Timestamp {
	return ts.TillNext(SECONDS_IN_HOUR)
}

// TillNext12h returns how much seconds (as Timestamp) must be passed
// until next half day (12h) (for the current Timestamp 'ts') will come.
func (ts Timestamp) TillNext12h() Timestamp {
	return ts.TillNext(SECONDS_IN_12H)
}

// TillNextNoon returns how much seconds (as Timestamp) must be passed
// until next noon (12.00 PM) (for the current Timestamp 'ts') will come.
func (ts Timestamp) TillNextNoon() Timestamp {
	var d = ts.TillNext(SECONDS_IN_DAY) + SECONDS_IN_12H
	if d >= SECONDS_IN_DAY {
		d -= SECONDS_IN_DAY
	}
	return d
}

// TillNextMidnight returns how much seconds (as Timestamp) must be passed
// until next midnight (12.00 AM) (for the current Timestamp 'ts') will come.
func (ts Timestamp) TillNextMidnight() Timestamp {
	return ts.TillNextDay()
}

// TillNextDay returns how much seconds (as Timestamp) must be passed
// until next day (for the current Timestamp 'ts') will come.
func (ts Timestamp) TillNextDay() Timestamp {
	return ts.TillNext(SECONDS_IN_DAY)
}

// TillNextMonth returns how much seconds (as Timestamp) must be passed
// until next month (for the current Timestamp 'ts') will come.
func (ts Timestamp) TillNextMonth() Timestamp {
	y, m, _ := dateFromUnix(ts)
	return ts.TillNext(InMonth(y, m))
}

// TillNextYear returns how much seconds (as Timestamp) must be passed
// until next year (for the current Timestamp 'ts') will come.
func (ts Timestamp) TillNextYear() Timestamp {
	return ts.TillNext(InYear(ts.Year()))
}

// TillNextMinute is the same as Timestamp.TillNextMinute()
// but for current time.
func TillNextMinute() Timestamp {
	return NewTimestampNow().TillNextMinute()
}

// TillNextHour is the same as Timestamp.TillNextHour() but for current time.
func TillNextHour() Timestamp {
	return NewTimestampNow().TillNextHour()
}

// TillNext12h is the same as Timestamp.TillNext12h() but for current time.
func TillNext12h() Timestamp {
	return NewTimestampNow().TillNext12h()
}

// TillNextNoon is the same as Timestamp.TillNextNoon() but for current time.
func TillNextNoon() Timestamp {
	return NewTimestampNow().TillNextNoon()
}

// TillNextMidnight is the same as Timestamp.TillNextMidnight()
// but for current time.
func TillNextMidnight() Timestamp {
	return NewTimestampNow().TillNextMidnight()
}

// TillNextDay is the same as Timestamp.TillNextDay() but for current time.
func TillNextDay() Timestamp {
	return NewTimestampNow().TillNextDay()
}

// TillNextMonth is the same as Timestamp.TillNextMonth() but for current time.
func TillNextMonth() Timestamp {
	return NewTimestampNow().TillNextMonth()
}

// TillNextYear is the same as Timestamp.TillNextYear() but for current time.
func TillNextYear() Timestamp {
	return NewTimestampNow().TillNextYear()
}
