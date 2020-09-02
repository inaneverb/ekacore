// Copyright © 2020. All rights reserved.
// Author: Ilya Stroy.
// Contacts: qioalice@gmail.com, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekatime

import "fmt"

type (
	// Hour is a special type that has enough space to store Hour's number.
	// Useless just by yourself but is a part of Time object.
	// Valid values: [0..23].
	Hour int8

	// Minute is a special type that has enough space to store Minute's number.
	// Useless just by yourself but is a part of Time object.
	// Valid values: [0..59].
	Minute int8

	// Second is a special type that has enough space to store Second's number.
	// Useless just by yourself but is a part of Time object.
	// Valid values: [0..59].
	Second int8

	// Time is a special object that has enough space to store a some Time (Clock)
	// (including Hour, Minute, Second) but does it most RAM efficient way
	// taking only 4 bytes.
	Time uint32
)

// Hour returns the hour number the current Time includes which.
//
// It guarantees that Hour() returns the valid hour number Time is of,
// only if Time has not been created manually but using constructors
// like NewTime(), Timestamp.Time(), Timestamp.Split(), etc.
func (t Time) Hour() Hour {
	return Hour(t >> _TIME_OFFSET_HOUR) & _TIME_MASK_HOUR
}

// Minute returns the minute number the current Time includes which.
//
// It guarantees that Minute() returns the valid minute number Time is of,
// only if Time has not been created manually but using constructors
// like NewTime(), Timestamp.Time(), Timestamp.Split(), etc.
func (t Time) Minute() Minute {
	return Minute(t >> _TIME_OFFSET_MINUTE) & _TIME_MASK_MINUTE
}

// Second returns the second number the current Time includes which.
//
// It guarantees that Second() returns the valid hour number Time is of,
// only if Time has not been created manually but using constructors
// like NewTime(), Timestamp.Time(), Timestamp.Split(), etc.
func (t Time) Second() Second {
	return Second(t >> _TIME_OFFSET_SECOND) & _TIME_MASK_SECOND
}

// Split returns the hour number, minutes number and seconds number the current Date
// includes which.
// It's just like a separate Hour(), Minute(), Second() calls.
func (t Time) Split() (h Hour, m Minute, s Second) {
	return t.Hour(), t.Minute(), t.Second()
}

// NewTime creates a new Time object using provided hour number, minute number,
// second number, normalizing these values and shifting time if it's required.
//
// Totally, it just stores the provided data if values are in their valid ranges,
// like: [0..23] for hour, [0..59] for minutes and seconds.
//
// If they are not, the time may be (will be) shifted. E.g:
// 21:02:64 (h == 21, m == 2, s == 64) -> 21:03:04 (h == 21, m == 3, s == 4).
func NewTime(h Hour, m Minute, s Second) Time {

	h, m, s = normalizeTime(h, m, s)

	// Do not forgot 'bitwise AND' between h, m, s and their bitmasks
	// if you will change the logic of getting valid h, m, s from invalid ones.
	// Now it's unnecessary (redundant).
	//h, m, s = h & _TIME_MASK_HOUR, m & _TIME_MASK_MINUTE, s & _TIME_MASK_SECOND

	return (Time(h) << _TIME_OFFSET_HOUR) |
		(Time(m) << _TIME_OFFSET_MINUTE) |
		(Time(s) << _TIME_OFFSET_SECOND)
}

// String returns the current Time's string representation in the following format:
// "hh:mm:ss".
func (t Time) String() string {
	return fmt.Sprintf("%02d:%02d:%02d", t.Hour(), t.Minute(), t.Second())
}
