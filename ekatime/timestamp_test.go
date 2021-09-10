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

// =========================================================================== //
// =========================================================================== //
// =========================================================================== //

func TestDate_Equal(t *testing.T) {
	d1 := ekatime.NewDate(2020, 9, 9)
	d2 := d1.WithTime(0, 0, 0).Date()
	require.False(t, d1 == d2)
	require.True(t, d1.ToCmp() == d2.ToCmp())
	require.True(t, d1.Equal(d2))
}

// =========================================================================== //
// =========================================================================== //
// =========================================================================== //

func Benchmark_BeginningOfYear_CachedYear(b *testing.B) {
	currYear := ekatime.NewTimestampNow().Year()
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ekatime.BeginningOfYear(currYear)
	}
}

func Benchmark_BeginningOfYear_NonCachedYear(b *testing.B) {
	notCachedYear := ekatime.NewTimestampNow().Year() + 20
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ekatime.BeginningOfYear(notCachedYear)
	}
}

func Benchmark_EndOfYear_CachedYear(b *testing.B) {
	currYear := ekatime.NewTimestampNow().Year()
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ekatime.EndOfYear(currYear)
	}
}

func Benchmark_EndOfYear_NonCachedYear(b *testing.B) {
	notCachedYear := ekatime.NewTimestampNow().Year() + 20
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ekatime.EndOfYear(notCachedYear)
	}
}

func Benchmark_BeginningAndEndOfYear_CachedYear(b *testing.B) {
	currYear := ekatime.NewTimestampNow().Year()
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ekatime.BeginningAndEndOfYear(currYear)
	}
}

func Benchmark_BeginningAndEndOfYear_NonCachedYear(b *testing.B) {
	notCachedYear := ekatime.NewTimestampNow().Year() + 20
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ekatime.BeginningAndEndOfYear(notCachedYear)
	}
}

// =========================================================================== //
// =========================================================================== //
// =========================================================================== //

func Benchmark_BeginningOfMonth_CachedYear(b *testing.B) {
	currYear := ekatime.NewTimestampNow().Year()
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ekatime.BeginningOfMonth(currYear, 12)
	}
}

func Benchmark_BeginningOfMonth_NonCachedYear(b *testing.B) {
	notCachedYear := ekatime.NewTimestampNow().Year() + 20
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ekatime.BeginningOfMonth(notCachedYear, 12)
	}
}

func Benchmark_EndOfMonth_CachedYear(b *testing.B) {
	currYear := ekatime.NewTimestampNow().Year()
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ekatime.EndOfMonth(currYear, 12)
	}
}

func Benchmark_EndOfMonth_NonCachedYear(b *testing.B) {
	notCachedYear := ekatime.NewTimestampNow().Year() + 20
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ekatime.EndOfMonth(notCachedYear, 12)
	}
}

func Benchmark_BeginningAndEndOfMonth_CachedYear(b *testing.B) {
	currYear := ekatime.NewTimestampNow().Year()
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ekatime.BeginningAndEndOfMonth(currYear, 12)
	}
}

func Benchmark_BeginningAndEndOfMonth_NonCachedYear(b *testing.B) {
	notCachedYear := ekatime.NewTimestampNow().Year() + 20
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ekatime.BeginningAndEndOfMonth(notCachedYear, 12)
	}
}

// =========================================================================== //
// =========================================================================== //
// =========================================================================== //

func TestTimestamp_String(t *testing.T) {
	ts1 := ekatime.NewTimestamp(2020, 9, 1,14, 2, 13)
	ts2 := ekatime.NewTimestamp(1812, 11, 24, 23, 59, 59)
	ts3 := ekatime.NewTimestamp(2100, 2, 15, 0, 0, 0)

	require.EqualValues(t, "2020/09/01 14:02:13", ts1.String())
	require.EqualValues(t, "1812/11/24 23:59:59", ts2.String())
	require.EqualValues(t, "2100/02/15 00:00:00", ts3.String())
}

func BenchmarkTimestamp_String_Cached(b *testing.B) {
	ts0 := ekatime.NewTimestamp(2020, 9, 1,14, 2, 13)
	b.ResetTimer(); b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = ts0.String()
	}
}

func BenchmarkTimestamp_String_FmtSprintf(b *testing.B) {
	dd, tt := ekatime.NewTimestamp(2020, 9, 1,14, 2, 13).Split()
	y, m, d := dd.Split()
	hh, mm, ss := tt.Split()
	b.ResetTimer(); b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = fmt.Sprintf("%04d/%02d/%02d %02d:%02d:%02d",y, m, d, hh, mm, ss)
	}
}

// =========================================================================== //
// =========================================================================== //
// =========================================================================== //

func TestTimestamp_ParseFrom(t *testing.T) {
	bd1 := []byte("1814/05/12 000000")   // valid
	bd2 := []byte("20200901")            // valid
	bd3 := []byte("20200901T23:59:59")   // valid
	bd4 := []byte("2020-09-01t12:13:14") // valid
	bd5 := []byte("2020-09-0112:13:14")  // invalid (w/o T)

	orig1 := ekatime.NewTimestamp(1814, 5, 12, 0, 0, 0)
	orig2 := ekatime.NewTimestamp(2020, 9, 1, 0, 0, 0)
	orig3 := ekatime.NewTimestamp(2020, 9, 1, 23, 59, 59)
	orig4 := ekatime.NewTimestamp(2020, 9, 1, 12, 13, 14)

	var (ts1, ts2, ts3, ts4, ts5 ekatime.Timestamp)

	err1 := ts1.ParseFrom(bd1)
	err2 := ts2.ParseFrom(bd2)
	err3 := ts3.ParseFrom(bd3)
	err4 := ts4.ParseFrom(bd4)
	err5 := ts5.ParseFrom(bd5)

	require.Equal(t, ts1, orig1)
	require.Equal(t, ts2, orig2)
	require.Equal(t, ts3, orig3)
	require.Equal(t, ts4, orig4)

	require.NoError(t, err1)
	require.NoError(t, err2)
	require.NoError(t, err3)
	require.NoError(t, err4)
	require.Error(t, err5)
}

func BenchmarkTimestamp_ParseFrom(b *testing.B) {
	bd0 := []byte("2020-09-01T12:13:14")
	var ts ekatime.Timestamp
	b.ResetTimer(); b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = ts.ParseFrom(bd0)
	}
}

func TestTimestamp_MarshalJSON(t *testing.T) {
	var ts = struct {
		TS ekatime.Timestamp `json:"ts"`
	}{
		TS: ekatime.NewTimestamp(2020, 9, 12, 13, 14, 15),
	}
	d, err := json.Marshal(&ts)

	require.NoError(t, err)
	require.EqualValues(t, string(d), `{"ts":"2020-09-12T13:14:15"}`)
}
