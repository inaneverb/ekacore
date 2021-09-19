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
	SECONDS_IN_MINUTE Timestamp = 60
	SECONDS_IN_HOUR   Timestamp = 3600
	SECONDS_IN_12H    Timestamp = 43200
	SECONDS_IN_DAY    Timestamp = 86400
	SECONDS_IN_WEEK   Timestamp = 604800

	SECONDS_IN_365_YEAR Timestamp = 31536000
	SECONDS_IN_366_YEAR Timestamp = 31622400
)

var (
	// Table of:
	//   The number of days in month (not leap year).
	//   Consumes 12 bytes in RAM totally.
	_Table0 = [12]Day{31, 28, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31}

	// Table of:
	//   The number of seconds in month (not leap year).
	//   Consumes 96 bytes in RAM totally.
	_Table1 = [12]Timestamp{
		2678400, 2419200, 2678400, 2592000, 2678400, 2592000,
		2678400, 2678400, 2592000, 2678400, 2592000, 2678400,
	}

	// Table of:
	//   The 1st day of each month in the year pov.
	//   Consumes 24 bytes in RAM totally.
	//   It assumes that February has 29 days.
	_Table2 = [12]Days{1, 32, 61, 92, 122, 153, 183, 214, 245, 275, 306, 336}

	// Table of:
	//   The beginning of month in unix timestamp for [1970..now+10] years.
	//   Consumes 96 byte in RAM per year.
	//   "The beginning of month" is 1st day with 00:00:01 AM clock.
	//
	// Assuming, it's 2020 year now, the upper bound will be is 2030 year the
	// table is being generated for.
	// So, table will contain 2030-1970 == 60 years and will consume
	//
	// 5856 bytes in RAM totally (about 5.7 KB).
	_Table5           [][12]Timestamp
	_Table5UpperBound Year // real year, not an index in '_Table5'
)

func initTable5() {
	_Table5UpperBound = Year(time.Now().Year()) + 10
	_Table5 = make([][12]Timestamp, _Table5UpperBound-1970+1)

	d := time.Unix(1, 0)
	for i, n := 0, len(_Table5); i < n; i++ {
		for j := 0; j < 12; j++ {
			_Table5[i][j] = Timestamp(d.Unix())
			d = d.AddDate(0, 1, 0)
		}
	}
}
