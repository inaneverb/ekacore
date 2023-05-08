// Copyright © 2020. All rights reserved.
// Author: Ilya Stroy.
// Contacts: iyuryevich@pm.me, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekastr_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/inaneverb/ekacore/ekastr/v4"
)

type testCaseInterpolatePartReport struct {
	part   string
	isVerb bool
}

var (
	itc = []struct {
		str      string
		expected []testCaseInterpolatePartReport
	}{
		{
			str: "Hello, {{name}}! Nice to meet you.",
			expected: []testCaseInterpolatePartReport{
				{"Hello, ", false},
				{"{{name}}", true},
				{"! Nice to meet you.", false},
			},
		},
		{
			str: "А вот и {{utf8}} текст с {{даже}} utf8 {{глаголом}}",
			expected: []testCaseInterpolatePartReport{
				{"А вот и ", false},
				{"{{utf8}}", true},
				{" текст с ", false},
				{"{{даже}}", true},
				{" utf8 ", false},
				{"{{глаголом}}", true},
			},
		},
	}
)

func TestInterpolateBytes(t *testing.T) {

	var gotParts = make([]testCaseInterpolatePartReport, 0, 32)
	var gotPartsPtr = &gotParts

	var gotVerb = func(verb []byte) {
		*gotPartsPtr = append(*gotPartsPtr, testCaseInterpolatePartReport{
			part:   string(verb),
			isVerb: true,
		})
	}

	var gotJustText = func(text []byte) {
		*gotPartsPtr = append(*gotPartsPtr, testCaseInterpolatePartReport{
			part: string(text),
		})
	}

	for i, testCase := range itc {
		gotParts = gotParts[:0]
		ekastr.InterpolateBytes([]byte(testCase.str), gotVerb, gotJustText)
		assert.Equal(t, testCase.expected, gotParts, "%i test case", i)
	}
}

func BenchmarkInterpolate(b *testing.B) {
	const S = "This is some {{kind}} of string that must be {{interpolated}}."
	b.ReportAllocs()
	var cb = func(_ string) {}
	for i := 0; i < b.N; i++ {
		ekastr.Interpolate(S, cb, cb)
	}
}

func BenchmarkInterpolateBytes(b *testing.B) {
	const S = "This is some {{kind}} of string that must be {{interpolated}}."
	b.ReportAllocs()
	var cb = func(_ []byte) {}
	for i := 0; i < b.N; i++ {
		arr := ekastr.ToBytes(S)
		ekastr.InterpolateBytes(arr, cb, cb)
	}
}

func BenchmarkPrintf(b *testing.B) {
	const S = "This is some %s of string that must be %s."
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = fmt.Sprintf(S, "kind", "interpolated")
	}
}
