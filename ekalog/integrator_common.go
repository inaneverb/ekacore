// Copyright © 2018-2021. All rights reserved.
// Author: Ilya Stroy.
// Contacts: qioalice@gmail.com, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekalog

import (
	"io"
	"sync"

	"github.com/qioalice/ekago/v2/ekatyp"

	"github.com/qioalice/ekago/v2/internal/ekaclike"
	"github.com/qioalice/ekago/v2/internal/ekaletter"
)

//noinspection GoSnakeCaseUsage
type (
	// CommonIntegrator is the implementation of Integrator interface.
	// It's SYNC Integrator, that calls io.Writer.Write() method one by one at the
	// same goroutine.
	//
	// This integrator covers 99% your cases. Why? Look:
	//  1. You can attach as many io.Writer as you want.
	//  2. You can define different encoders for different io.Writer.
	//  3. You can specify different minimum enabled levels for different io.Writer.
	//  4. You can specify different minimum levels for stacktrace for different io.Writer.
	//
	// Yes. You can do something like this:
	//  - Handle all entries, encode them to JSON, and write to os.Stdout and file';
	//  - Handle only LEVEL_WARNING and more dangerous log entries,
	//    encode them to YAML, and write to 'file2',
	//  - Handle only LEVEL_ERROR and more dangerous log entries,
	//    encode them just as is and write them to the 'os.Stderr'.
	//
	// And all of these things can be done using only one CommonIntegrator object!
	//
	// How? Look:
	// 		ig := new(CommonIntegrator)
	// 		ig = ig.WithEncoder(encoder1).      // there is 1st writer registration begin
	// 		        WithMinLevel(LEVEL_DEBUG).
	// 		        WriteTo(os.Stdout, file1).  // there is 1st writer registration end
	// 		        WithEncoder(encoder2).      // there is 2nd writer registration begin
	// 		        WithMinLevel(LEVEL_WARNING).
	// 		        WriteTo(file2).             // there is 2nd writer registration end
	// 		        WithEncoder(encoder3).      // there is 3rd writer registration begin
	// 		        WithMinLevel(LEVEL_ERROR).
	//              WriteTo(os.Stderr)          // there is 3rd writer registration end
	// And there is!
	// Now you have 'ig' object and it is your integrator.
	//
	// Guess you noticed that you must pass something to WithEncoder() method.
	// You must create an encoder, specify its behaviour and pass them into.
	//
	// You may use one of default encoder, ekalog provides:
	// - JSON encoder: CI_JSONEncoder class.
	// - Plain text encoder (w/ TTY coloring supporting): CI_ConsoleEncoder class.
	// You can register them using WithEncoder().
	//
	// NOTICE.
	// CommonIntegrator IS A SYNC INTEGRATOR, MEANING THAT ALL FUNCTION CALLS
	// WILL BE DONE WITH THE SAME GOROUTINE, THE CALLER REQUEST LOGGING
	// IS RUNNING IN.
	//
	// WARNING.
	// CommonIntegrator is thread-unsafety at the initialization*
	// and unavailable to be modified after it is registered** with some Logger.
	// If you're trying to modify CommonIntegrator after it was registered
	// it will do nothing.
	// If some CommonIntegrator was once registered it's unavailable
	// to be registered again explicitly. You should create a new CommonIntegrator,
	// initialize it and then use it.
	//
	// ----------
	//
	//  * "Initialization" is calling methods WithEncoder(), WithMinLevel(),
	//    WriteTo().
	//
	//  ** "Registering" is linking CommonIntegrator with some Logger
	//    using Logger.ReplaceIntegrator()
	//    or package's functions WithIntegrator(), ReplaceIntegrator(),
	//
	CommonIntegrator struct {

		// How it works?
		//
		// When you call WithMinLevel, WithMinLevelForStackTrace, WithEncoder
		// it saves min levels or encoder to output[i], where i is idx.
		//
		// When you call WriteTo it finishes output[i] (where i is idx),
		// and thus output[i] considered completed, the i counter increases
		// and then go to new _CI_Output.

		// -------------------------------------------------------------------- //

		// mu is an CommonIntegrator initialization helper.
		// It allows CommonIntegrator to be thread-safety.
		//
		// status doesn't handle a case when:
		// 1. G1 reads status (0), can modify CommonIntegrator
		// 2. G2 changes status (0->1), first use of CommonIntegrator
		// 3. Now it's data race between G1, G2.
		mu sync.Mutex

		// isRegistered is true if CommonIntegrator has been registered
		// with some Logger.
		isRegistered bool

		// -------------------------------------------------------------------- //

		// oll is the lowest Level among all output's minimum Level enabled.
		// uint32 because of sync/atomic operations, casted to Level later.
		oll Level

		// stll is the lowest level among all output's minimum Level
		// for stacktrace being generated for.
		stll Level

		// output contains an outputs that are under registration
		// or already approved.
		output []_CI_Output

		// idx is an index of output to object that is under initialization
		// right now.
		idx int
	}

	// CI_Encoder is an interface that types must implement to be allowed
	// for being register with CommonIntegrator as one of encoders.
	CI_Encoder interface {

		// PreEncodeField must encode passed ekaletter.LetterField,
		// save its encoded raw data in the internal parts and then use it
		// for each Entry passed to EncodeEntry() method.
		PreEncodeField(f ekaletter.LetterField)

		// EncodeEntry must encode passed Entry by its own way and return
		// an encoded raw data of Entry.
		// Error handling is on implementation's shoulders.
		EncodeEntry(e *Entry) []byte
	}
)


// --------------------- IMPLEMENT Integrator INTERFACE ----------------------- //
// ---------------------------------------------------------------------------- //


// MinLevelEnabled returns minimum level an Integrator will handle Entry with.
// E.g. if minimum level is LEVEL_WARNING then LEVEL_DEBUG, LEVEL_INFO logs will be dropped.
//
// This method is used by internal Logger's part and this level may be set by
// WithMinLevel().
//
// MinLevelEnabled() must not be used by caller before CommonIntegrator
// is registered with some logger u
// Otherwise it may return some strange results.
func (ci *CommonIntegrator) MinLevelEnabled() Level {
	ci.assertNil()
	return ci.oll
}

// MinLevelForStackTrace returns a minimum level starting with an Entry
// must generate and attach a stacktrace.
// This method allows Logger not to generate stacktrace if there's no attached
// ekaerr.Error if Entry's Level less than that.
//
// This method is used by internal Logger's part and this level may be set by
// WithMinLevelForStackTrace().
//
// WARNING!
// DOES NOT IMPACT FOR EKAERR ERROR OBJECTS THAT ATTACHED TO THE LOG ENTRY.
// THEY ALSO MAY DON'T HAVE A STACKTRACE BUT THEY HAVE ITS OWN RULES OF THAT.
func (ci *CommonIntegrator) MinLevelForStackTrace() Level {
	ci.assertNil()
	return ci.stll
}

// PreEncodeField passes presented ekaletter.LetterField to all registered
// CI_Encoder objects, saving it as encoded RAW data inside them to attach them later
// to each Entry that must be logged.
//
// PreEncodeField is for internal purposes only and MUST NOT be called directly.
// UB otherwise, may panic.
//
// If PreEncodeField is called directly, the caller MUST call it only AFTER
// registration of CommonIntegration with some Logger done.
func (ci *CommonIntegrator) PreEncodeField(f ekaletter.LetterField) {

	ci.assertNil()

	for i, n := 0, len(ci.output); i < n; i++ {
		ci.output[i].encoder.PreEncodeField(f)
	}
}

// EncodeAndWrite encodes Entry using registered CI_Encoder objects and then writes
// obtained RAW data ([]byte) to correspondent io.Writer objects.
//
// According with Logger's code, Entry won't be here if it's Level < MinLevelEnabled().
//
// EncodeAndWrite is for internal purposes only and MUST NOT be called directly.
// UB otherwise, may panic.
func (ci *CommonIntegrator) EncodeAndWrite(entry *Entry) {

	ci.assertNil()

	// it guarantees that ci.output is not empty,
	// because each CommonIntegrator object is checked by tryToBuild().

	for _, output := range ci.output {

		// maybe we must remove stacktrace?
		logStacktraceBak := entry.LogLetter.StackTrace
		if output.stacktraceMinLevel > entry.Level {
			entry.LogLetter.StackTrace = nil
		}

		encodedEntry := output.encoder.EncodeEntry(entry)

		// restore stacktrace
		entry.LogLetter.StackTrace = logStacktraceBak

		for _, destination := range output.writers {
			_, _ = destination.Write(encodedEntry)
		}
	}
}

// Sync flushes all pending log entries to all registered destinations,
// until error occurred. If that will happen, the process will stop and error
// will be returned.
func (ci *CommonIntegrator) Sync() error {

	ci.assertNil()

	ci.mu.Lock()
	defer ci.mu.Unlock()

	if !ci.isRegistered {
		return nil
	}

	for _, output := range ci.output {
		for _, destination := range output.writers {
			if syncer, ok := destination.(ekatyp.Syncer); ok {
				if err := syncer.Sync(); err != nil {
					return err
				}
			}
		}
	}

	return nil
}


// -------------------- CommonIntegrator BUILDING METHODS --------------------- //
// ---------------------------------------------------------------------------- //


// WithEncoder marks that all next registered writers by WriteTo() method
// will be associated with passed CI_Encoder encoder.
func (ci *CommonIntegrator) WithEncoder(enc CI_Encoder) *CommonIntegrator {

	ci.assertWithLock()
	defer ci.mu.Unlock()

	// encAddr == nil if encoder == nil
	switch encAddr := ekaclike.TakeRealAddr(enc); {

	case encAddr == nil && len(ci.output) == 0:
		ci.output = append(ci.output, _CI_Output{
			encoder: defaultConsoleEncoder,
		})
		// ci.idx == 0 already (because len(ci.output) == 0)

	case encAddr == nil:
		// do nothing.
		// if there was a prev encoder, next writers will be output to.
		// otherwise defaultConsoleEncoder will be used.
	}

	// Now we know that CI_Encoder is not nil and we need to add it somewhere.
	// Encoders might be CI_ConsoleEncoder or CI_JSONEncoder that must be built.

	switch encTyped := enc.(type) {

	case *CI_ConsoleEncoder:
		encTyped.doBuild()

	case *CI_JSONEncoder:
		encTyped.doBuild()
	}

	// Final step.
	// Determine to which _CI_Output a CI_Encoder will be added.

	switch {

	case len(ci.output) == 0:
		ci.output = append(ci.output, _CI_Output{
			encoder: enc,
		})

	case len(ci.output[ci.idx].writers) == 0:
		// maybe writers of prev encoder are empty?
		ci.output[ci.idx].encoder = enc

	default:
		ci.output = append(ci.output, _CI_Output{
			encoder: enc,
		})
		ci.idx++
	}

	return ci
}

// WithMinLevel changes minimum level log's Entry to be processed for next
// registered writers by WriteTo() method.
func (ci *CommonIntegrator) WithMinLevel(minLevel Level) *CommonIntegrator {

	ci.assertWithLock()
	defer ci.mu.Unlock()

	if len(ci.output) == 0 {
		// only in that case ci.idx == 0,
		// it was a direct call WithMinLevel(), even w/o WithEncoder() before.
		ci.WithEncoder(nil) // then here will no SEGFAULT
	}

	ci.output[ci.idx].minLevel = minLevel
	return ci
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
func (ci *CommonIntegrator) WithMinLevelForStackTrace(minLevel Level) *CommonIntegrator {

	ci.assertWithLock()
	defer ci.mu.Unlock()

	if len(ci.output) == 0 {
		// only in that case ci.idx == 0,
		// it was a direct call WithMinLevel(), even w/o WithEncoder() before.
		ci.WithEncoder(nil) // then here will no SEGFAULT
	}

	ci.output[ci.idx].stacktraceMinLevel = minLevel
	return ci
}

// WriteTo registers all passed io.Writer as CommonIntegrator destinations
// for the CI_Encoder that has been specified using last WithEncoder() call
// before this WriteTo() call.
func (ci *CommonIntegrator) WriteTo(writers ...io.Writer) *CommonIntegrator {

	ci.assertWithLock()
	defer ci.mu.Unlock()

	switch {
	case len(writers) == 0 || len(writers) == 1 && writers[0] == nil:
		return ci

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
			return ci
		}
	}

	if len(ci.output) == 0 {
		// only in that case ci.idx == 0,
		// it was a direct call WriteTo(), even w/o WithEncoder() before.
		ci.WithEncoder(nil) // otherwise there will be SEGFAULT
	}

	ci.output[ci.idx].writers = append(ci.output[ci.idx].writers, writers...)
	return ci
}
