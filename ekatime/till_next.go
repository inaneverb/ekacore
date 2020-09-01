// Copyright © 2020. All rights reserved.
// Author: Ilya Stroy.
// Contacts: qioalice@gmail.com, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekatime

import (
	"time"
)

// TillNextMinute returns how much ns (as time.Duration) must be passed until
// next minute (for the current Timestamp 'ts') will came.
func (ts Timestamp) TillNextMinute() time.Duration {
	return ts.tillNext(60)
}

// TillNextHour returns how much ns (as time.Duration) must be passed until
// next hour (for the current Timestamp 'ts') will came.
func (ts Timestamp) TillNextHour() time.Duration {
	return ts.tillNext(3600)
}

// TillNext12h returns how much ns (as time.Duration) must be passed until
// next half day (12h) (for the current Timestamp 'ts') will came.
func (ts Timestamp) TillNext12h() time.Duration {
	return ts.tillNext(SECONDS_IN_12H)
}

// TillNextNoon returns how much ns (as time.Duration) must be passed until
// next noon (12.00 PM) (for the current Timestamp 'ts') will came.
func (ts Timestamp) TillNextNoon() time.Duration {
	d := ts.tillNext(SECONDS_IN_DAY) + 12 * time.Hour
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
	return ts.tillNext(SECONDS_IN_DAY)
}

// TillNextMonth returns how much ns (as time.Duration) must be passed until
// next month (for the current Timestamp 'ts') will came.
func (ts Timestamp) TillNextMonth() time.Duration {
	y, m, _ := dateFromUnix(ts)
	return ts.tillNext(InMonth(y, m))
}

// TillNextYear returns how much ns (as time.Duration) must be passed until
// next year (for the current Timestamp 'ts') will came.
func (ts Timestamp) TillNextYear() time.Duration {
	return ts.tillNext(InYear(ts.Year()))
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