// Copyright © 2020-2023. All rights reserved.
// Author: Ilya Stroy.
// Contacts: iyuryevich@pm.me, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekastr_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/qioalice/ekago/v4/ekastr"
)

func TestToBytes(t *testing.T) {
	var s = "string"
	require.EqualValues(t, []byte("string"), ekastr.ToBytes(s))
}

func TestFromBytes(t *testing.T) {
	var b = []byte("byte array")
	require.EqualValues(t, "byte array", ekastr.FromBytes(b))
}

func BenchmarkFromBytes(b *testing.B) {
	b.ReportAllocs()

	var ba = []byte("byte array")
	for i := 0; i < b.N; i++ {
		_ = ekastr.FromBytes(ba)
	}
}

func BenchmarkToBytes(b *testing.B) {
	b.ReportAllocs()

	var s = "string"
	for i := 0; i < b.N; i++ {
		_ = ekastr.ToBytes(s)
	}
}
