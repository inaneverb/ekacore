// Copyright © 2021. All rights reserved.
// Author: Ilya Stroy.
// Contacts: qioalice@gmail.com, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekatime_test

import (
	"testing"

	"github.com/qioalice/ekago/v3/ekatime"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCalendar2015(t *testing.T) {

	const YEAR = ekatime.Year(2015)

	cal := ekatime.NewCalendar(YEAR, true, false)

	type (
		T1 struct {
			d        ekatime.Day
			isDayOff bool
		}
		T0 struct {
			m  ekatime.Month
			ds []T1
		}
	)

	for _, r := range []T0{
		{ekatime.MONTH_JANUARY, []T1{{1, true}, {2, true}, {5, true}, {6, true}, {7, true}, {8, true}, {9, true}}},
		{ekatime.MONTH_FEBRUARY, []T1{{23, true}}},
		{ekatime.MONTH_MARCH, []T1{{9, true}}},
		{ekatime.MONTH_MAY, []T1{{1, true}, {4, true}, {11, true}}},
		{ekatime.MONTH_JUNE, []T1{{12, true}}},
		{ekatime.MONTH_NOVEMBER, []T1{{4, true}}},
	} {
		for _, d := range r.ds {
			dd := ekatime.NewDate(YEAR, r.m, d.d)
			cal.OverrideDate(dd, d.isDayOff)

			assert.True(t, cal.IsDayOff(dd) == d.isDayOff,
				"DD: %s, DayOff: %t", dd.String(), d.isDayOff)
		}
	}

	workDays := []ekatime.Days{15, 19, 21, 22, 18, 21, 23, 21, 22, 22, 20, 23}
	for i, workDays := range workDays {
		m := ekatime.Month(i + 1)
		assert.Equal(t, workDays, cal.WorkDaysCount(m),
			"M: %s, WorkDays: %v", m.String(), cal.WorkDays(m))
	}

	daysOff := []ekatime.Days{16, 9, 10, 8, 13, 9, 8, 10, 8, 9, 10, 8}
	for i, daysOff := range daysOff {
		m := ekatime.Month(i + 1)
		assert.Equal(t, daysOff, cal.DaysOffCount(m),
			"M: %s, DaysOff: %v", m.String(), cal.DaysOff(m))
	}
}

func TestCalendar2020(t *testing.T) {

	const YEAR = ekatime.Year(2020)

	cal := ekatime.NewCalendar(YEAR, true, false)

	type (
		T1 struct {
			d        ekatime.Day
			isDayOff bool
		}
		T0 struct {
			m  ekatime.Month
			ds []T1
		}
	)

	for _, r := range []T0{
		{ekatime.MONTH_JANUARY, []T1{{1, true}, {2, true}, {3, true}, {6, true}, {7, true}, {8, true}}},
		{ekatime.MONTH_FEBRUARY, []T1{{24, true}}},
		{ekatime.MONTH_MARCH, []T1{{9, true}, {30, true}, {31, true}}},
		{ekatime.MONTH_MAY, []T1{{1, true}, {4, true}, {5, true}, {6, true}, {7, true}, {8, true}, {11, true}}},
		{ekatime.MONTH_JUNE, []T1{{12, true}, {24, true}}},
		{ekatime.MONTH_JULY, []T1{{1, true}}},
		{ekatime.MONTH_NOVEMBER, []T1{{4, true}}},
	} {
		for _, d := range r.ds {
			dd := ekatime.NewDate(YEAR, r.m, d.d)
			cal.OverrideDate(dd, d.isDayOff)

			assert.True(t, cal.IsDayOff(dd) == d.isDayOff,
				"DD: %s, DayOff: %t", dd.String(), d.isDayOff)
		}
	}

	for d, n := ekatime.Day(1), ekatime.Day(30); d <= n; d++ {
		dd := ekatime.NewDate(YEAR, ekatime.MONTH_APRIL, d)
		cal.OverrideDate(dd, true)

		assert.True(t, cal.IsDayOff(dd),
			"DD: %s, DayOff: %t", dd.String(), true)
	}

	workDays := []ekatime.Days{17, 19, 19, 0, 14, 20, 22, 21, 22, 22, 20, 23}
	for i, workDays := range workDays {
		m := ekatime.Month(i + 1)
		assert.Equal(t, workDays, cal.WorkDaysCount(m),
			"M: %s, WorkDays: %v", m.String(), cal.WorkDays(m))
	}

	daysOff := []ekatime.Days{14, 10, 12, 30, 17, 10, 9, 10, 8, 9, 10, 8}
	for i, daysOff := range daysOff {
		m := ekatime.Month(i + 1)
		assert.Equal(t, daysOff, cal.DaysOffCount(m),
			"M: %s, DaysOff: %v", m.String(), cal.DaysOff(m))
	}
}

func TestCalendar2021(t *testing.T) {

	const YEAR = ekatime.Year(2021)

	cal := ekatime.NewCalendar(YEAR, true, false)

	type (
		T1 struct {
			d        ekatime.Day
			isDayOff bool
		}
		T0 struct {
			m  ekatime.Month
			ds []T1
		}
	)

	for _, r := range []T0{
		{ekatime.MONTH_JANUARY, []T1{{1, true}, {4, true}, {5, true}, {6, true}, {7, true}, {8, true}}},
		{ekatime.MONTH_FEBRUARY, []T1{{20, false}, {22, true}, {23, true}}},
		{ekatime.MONTH_MARCH, []T1{{8, true}}},
		{ekatime.MONTH_MAY, []T1{{3, true}, {4, true}, {5, true}, {6, true}, {7, true}, {10, true}}},
		{ekatime.MONTH_JUNE, []T1{{14, true}}},
		{ekatime.MONTH_NOVEMBER, []T1{{4, true}, {5, true}}},
		{ekatime.MONTH_DECEMBER, []T1{{31, true}}},
	} {
		for _, d := range r.ds {
			dd := ekatime.NewDate(YEAR, r.m, d.d)
			cal.OverrideDate(dd, d.isDayOff)

			assert.True(t, cal.IsDayOff(dd) == d.isDayOff,
				"DD: %s, DayOff: %t", dd.String(), d.isDayOff)
		}
	}

	encodedBinary, err := cal.MarshalBinary()
	require.NoError(t, err)

	cal = ekatime.NewCalendar(1991, true, false)
	err = cal.UnmarshalBinary(encodedBinary)
	require.NoError(t, err)

	workDays := []ekatime.Days{15, 19, 22, 22, 15, 21, 22, 22, 22, 21, 20, 22}
	for i, workDays := range workDays {
		m := ekatime.Month(i + 1)
		assert.Equal(t, workDays, cal.WorkDaysCount(m),
			"M: %s, WorkDays: %v", m.String(), cal.WorkDays(m))
	}

	encodedText, err := cal.MarshalText()
	require.NoError(t, err)

	cal = ekatime.NewCalendar(2031, true, false)
	err = cal.UnmarshalText(encodedText)
	require.NoError(t, err)

	daysOff := []ekatime.Days{16, 9, 9, 8, 16, 9, 9, 9, 8, 10, 10, 9}
	for i, daysOff := range daysOff {
		m := ekatime.Month(i + 1)
		assert.Equal(t, daysOff, cal.DaysOffCount(m),
			"M: %s, DaysOff: %v", m.String(), cal.DaysOff(m))
	}
}
