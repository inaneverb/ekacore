// Copyright © 2020. All rights reserved.
// Author: Ilya Stroy.
// Contacts: qioalice@gmail.com, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekalog

import (
	"io"
	"unsafe"

	"github.com/qioalice/ekago/v2/ekadanger"
	"github.com/qioalice/ekago/v2/ekatyp"
)

//noinspection GoSnakeCaseUsage
type (
	// CommonIntegrator is the implementation of Integrator interface.
	// It's SYNC Integrator, that calls 'io.Writer.Write' method one by one at the
	// same goroutine.
	//
	// This integrator covers 99% your cases. Why? Listen:
	// 1. You can attach as many 'io.Writer's as you want.
	// 2. You can define different encoders to different 'io.Writer's.
	// 3. You can specify different minimum enabled levels for different 'io.Writer's.
	// 4. You can specify different minimum levels for stacktrace for different 'io.Writer's.
	//
	// Yes. You can do something like this:
	// - Handle all entries, encode them to JSON, and write to 'os.Stdout' and 'file1';
	// - Handle log entries only with 'Warning' and more dangerous log levels,
	//   encode them to YAML, and write to 'file2',
	// - Handle log only 'Error' log entries and, encode them just as is and write them
	//   to the 'os.Stderr'. ...
	// ... and all of these things you can done using only one CommonIntegrator object!
	//
	// How? Look:
	// 		ig := new(CommonIntegrator)
	// 		ig = ig.WithEncoder(encoder1).      // there is 1st dest registration begin
	// 		        WithMinLevel(LEVEL_DEBUG).
	// 		        WriteTo(os.Stdout, file1).  // there is 1st dest registration end
	// 		        WithEncoder(encoder2).      // there is 2nd dest registration begin
	// 		        WithMinLevel(LEVEL_WARNING).
	// 		        WriteTo(file2).             // there is 2nd dest registration end
	// 		        WithEncoder(encoder3).      // there is 3rd dest registration begin
	// 		        WithMinLevel(LEVEL_ERROR).
	//              WriteTo(os.Stderr)          // there is 3rd dest registration end
	// And there is!
	// Now you have 'ig' object and it is your integrator.
	//
	// Guess you noticed that you must pass something to WithEncoder() method.
	// You must create an encoder, specify its behaviour and pass them into.
	//
	// You may use one of default encoder, ekalog provides:
	// - JSON encoder: CI_JSONEncoder class.
	// - Plain text encoder (w/ TTY coloring supporting): CI_ConsoleEncoder class.
	//
	// You may instantiate them, specify format of encoding, call then
	// FreezeAndGetEncoder() to get built encoder and there is!
	// You can register it now using WithEncoder().
	CommonIntegrator struct {

		// How it works?
		// When you call WithMinLevel, WithMinLevelForStackTrace, WithEncoder
		//   it saves min levels or encoder to output[i], where i is idx.
		// When you call WriteTo it finishes output[i] (where i is idx),
		//   and thus output[i] considered completed, the i counter increases
		//   and then go to new _CI_Output.

		output []_CI_Output // registered and approved destinations.
		oll    Level        // the lowest level among all output's levels.
		stll   Level        // the lowest level of stacktrace generating among all output's levels.
		idx    int          // idx of current object in 'output' being registered.
	}

	// _CI_Output is a CommonIntegrator part that contains encoder
	// and destination 'io.Writer's, Logger's Entry will be written to.
	_CI_Output struct {
		ml   Level       // minimum level log entry should have to be processed
		stml Level       // minimum level starting with stacktrace must be added to the entry
		enc  CI_Encoder  // func that encoders 'Entry' object to '[]byte'
		dest []io.Writer // slice of 'io.Writer's, log entry will be written to
	}

	// CI_Encoder is the Common Integrator's Encoder and it's an alias
	// to the function that takes an one log message (as Entry type),
	// encodes it (by some rules) and returns the encoded data as RAW bytes.
	CI_Encoder func(e *Entry) []byte

	// An interface for default encoders.
	_CI_EncoderGenerator interface {

		// FreezeAndGetEncoder must perform all preparations, allocations,
		// something else that may allow to increase runtime encoding.
		//
		// IT MUST BE NOT POSSIBLE TO CHANGE ENCODER BEHAVIOUR AFTER
		// THIS METHOD IS CALLED AND SOME BUILT ENCODER HAS BEEN OBTAINED!
		FreezeAndGetEncoder() CI_Encoder
	}
)

// MinLevelEnabled returns minimum level an Integrator will handle Logger's Entries with.
// E.g. if minimum level is LEVEL_WARNING then LEVEL_DEBUG, LEVEL_INFO logs will be dropped.
func (bi *CommonIntegrator) MinLevelEnabled() Level {
	return bi.oll
}

// MinLevelForStackTrace returns a minimum level starting with a Logger's Entry
// must generate and attach a stacktrace.
//
// This method is used by internal Logger's part and this level may be set by
// WithMinLevelForStackTrace().
//
// WARNING!
// DOES NOT IMPACT FOR EKAERR ERROR OBJECTS THAT ATTACHED TO THE LOG ENTRY.
// THEY WILL HAVE STACKTRACE ANYWAY.
func (bi *CommonIntegrator) MinLevelForStackTrace() Level {
	return bi.stll
}

// Write writes log entry to all registered destinations.
func (bi *CommonIntegrator) Write(entry *Entry) {

	// it guarantees that bi.output is not empty,
	// because each CommonIntegrator object is checked by tryToBuild().

	for _, output := range bi.output {

		// maybe we must remove stacktrace?
		logStacktraceBak := entry.LogLetter.StackTrace
		if output.stml > entry.Level {
			entry.LogLetter.StackTrace = nil
		}

		encodedEntry := output.enc(entry)

		// restore stacktrace
		entry.LogLetter.StackTrace = logStacktraceBak

		for _, destination := range output.dest {
			_, _ = destination.Write(encodedEntry)
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
		for _, destination := range output.dest {
			if syncer, ok := destination.(ekatyp.Syncer); ok {
				if err := syncer.Sync(); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

// IsAsync always returns false, cause CommonIntegrator is a SYNCHRONOUS integrator.
func (bi *CommonIntegrator) IsAsync() bool {
	return false
}

// WithEncoder marks that all next registered writers by WriteTo() method
// will be associated with 'enc' encoder.
func (bi *CommonIntegrator) WithEncoder(enc CI_Encoder) *CommonIntegrator {

	if bi == nil {
		return bi
	}

	// encAddr == nil if enc == nil
	switch encAddr := ekadanger.TakeRealAddr(enc); {

	case encAddr == nil && len(bi.output) == 0:
		bi.output = append(bi.output, _CI_Output{
			enc: defaultConsoleEncoder,
		})
		// bi.idx == 0 already (because len(bi.output) == 0)

	case encAddr == nil:
		// do nothing.
		// if there was a prev enc, next writers will be output to.
		// otherwise defaultConsoleEncoder will be used.

	case len(bi.output) == 0:
		bi.output = append(bi.output, _CI_Output{
			enc: enc,
		})

	default:
		// maybe writers of prev encoder are empty?
		if len(bi.output[bi.idx].dest) == 0 {
			bi.output[bi.idx].enc = enc
			return bi
		}
		bi.output = append(bi.output, _CI_Output{
			enc: enc,
		})
		bi.idx++
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
		// only in that case bi.idx == 0,
		// it was a direct call WithMinLevel(), even w/o WithEncoder() before.
		bi.WithEncoder(nil) // then here will no SEGFAULT
	}

	bi.output[bi.idx].ml = minLevel
	return bi
}

// WithMinLevelForStackTrace changes a minimum level log's Entry stacktrace being
// generated for and saves it for next registered writers by WriteTo() method.
//
// If it's less than level, registered by WithMinLevel(),
// it will be overwritten.
//
// WARNING!
// DOES NOT IMPACT FOR EKAERR ERROR OBJECTS THAT ATTACHED TO THE LOG ENTRY.
// THEY WILL HAVE STACKTRACE ANYWAY.
func (bi *CommonIntegrator) WithMinLevelForStackTrace(minLevel Level) *CommonIntegrator {

	if bi == nil {
		return nil
	}

	if len(bi.output) == 0 {
		// only in that case bi.idx == 0,
		// it was a direct call WithMinLevel(), even w/o WithEncoder() before.
		bi.WithEncoder(nil) // then here will no SEGFAULT
	}

	bi.output[bi.idx].stml = minLevel
	return bi
}

// WriteTo registers all passed 'io.Writer's as logger's destinations.
// All previous WithEncoder(), WithMinLevel() calls will be applied to these writers.
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
		// only in that case bi.idx == 0,
		// it was a direct call WriteTo(), even w/o WithEncoder() before.
		bi.WithEncoder(nil) // otherwise there will be SEGFAULT
	}

	bi.output[bi.idx].dest =
		append(bi.output[bi.idx].dest, writers...)

	return bi
}

// tryToBuild tries to "build" CommonIntegrator objects:
// - drop all barely registered _CI_Output objects,
// - calculate lowest levels of all _CI_Output s.
//
// Returns 'false' only if bi == nil or there is no registered writers.
// Otherwise always 'true' is returned.
func (bi *CommonIntegrator) tryToBuild() (wasBuilt bool) {

	if bi == nil || len(bi.output) == 0 {
		return false
	}

	// only last one bi.output could have empty writers set.
	// fix it then
	if len(bi.output[bi.idx].dest) == 0 {
		if bi.output = bi.output[:bi.idx]; len(bi.output) == 0 {
			// no valid enc/w after cut empty one
			return false
		}
	}

	bi.oll = Level(0xFF)
	bi.stll = Level(0xFF)

	for _, output := range bi.output {
		if output.ml < bi.oll {
			bi.oll = output.ml
		}
		if output.stml < bi.stll {
			bi.stll = output.stml
		}
	}

	return true
}
