// Copyright © 2020. All rights reserved.
// Author: Ilya Stroy.
// Contacts: iyuryevich@pm.me, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekaclike_test

import (
	"runtime"
	"testing"
	"unsafe"

	"github.com/qioalice/ekago/v3/internal/ekaclike"

	"github.com/stretchr/testify/assert"
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
		o := newT(10)

		var (
			typedFoo *typeFoo
			typedBar *typeBar
		)
		{
			addrFoo := o.Foo
			typedFoo = &addrFoo
			addrBar := o.Bar
			typedBar = &addrBar
		}

		runtime.GC()

		assert.Equal(t, 10, (*(*typeFoo)(typedFoo))())
		assert.Equal(t, 20, (*(*typeBar)(typedBar))(20))
		assert.Equal(t, 20, (*(*typeFoo)(typedFoo))())

		o.Bar(10)

		ptrFoo = ekaclike.TakeCallableAddr(o.Foo)
		ptrBar = ekaclike.TakeCallableAddr(o.Bar)
	}

	runtime.GC()

	assert.NotNil(t, ptrFoo)
	assert.NotNil(t, ptrBar)

	assert.Equal(t, 10, (*(*typeFoo)(ptrFoo))())
	assert.Equal(t, 20, (*(*typeBar)(ptrBar))(20))
	assert.Equal(t, 20, (*(*typeFoo)(ptrFoo))())
}

func TestTakeCallableAddr2(t *testing.T) {

	o := newT(10)
	foo := o.Foo
	bar := o.Bar

	runtime.GC()

	assert.Equal(t, 10, foo())
	assert.Equal(t, 20, bar(20))
	assert.Equal(t, 20, foo())
}

type CustomError struct{}

func (_ *CustomError) Error() string { return "<custom error>" }

func TestTakeRealAddrForError(t *testing.T) {

	customNilError := (*CustomError)(nil)
	customNotNilError := new(CustomError)

	var legacyNilError error = customNilError
	var legacyNotNilError error = customNotNilError

	assert.True(t, ekaclike.TakeRealAddr(customNilError) == nil)
	assert.True(t, ekaclike.TakeRealAddr(legacyNilError) == nil)

	assert.True(t, ekaclike.TakeRealAddr(customNotNilError) != nil)
	assert.True(t, ekaclike.TakeRealAddr(legacyNotNilError) != nil)

	// This is why this test exists:
	assert.True(t, legacyNilError != nil)
}

func TestTakeCallableAddr3(t *testing.T) {

	f := func(x int32) int32 {
		return x
	}

	z := (*(*func(float32) int32)(ekaclike.TakeCallableAddr(f)))(-1 * 12345678e-4)
	assert.Equal(t, int32(-996519381), z)
}
