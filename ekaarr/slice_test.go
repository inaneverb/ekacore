package ekaarr_test

import (
	"reflect"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/require"

	"github.com/inaneverb/ekacore/ekaarr/v4"
)

func slicePtr[T comparable](in []T) uintptr {
	return (*reflect.SliceHeader)(unsafe.Pointer(&in)).Data
}

func TestFilter(t *testing.T) {

	var in = []int{1, 2, 3, 4, -1, 0, -2, -5}
	var out = []int{1, 2, 3, 4}
	var got = ekaarr.Filter(in, func(n int) bool { return n > 0 })

	require.Equal(t, out, got)
	require.Equal(t, slicePtr(in), slicePtr(got))

	in = []int{1, 2, 3, 4}
	out = []int{1, 2, 3, 4}
	got = ekaarr.Filter(in, func(n int) bool { return n > 0 })

	require.Equal(t, out, got)
	require.Equal(t, slicePtr(in), slicePtr(got))

	in = []int{-1, -2, -3, 0, -3, -2, -1}
	out = []int{}
	got = ekaarr.Filter(in, func(n int) bool { return n > 0 })

	require.Equal(t, out, got)

	in = []int{1, -1, 2, 3, -1, -2, -3, 4, -5, -6, 5, 6, 7}
	out = []int{1, 2, 3, 4, 5, 6, 7}
	got = ekaarr.Filter(in, func(n int) bool { return n > 0 })

	require.Equal(t, out, got)
}

func BenchmarkFilter(b *testing.B) {

	var bgFilter = func(in []int, cb func(int) bool) func(*testing.B) {
		return func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				_ = ekaarr.Filter(in, cb)
			}
		}
	}

	var bgIter = func(in []int) func(*testing.B) {
		return func(b *testing.B) {
			b.ReportAllocs()
			var out = make([]int, 0, len(in))
			for i, n := 0, len(in); i < n; i++ {
				if in[i] > 0 {
					out = append(out, in[i])
				}
			}
		}
	}

	var f = func(n int) bool { return n > 0 }

	var bg = func(b *testing.B, name string, in []int) {
		b.Run(name, bgFilter(in, f))
		b.Run(name+"-iter", bgIter(in))
	}

	var in1 = []int{1, 2, 3, 4, -1, 0, -2, -5}
	var in2 = []int{1, 2, 3, 4}
	var in3 = []int{-1, -2, -3, 0, -3, -2, -1}
	var in4 = []int{1, -1, 2, 3, -1, -2, -3, 4, -5, -6, 5, 6, 7}

	var mul = func(in []int, n int) []int {
		var out = make([]int, 0, len(in)*n)
		for i := 0; i < n; i++ {
			out = append(out, in...)
		}
		return out
	}

	var sort = func(in []int) []int {
		var pos []int
		var neg []int
		for i, n := 0, len(in); i < n; i++ {
			if in[i] > 0 {
				pos = append(pos, in[i])
			} else {
				neg = append(neg, in[i])
			}
		}
		return append(pos, neg...)
	}

	bg(b, "TrimRight", in1)
	bg(b, "NoFilter", in2)
	bg(b, "FilterAll", in3)
	bg(b, "Common", in4)
	bg(b, "TrimRight-x10", sort(mul(in1, 10)))
	bg(b, "NoFilter-x10", mul(in2, 10))
	bg(b, "FilterAll-x10", mul(in3, 10))
	bg(b, "Common-x10", mul(in4, 10))
	bg(b, "TrimRight-x100", sort(mul(in1, 100)))
	bg(b, "NoFilter-x100", mul(in2, 100))
	bg(b, "FilterAll-x100", mul(in3, 100))
	bg(b, "Common-x100", mul(in4, 100))
}

func TestContains(t *testing.T) {

	for _, tc := range []struct {
		In1, In2 []int
		Any, All bool
	}{
		{[]int{1, 2, 3, 4}, []int{1}, true, true},
		{[]int{1, 2, 3, 4}, []int{2, 5}, true, false},
		{[]int{1, 2, 3, 4}, []int{}, false, false},
		{[]int{1, 2, 3, 4}, []int{1, 4}, true, true},
		{[]int{1, 2, 3, 4}, []int{2, 3}, true, true},
		{[]int{1, 2, 3, 4}, []int{0, 1}, true, false},
		{[]int{1, 2, 3, 4}, []int{4, 5}, true, false},
		{[]int{1, 2, 3, 4}, []int{1, 2, 3, 4}, true, true},
		{[]int{}, []int{1, 2}, false, false},
		{[]int{}, []int{}, false, false},
		{[]int{1}, []int{1}, true, true},
		{[]int{1}, []int{1, 2}, true, false},
		{[]int{1, 2}, []int{1}, true, true},
		{[]int{1, 2}, []int{2}, true, true},
		{[]int{1, 2}, []int{3}, false, false},
	} {
		const F1 = "[ANY] In1: %v, In2: %v\n"
		const F2 = "[ALL] In1: %v, In2: %v\n"

		var r = ekaarr.ContainsAny(tc.In1, tc.In2...)
		require.Equalf(t, tc.Any, r, F1, tc.In1, tc.In2)

		r = ekaarr.ContainsAll(tc.In1, tc.In2...)
		require.Equalf(t, tc.All, r, F2, tc.In1, tc.In2)
	}
}

func TestReduce(t *testing.T) {

	var in = []int{1, 2, 3, 4, 5}
	var sum = 0
	var mul = 1

	for i, n := 0, len(in); i < n; i++ {
		sum += in[i]
		mul *= in[i]
	}

	var sumCb = func(acc int, v int, _ int, _ []int) int { return acc + v }
	var mulCb = func(acc int, v int, _ int, _ []int) int { return acc * v }

	require.Equal(t, sum, ekaarr.Reduce(in, 0, sumCb))
	require.Equal(t, mul, ekaarr.Reduce(in, 1, mulCb))
}

func TestRemove(t *testing.T) {

	for _, tc := range []struct {
		In1, In2, Out []int
	}{
		{[]int{1, 2, 3, 4}, []int{1, 2}, []int{3, 4}},
		{[]int{1, 2, 3, 4}, []int{1, 2, 3, 4}, []int{}},
		{[]int{1, 2, 3, 4}, []int{}, []int{1, 2, 3, 4}},
		{[]int{1, 2, 3, 4}, []int{1}, []int{2, 3, 4}},
	} {
		const F = "In1: %v, In2: %v\n"

		var r = ekaarr.Remove(tc.In1, tc.In2...)
		require.Equalf(t, tc.Out, r, F, tc.In1, tc.In2)
	}
}

func TestUnique(t *testing.T) {

	for _, tc := range []struct {
		In, Out []int
	}{
		{[]int{1, 2, 3, 4, 5}, []int{1, 2, 3, 4, 5}},
		{[]int{1, 2}, []int{1, 2}},
		{[]int{1}, []int{1}},
		{[]int{1, 1}, []int{}},
		{[]int{1, 2, 3, 2}, []int{1, 3}},
		{[]int{1, 2, 2, 1}, []int{}},
		{[]int{1, 2, 3, 1}, []int{2, 3}},
		{[]int{1, 1, 1, 1, 4}, []int{4}},
	} {
		const F = "In1: %v\n"

		var in = append(tc.In[:0:0], tc.In...)
		require.ElementsMatchf(t, tc.Out, ekaarr.Unique(tc.In), F, in)
	}
}

func TestDistinct(t *testing.T) {

	for _, tc := range []struct {
		In, Out []int
	}{
		{[]int{1, 2, 3, 4, 5}, []int{1, 2, 3, 4, 5}},
		{[]int{1, 2}, []int{1, 2}},
		{[]int{1}, []int{1}},
		{[]int{1, 1}, []int{1}},
		{[]int{1, 2, 3, 2}, []int{1, 2, 3}},
		{[]int{1, 2, 2, 1}, []int{1, 2}},
		{[]int{1, 2, 3, 1}, []int{1, 2, 3}},
		{[]int{1, 1, 1, 1, 4}, []int{1, 4}},
	} {
		const F = "In1: %v\n"

		require.ElementsMatchf(t, tc.Out, ekaarr.Distinct(tc.In), F, tc.In)
	}
}

var tdMinMax = []int{
	31, 25, -86, 9, -72, -78, 61, -58, 60, 36, 84, -46, 95, 86, 100, 18,
	63, 52, -13, -89, -71, -27, 91, 81, -9, -59, 80, 14, 99, 94, 12, -60,
	72, 4, -80, -84, -61, -51, 43, -53, 8, 62, -22, -94, -79, 28, 67, -3,
	-68, 73, -87, 65, 66, 0, -8, -69, 32, -97, -2, 90, -38, 96, -32, 59,
	40, -7, -98, -24, 58, -55, -74, 39, 87, 37, 74, -82, 2, -50, -21, -90,
	-100, 57, -77, 30, 42, 71, 27, 33, 3, -10, -20, -92, -25, 23, -43, -75,
	50, 13, 22, -49,
}

func TestMinMax(t *testing.T) {

	t.Run("Min", func(t *testing.T) {
		require.EqualValues(t, -100, ekaarr.Min(tdMinMax))
	})
	t.Run("Max", func(t *testing.T) {
		require.EqualValues(t, 100, ekaarr.Max(tdMinMax))
	})

	t.Run("MinElems", func(t *testing.T) {
		require.EqualValues(t, -100, ekaarr.MinElems(tdMinMax...))
	})
	t.Run("MaxElems", func(t *testing.T) {
		require.EqualValues(t, 100, ekaarr.MaxElems(tdMinMax...))
	})
}

func BenchmarkMinMax(b *testing.B) {

	b.Run("Min", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_ = ekaarr.Min(tdMinMax)
		}
	})

	b.Run("Max", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_ = ekaarr.Max(tdMinMax)
		}
	})

	b.Run("MinElems", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_ = ekaarr.MinElems(tdMinMax...)
		}
	})

	b.Run("MaxElems", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_ = ekaarr.MaxElems(tdMinMax...)
		}
	})
}
