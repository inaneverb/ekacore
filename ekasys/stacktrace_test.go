// Copyright © 2019. All rights reserved.
// Author: Ilya Yuryevich.
// Contacts: qioalice@gmail.com, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekasys

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// benchGetStackTraceCommonDepth aux bench func that starts
// GetStackTrace bench with specified 'depth' arg.
func benchGetStackTraceCommonDepth(b *testing.B, depth int) {

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = GetStackTrace(0, depth)
	}
}

// benchGetStackTraceSyntheticDepth increases stack depth level artificially
// by 'createDepth' value and then starts GetStackTrace bench with
// specified 'depth' arg.
func benchGetStackTraceSyntheticDepth(b *testing.B, depth, createDepth int) {

	type tF func(int, int) StackTrace

	wrapper := func(f tF) tF {
		return func(i1 int, i2 int) StackTrace {
			return f(i1, i2)
		}
	}

	var f tF = GetStackTrace

	for i := 0; i < createDepth; i++ {
		f = wrapper(f)
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = f(0, depth)
	}
}

// Bench GetStackTrace with 'skip' == 0, 'depth' == 1 on common stack.
func Benchmark_GetStackTrace_CommonDepth_1(b *testing.B) {
	benchGetStackTraceCommonDepth(b, 1)
}

// Bench GetStackTrace with 'skip' == 0, 'depth' == 10 on common stack.
func Benchmark_GetStackTrace_CommonDepth_10(b *testing.B) {
	benchGetStackTraceCommonDepth(b, 10)
}

// Bench GetStackTrace with 'skip' == 0, 'depth' == -1 (full depth)
// on common stack.
func Benchmark_GetStackTrace_CommonDepth_Full(b *testing.B) {
	benchGetStackTraceCommonDepth(b, -1)
}

// Bench GetStackTrace with 'skip' == 0, 'depth' == 1
// on artificially enlarged stack by 10.
func Benchmark_GetStackTrace_SyntheticDepth_1_of_10(b *testing.B) {
	benchGetStackTraceSyntheticDepth(b, 1, 10)
}

// Bench GetStackTrace with 'skip' == 0, 'depth' == 10
// on artificially enlarged stack by 10.
func Benchmark_GetStackTrace_SyntheticDepth_10_of_10(b *testing.B) {
	benchGetStackTraceSyntheticDepth(b, 10, 10)
}

// Bench GetStackTrace with 'skip' == 0, 'depth' == -1 (full depth)
// on artificially enlarged stack by 10.
func Benchmark_GetStackTrace_SyntheticDepth_Full_of_10(b *testing.B) {
	benchGetStackTraceSyntheticDepth(b, -1, 10)
}

// Bench GetStackTrace with 'skip' == 0, 'depth' == 10
// on artificially enlarged stack by 10.
func Benchmark_GetStackTrace_SyntheticDepth_10_of_20(b *testing.B) {
	benchGetStackTraceSyntheticDepth(b, 10, 20)
}

// Bench GetStackTrace with 'skip' == 0, 'depth' == -1 (full depth)
// on artificially enlarged stack by 10.
func Benchmark_GetStackTrace_SyntheticDepth_Full_of_20(b *testing.B) {
	benchGetStackTraceSyntheticDepth(b, -1, 20)
}

// ---------------------------------------------------------------------------- //

// benchGetStackFramePointsCommonDepth aux bench func that starts
// getStackFramePoints bench with specified 'depth' arg.
func benchGetStackFramePointsCommonDepth(b *testing.B, depth int) {

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = getStackFramePoints(0, depth)
	}
}

// benchGetStackFramePointsSyntheticDepth increases stack depth level artificially
// by 'createDepth' value and then starts getStackFramePoints bench with
// specified 'depth' arg.
func benchGetStackFramePointsSyntheticDepth(b *testing.B, depth, createDepth int) {

	type tF func(int, int) []uintptr

	wrapper := func(f tF) tF {
		return func(i1 int, i2 int) []uintptr {
			return f(i1, i2)
		}
	}

	var f tF = getStackFramePoints

	for i := 0; i < createDepth; i++ {
		f = wrapper(f)
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = f(0, depth)
	}
}

// Bench getStackFramePoints with 'skip' == 0, 'depth' == 1 on common stack.
func Benchmark_getStackFramePoints_CommonDepth_1(b *testing.B) {
	benchGetStackFramePointsCommonDepth(b, 1)
}

// Bench getStackFramePoints with 'skip' == 0, 'depth' == 10 on common stack.
func Benchmark_getStackFramePoints_CommonDepth_10(b *testing.B) {
	benchGetStackFramePointsCommonDepth(b, 10)
}

// Bench getStackFramePoints with 'skip' == 0, 'depth' == -1 (full depth)
// on common stack.
func Benchmark_getStackFramePoints_CommonDepth_Full(b *testing.B) {
	benchGetStackFramePointsCommonDepth(b, -1)
}

// Bench getStackFramePoints with 'skip' == 0, 'depth' == 1
// on artificially enlarged stack by 10.
func Benchmark_getStackFramePoints_SyntheticDepth_1_of_10(b *testing.B) {
	benchGetStackFramePointsSyntheticDepth(b, 1, 10)
}

// Bench getStackFramePoints with 'skip' == 0, 'depth' == 10
// on artificially enlarged stack by 10.
func Benchmark_getStackFramePoints_SyntheticDepth_10_of_10(b *testing.B) {
	benchGetStackFramePointsSyntheticDepth(b, 10, 10)
}

// Bench getStackFramePoints with 'skip' == 0, 'depth' == -1 (full depth)
// on artificially enlarged stack by 10.
func Benchmark_getStackFramePoints_SyntheticDepth_Full_of_10(b *testing.B) {
	benchGetStackFramePointsSyntheticDepth(b, -1, 10)
}

// Bench getStackFramePoints with 'skip' == 0, 'depth' == 10
// on artificially enlarged stack by 10.
func Benchmark_getStackFramePoints_SyntheticDepth_10_of_20(b *testing.B) {
	benchGetStackFramePointsSyntheticDepth(b, 10, 20)
}

// Bench getStackFramePoints with 'skip' == 0, 'depth' == -1 (full depth)
// on artificially enlarged stack by 10.
func Benchmark_getStackFramePoints_SyntheticDepth_Full_of_20(b *testing.B) {
	benchGetStackFramePointsSyntheticDepth(b, -1, 20)
}

// ---------------------------------------------------------------------------- //

// Test GetStackTrace with 'skip' == 0, 'depth' == 1,
// tests:
// - GetStackTrace returns slice with len == 1 (as depth)
// - frame.Function contains current test name
func TestGetStackTraceCommonDepth1(t *testing.T) {

	frames := GetStackTrace(0, 1)

	assert.Len(t, frames, 1, "invalid len of frames")
	assert.Contains(t, frames[0].Function, "TestGetStackTraceCommonDepth1",
		"wrong function name")
}

// Test GetStackTrace with 'skip' == -3 (include hidden frames),
// 'depth' == -1 (full depth) tests:
// - GetStackTrace returns slice with len >= 3
// (at least hidden frames were included to the output)
// - first three returned frames have valid function names
func TestGetStackTraceCommonDepthAbsolutelyFull(t *testing.T) {

	frames := GetStackTrace(-3, -1)

	assert.True(t, len(frames) >= 3, "invalid len of frames")

	funcNames := []string{
		"runtime.Callers", "getStackFramePoints", "GetStackTrace",
	}

	for i := 0; i < len(funcNames) && i < len(frames); i++ {
		assert.Contains(t, frames[i].Function, funcNames[i],
			"wrong function name")
	}

	frames.Print(nil)
}

// ---------------------------------------------------------------------------- //

type T struct{}

func (T) foo() StackFrame {
	return GetStackTrace(0, 1)[0]
}

// TestStackFrame_DoFormat just see what StackFrame.DoFormat generates.
func TestStackFrame_DoFormat(t *testing.T) {

	frame := GetStackTrace(0, 1)[0]
	fmt.Println(frame.doFormat())

	frame = new(T).foo()
	fmt.Println(frame.doFormat())
}

// Bench StackFrame.doFormat func (generating readable output of stack frame).
func BenchmarkStackFrame_DoFormat(b *testing.B) {

	b.ReportAllocs()
	b.StopTimer()

	frame := GetStackTrace(0, 1)[0]

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		frame.doFormat()
	}
}
