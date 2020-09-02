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
