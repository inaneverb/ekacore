// Copyright © 2021. All rights reserved.
// Author: Ilya Stroy.
// Contacts: iyuryevich@pm.me, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekatime_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/qioalice/ekago/v4/ekatime"
)

func TestTimestamp_TillNext(t *testing.T) {
	ts := ekatime.NewDate(2012, ekatime.MONTH_JANUARY, 12).WithTime(13, 14, 15)
	for _, n := range []struct {
		expected time.Duration
		range_   ekatime.Timestamp
	}{
		{45*time.Minute + 45*time.Second, 2 * ekatime.SECONDS_IN_HOUR},
		{105*time.Minute + 45*time.Second, 3 * ekatime.SECONDS_IN_HOUR},
		{15*time.Minute + 45*time.Second, 30 * ekatime.SECONDS_IN_MINUTE},
	} {
		require.Equal(t, n.expected, ts.TillNext(n.range_))
	}
}
