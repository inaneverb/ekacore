// Copyright © 2020-2022. All rights reserved.
// Author: Ilya Stroy.
// Contacts: iyuryevich@pm.me, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekamath_test

import (
	"strconv"
	"testing"

	"github.com/qioalice/ekago/v4/ekamath"
)

func BenchmarkStat(b *testing.B) {

	var statWeight = []int{100, 1_000, 10_000, 100_000}

	var genBench = func(n int, s ekamath.Stat[int]) func(*testing.B) {
		return func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				s.Clear()
				for j := 0; j < n; j++ {
					s.Count(j)
				}
				_ = s.Avg()
			}
		}
	}

	var sc = ekamath.NewStatCumulative[int]()
	var si = ekamath.NewStatIterative[int]()

	for i, n := 0, len(statWeight); i < n; i++ {
		var weightStr = strconv.Itoa(statWeight[i])

		b.Run("Cumulative"+"-"+weightStr, genBench(statWeight[i], sc))
		b.Run("Iterative"+"-"+weightStr, genBench(statWeight[i], si))
	}
}
