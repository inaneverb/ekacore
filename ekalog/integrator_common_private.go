// Copyright © 2020. All rights reserved.
// Author: Ilya Stroy.
// Contacts: iyuryevich@pm.me, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekalog

import (
	"io"
)

//noinspection GoSnakeCaseUsage
type (
	// _CI_Output is a CommonIntegrator part that contains encoder
	// and destination io.Writer set, encoded Entry will be written to.
	//
	// It used at the CommonIntegrator building procedure.
	_CI_Output struct {
		minLevel           Level       // minimum level log entry should have to be processed
		stacktraceMinLevel Level       // minimum level starting with stacktrace must be added to the entry
		encoder            CI_Encoder  // func that encoders Entry object to []byte
		writers            []io.Writer // slice of io.Writer, log entry will be written to
		preEncodedFields   []byte      // raw data of pre-encoded fields
	}
)

// assertNil panics if current CommonIntegrator is nil.
func (ci *CommonIntegrator) assertNil() {
	if ci == nil {
		panic("CommonIntegrator is nil object")
	}
}

// assertWithLock calls assertNil(), locks then CommonIntegrator and checks,
// whether current CommonIntegrator is not registered before and has at least 1
// valid io.Writer.
// The caller must take a responsibility about unlocking CommonIntegrator.
func (ci *CommonIntegrator) assertWithLock() {
	ci.assertNil()
	ci.mu.Lock()
	switch {
	case ci.isRegistered:
		panic("Failed to build CommonIntegrator. It has already registered.")
	}
}

// build tries to prepare CommonIntegrator to be used with Logger:
//  - Drops last not fully registered _CI_Output object,
//  - Calculates lowest levels of all _CI_Output.
//
// Requirements:
//  - Integrator must not be nil (even typed nil), panic otherwise;
//  - If Integrator is CommonIntegrator
//    it must not be registered with some Logger before, panic otherwise;
//  - If Integrator is CommonIntegrator
//    it must have at least 1 registered io.Writer, panic otherwise.
func (ci *CommonIntegrator) build() {

	// build() cannot be called if CommonIntegrator is nil.
	// So, another one nil check is redundant.

	ci.assertWithLock()
	defer ci.mu.Unlock()

	if len(ci.output) == 0 {
		panic("Failed to build CommonIntegrator. There is no valid io.Writer.")
	}

	// Only last one ci.output could have empty writers set. Fix it then.
	if !ci.isRegistered && len(ci.output[ci.idx].writers) == 0 {
		ci.output = ci.output[:ci.idx]
	}

	ci.oll = LEVEL_WARNING
	ci.stll = LEVEL_WARNING

	for _, output := range ci.output {
		if output.minLevel > ci.oll {
			ci.oll = output.minLevel
		}
		if output.stacktraceMinLevel > ci.stll {
			ci.stll = output.stacktraceMinLevel
		}
	}

	ci.isRegistered = true
}
