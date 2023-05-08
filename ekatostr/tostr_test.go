// Copyright © 2023. All rights reserved.
// Author: Ilya Stroy.
// Contacts: inaneverb@pm.me, https://github.com/inaneverb
// License: https://opensource.org/licenses/MIT

package ekatostr_test

import (
	"bytes"
	"reflect"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/inaneverb/ekacore/ekatostr/v4"
)

type TestStruct struct {
	Empty  struct{}
	I      int
	S      string
	hidden bool
	V      any
}

func NewTestStruct(i int, s string, b bool, v any) *TestStruct {
	return &TestStruct{I: i, S: s, hidden: b, V: v}
}

func TestToStr(t *testing.T) {

	const GT64 = "hellohellohellohellohellohellohellohellohellohellohellohellohello"
	const GT64HASH = "sha256:79e26b80bd90f047f74914b2cda508d578d66c8d335456eb4875ba7fe2dabc4a"

	for _, tc := range []struct {
		V   any
		Exp string
	}{
		{true, `true`},
		{"str", `"str"`},
		{42, `42`},
		{3.14, `3.14`},
		{[]int{1, 2, 3}, `[<~3> 1,2,3]`},
		{NewTestStruct(42, "hello", true, []int8{1, 2, 3}), `{I=42, S="hello", V=[<~3> 1,2,3]}`},
		{[]byte(GT64), GT64HASH},
		{bytes.NewBufferString(GT64), GT64HASH},
		{[]byte("hello"), `68656c6c6f`},
		{bytes.NewBufferString("hello"), `68656c6c6f`},
	} {
		assert.Equal(t, tc.Exp, ekatostr.ToStr(tc.V), "[%[1]T] %+[1]v", tc.V)
	}
}

func BenchmarkToStr(b *testing.B) {

	//goland:noinspection GoSnakeCaseUsage
	const BH_OFF = ekatostr.TOSTR_BH_LOW_JSON | ekatostr.TOSTR_BH_SKIP_ZERO

	var g = func(v any, bufSize int) (string, func(b *testing.B)) {
		var name = reflect.TypeOf(v).String() + "/" + strconv.Itoa(bufSize)
		return name, func(b *testing.B) {
			var out = make([]byte, 0, bufSize)
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				out, _ = ekatostr.ToStrTo(out, v, ^(uint8(0))&^BH_OFF)
				out = out[:0:bufSize]
			}
		}
	}

	b.Run(g(42, 256))
	b.Run(g("test string", 256))
	b.Run(g(NewTestStruct(42, "hello", true, []int8{1, 2, 3}), 256))
}
