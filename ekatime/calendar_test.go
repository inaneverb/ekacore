// Copyright © 2020. All rights reserved.
// Author: Ilya Stroy.
// Contacts: qioalice@gmail.com, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekatime_test

import (
	"fmt"
	"testing"

	"github.com/qioalice/ekago/v2/ekatime"

	"github.com/json-iterator/go"
)

func TestCalendar_Today(t *testing.T) {

	c := new(ekatime.Calendar).
		RegJsonEncoder(func(today *ekatime.Today) []byte {
			data, _ := jsoniter.Marshal(today)
			return data
		})

	c.EventAdd(ekatime.NewEvent(ekatime.NewDate(2020, 9, 1), 1, true))
	c.Run()
	fmt.Println(string(c.Today().AsJson))
}
