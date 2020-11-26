// Copyright © 2020. All rights reserved.
// Author: Ilya Stroy.
// Contacts: qioalice@gmail.com, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekalog_test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"strings"
	"testing"
	"time"

	"github.com/qioalice/ekago/v2/ekadeath"
	"github.com/qioalice/ekago/v2/ekaerr"
	"github.com/qioalice/ekago/v2/ekalog"
)

func foo() *ekaerr.Error {
	return ekaerr.Interrupted.
		New("fwefwf").
		AddFields("test", 42).
		Throw()
}

func TestLog(t *testing.T) {

	consoleEncoder := new(ekalog.CI_JSONEncoder)
	b := bytes.NewBuffer(nil)

	stdoutConsoleIntegrator := new(ekalog.CommonIntegrator).
		WithEncoder(consoleEncoder.FreezeAndGetEncoder()).
		WithMinLevel(ekalog.LEVEL_DEBUG).
		WriteTo(b)

	ekadeath.Reg(func() {
		str := b.String()
		//str = strings.ReplaceAll(str, "\033", "\\033")
		_ = strings.ReplaceAll
		fmt.Println(str)
	})

	ekalog.ReplaceIntegrator(stdoutConsoleIntegrator)

	ekaerr.Interrupted.New("test").LogAsWarn("", "key", "value")

	ekalog.Debug("test %s %d", "hello", 25, "arg", 42)

	ekalog.Debug("test", "field1", 42, "field2", nil)
	ekalog.Info("test", "dur", time.Minute * 20 + time.Second * 12, "i64", int64(3234234))
	ekalog.Warn("test", "time", time.Now())
	ekalog.Error("test")

	foo().LogAsFatal()
}

func BenchmarkLog(b *testing.B) {
	b.StopTimer()
	b.ReportAllocs()

	devNullIntegrator := new(ekalog.CommonIntegrator).
		WithEncoder(new(ekalog.CI_ConsoleEncoder).FreezeAndGetEncoder()).
		WithMinLevel(ekalog.LEVEL_DEBUG).
		WriteTo(ioutil.Discard)

	ekalog.ReplaceIntegrator(devNullIntegrator)

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		ekalog.Debug("test", "field", 41, "field2", nil)
	}
}
