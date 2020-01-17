// Copyright © 2019. All rights reserved.
// Author: Ilya Yuryevich.
// Contacts: qioalice@gmail.com, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package sys

import "runtime"

// StackTrace is the slice of StackFrames, nothing more.
// Each stack level described separately.
type StackTrace []StackFrame

// StackFrame represents one stack level (frame/item).
// It general purpose is runtime.Frame type extending.
type StackFrame struct {
	runtime.Frame
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

	return framePoints[:framePointsLen]
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