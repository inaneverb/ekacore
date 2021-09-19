// Copyright © 2021. All rights reserved.
// Author: Ilya Stroy.
// Contacts: qioalice@gmail.com, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekatime_test

import (
	"testing"

	"github.com/qioalice/ekago/v3/ekatime"

	"github.com/stretchr/testify/require"
)

func BenchmarkNewWorkCalendar(b *testing.B) {
	b.ReportAllocs()
	b.StopTimer()
	wc := ekatime.NewWorkCalendar(2021, true, false)
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		wc.DoSaturdayAndSundayDayOff()
	}
}

func foo() {
	_ = require.New
}