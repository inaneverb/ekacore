// Copyright © 2020. All rights reserved.
// Author: Ilya Yuryevich.
// Contacts: iyuryevich@pm.me, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekastr_test

import (
	"strconv"
	"testing"

	"github.com/inaneverb/ekacore/ekastr/v4"
)

func BenchmarkRequiredForI64(b *testing.B) {
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_ = ekastr.RequiredForI64(int64(i))
	}
}

func BenchmarkBItoa64(b *testing.B) {
	b.ReportAllocs()

	var buf = make([]byte, 20)

	for i := 0; i < b.N; i++ {
		_ = ekastr.BItoa64(buf, int64(i))
	}
}

func BenchmarkStdItoa(b *testing.B) {
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_ = strconv.Itoa(i)
	}
}
