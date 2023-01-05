// Copyright © 2020-2023. All rights reserved.
// Author: Ilya Stroy.
// Contacts: iyuryevich@pm.me, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekastr_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/qioalice/ekago/ekastr/v4"
)

func TestS2B(t *testing.T) {
	var s = "string"
	require.EqualValues(t, []byte("string"), ekastr.S2B(s))
}

func TestB2S(t *testing.T) {
	var b = []byte("byte array")
	require.EqualValues(t, "byte array", ekastr.B2S(b))
}

func BenchmarkB2S(b *testing.B) {
	b.ReportAllocs()

	var ba = []byte("byte array")
	for i := 0; i < b.N; i++ {
		_ = ekastr.B2S(ba)
	}
}

func BenchmarkS2B(b *testing.B) {
	b.ReportAllocs()

	var s = "string"
	for i := 0; i < b.N; i++ {
		_ = ekastr.S2B(s)
	}
}
