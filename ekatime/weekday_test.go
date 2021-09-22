// Copyright © 2020-2021. All rights reserved.
// Author: Ilya Stroy.
// Contacts: iyuryevich@pm.me, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekatime_test

import (
	"testing"

	"github.com/qioalice/ekago/v3/ekatime"

	"github.com/stretchr/testify/require"
)

func TestWeekday_To06(t *testing.T) {
	for _, n := range []struct {
		wd         ekatime.Weekday
		expected06 int8
	}{
		{ekatime.WEEKDAY_MONDAY, 1},
		{ekatime.WEEKDAY_TUESDAY, 2},
		{ekatime.WEEKDAY_WEDNESDAY, 3},
		{ekatime.WEEKDAY_THURSDAY, 4},
		{ekatime.WEEKDAY_FRIDAY, 5},
		{ekatime.WEEKDAY_SATURDAY, 6},
		{ekatime.WEEKDAY_SUNDAY, 0},
	} {
		require.Equal(t, n.expected06, n.wd.To06())
	}
}

func TestWeekday_From06(t *testing.T) {
	for _, n := range []struct {
		_06             int8
		expectedWeekday ekatime.Weekday
	}{
		{1, ekatime.WEEKDAY_MONDAY},
		{2, ekatime.WEEKDAY_TUESDAY},
		{3, ekatime.WEEKDAY_WEDNESDAY},
		{4, ekatime.WEEKDAY_THURSDAY},
		{5, ekatime.WEEKDAY_FRIDAY},
		{6, ekatime.WEEKDAY_SATURDAY},
		{0, ekatime.WEEKDAY_SUNDAY},
	} {
		require.Equal(t, n.expectedWeekday, *new(ekatime.Weekday).From06(n._06))
	}
}

func TestWeekdayJan1(t *testing.T) {
	for _, n := range []struct {
		y               ekatime.Year
		expectedWeekday ekatime.Weekday
	}{
		{1500, ekatime.WEEKDAY_MONDAY},
		{1600, ekatime.WEEKDAY_SATURDAY},
		{1700, ekatime.WEEKDAY_FRIDAY},
		{1800, ekatime.WEEKDAY_WEDNESDAY},
		{1900, ekatime.WEEKDAY_MONDAY},
		{2000, ekatime.WEEKDAY_SATURDAY},
		{2100, ekatime.WEEKDAY_FRIDAY},
		{2200, ekatime.WEEKDAY_WEDNESDAY},
		{2300, ekatime.WEEKDAY_MONDAY},
		{2400, ekatime.WEEKDAY_SATURDAY},
		{2500, ekatime.WEEKDAY_FRIDAY},
		{2600, ekatime.WEEKDAY_WEDNESDAY},
		{2700, ekatime.WEEKDAY_MONDAY},
		{2800, ekatime.WEEKDAY_SATURDAY},
		{2900, ekatime.WEEKDAY_FRIDAY},
		{1926, ekatime.WEEKDAY_FRIDAY},
		{1927, ekatime.WEEKDAY_SATURDAY},
		{1928, ekatime.WEEKDAY_SUNDAY},
		{1929, ekatime.WEEKDAY_TUESDAY},
		{1930, ekatime.WEEKDAY_WEDNESDAY},
		{1931, ekatime.WEEKDAY_THURSDAY},
		{1932, ekatime.WEEKDAY_FRIDAY},
		{1933, ekatime.WEEKDAY_SUNDAY},
		{1934, ekatime.WEEKDAY_MONDAY},
		{1935, ekatime.WEEKDAY_TUESDAY},
		{1936, ekatime.WEEKDAY_WEDNESDAY},
		{1937, ekatime.WEEKDAY_FRIDAY},
		{1938, ekatime.WEEKDAY_SATURDAY},
		{1939, ekatime.WEEKDAY_SUNDAY},
		{1940, ekatime.WEEKDAY_MONDAY},
		{1941, ekatime.WEEKDAY_WEDNESDAY},
		{1942, ekatime.WEEKDAY_THURSDAY},
		{1943, ekatime.WEEKDAY_FRIDAY},
		{1944, ekatime.WEEKDAY_SATURDAY},
		{1945, ekatime.WEEKDAY_MONDAY},
		{1946, ekatime.WEEKDAY_TUESDAY},
		{1947, ekatime.WEEKDAY_WEDNESDAY},
		{1948, ekatime.WEEKDAY_THURSDAY},
		{1949, ekatime.WEEKDAY_SATURDAY},
		{1950, ekatime.WEEKDAY_SUNDAY},
		{1951, ekatime.WEEKDAY_MONDAY},
		{1952, ekatime.WEEKDAY_TUESDAY},
		{1953, ekatime.WEEKDAY_THURSDAY},
		{1954, ekatime.WEEKDAY_FRIDAY},
		{1955, ekatime.WEEKDAY_SATURDAY},
		{1956, ekatime.WEEKDAY_SUNDAY},
		{1957, ekatime.WEEKDAY_TUESDAY},
		{1958, ekatime.WEEKDAY_WEDNESDAY},
		{1959, ekatime.WEEKDAY_THURSDAY},
		{1960, ekatime.WEEKDAY_FRIDAY},
		{1961, ekatime.WEEKDAY_SUNDAY},
		{1962, ekatime.WEEKDAY_MONDAY},
		{1963, ekatime.WEEKDAY_TUESDAY},
		{1964, ekatime.WEEKDAY_WEDNESDAY},
		{1965, ekatime.WEEKDAY_FRIDAY},
		{1966, ekatime.WEEKDAY_SATURDAY},
		{1967, ekatime.WEEKDAY_SUNDAY},
		{1968, ekatime.WEEKDAY_MONDAY},
		{1969, ekatime.WEEKDAY_WEDNESDAY},
		{1970, ekatime.WEEKDAY_THURSDAY},
		{1971, ekatime.WEEKDAY_FRIDAY},
		{1972, ekatime.WEEKDAY_SATURDAY},
		{1973, ekatime.WEEKDAY_MONDAY},
		{1974, ekatime.WEEKDAY_TUESDAY},
		{1975, ekatime.WEEKDAY_WEDNESDAY},
		{1976, ekatime.WEEKDAY_THURSDAY},
		{1977, ekatime.WEEKDAY_SATURDAY},
		{1978, ekatime.WEEKDAY_SUNDAY},
		{1979, ekatime.WEEKDAY_MONDAY},
		{1980, ekatime.WEEKDAY_TUESDAY},
		{1981, ekatime.WEEKDAY_THURSDAY},
		{1982, ekatime.WEEKDAY_FRIDAY},
		{1983, ekatime.WEEKDAY_SATURDAY},
		{1984, ekatime.WEEKDAY_SUNDAY},
		{1985, ekatime.WEEKDAY_TUESDAY},
		{1986, ekatime.WEEKDAY_WEDNESDAY},
		{1987, ekatime.WEEKDAY_THURSDAY},
		{1988, ekatime.WEEKDAY_FRIDAY},
		{1989, ekatime.WEEKDAY_SUNDAY},
		{1990, ekatime.WEEKDAY_MONDAY},
		{1991, ekatime.WEEKDAY_TUESDAY},
		{1992, ekatime.WEEKDAY_WEDNESDAY},
		{1993, ekatime.WEEKDAY_FRIDAY},
		{1994, ekatime.WEEKDAY_SATURDAY},
		{1995, ekatime.WEEKDAY_SUNDAY},
		{1996, ekatime.WEEKDAY_MONDAY},
		{1997, ekatime.WEEKDAY_WEDNESDAY},
		{1998, ekatime.WEEKDAY_THURSDAY},
		{1999, ekatime.WEEKDAY_FRIDAY},
		{2000, ekatime.WEEKDAY_SATURDAY},
		{2001, ekatime.WEEKDAY_MONDAY},
		{2002, ekatime.WEEKDAY_TUESDAY},
		{2003, ekatime.WEEKDAY_WEDNESDAY},
		{2004, ekatime.WEEKDAY_THURSDAY},
		{2005, ekatime.WEEKDAY_SATURDAY},
		{2006, ekatime.WEEKDAY_SUNDAY},
		{2007, ekatime.WEEKDAY_MONDAY},
		{2008, ekatime.WEEKDAY_TUESDAY},
		{2009, ekatime.WEEKDAY_THURSDAY},
		{2010, ekatime.WEEKDAY_FRIDAY},
		{2011, ekatime.WEEKDAY_SATURDAY},
		{2012, ekatime.WEEKDAY_SUNDAY},
		{2013, ekatime.WEEKDAY_TUESDAY},
		{2014, ekatime.WEEKDAY_WEDNESDAY},
		{2015, ekatime.WEEKDAY_THURSDAY},
		{2016, ekatime.WEEKDAY_FRIDAY},
		{2017, ekatime.WEEKDAY_SUNDAY},
		{2018, ekatime.WEEKDAY_MONDAY},
		{2019, ekatime.WEEKDAY_TUESDAY},
		{2020, ekatime.WEEKDAY_WEDNESDAY},
		{2021, ekatime.WEEKDAY_FRIDAY},
		{2022, ekatime.WEEKDAY_SATURDAY},
		{2023, ekatime.WEEKDAY_SUNDAY},
		{2024, ekatime.WEEKDAY_MONDAY},
		{2025, ekatime.WEEKDAY_WEDNESDAY},
		{2026, ekatime.WEEKDAY_THURSDAY},
	} {
		require.Equal(t, n.expectedWeekday, ekatime.WeekdayJan1(n.y))
	}
}
