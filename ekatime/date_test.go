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
