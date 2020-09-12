// Copyright © 2020. All rights reserved.
// Author: Ilya Stroy.
// Contacts: qioalice@gmail.com, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekatime_test

import (
	"fmt"
	"testing"

	"github.com/qioalice/ekago/v2/ekatime"

	"github.com/stretchr/testify/require"
)

func TestTime_String(t *testing.T) {
	t1 := ekatime.NewTime(14, 2, 13)
	t2 := ekatime.NewTime(23, 59, 59)
	t3 := ekatime.NewTime(0, 0, 0)

	require.EqualValues(t, "14:02:13", t1.String())
	require.EqualValues(t, "23:59:59", t2.String())
	require.EqualValues(t, "00:00:00", t3.String())
}

func BenchmarkTime_String_Cached(b *testing.B) {
	t0 := ekatime.NewTime(14, 2, 13)
	b.ResetTimer(); b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = t0.String()
	}
}

func BenchmarkTime_String_FmtSprintf(b *testing.B) {
	hh, mm, ss := ekatime.NewTime(14, 2, 13).Split()
	b.ResetTimer(); b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = fmt.Sprintf("%02d:%02d:%02d", hh, mm, ss)
	}
}

func TestTime_ParseFrom(t *testing.T) {
	bt1 := []byte("000000")   // valid
	bt2 := []byte("23:59:59") // valid
	bt3 := []byte("23:59:5")  // invalid
	bt4 := []byte("23:5959")  // invalid
	bt5 := []byte("2359:5")   // invalid

	var (t1, t2, t3, t4, t5 ekatime.Time)

	err1 := t1.ParseFrom(bt1)
	err2 := t2.ParseFrom(bt2)
	err3 := t3.ParseFrom(bt3)
	err4 := t4.ParseFrom(bt4)
	err5 := t5.ParseFrom(bt5)

	require.Equal(t, ekatime.NewTime(0, 0, 0), t1)
	require.Equal(t, ekatime.NewTime(23, 59, 59), t2)

	require.Nil(t, err1)
	require.Nil(t, err2)
	require.NotNil(t, err3)
	require.NotNil(t, err4)
	require.NotNil(t, err5)
}

func BenchmarkTime_ParseFrom(b *testing.B) {
	bt0 := []byte("23:59:59")
	var tt ekatime.Time
	b.ResetTimer(); b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = tt.ParseFrom(bt0)
	}
}

func TestTime_MarshalJSON(t *testing.T) {

}