// Copyright © 2019. All rights reserved.
// Author: Ilya Yuryevich.
// Contacts: qioalice@gmail.com, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package syse

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// benchGetStackTraceCommonDepth aux bench func that starts
// getStackTrace bench with specified 'depth' arg.
func benchGetStackTraceCommonDepth(b *testing.B, depth int) {

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		GetStackTrace(0, depth)
	}
}

// benchGetStackTraceSyntheticDepth increases stack depth level artificially
// by 'createDepth' value and then starts getStackTrace bench with
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
		f(0, depth)
	}
}

// Bench getStackTrace with 'skip' == 0, 'depth' == 1 on common stack.
func BenchmarkGetStackTraceCommonDepth1(b *testing.B) {
	benchGetStackTraceCommonDepth(b, 1)
}

// Bench getStackTrace with 'skip' == 0, 'depth' == 10 on common stack.
func BenchmarkGetStackTraceCommonDepth10(b *testing.B) {
	benchGetStackTraceCommonDepth(b, 10)
}

// Bench getStackTrace with 'skip' == 0, 'depth' == -1 (full depth)
// on common stack.
func BenchmarkGetStackTraceCommonDepthFull(b *testing.B) {
	benchGetStackTraceCommonDepth(b, -1)
}

// Bench getStackTrace with 'skip' == 0, 'depth' == 1
// on artificially enlarged stack by 10.
func BenchmarkGetStackTraceSyntheticDepth1(b *testing.B) {
	benchGetStackTraceSyntheticDepth(b, 1, 10)
}

// Bench getStackTrace with 'skip' == 0, 'depth' == 10
// on artificially enlarged stack by 10.
func BenchmarkGetStackTraceSyntheticDepth10(b *testing.B) {
	benchGetStackTraceSyntheticDepth(b, 10, 10)
}

// Bench getStackTrace with 'skip' == 0, 'depth' == -1 (full depth)
// on artificially enlarged stack by 10.
func BenchmarkGetStackTraceSyntheticDepthFull(b *testing.B) {
	benchGetStackTraceSyntheticDepth(b, -1, 10)
}

// Test getStackTrace with 'skip' == 0, 'depth' == 1,
// tests:
// - getStackTrace returns slice with len == 1 (as depth)
// - frame.Function contains current test name
func TestGetStackTraceCommonDepth1(t *testing.T) {

	frames := GetStackTrace(0, 1)

	assert.Len(t, frames, 1, "invalid len of frames")
	assert.Contains(t, frames[0].Function, "TestGetStackTraceCommonDepth1",
		"wrong function name")
}

// Test getStackTrace with 'skip' == -3 (include hidden frames),
// 'depth' == -1 (full depth) tests:
// - getStackTrace returns slice with len >= 3
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
}
