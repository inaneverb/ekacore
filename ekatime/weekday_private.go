// Copyright © 2020-2021. All rights reserved.
// Author: Ilya Stroy.
// Contacts: qioalice@gmail.com, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekatime

import (
	"bytes"
)

var (
	// _WeekdayStr is just English names of days of week.
	_WeekdayStr = [...]string {
		"Unknown",
		"Wednesday",
		"Thursday",
		"Friday",
		"Saturday",
		"Sunday",
		"Monday",
		"Tuesday",
	}

	_WeekdayBytes = [len(_WeekdayStr)][]byte{}
)

// asPartOfDate returns current Weekday but as part of Date object
// (casted, bit shifted, ready for being bit added to some Date object).
func (w Weekday) asPartOfDate() Date {
	return Date(w+1) << _DATE_OFFSET_WEEKDAY
}

// byteSliceEncode returns the current weekday's []byte representation,
// the same as String() returns but double quoted.
func (w Weekday) byteSliceEncode() []byte {
	if w < 0 || w > 6 {
		//noinspection GoAssignmentToReceiver
		w = -1
	}
	return _WeekdayBytes[w+1]
}

// byteSliceDecode decodes the weekday's value from 'data',
// saving it into the current weekday's object. Always returns nil.
// Saves -1 if 'data' does not contain valid weekday.
func (w *Weekday) byteSliceDecode(data []byte) error {
	if data != nil {
		for i, n := 1, len(_WeekdayBytes); i < n; i++ {
			if bytes.Equal(data, _WeekdayBytes[i]) {
				*w = Weekday(i-1)
				return nil
			}
		}
	}
	*w = -1
	return nil
}

// initWeekday fills weekdays []byte representation for fast encoding/decoding.
func initWeekday() {
	for i, n := 0, len(_WeekdayStr); i < n; i++ {
		_WeekdayBytes[i] = []byte("\"" + _WeekdayStr[i] + "\"")
	}
}
