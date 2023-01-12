// Copyright © 2020-2021. All rights reserved.
// Author: Ilya Stroy.
// Contacts: iyuryevich@pm.me, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekaerr_test

import (
	"errors"
	"fmt"
	"testing"
	"unsafe"

	"github.com/qioalice/ekago/v4/ekaerr"
	"github.com/qioalice/ekago/v4/ekalog"
	"github.com/qioalice/ekago/v4/internal/ekaletter"
)

//go:noinline
func wrapStdErrors(layer int) error {

	if layer == 0 {
		//goland:noinspection GoErrorStringFormat
		return errors.New("Something goes wrong")
	}

	err := wrapStdErrors(layer - 1)
	if layer&1 == 0 {
		return fmt.Errorf("wrap: %w, add message", err)
	} else {
		return err
	}
}

//go:noinline
func wrapEkagoErr(layer int) *ekaerr.Error {

	if layer == 0 {
		return ekaerr.IllegalArgument.New("Something goes wrong")
	}

	err := wrapEkagoErr(layer - 1)
	if layer&1 == 0 {
		return err.AddMessage("Add message").Throw()
	} else {
		return err.Throw()
	}
}

//go:noinline
func wrapEkagoErrLightweight(layer int) *ekaerr.Error {

	if layer == 0 {
		return ekaerr.IllegalArgument.NewLightweight("Something goes wrong")
	}

	err := wrapEkagoErrLightweight(layer - 1)
	if layer&1 == 0 {
		return err.AddMessage("Add message").WithInt("layer", layer+1).Throw()
	} else {
		return err.Throw()
	}
}

func TestFoo(t *testing.T) {

	err := wrapEkagoErr(16)
	l := ekaletter.BridgeErrorGetLetter(unsafe.Pointer(err))

	fmt.Printf("Messages len = %d (%d), Fields len = %d (%d), Stacktrace len = %d (%d)\n",
		len(l.Messages), cap(l.Messages),
		len(l.Fields), cap(l.Fields),
		len(l.StackTrace), cap(l.StackTrace),
	)

	ekalog.ReplaceEncoder(new(ekalog.CI_JSONEncoder).SetOneDepthLevel(true))
	ekalog.Errore("", err)
}

var cases = []struct {
	layers int
}{
	{1},
	{16},
	{32},
	{64},
	{128},
}

func BenchmarkWrap(b *testing.B) {
	for _, tc := range cases {
		b.Run(fmt.Sprintf("std errors %v layers", tc.layers), func(b *testing.B) {
			b.ReportAllocs()
			for n := 0; n < b.N; n++ {
				_ = wrapStdErrors(tc.layers)
			}
		})

		b.Run(fmt.Sprintf("ekaerr %v layers", tc.layers), func(b *testing.B) {
			b.ReportAllocs()
			for n := 0; n < b.N; n++ {
				err := wrapEkagoErr(tc.layers)
				ekaerr.ReleaseError(err)
			}
		})

		b.Run(fmt.Sprintf("ekaerr LIGHT %v layers", tc.layers), func(b *testing.B) {
			b.ReportAllocs()
			for n := 0; n < b.N; n++ {
				err := wrapEkagoErrLightweight(tc.layers)
				ekaerr.ReleaseError(err)
			}
		})
	}
}
