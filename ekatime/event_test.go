// Copyright © 2020. All rights reserved.
// Author: Ilya Stroy.
// Contacts: iyuryevich@pm.me, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekatime_test

import (
	"testing"

	"github.com/qioalice/ekago/v3/ekatime"

	"github.com/stretchr/testify/require"
)

func TestEvent_String(t *testing.T) {
	e := ekatime.NewEvent(ekatime.NewDate(2020, 12, 31), 1, true)
	require.EqualValues(t, "2020/12/31 [Dayoff] ID: 1", e.String())
}
