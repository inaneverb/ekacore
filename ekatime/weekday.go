package ekatime

type (
	// Weekday is a special type that has enough space to store weekday's number.
	// Useless just by yourself but is a part of Date object.
	// Valid values: [0..6]. Use predefined constants to make it clear.
	Weekday int8
)

//noinspection GoSnakeCaseUsage,GoUnusedConst
const (
	WEEKDAY_WEDNESDAY Weekday = iota
	WEEKDAY_THURSDAY
	WEEKDAY_FRIDAY
	WEEKDAY_SATURDAY
	WEEKDAY_SUNDAY
	WEEKDAY_MONDAY
	WEEKDAY_TUESDAY
)

// WeekdayJan1 returns a Weekday of 1 January of requested Year.
// Returns -1 if passed Year is not in the range [1..4095].
//
// Yes, although other functions and methods declares 1900 as lower bound of Year,
// here you can use even lower than 1900.
func WeekdayJan1(y Year) Weekday {

	// Algorithm from:
	// https://people.cs.nctu.edu.tw/~tsaiwn/sisc/runtime_error_200_div_by_0/www.merlyn.demon.co.uk/weekcalc.htm
	// (if not available, then thank you, Mr Zeller)
	// https://en.wikipedia.org/wiki/Christian_Zeller

	if y <= 0 || y > _YEAR_MAX {
		return -1
	} else if y >= 1901 && y <= 2099 {
		yf := float32(y)
		return WeekdayFrom06(int8(int16(1.25 * (yf-1)) % 7))
	} else {
		y--
		yf := float32(y)
		return WeekdayFrom06(int8(int16((1.25 * yf) - (yf/100) + (yf/400) +1) % 7))
	}
}

// Next returns the next day of week following the current one.
func (w Weekday) Next() Weekday {
	if w == WEEKDAY_TUESDAY {
		return WEEKDAY_WEDNESDAY
	}
	return w+1
}

// Prev returns the prev day of week before the current one.
func (w Weekday) Prev() Weekday {
	if w == WEEKDAY_WEDNESDAY {
		return WEEKDAY_TUESDAY
	}
	return w-1
}

// IsDayOff reports whether current day of week is Saturday or Sunday.
func (w Weekday) IsDayOff() bool {
	return w == WEEKDAY_SATURDAY || w == WEEKDAY_SUNDAY
}

// String returns the current Weekday's string representation as full day of week
// name, capitalized.
func (w Weekday) String() string {
	if w < 0 || w > 6 {
		//noinspection GoAssignmentToReceiver
		w = -1
	}
	return _WeekdayStr[w+1]
}

// MarshalJSON returns the current Weekday's string representation and nil (always).
// It will be the same as String() returns (but double quoted).
func (w Weekday) MarshalJSON() ([]byte, error) {
	return w.byteSliceEncode(), nil
}

// UnmarshalJSON tries to decode 'data' and treat it as encoded weekday.
// If it's so, saves the weekday into the current Weekday's object.
// If 'data' contains an invalid data, -1 will be stored into the current Weekday.
// Always returns nil.
func (w *Weekday) UnmarshalJSON(data []byte) error {
	return w.byteSliceDecode(data)
}

// From06 parses passed i8, transforms it to the Weekday, saves to the current object,
// and returns it.
// Requires i8 be in the range [0..6], where 0 - Sunday, 6 - Saturday.
// Does nothing if w == nil, but returns it.
func (w *Weekday) From06(i8 int8) *Weekday {
	if w != nil && i8 >= 0 && i8 <= 6 {
		*w = WeekdayFrom06(i8)
	}
	return w
}

// WeekdayFrom06 parses passed i8, transforms it to the Weekday and returns.
// Requires i8 be in the range [0..6], where 0 - Sunday, 6 - Saturday.
// Returns -1 if i8 is not in the allowed range.
func WeekdayFrom06(i8 int8) Weekday {
	if i8 >= 0 && i8 <= 6 {
		return Weekday(i8 + 4) % 7
	} else {
		return -1
	}
}

// To06 returns a uint8 representation of the current Weekday,
// where 0 - Sunday, 6 - Saturday.
func (w Weekday) To06() int8 {
	return int8(w + 3) % 7
}
