// Copyright © 2020. All rights reserved.
// Author: Ilya Stroy.
// Contacts: qioalice@gmail.com, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekalog_encoder_datadog_test

import (
	"bytes"
	"fmt"
	"os"
	"testing"
	"time"
	"unsafe"

	"github.com/qioalice/ekago/v2/ekadeath"
	"github.com/qioalice/ekago/v2/ekaerr"
	"github.com/qioalice/ekago/v2/ekalog"

	"github.com/qioalice/ekago/v2/ekalog/encoders/datadog"
	"github.com/qioalice/ekago/v2/ekalog/writers/http"
)

type bw bytes.Buffer

func (bw *bw) Write(p []byte) (n int, err error) {
	b := (*bytes.Buffer)(unsafe.Pointer(bw))
	n, err = b.Write(p)
	if err != nil {
		return n, err
	}
	return n+1, b.WriteByte('\n')
}

func foo() *ekaerr.Error {
	return ekaerr.Interrupted.New("some message").W("test", 42).Throw()
}

func log() {
	ekalog.WithThis("sys.dd.service", "for_datadog_service", "sys.dd.ddtags", "env:dev")

	ekaerr.Interrupted.New("test (this should be an interrupted error)").
		LogAsWarn("", "key", "value")

	ekalog.Debug("test with printf and args %s %d", "hello", 25, "arg", 42)

	ekalog.Debug("test debug", "field1", 42, "field2", nil)
	ekalog.Info("test info", "dur", time.Minute * 20 + time.Second * 12, "i64", int64(3234234))
	ekalog.Warn("test warn", "time", time.Now())
	ekalog.Error("test err")

	foo().
		LogAsFatal()
}

func TestExampleLog(_ *testing.T) {

	// WARNING! REMINDER!
	//
	// BY DEFAULT CI_DatadogEncoder DOES NOT ADDS ENDING '\n' TO THE END OF JSON DATA.
	// BUT HERE USED A CUSTOM WRITER (bw TYPE) THAT DOES IT TO PROVIDE TO YOU
	// INCREASED READABILITY.

	b := (*bw)(unsafe.Pointer(bytes.NewBuffer(nil)))

	stdoutIntegrator := new(ekalog.CommonIntegrator).
		WithEncoder(new(ekalog_encoder_datadog.CI_DatadogEncoder).FreezeAndGetEncoder()).
		WithMinLevel(ekalog.LEVEL_DEBUG).
		WriteTo(b)

	ekadeath.Reg(func() {
		str := (*bytes.Buffer)(unsafe.Pointer(b))
		fmt.Println(str.String())
	})

	ekalog.ReplaceIntegrator(stdoutIntegrator)
	log()
}

func TestRealLog(_ *testing.T) {
	//noinspection GoSnakeCaseUsage
	const DATADOG_TOKEN = "<place_token_here>"

	ddWriter := new(ekalog_writer_http.CI_WriterHttp).
		UseProviderDataDog(ekalog_writer_http.DATADOG_ADDR_EU, DATADOG_TOKEN)

	if err := ddWriter.Ping(); err.IsNotNil() {
		err.LogAsFatal()
	}

	stdoutIntegrator := new(ekalog.CommonIntegrator).
		WithEncoder(new(ekalog_encoder_datadog.CI_DatadogEncoder).FreezeAndGetEncoder()).
		WithMinLevel(ekalog.LEVEL_DEBUG).
		WriteTo(ddWriter, os.Stdout)

	ekalog.ReplaceIntegrator(stdoutIntegrator)
	log()
}

