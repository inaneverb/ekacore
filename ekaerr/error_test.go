// Copyright © 2020. All rights reserved.
// Author: Ilya Stroy.
// Contacts: qioalice@gmail.com, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekaerr_test

import (
	"fmt"
	"io/ioutil"
	"runtime"
	"testing"

	"github.com/qioalice/ekago/v2/ekaerr"
	"github.com/qioalice/ekago/v2/ekalog"
)

type T struct{}

func (T) String() string {
	return "stringer"
}

func foo() *ekaerr.Error {
	return foo1().
		AddMessage("foo bad").
		AddFields("foo_arg", T{}).
		Throw()
}

func foo1() *ekaerr.Error {
	return foo2().
		AddMessage("foo1 bad").
		AddFields("foo1_arg", 23).
		Throw()
}

func foo2() *ekaerr.Error {
	return foo3().
		AddMessage("foo2 bad").
		AddFields("foo2_arg?", "").
		Throw()
}

func foo3() *ekaerr.Error {
	return ekaerr.IllegalState.
		New("what??", "arg1?", nil).
		Throw()
}

func TestError(t *testing.T) {
	foo().LogAsWarn()
}

func BenchmarkError(b *testing.B) {
	b.StopTimer()
	b.ReportAllocs()

	devNullJSONIntegrator := new(ekalog.CommonIntegrator).
		WithEncoder(new(ekalog.CI_JSONEncoder).FreezeAndGetEncoder()).
		WithMinLevel(ekalog.LEVEL_DEBUG).
		WriteTo(ioutil.Discard)

	ekalog.ReplaceIntegrator(devNullJSONIntegrator)

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		foo2().LogAsWarn()
	}

	defer func() {
		runtime.GC()
		fmt.Printf("EKAERR: %+v\n", ekaerr.EPS())
		fmt.Printf("EKALOG: %+v\n", ekalog.EPS())
		fmt.Println()
	}()
}
