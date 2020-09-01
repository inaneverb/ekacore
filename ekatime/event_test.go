// Copyright © 2020. All rights reserved.
// Author: Ilya Stroy.
// Contacts: qioalice@gmail.com, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekatime_test

import (
	"fmt"
	"testing"

	"github.com/qioalice/ekago/v2/ekatime"
)

func TestEvent_String(t *testing.T) {
	fmt.Println(ekatime.NewEvent(ekatime.NewDate(2020, 12, 31), 1, true).String())
}
