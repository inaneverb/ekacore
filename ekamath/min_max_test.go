package ekamath_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/inaneverb/ekacore/ekaext/v4"
	"github.com/inaneverb/ekacore/ekamath/v4"
)

func TestMinMax(t *testing.T) {

	for _, tc := range []struct {
		A, B, Exp int32
		IsMin     bool
	}{
		{1, 1, 1, true},
		{-1, -1, -1, true},
		{1, 1, 1, false},
		{-1, -1, -1, false},
		{1, 2, 1, true},
		{2, 1, 1, true},
		{1, 2, 2, false},
		{2, 1, 2, false},
		{-1, -1, -1, true},
		{-1, -1, -1, false},
		{-1, -2, -2, true},
		{-2, -1, -2, true},
		{-1, -2, -1, false},
		{-2, -1, -1, false},
		{-1, 1, -1, true},
		{1, -1, -1, true},
		{-1, 1, 1, false},
		{1, -1, 1, false},
	} {
		var a, b = tc.A, tc.B
		var got = ekaext.If(tc.IsMin, ekamath.Min(a, b), ekamath.Max(a, b))

		require.EqualValuesf(t, tc.Exp, got, "IsMin: %t, A: %d, B: %d", tc.IsMin, a, b)
	}
}

func TestClamp(t *testing.T) {

	for _, tc := range []struct {
		A, B, V, Exp int32
	}{
		{-2, 2, 1, 1},
		{-2, 2, 3, 2},
		{-2, 2, -3, -2},
		{2, -2, 1, 1},
		{2, -2, 3, 2},
		{2, -2, -3, -2},
		{1, 5, 4, 4},
		{3, 5, 10, 5},
		{5, 3, 1, 3},
	} {
		var got = ekamath.Clamp(tc.V, tc.A, tc.B)
		require.EqualValuesf(t, tc.Exp, got, "A: %d, B: %d, V: %d", tc.A, tc.B, tc.V)
	}
}
