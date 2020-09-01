// Copyright © 2020. All rights reserved.
// Author: Ilya Stroy.
// Contacts: qioalice@gmail.com, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekatime_test

import (
	"testing"

	"github.com/qioalice/ekago/v2/ekatime"

	"github.com/stretchr/testify/require"
)

// =========================================================================== //
// =========================================================================== //
// =========================================================================== //

func TestDate_Equal(t *testing.T) {
	d1 := ekatime.NewDate(2020, 9, 9)
	d2 := ekatime.UnixFrom(d1, 0).Date()
	require.False(t, d1 == d2)
	require.True(t, d1.ToCmp() == d2.ToCmp())
	require.True(t, d1.Equal(d2))
}

// =========================================================================== //
// =========================================================================== //
// =========================================================================== //

func Benchmark_BeginningOfYear_CachedYear(b *testing.B) {
	currYear := ekatime.Now().Year()
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ekatime.BeginningOfYear(currYear)
	}
}

func Benchmark_BeginningOfYear_NonCachedYear(b *testing.B) {
	notCachedYear := ekatime.Now().Year() + 20
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ekatime.BeginningOfYear(notCachedYear)
	}
}

func Benchmark_EndOfYear_CachedYear(b *testing.B) {
	currYear := ekatime.Now().Year()
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ekatime.EndOfYear(currYear)
	}
}

func Benchmark_EndOfYear_NonCachedYear(b *testing.B) {
	notCachedYear := ekatime.Now().Year() + 20
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ekatime.EndOfYear(notCachedYear)
	}
}

func Benchmark_BeginningAndEndOfYear_CachedYear(b *testing.B) {
	currYear := ekatime.Now().Year()
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ekatime.BeginningAndEndOfYear(currYear)
	}
}

func Benchmark_BeginningAndEndOfYear_NonCachedYear(b *testing.B) {
	notCachedYear := ekatime.Now().Year() + 20
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
	currYear := ekatime.Now().Year()
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ekatime.BeginningOfMonth(currYear, 12)
	}
}

func Benchmark_BeginningOfMonth_NonCachedYear(b *testing.B) {
	notCachedYear := ekatime.Now().Year() + 20
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ekatime.BeginningOfMonth(notCachedYear, 12)
	}
}

func Benchmark_EndOfMonth_CachedYear(b *testing.B) {
	currYear := ekatime.Now().Year()
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ekatime.EndOfMonth(currYear, 12)
	}
}

func Benchmark_EndOfMonth_NonCachedYear(b *testing.B) {
	notCachedYear := ekatime.Now().Year() + 20
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ekatime.EndOfMonth(notCachedYear, 12)
	}
}

func Benchmark_BeginningAndEndOfMonth_CachedYear(b *testing.B) {
	currYear := ekatime.Now().Year()
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ekatime.BeginningAndEndOfMonth(currYear, 12)
	}
}

func Benchmark_BeginningAndEndOfMonth_NonCachedYear(b *testing.B) {
	notCachedYear := ekatime.Now().Year() + 20
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ekatime.BeginningAndEndOfMonth(notCachedYear, 12)
	}
}

// =========================================================================== //
// =========================================================================== //
// =========================================================================== //