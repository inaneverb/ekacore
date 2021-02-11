// Copyright © 2020. All rights reserved.
// Author: Ilya Stroy.
// Contacts: qioalice@gmail.com, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekatime_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/qioalice/ekago/v2/ekatime"

	"github.com/stretchr/testify/require"
)

func TestDate_String(t *testing.T) {
	d1 := ekatime.NewDate(2020, 9, 1)
	d2 := ekatime.NewDate(1812, 11, 24)
	d3 := ekatime.NewDate(2100, 2, 15)

	require.EqualValues(t, "2020/09/01", d1.String())
	require.EqualValues(t, "1812/11/24", d2.String())
	require.EqualValues(t, "2100/02/15", d3.String())
}

func BenchmarkDate_String_CachedYear(b *testing.B) {
	d := ekatime.NewDate(2020, 9, 1)
	b.ResetTimer(); b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = d.String()
	}
}

func BenchmarkDate_String_GeneratedYear(b *testing.B) {
	d := ekatime.NewDate(2212, 8, 8)
	b.ResetTimer(); b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = d.String()
	}
}

func BenchmarkDate_FmtSprintf(b *testing.B) {
	y, m, d := ekatime.NewDate(1914, 3, 9).Split()
	b.ResetTimer(); b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = fmt.Sprintf("%04d/%02d/%02d", y, m, d)
	}
}

func TestDate_ParseFrom(t *testing.T) {
	bd1 := []byte("1814/05/12") // valid
	bd2 := []byte("20200901")   // valid
	bd3 := []byte("2020-1120")  // invalid
	bd4 := []byte("202011-2")   // invalid
	bd5 := []byte("202011-20")  // invalid

	var (d1, d2, d3, d4, d5 ekatime.Date)

	err1 := d1.ParseFrom(bd1)
	err2 := d2.ParseFrom(bd2)
	err3 := d3.ParseFrom(bd3)
	err4 := d4.ParseFrom(bd4)
	err5 := d5.ParseFrom(bd5)

	require.Equal(t, ekatime.NewDate(1814, 5, 12).ToCmp(), d1.ToCmp())
	require.Equal(t, ekatime.NewDate(2020, 9, 1).ToCmp(), d2.ToCmp())

	require.NoError(t, err1)
	require.NoError(t, err2)
	require.Error(t, err3)
	require.Error(t, err4)
	require.Error(t, err5)
}

func BenchmarkDate_ParseFrom(b *testing.B) {
	bd0 := []byte("2020/09/01")
	var d ekatime.Date
	b.ResetTimer(); b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = d.ParseFrom(bd0)
	}
}

func TestDate_MarshalJSON(t *testing.T) {
	var dt = struct {
		D ekatime.Date `json:"d"`
	}{
		D: ekatime.NewDate(2020, 9, 12),
	}
	d, err := json.Marshal(&dt)

	require.NoError(t, err)
	require.EqualValues(t, string(d), `{"d":"2020-09-12"}`)
}

func TestDate_Replace(t *testing.T) {
	d := ekatime.NewDate(2021, ekatime.MONTH_FEBRUARY, 10)

	d = d.Replace(2013, ekatime.MONTH_JANUARY, 2)
	require.Equal(t, ekatime.NewDate(2013, ekatime.MONTH_JANUARY, 02).ToCmp(), d.ToCmp())

	d = d.Replace(2020, -2, 4)
	require.Equal(t, ekatime.NewDate(2020, ekatime.MONTH_JANUARY, 4).ToCmp(), d.ToCmp())

	d = d.Replace(0, 0, -50)
	require.Equal(t, ekatime.NewDate(2020, ekatime.MONTH_JANUARY, 4).ToCmp(), d.ToCmp())

	d = d.Replace(1899, 61, 31)
	require.Equal(t, ekatime.NewDate(2020, ekatime.MONTH_JANUARY, 31).ToCmp(), d.ToCmp())

	d = d.Replace(2013, ekatime.MONTH_FEBRUARY, -2)
	require.Equal(t, ekatime.NewDate(2013, ekatime.MONTH_JANUARY, 31).ToCmp(), d.ToCmp())

	d = d.Replace(2014, ekatime.MONTH_FEBRUARY, 30)
	require.Equal(t, ekatime.NewDate(2014, ekatime.MONTH_JANUARY, 30).ToCmp(), d.ToCmp())

	d = d.Replace(2000, ekatime.MONTH_FEBRUARY, 29)
	require.Equal(t, ekatime.NewDate(2000, ekatime.MONTH_FEBRUARY, 29).ToCmp(), d.ToCmp())

	d = d.Replace(2001, -1, -1)
	require.Equal(t, ekatime.NewDate(2000, ekatime.MONTH_FEBRUARY, 29).ToCmp(), d.ToCmp())

	d = d.Replace(0, 0, 13)
	require.Equal(t, ekatime.NewDate(2000, ekatime.MONTH_FEBRUARY, 13).ToCmp(), d.ToCmp())
}

func TestDate_Add(t *testing.T) {
	d := ekatime.NewDate(2021, ekatime.MONTH_FEBRUARY, 10)

	d = d.Add(1, 2, 3)
	require.Equal(t, ekatime.NewDate(2022, ekatime.MONTH_APRIL, 13).ToCmp(), d.ToCmp())

	d = d.Add(0, -1, -13)
	require.Equal(t, ekatime.NewDate(2022, ekatime.MONTH_FEBRUARY, 28).ToCmp(), d.ToCmp())

	d = d.Add(0, 1, 3)
	require.Equal(t, ekatime.NewDate(2022, ekatime.MONTH_MARCH, 31).ToCmp(), d.ToCmp())

	d = d.Add(0, 1, 0)
	require.Equal(t, ekatime.NewDate(2022, ekatime.MONTH_MAY, 1).ToCmp(), d.ToCmp())

	d = d.Add(0, 127, 0)
	require.Equal(t, ekatime.NewDate(2032, ekatime.MONTH_DECEMBER, 1).ToCmp(), d.ToCmp())
}

func TestDate_Days(t *testing.T) {
	var d ekatime.Date

	d = ekatime.NewDate(2021, ekatime.MONTH_FEBRUARY, 11)
	require.EqualValues(t, 42, d.Days())

	d = ekatime.NewDate(2021, ekatime.MONTH_DECEMBER, 31)
	require.EqualValues(t, 365, d.Days())

	d = ekatime.NewDate(2021, ekatime.MONTH_JANUARY, 1)
	require.EqualValues(t, 1, d.Days())

	d = ekatime.NewDate(2021, ekatime.MONTH_JANUARY, 0)
	require.EqualValues(t, 366, d.Days())
}
