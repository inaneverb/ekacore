// Copyright © 2021. All rights reserved.
// Author: Ilya Stroy.
// Contacts: qioalice@gmail.com, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekaerr

import (
	"fmt"
	"io/ioutil"
	"runtime"
	"testing"

	"github.com/qioalice/ekago/v2/ekalog"
)

type T struct{}

func (T) String() string {
	return "stringer"
}

func foo() *Error {
	return foo1().
		AddMessage("foo bad").
		AddFields("foo_arg", T{}).
		Throw()
}

func foo1() *Error {
	return foo2().
		AddMessage("foo1 bad").
		AddFields("foo1_arg", 23).
		Throw()
}

func foo2() *Error {
	return foo3().
		AddMessage("foo2 bad").
		AddFields("foo2_arg?", "").
		Throw()
}

func foo3() *Error {
	return IllegalState.
		New("what??", "arg1?", nil).
		Throw()
}

func TestError(t *testing.T) {
	foo().LogAsWarn()
}

func BenchmarkErrorAllocate(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = allocError()
	}
}

func BenchmarkErrorAllocateAndRelease(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		releaseError(allocError().(*Error))
	}
}

func BenchmarkErrorAcquire(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = acquireError()
	}
}

func BenchmarkErrorAcquireAndRelease(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		releaseError(acquireError())
	}
}

//goland:noinspection GoSnakeCaseUsage
const (
	BENCHMARK_PRINT_ERROR_POOL_INFO = false
)

func benchError(b *testing.B) {
	b.ReportAllocs()

	devNullJSONIntegrator := new(ekalog.CommonIntegrator).
		WithEncoder(new(ekalog.CI_JSONEncoder).FreezeAndGetEncoder()).
		WithMinLevel(ekalog.LEVEL_DEBUG).
		WriteTo(ioutil.Discard)

	ekalog.ReplaceIntegrator(devNullJSONIntegrator)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		foo().LogAsWarn()
	}
	b.StopTimer()

	if BENCHMARK_PRINT_ERROR_POOL_INFO {
		defer func() {
			runtime.GC()
			fmt.Printf("EKAERR: %+v\n", EPS())
			fmt.Printf("EKALOG: %+v\n", ekalog.EPS())
			fmt.Println()
		}()
	}
}

func BenchmarkErrorGetStackTraceVer1(b *testing.B) {
	benchError(b)
}

func BenchmarkErrorGetStackTraceVer2(b *testing.B) {
	_USE_GET_STACKTRACE_VER2 = true
	benchError(b)
}
