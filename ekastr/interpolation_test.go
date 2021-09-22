// Copyright © 2020. All rights reserved.
// Author: Ilya Stroy.
// Contacts: iyuryevich@pm.me, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekastr_test

import (
	"fmt"
	"testing"

	"github.com/qioalice/ekago/v3/ekastr"

	"github.com/stretchr/testify/assert"
)

type (
	/*
		_TITC is Interpolation's test cases
	*/
	_TITC struct {
		str      string
		expected []_TITCSP
	}

	/*
		_TITCSP is Interpolation test cases' string parts
	*/
	_TITCSP struct {
		part   string
		isVerb bool
	}
)

var (
	itc = []_TITC{
		{
			str: "Hello, {{name}}! Nice to meet you.",
			expected: []_TITCSP{
				{"Hello, ", false},
				{"{{name}}", true},
				{"! Nice to meet you.", false},
			},
		},
		{
			str: "А вот и {{utf8}} текст с {{даже}} utf8 {{глаголом}}",
			expected: []_TITCSP{
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

func TestInterpolateb(t *testing.T) {

	var (
		gotParts    = make([]_TITCSP, 0, 32)
		gotPartsPtr = &gotParts
	)

	gotVerb := func(verb []byte) {
		*gotPartsPtr = append(*gotPartsPtr, _TITCSP{
			part:   string(verb),
			isVerb: true,
		})
	}

	gotJustText := func(text []byte) {
		*gotPartsPtr = append(*gotPartsPtr, _TITCSP{
			part: string(text),
		})
	}

	for i, testCase := range itc {
		gotParts = gotParts[:0]
		ekastr.Interpolateb([]byte(testCase.str), gotVerb, gotJustText)
		assert.Equal(t, testCase.expected, gotParts, "%i test case", i)
	}
}

func BenchmarkInterpolate(b *testing.B) {
	const S = "This is some {{kind}} of string that must be {{interpolated}}."
	b.ReportAllocs()
	cb := func(_ string) {}
	for i := 0; i < b.N; i++ {
		ekastr.Interpolate(S, cb, cb)
	}
}

func BenchmarkInterpolateb(b *testing.B) {
	const S = "This is some {{kind}} of string that must be {{interpolated}}."
	b.ReportAllocs()
	cb := func(_ []byte) {}
	for i := 0; i < b.N; i++ {
		arr := ekastr.S2B(S)
		ekastr.Interpolateb(arr, cb, cb)
	}
}

func BenchmarkPrintf(b *testing.B) {
	const S = "This is some %s of string that must be %s."
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = fmt.Sprintf(S, "kind", "interpolated")
	}
}
