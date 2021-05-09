// Copyright © 2020. All rights reserved.
// Author: Ilya Stroy.
// Contacts: qioalice@gmail.com, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekatime_test

import (
	"fmt"
	"testing"

	"github.com/qioalice/ekago/v3/ekaerr"
	"github.com/qioalice/ekago/v3/ekatime"

	"github.com/json-iterator/go"
	"github.com/stretchr/testify/require"
)

func TestCalendar_Today(t *testing.T) {

	c := new(ekatime.Calendar).
		RegJsonEncoder(func(today *ekatime.Today) ([]byte, *ekaerr.Error) {
			data, _ := jsoniter.Marshal(today)
			return data, nil
		})

	c.EventAdd(ekatime.NewEvent(ekatime.NewDate(2020, 9, 1), 1, true))
	c.RunAsync()
	fmt.Println(string(c.Today().AsJson))
}

func TestCalendar_WorkdaysFor(t *testing.T) {

	c := new(ekatime.Calendar).
		EventAdd(ekatime.NewEvent(ekatime.NewDate(2020, 9, 15), 1, true)).
		EventAdd(ekatime.NewEvent(ekatime.NewDate(2020, 9, 16), 1, true)).
		EventAdd(ekatime.NewEvent(ekatime.NewDate(2020, 9, 20), 1, false))

	c.RunAsync()

	current, total :=
		c.WorkdaysFor(ekatime.NewDate(2020, 9, 12), 21)

	require.EqualValues(t, 5, current)
	require.EqualValues(t, 12, total)
}