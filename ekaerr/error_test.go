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

	"github.com/qioalice/ekago/ekaerr"
	"github.com/qioalice/ekago/ekalog"
)

func foo() *ekaerr.Error {
	return foo1().S("foo bad").W("foo_arg", 131).Throw()
}

func foo1() *ekaerr.Error {
	return foo2().S("foo1 bad").W("foo1_arg", 23).Throw()
}

func foo2() *ekaerr.Error {
	return foo3().S("foo2 bad").W("foo2_arg", 3243).Throw()
}

func foo3() *ekaerr.Error {
	return ekaerr.IllegalState.New("what??", "arg1", nil).Throw()
}

func TestError(t *testing.T) {
	foo2().LogAsWarning()
}

func BenchmarkError(b *testing.B) {
	b.StopTimer()
	b.ReportAllocs()

	devNullJSONIntegrator := new(ekalog.CommonIntegrator).
		WithEncoder(new(ekalog.JSONEncoder).FreezeAndGetEncoder()).
		WithMinLevel(ekalog.LEVEL_DEBUG).
		WriteTo(ioutil.Discard)

	ekalog.ReplaceIntegrator(devNullJSONIntegrator)

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		foo2().LogAsWarning()
	}

	defer func() {
		runtime.GC()
		fmt.Printf("EKAERR: %+v\n", ekaerr.EPS())
		fmt.Printf("EKALOG: %+v\n", ekalog.EPS())
		fmt.Println()
	}()
}
