package ekatime

type (
	//
	Today struct {

		c                *Calendar   `json:"-"`

		Timestamp        Timestamp   `json:"ts"`

		Time             Time        `json:"-"`
		Date             Date        `json:"-"`

		Year             Year        `json:"year"`
		Month            Month       `json:"month"`
		Day              Day         `json:"day"`

		Weekday          Weekday     `json:"weekday"`

		Hour             Hour        `json:"-"`
		Minute           Minute      `json:"-"`
		Second           Second      `json:"-"`

		WorkDayCurrent   Day         `json:"work_day_current"`
		IsDayOff         bool        `json:"is_dayoff"`

		DaysInMonth      Day         `json:"days_in_month"`
		WorkDayTotal     Day         `json:"work_day_total"`
		DayOffTotal      Day         `json:"dayoff_total"`

		AsJson           []byte      `json:"-"`
		AsYourOwn1       []byte      `json:"-"`
	}
)

// Copy returns the current Today's object copy but with the same encoded data
// (only pointers are copied, not internal data).
// If you want to copy them too, use CopyWithEncodedData() instead.
func (t *Today) Copy() *Today {
	cp := *t
	return &cp
}

// CopyWithEncodedData returns the current Today's object copy even with encoded data.
// Yes, a new buffers are allocated and the encoded data explicitly copied.
func (t *Today) CopyWithEncodedData() *Today {
	cp := t.Copy()
	// https://github.com/go101/go101/wiki/How-to-perfectly-clone-a-slice
	cp.AsJson = append(t.AsJson[:0:0], t.AsJson...)
	cp.AsYourOwn1 = append(t.AsYourOwn1[:0:0], t.AsYourOwn1...)
	return cp
}
