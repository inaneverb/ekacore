// Copyright © 2020. All rights reserved.
// Author: Ilya Stroy.
// Contacts: iyuryevich@pm.me, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekaunsafe_test

import (
	"runtime"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"

	"github.com/inaneverb/ekacore/ekaunsafe/v4"
)

type t1 struct {
	i int
}

func (v *t1) Foo() int         { return v.i }
func (v *t1) Bar(newI int) int { v.i = newI; return v.i }

func newT(i int) t1 { return t1{i: i} }

func TestTakeCallableAddr(t *testing.T) {

	type (
		typeFoo = func() int
		typeBar = func(int) int
	)

	var ptrFoo, ptrBar unsafe.Pointer
	{
		var o = newT(10)

		var (
			typedFoo *typeFoo
			typedBar *typeBar
		)
		{
			var addrFoo = o.Foo
			typedFoo = &addrFoo
			var addrBar = o.Bar
			typedBar = &addrBar
		}

		runtime.GC()

		assert.Equal(t, 10, (*(*typeFoo)(typedFoo))())
		assert.Equal(t, 20, (*(*typeBar)(typedBar))(20))
		assert.Equal(t, 20, (*(*typeFoo)(typedFoo))())

		o.Bar(10)

		ptrFoo = ekaunsafe.TakeCallableAddr(o.Foo)
		ptrBar = ekaunsafe.TakeCallableAddr(o.Bar)
	}

	runtime.GC()

	assert.NotNil(t, ptrFoo)
	assert.NotNil(t, ptrBar)

	assert.Equal(t, 10, (*(*typeFoo)(ptrFoo))())
	assert.Equal(t, 20, (*(*typeBar)(ptrBar))(20))
	assert.Equal(t, 20, (*(*typeFoo)(ptrFoo))())
}

func TestTakeCallableAddr2(t *testing.T) {

	var o = newT(10)
	var foo = o.Foo
	var bar = o.Bar

	runtime.GC()

	assert.Equal(t, 10, foo())
	assert.Equal(t, 20, bar(20))
	assert.Equal(t, 20, foo())
}

type CustomError struct{}

func (_ *CustomError) Error() string { return "<custom error>" }

func TestTakeRealAddrForError(t *testing.T) {

	var customNilError = (*CustomError)(nil)
	var customNotNilError = new(CustomError)

	var legacyNilError error = customNilError
	var legacyNotNilError error = customNotNilError

	assert.True(t, ekaunsafe.TakeRealAddr(customNilError) == nil)
	assert.True(t, ekaunsafe.TakeRealAddr(legacyNilError) == nil)

	assert.True(t, ekaunsafe.TakeRealAddr(customNotNilError) != nil)
	assert.True(t, ekaunsafe.TakeRealAddr(legacyNotNilError) != nil)

	// This is why this test exists:
	assert.True(t, legacyNilError != nil)
}

func TestTakeCallableAddr3(t *testing.T) {

	const U8 = 0xAB

	type tf1 = func(uint16) uint16
	type tf2 = func(uint8, uint8) (uint8, uint8)

	var f1 = tf1(func(u uint16) uint16 { return u })
	var f2 = *(*tf2)(ekaunsafe.TakeCallableAddr(f1))

	var r1, r2 = f2(U8, U8)

	assert.EqualValues(t, U8, r1)
	assert.EqualValues(t, U8, r2)
}

func BenchmarkTakeCallableAddr3(b *testing.B) {
	b.ReportAllocs()

	const U8 = 0xAB

	type tf1 = func(uint16) uint16
	type tf2 = func(uint8, uint8) (uint8, uint8)

	var f1 = tf1(func(u uint16) uint16 { return u })

	for i := 0; i < b.N; i++ {
		var f2 = *(*tf2)(ekaunsafe.TakeCallableAddr(f1))
		_, _ = f2(U8, U8)
	}
}
