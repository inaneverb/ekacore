// Copyright © 2021. All rights reserved.
// Author: Ilya Stroy.
// Contacts: qioalice@gmail.com, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekatyp_test

import (
	"fmt"
	"testing"

	"github.com/json-iterator/go"

	"github.com/qioalice/ekago/v2/ekatyp"
)

type (
	T struct {
		ULID ekatyp.ULID `json:"ulid"`
	}
)

func TestULID_MarshalJSON(t *testing.T) {
	x := T{ ULID: ekatyp.ULID_New_OrPanic() }
	encoded, err := jsoniter.Marshal(x)
	fmt.Println(string(encoded), err)
}

func BenchmarkULID_New(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = ekatyp.ULID_New_OrPanic()
	}
}
