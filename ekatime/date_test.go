// Copyright © 2020. All rights reserved.
// Author: Ilya Stroy.
// Contacts: qioalice@gmail.com, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekatime_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/qioalice/ekago/v3/ekatime"

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

	d = ekatime.NewDate(2021, ekatime.MONTH_SEPTEMBER, 10)
	require.EqualValues(t, 253, d.Days())
}

var TestDateWeekNumber = []struct{
	y ekatime.Year
	m ekatime.Month
	d ekatime.Day
	expectedWeekNumber ekatime.WeekNumber
}{
	{ 2004, 12, 27, 53 },
	{ 2004, 12, 28, 53 },
	{ 2004, 12, 29, 53 },
	{ 2004, 12, 30, 53 },
	{ 2004, 12, 31, 53 },
	{ 2005, 1, 1, 53 },
	{ 2005, 1, 2, 53 },
	{ 2005, 1, 3, 1 },
	{ 2005, 1, 4, 1 },
	{ 2005, 1, 5, 1 },
	{ 2005, 1, 6, 1 },
	{ 2005, 12, 27, 52 },
	{ 2005, 12, 28, 52 },
	{ 2005, 12, 29, 52 },
	{ 2005, 12, 30, 52 },
	{ 2005, 12, 31, 52 },
	{ 2006, 1, 1, 52 },
	{ 2006, 1, 2, 1 },
	{ 2006, 1, 3, 1 },
	{ 2006, 1, 4, 1 },
	{ 2006, 1, 5, 1 },
	{ 2006, 1, 6, 1 },
	{ 2006, 12, 27, 52 },
	{ 2006, 12, 28, 52 },
	{ 2006, 12, 29, 52 },
	{ 2006, 12, 30, 52 },
	{ 2006, 12, 31, 52 },
	{ 2007, 1, 1, 1 },
	{ 2007, 1, 2, 1 },
	{ 2007, 1, 3, 1 },
	{ 2007, 1, 4, 1 },
	{ 2007, 1, 5, 1 },
	{ 2007, 1, 6, 1 },
	{ 2007, 12, 27, 52 },
	{ 2007, 12, 28, 52 },
	{ 2007, 12, 29, 52 },
	{ 2007, 12, 30, 52 },
	{ 2007, 12, 31, 1 },
	{ 2008, 1, 1, 1 },
	{ 2008, 1, 2, 1 },
	{ 2008, 1, 3, 1 },
	{ 2008, 1, 4, 1 },
	{ 2008, 1, 5, 1 },
	{ 2008, 1, 6, 1 },
	{ 2008, 12, 26, 52 },
	{ 2008, 12, 27, 52 },
	{ 2008, 12, 28, 52 },
	{ 2008, 12, 29, 1 },
	{ 2008, 12, 30, 1 },
	{ 2008, 12, 31, 1 },
	{ 2009, 1, 1, 1 },
	{ 2009, 1, 2, 1 },
	{ 2009, 1, 3, 1 },
	{ 2009, 1, 4, 1 },
	{ 2009, 1, 5, 2 },
	{ 2009, 1, 6, 2 },
	{ 2009, 12, 27, 52 },
	{ 2009, 12, 28, 53 },
	{ 2009, 12, 29, 53 },
	{ 2009, 12, 30, 53 },
	{ 2009, 12, 31, 53 },
	{ 2010, 1, 1, 53 },
	{ 2010, 1, 2, 53 },
	{ 2010, 1, 3, 53 },
	{ 2010, 1, 4, 1 },
	{ 2010, 1, 5, 1 },
	{ 2010, 1, 6, 1 },
	{ 2010, 12, 27, 52 },
	{ 2010, 12, 28, 52 },
	{ 2010, 12, 29, 52 },
	{ 2010, 12, 30, 52 },
	{ 2010, 12, 31, 52 },
	{ 2011, 1, 1, 52 },
	{ 2011, 1, 2, 52 },
	{ 2011, 1, 3, 1 },
	{ 2011, 1, 4, 1 },
	{ 2011, 1, 5, 1 },
	{ 2011, 1, 6, 1 },
	{ 2011, 12, 27, 52 },
	{ 2011, 12, 28, 52 },
	{ 2011, 12, 29, 52 },
	{ 2011, 12, 30, 52 },
	{ 2011, 12, 31, 52 },
	{ 2012, 1, 1, 52 },
	{ 2012, 1, 2, 1 },
	{ 2012, 1, 3, 1 },
	{ 2012, 1, 4, 1 },
	{ 2012, 1, 5, 1 },
	{ 2012, 1, 6, 1 },
	{ 2012, 12, 26, 52 },
	{ 2012, 12, 27, 52 },
	{ 2012, 12, 28, 52 },
	{ 2012, 12, 29, 52 },
	{ 2012, 12, 30, 52 },
	{ 2012, 12, 31, 1 },
	{ 2013, 1, 1, 1 },
	{ 2013, 1, 2, 1 },
	{ 2013, 1, 3, 1 },
	{ 2013, 1, 4, 1 },
	{ 2013, 1, 5, 1 },
	{ 2013, 1, 6, 1 },
	{ 2013, 12, 27, 52 },
	{ 2013, 12, 28, 52 },
	{ 2013, 12, 29, 52 },
	{ 2013, 12, 30, 1 },
	{ 2013, 12, 31, 1 },
	{ 2014, 1, 1, 1 },
	{ 2014, 1, 2, 1 },
	{ 2014, 1, 3, 1 },
	{ 2014, 1, 4, 1 },
	{ 2014, 1, 5, 1 },
	{ 2014, 1, 6, 2 },
	{ 2014, 12, 27, 52 },
	{ 2014, 12, 28, 52 },
	{ 2014, 12, 29, 1 },
	{ 2014, 12, 30, 1 },
	{ 2014, 12, 31, 1 },
	{ 2015, 1, 1, 1 },
	{ 2015, 1, 2, 1 },
	{ 2015, 1, 3, 1 },
	{ 2015, 1, 4, 1 },
	{ 2015, 1, 5, 2 },
	{ 2015, 1, 6, 2 },
	{ 2015, 12, 27, 52 },
	{ 2015, 12, 28, 53 },
	{ 2015, 12, 29, 53 },
	{ 2015, 12, 30, 53 },
	{ 2015, 12, 31, 53 },
	{ 2016, 1, 1, 53 },
	{ 2016, 1, 2, 53 },
	{ 2016, 1, 3, 53 },
	{ 2016, 1, 4, 1 },
	{ 2016, 1, 5, 1 },
	{ 2016, 1, 6, 1 },
	{ 2016, 12, 26, 52 },
	{ 2016, 12, 27, 52 },
	{ 2016, 12, 28, 52 },
	{ 2016, 12, 29, 52 },
	{ 2016, 12, 30, 52 },
	{ 2016, 12, 31, 52 },
	{ 2017, 1, 1, 52 },
	{ 2017, 1, 2, 1 },
	{ 2017, 1, 3, 1 },
	{ 2017, 1, 4, 1 },
	{ 2017, 1, 5, 1 },
	{ 2017, 1, 6, 1 },
	{ 2017, 12, 27, 52 },
	{ 2017, 12, 28, 52 },
	{ 2017, 12, 29, 52 },
	{ 2017, 12, 30, 52 },
	{ 2017, 12, 31, 52 },
	{ 2018, 1, 1, 1 },
	{ 2018, 1, 2, 1 },
	{ 2018, 1, 3, 1 },
	{ 2018, 1, 4, 1 },
	{ 2018, 1, 5, 1 },
	{ 2018, 1, 6, 1 },
	{ 2018, 12, 27, 52 },
	{ 2018, 12, 28, 52 },
	{ 2018, 12, 29, 52 },
	{ 2018, 12, 30, 52 },
	{ 2018, 12, 31, 1 },
	{ 2019, 1, 1, 1 },
	{ 2019, 1, 2, 1 },
	{ 2019, 1, 3, 1 },
	{ 2019, 1, 4, 1 },
	{ 2019, 1, 5, 1 },
	{ 2019, 1, 6, 1 },
	{ 2019, 12, 27, 52 },
	{ 2019, 12, 28, 52 },
	{ 2019, 12, 29, 52 },
	{ 2019, 12, 30, 1 },
	{ 2019, 12, 31, 1 },
	{ 2020, 1, 1, 1 },
	{ 2020, 1, 2, 1 },
	{ 2020, 1, 3, 1 },
	{ 2020, 1, 4, 1 },
	{ 2020, 1, 5, 1 },
	{ 2020, 1, 6, 2 },
	{ 2020, 12, 26, 52 },
	{ 2020, 12, 27, 52 },
	{ 2020, 12, 28, 53 },
	{ 2020, 12, 29, 53 },
	{ 2020, 12, 30, 53 },
	{ 2020, 12, 31, 53 },
	{ 2021, 1, 1, 53 },
	{ 2021, 1, 2, 53 },
	{ 2021, 1, 3, 53 },
	{ 2021, 1, 4, 1 },
	{ 2021, 1, 5, 1 },
	{ 2021, 1, 6, 1 },
	{ 2021, 12, 27, 52 },
	{ 2021, 12, 28, 52 },
	{ 2021, 12, 29, 52 },
	{ 2021, 12, 30, 52 },
	{ 2021, 12, 31, 52 },
	{ 2022, 1, 1, 52 },
	{ 2022, 1, 2, 52 },
	{ 2022, 1, 3, 1 },
	{ 2022, 1, 4, 1 },
	{ 2022, 1, 5, 1 },
	{ 2022, 1, 6, 1 },
	{ 2022, 12, 27, 52 },
	{ 2022, 12, 28, 52 },
	{ 2022, 12, 29, 52 },
	{ 2022, 12, 30, 52 },
	{ 2022, 12, 31, 52 },
	{ 2023, 1, 1, 52 },
	{ 2023, 1, 2, 1 },
	{ 2023, 1, 3, 1 },
	{ 2023, 1, 4, 1 },
	{ 2023, 1, 5, 1 },
	{ 2023, 1, 6, 1 },
	{ 2023, 12, 27, 52 },
	{ 2023, 12, 28, 52 },
	{ 2023, 12, 29, 52 },
	{ 2023, 12, 30, 52 },
	{ 2023, 12, 31, 52 },
	{ 2024, 1, 1, 1 },
	{ 2024, 1, 2, 1 },
	{ 2024, 1, 3, 1 },
	{ 2024, 1, 4, 1 },
	{ 2024, 1, 5, 1 },
	{ 2024, 1, 6, 1 },
	{ 2024, 12, 26, 52 },
	{ 2024, 12, 27, 52 },
	{ 2024, 12, 28, 52 },
	{ 2024, 12, 29, 52 },
	{ 2024, 12, 30, 1 },
	{ 2024, 12, 31, 1 },
	{ 2025, 1, 1, 1 },
	{ 2025, 1, 2, 1 },
	{ 2025, 1, 3, 1 },
	{ 2025, 1, 4, 1 },
	{ 2025, 1, 5, 1 },
	{ 2025, 1, 6, 2 },
	{ 2025, 12, 27, 52 },
	{ 2025, 12, 28, 52 },
	{ 2025, 12, 29, 1 },
	{ 2025, 12, 30, 1 },
	{ 2025, 12, 31, 1 },
	{ 2026, 1, 1, 1 },
	{ 2026, 1, 2, 1 },
	{ 2026, 1, 3, 1 },
	{ 2026, 1, 4, 1 },
	{ 2026, 1, 5, 2 },
	{ 2026, 1, 6, 2 },
	{ 2026, 12, 27, 52 },
	{ 2026, 12, 28, 53 },
	{ 2026, 12, 29, 53 },
	{ 2026, 12, 30, 53 },
	{ 2026, 12, 31, 53 },
	{ 2027, 1, 1, 53 },
	{ 2027, 1, 2, 53 },
	{ 2027, 1, 3, 53 },
	{ 2027, 1, 4, 1 },
	{ 2027, 1, 5, 1 },
	{ 2027, 1, 6, 1 },
	{ 2027, 12, 27, 52 },
	{ 2027, 12, 28, 52 },
	{ 2027, 12, 29, 52 },
	{ 2027, 12, 30, 52 },
	{ 2027, 12, 31, 52 },
	{ 2028, 1, 1, 52 },
	{ 2028, 1, 2, 52 },
	{ 2028, 1, 3, 1 },
	{ 2028, 1, 4, 1 },
	{ 2028, 1, 5, 1 },
	{ 2028, 1, 6, 1 },
	{ 2028, 12, 26, 52 },
	{ 2028, 12, 27, 52 },
	{ 2028, 12, 28, 52 },
	{ 2028, 12, 29, 52 },
	{ 2028, 12, 30, 52 },
	{ 2028, 12, 31, 52 },
	{ 2029, 1, 1, 1 },
	{ 2029, 1, 2, 1 },
	{ 2029, 1, 3, 1 },
	{ 2029, 1, 4, 1 },
	{ 2029, 1, 5, 1 },
	{ 2029, 1, 6, 1 },
	{ 2029, 12, 27, 52 },
	{ 2029, 12, 28, 52 },
	{ 2029, 12, 29, 52 },
	{ 2029, 12, 30, 52 },
	{ 2029, 12, 31, 1 },
	{ 2030, 1, 1, 1 },
	{ 2030, 1, 2, 1 },
	{ 2030, 1, 3, 1 },
	{ 2030, 1, 4, 1 },
	{ 2030, 1, 5, 1 },
	{ 2030, 1, 6, 1 },
	{ 2030, 12, 27, 52 },
	{ 2030, 12, 28, 52 },
	{ 2030, 12, 29, 52 },
	{ 2030, 12, 30, 1 },
	{ 2030, 12, 31, 1 },
	{ 2031, 1, 1, 1 },
	{ 2031, 1, 2, 1 },
	{ 2031, 1, 3, 1 },
	{ 2031, 1, 4, 1 },
	{ 2031, 1, 5, 1 },
	{ 2031, 1, 6, 2 },
	{ 2031, 12, 27, 52 },
	{ 2031, 12, 28, 52 },
	{ 2031, 12, 29, 1 },
	{ 2031, 12, 30, 1 },
	{ 2031, 12, 31, 1 },
	{ 2032, 1, 1, 1 },
	{ 2032, 1, 2, 1 },
	{ 2032, 1, 3, 1 },
	{ 2032, 1, 4, 1 },
	{ 2032, 1, 5, 2 },
	{ 2032, 1, 6, 2 },
	{ 2032, 12, 26, 52 },
	{ 2032, 12, 27, 53 },
	{ 2032, 12, 28, 53 },
	{ 2032, 12, 29, 53 },
	{ 2032, 12, 30, 53 },
	{ 2032, 12, 31, 53 },
	{ 2033, 1, 1, 53 },
	{ 2033, 1, 2, 53 },
	{ 2033, 1, 3, 1 },
	{ 2033, 1, 4, 1 },
	{ 2033, 1, 5, 1 },
	{ 2033, 1, 6, 1 },
	{ 2033, 12, 27, 52 },
	{ 2033, 12, 28, 52 },
	{ 2033, 12, 29, 52 },
	{ 2033, 12, 30, 52 },
	{ 2033, 12, 31, 52 },
	{ 2034, 1, 1, 52 },
	{ 2034, 1, 2, 1 },
	{ 2034, 1, 3, 1 },
	{ 2034, 1, 4, 1 },
	{ 2034, 1, 5, 1 },
	{ 2034, 1, 6, 1 },
	{ 2034, 12, 27, 52 },
	{ 2034, 12, 28, 52 },
	{ 2034, 12, 29, 52 },
	{ 2034, 12, 30, 52 },
	{ 2034, 12, 31, 52 },
	{ 2035, 1, 1, 1 },
	{ 2035, 1, 2, 1 },
	{ 2035, 1, 3, 1 },
	{ 2035, 1, 4, 1 },
	{ 2035, 1, 5, 1 },
	{ 2035, 1, 6, 1 },
}

func TestDate_ISOWeek(t *testing.T) {
	for _, n := range TestDateWeekNumber {
		require.Equal(t, n.expectedWeekNumber, ekatime.NewDate(n.y, n.m, n.d).ISOWeek())
	}
}

func TestNewDateFromDays(t *testing.T) {
	require.Equal(t,
		ekatime.NewDate(2021, ekatime.MONTH_SEPTEMBER, 10).ToCmp(),
		ekatime.NewDateFromDays(2021, 253).ToCmp(),
	)
	require.EqualValues(t,
		253,
		ekatime.NewDateFromDays(2021, 253).Days(),
	)
}
