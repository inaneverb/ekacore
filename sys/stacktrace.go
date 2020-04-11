// Copyright © 2019. All rights reserved.
// Author: Ilya Yuryevich.
// Contacts: qioalice@gmail.com, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package sys

import (
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

// StackTrace is the slice of StackFrames, nothing more.
// Each stack level described separately.
type StackTrace []StackFrame

// StackFrame represents one stack level (frame/item).
// It general purpose is runtime.Frame type extending.
type StackFrame struct {
	runtime.Frame
	Format               string
	FormatFileOffset     int
	FormatFullPathOffset int
}

// DoFormat represents stack frame in the following string format:
// "<package>/<func> (<short_file>:<file_line>) <full_package_path>".
// Also saves it to the Format field. And does not regenerate it if it's not empty.
func (f *StackFrame) DoFormat() string {

	if f.Format == "" {

		fullPackage, fn := filepath.Split(f.Function)
		_, file := filepath.Split(f.File)

		// we need last package from the fullPackage
		lastPackage := filepath.Base(fullPackage)

		// need remove last package from fullPackage
		if len(lastPackage)+2 <= len(fullPackage) && lastPackage != "." {
			fullPackage = fullPackage[:len(fullPackage)-len(lastPackage)-2]
		}

		f.Format += lastPackage + "/" + fn

		f.FormatFileOffset = len(f.Format)
		f.Format += " (" + file + ":" + strconv.Itoa(f.Line) + ")"

		f.FormatFullPathOffset = len(f.Format)
		f.Format += " " + fullPackage
	}

	return f.Format
}

// getStackFramePoints returns the stack trace point's slice
// that contains 'count' points and starts from 'skip' depth level.
//
// You can pass any value <= 0 as 'count' to get full stack trace points.
func getStackFramePoints(skip, count int) (framePoints []uintptr) {

	// allow to get absolutely full stack trace
	// (include 'getStackFramePoints' and 'runtime.Callers' functions)
	if skip < -2 {
		skip = -2
	}

	// but by default (if 0 has been passed as 'skip', it means to skip
	// these functions ('getStackFramePoints' and 'runtime.Callers'),
	// Thus:
	// skip < -2 => skip = 0 => with these functions
	// skip == 0 => skip = 2 => w/o these functions
	skip += 2

	// get exactly as many stack trace frames as 'count' is
	// only if 'count' is more than zero
	if count > 0 {
		framePoints = make([]uintptr, count)
		return framePoints[:runtime.Callers(skip, framePoints)]
	}

	const (
		// how much frame points requested first time
		baseFullStackFramePointsLen int = 16

		// maximum requested frame points
		maxFullStackFramePointsLen int = 128
	)

	// runtime.Callers only fills slice we provide which
	// so, if slice is full, reallocate mem and try to request frames again
	framePointsLen := 0
	for count = baseFullStackFramePointsLen; ; count <<= 1 {

		framePoints = make([]uintptr, count)
		framePointsLen = runtime.Callers(skip, framePoints)

		if framePointsLen < count || count == maxFullStackFramePointsLen {
			break
		}
	}

	framePoints = framePoints[:framePointsLen]
	return framePoints[:len(framePoints)-1] // ignore Go internal functions
}

// GetStackTrace returns the stack trace as StackFrame object's slice,
// that have specified 'depth' and starts from 'skip' depth level.
// Each StackFrame object represents an one stack trace depth-level.
//
// You can pass any value <= 0 as 'depth' to get full stack trace.
func GetStackTrace(skip, depth int) (stacktrace StackTrace) {

	// see the same code section in 'getStackFramePoints'
	// to more details what happening here with 'skip' arg
	if skip < -3 {
		skip = -3
	}
	skip++

	// prepare to get runtime.Frame objects:
	// - get stack trace frame points,
	// - create runtime.Frame iterator by frame points from prev step
	framePoints := getStackFramePoints(skip, depth)
	framePointsLen := len(framePoints)
	frameIterator := runtime.CallersFrames(framePoints)

	// alloc mem for slice that will have as many 'runtime.Frame' objects
	// as many frame points we got
	stacktrace = make([]StackFrame, framePointsLen)

	// i in func scope (not in loop's) because it will be 'frames' len
	i := 0
	for more := true; more && i < framePointsLen; i++ {
		stacktrace[i].Frame, more = frameIterator.Next()
	}

	// but frameIterator can provide less 'runtime.Frame' objects
	// than we requested -> should fix 'frames' len w/o reallocate
	return stacktrace[:i]
}

// ExcludeInternal returns stacktrace based on current but with excluded all
// Golang internal functions such as runtime.doInit, runtime.main, etc.
func (s StackTrace) ExcludeInternal() StackTrace {

	// because some internal golang functions (such as runtime.gopanic)
	// could be embedded to user's function stacktrace's part,
	// we can't just cut and drop last part of stacktrace when we found
	// a function with a "runtime." prefix from the beginning to end.
	// instead, we starting from the end and generating "ignore list" -
	// a special list of stack frames that won't be included to the result set.

	idx := len(s) - 1
	for idx > 0 && strings.HasPrefix(s[idx].Function, "runtime.") {
		idx--
	}
	return s[:idx+1]
}

// Write writes generated stacktrace to the w or to the stdout if w == nil.
func (s StackTrace) Write(w io.Writer) (n int, err error) {

	if w == nil {
		w = os.Stdout
	}

	for _, frame := range s {

		nn, err_ := w.Write([]byte(frame.DoFormat()))
		if err_ != nil {
			return n, err_
		}
		n += nn

		// write \n
		if _, err_ := w.Write([]byte{'\n'}); err_ != nil {
			return n, err_
		}
		n += 1
	}

	return n, nil
}

// Print prints generated stacktrace to the w or to the stdout if w == nil.
// Ignores all errors. To write with error tracking use Write method.
func (s StackTrace) Print(w io.Writer) {
	_, _ = s.Write(w)
}
