// Copyright © 2021. All rights reserved.
// Author: Ilya Stroy.
// Contacts: iyuryevich@pm.me, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekasys

import (
	"testing"
)

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

	var wrapper = func(f tF) tF {
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
