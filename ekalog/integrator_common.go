// Copyright © 2020. All rights reserved.
// Author: Ilya Stroy.
// Contacts: qioalice@gmail.com, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekalog

import (
	"io"
	"unsafe"

	"github.com/qioalice/ekago/ekadanger"
	"github.com/qioalice/ekago/ekatyp"
)

// CommonIntegrator is the implementation way of Integrator interface.
// It's SYNC Integrator, that calls 'io.Writer.Write' method one by one at the
// same goroutine.
//
// This integrator covers 99% your cases. Why? Listen:
// 1. You can attach many 'io.Writer's.
// 2. And you can define different encoders to different 'io.Writer's.
// 3. And you can declare different minimum enabled levels for different 'io.Writer's.
//
// Yes. You can do something like this:
// - Handle all entries, encode them to JSON, and write to to 'os.Stdout' and 'file1'
// - Handle log entries only with 'Warning' and more dangerous log levels,
//   encode them to YAML, and write to 'file2',
// - Handle log only 'Error' log entries and, encode them just as is and write them
//   to the 'os.Stderr'. ...
// ... and all of these things you can done using only one CommonIntegrator object!
//
// How? Look:
// 		ig := new(CommonIntegrator)
// 		ig = ig.WithEncoder(encoder1). // there is 1st dest registration begin
// 		        WithMinLevel(Level.Debug).
// 		        WriteTo(os.Stdout). // there is 1st dest registration end
// 		        WithEncoder(encoder2). // there is 2nd dest registration begin
// 		        WithMinLevel(Level.Warn).
// 		        WriteTo(os.Stderr). // there is 2nd dest registration end
// And there is!
// Now you have 'ig' object and it is your integrator.
type CommonIntegrator struct {

	// How it works? When you call WithMinLevel, WithEncoder it saves min level
	// or encoder to output[i], where i is outputRegisterIdx.
	// When you call WriteTo it finishes output[i] (where i is outputRegisterIdx),
	// and thus output[i] considered completed, the i counter increases
	// and then go to new commonIntegratorOutput.

	// registered and approved commonIntegratorOutput objects
	output []commonIntegratorOutput

	// the lowest level among all commonIntegratorOutput levels.
	outputLowestLevel Level

	// idx of current commonIntegratorOutput object in 'output' being registered
	outputRegisterIdx int
}

// CommonIntegratorEncoder is the alias to the function that takes an one log message
// (as Entry type), encodes it (by some rules) and returns the encoded
// data as RAW bytes.
type CommonIntegratorEncoder func(e *Entry) []byte

type commonIntegratorEncoderGenerator interface {
	FreezeAndGetEncoder() CommonIntegratorEncoder
}

// commonIntegratorOutput is a CommonIntegrator part that contains encoder
// and destination 'io.Writer's, log entry will be written to.
type commonIntegratorOutput struct {

	// minimum level log entry should have to be processed
	MinLevel Level

	// func that converts 'Entry' object to '[]byte' RAW data
	Encoder CommonIntegratorEncoder

	// unique addr of 'Encoder', allows to compare encoder functions to avoid
	// saving more than one commonIntegratorOutput objects with the same encoders
	EncoderAddr unsafe.Pointer

	// slice of 'io.Writer's, log entry will be written to
	Destinations []io.Writer
}

// MinLevelEnabled returns minimum log's Level an Integrator will handle
// log entries with.
// E.g. if minimum level is 'Warning', 'Debug' logs will be dropped.
func (bi *CommonIntegrator) MinLevelEnabled() Level {
	return bi.outputLowestLevel
}

// Write writes log entry to all registered destinations.
func (bi *CommonIntegrator) Write(entry *Entry) {

	// it guarantees that bi.output is not empty,
	// because each CommonIntegrator object is checked by tryToBuild().

	for _, output := range bi.output {
		rawData := output.Encoder(entry)
		for _, destination := range output.Destinations {
			_, _ = destination.Write(rawData)
		}
	}
}

// Sync flushes all pending log entries to all registered destinations,
// until error occurred. If that will happen, the process will stop and error
// will be returned.
func (bi *CommonIntegrator) Sync() error {

	// it guarantees that bi.output is not empty,
	// because each CommonIntegrator object is checked by tryToBuild().

	for _, output := range bi.output {
		for _, destination := range output.Destinations {
			if syncer, ok := destination.(ekatyp.Syncer); ok {
				if err := syncer.Sync(); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

//
func (bi *CommonIntegrator) IsAsync() bool {
	return false
}

// WithEncoder marks that all next registered writers by WriteTo() method
// will be associated with 'enc' encoder.
func (bi *CommonIntegrator) WithEncoder(enc CommonIntegratorEncoder) *CommonIntegrator {

	if bi == nil {
		return bi
	}

	encAddr := ekadanger.TakeRealAddr(enc) // encAddr == nil if enc == nil

	switch {
	case encAddr == nil && len(bi.output) == 0:
		bi.output = append(bi.output, commonIntegratorOutput{
			EncoderAddr: jsonEncoderAddr,
			Encoder:     jsonEncoder,
		})
		// bi.outputRegisterIdx == 0 already (because len(bi.output) == 0)

	case encAddr == nil:
		// do nothing.
		// if there was a prev enc, next writers will be output to.
		// otherwise consoleEncoder will be used.

	case len(bi.output) == 0:
		bi.output = append(bi.output, commonIntegratorOutput{
			EncoderAddr: encAddr,
			Encoder:     enc,
		})

	default:
		// maybe bi.output already has a commonIntegratorOutput object with
		// the same 'enc' as encAddr?
		for i, alreadyRegistered := range bi.output {
			if alreadyRegistered.EncoderAddr == encAddr {
				// current's w can not be empty only if it not last
				// but if it so, will be initialized by WriteTo() or reused
				// by next call of WithEncoder().
				bi.outputRegisterIdx = i
				return bi
			}
		}
		// maybe writers of prev encoder are empty?
		if len(bi.output[bi.outputRegisterIdx].Destinations) == 0 {
			bi.output[bi.outputRegisterIdx].EncoderAddr = encAddr
			bi.output[bi.outputRegisterIdx].Encoder = enc
			return bi
		}
		bi.output = append(bi.output, commonIntegratorOutput{
			EncoderAddr: encAddr,
			Encoder:     enc,
		})
		bi.outputRegisterIdx++
	}

	return bi
}

// WithMinLevel changes minimum level log's Entry to be processed for next
// registered writers by WriteTo() method.
func (bi *CommonIntegrator) WithMinLevel(minLevel Level) *CommonIntegrator {

	if bi == nil {
		return nil
	}

	if len(bi.output) == 0 {
		// only in that case bi.outputRegisterIdx == 0,
		// it was a direct call WithMinLevel(), even w/o WithEncoder() before.
		bi.WithEncoder(nil) // then here will no SEGFAULT
	}

	bi.output[bi.outputRegisterIdx].MinLevel = minLevel
	return bi
}

// WriteTo registers all passed 'io.Writer's as logger's destinations.
// All previous WithEncoder(), WithMinLevel() calls will be applied to these
// writers.
func (bi *CommonIntegrator) WriteTo(writers ...io.Writer) *CommonIntegrator {

	if bi == nil {
		return bi
	}

	switch {
	case len(writers) == 0 || len(writers) == 1 && writers[0] == nil:
		return bi

	default:
		// keep only not nil writers
		// https://github.com/golang/go/wiki/SliceTricks#filtering-without-allocating
		notNilWriters := writers[:0]
		for _, writer := range writers {
			if writer != nil {
				notNilWriters = append(notNilWriters, writer)
			}
		}
		if writers = notNilWriters; len(writers) == 0 {
			return bi
		}
	}

	if len(bi.output) == 0 {
		// only in that case bi.outputRegisterIdx == 0,
		// it was a direct call WriteTo(), even w/o WithEncoder() before.
		bi.WithEncoder(nil) // otherwise there will be SEGFAULT
	}

	bi.output[bi.outputRegisterIdx].Destinations =
		append(bi.output[bi.outputRegisterIdx].Destinations, writers...)

	return bi
}

// tryToBuild tries to "build" CommonIntegrator objects:
// - drop all barely registered commonIntegratorOutput objects,
// - calculate lowest output level of all commonIntegratorOutput s.
//
// Returns 'false' only if bi == nil or there is no registered writers.
// Otherwise always 'true' is returned.
func (bi *CommonIntegrator) tryToBuild() (wasBuilt bool) {

	if bi == nil || len(bi.output) == 0 {
		return false
	}

	// only last one bi.output could have empty writers set. fix it then
	if len(bi.output[bi.outputRegisterIdx].Destinations) == 0 {
		if bi.output = bi.output[:bi.outputRegisterIdx]; len(bi.output) == 0 {
			// no valid Encoder/w after cut empty one
			return false
		}
	}

	bi.outputLowestLevel = Level(0xFF)

	for _, output := range bi.output {

		if output.MinLevel < bi.outputLowestLevel {
			bi.outputLowestLevel = output.MinLevel
		}

	}

	return true
}
