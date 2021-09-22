// Copyright © 2021. All rights reserved.
// Author: Ilya Stroy.
// Contacts: iyuryevich@pm.me, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekaerr

import (
	"testing"
)

func BenchmarkErrorAllocate(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = allocError()
	}
}

func BenchmarkErrorAllocateAndRelease(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		releaseError(allocError().(*Error))
	}
}

func BenchmarkErrorAcquire(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = acquireError()
	}
}

func BenchmarkErrorAcquireAndRelease(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		releaseError(acquireError())
	}
}
