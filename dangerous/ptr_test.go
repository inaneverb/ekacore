// Copyright © 2019. All rights reserved.
// Author: Alice Qio.
// Contacts: <qioalice@gmail.com>.
// License: https://opensource.org/licenses/MIT

package dangerous

import "unsafe"
import "runtime"
import "testing"

import "github.com/stretchr/testify/assert"

type t1 struct {
	i int
}

func (v t1) Foo() int          { return v.i }
func (v *t1) Bar(newI int) int { v.i = newI; return v.i }

func newT(i int) t1 { return t1{i: i} }

//
func TestTakeCallableAddr(t *testing.T) {

	type (
		typeFoo = func() int
		typeBar = func(int) int
	)

	var ptrFoo, ptrBar unsafe.Pointer
	{
		o := newT(10)
		ptrFoo = TakeCallableAddr(o.Foo)
		ptrBar = TakeCallableAddr(o.Bar)
	}

	runtime.GC()

	assert.NotNil(t, ptrFoo)
	assert.NotNil(t, ptrBar)

	assert.Equal(t, 10, (*(*typeFoo)(ptrFoo))())
	assert.Equal(t, 20, (*(*typeBar)(ptrBar))(20))
	assert.Equal(t, 20, (*(*typeFoo)(ptrFoo))())
}

//
func TestTakeCallableAddr2(t *testing.T) {

	o := newT(10)
	foo := o.Foo
	bar := o.Bar

	runtime.GC()

	assert.Equal(t, 10, foo())
	assert.Equal(t, 20, bar(20))
	assert.Equal(t, 20, foo())
}
