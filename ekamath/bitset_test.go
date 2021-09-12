// Copyright © 2021. All rights reserved.
// Author: Ilya Stroy.
// Contacts: qioalice@gmail.com, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekamath_test

import (
	"testing"

	"github.com/qioalice/ekago/v3/ekamath"

	"github.com/stretchr/testify/require"
)

func TestBitSet(t *testing.T) {

	bs := ekamath.NewBitSet(32)
	t.Cleanup(func() {
		bs.DebugDump()
	})

	//goland:noinspection GoSnakeCaseUsage
	const (
		SET_1 = 2
		SET_2 = 10
		SET_3 = 1064
	)

	require.True(t, bs.IsEmpty())

	bs.Set(SET_1, true)
	bs.Set(SET_2, true)
	bs.Set(SET_3, true)

	require.False(t, bs.IsEmpty())
	require.EqualValues(t, 3, bs.Count())

	require.True(t, bs.IsSet(SET_1))
	require.True(t, bs.IsSet(SET_2))
	require.True(t, bs.IsSet(SET_3))

	require.False(t, bs.IsSet(SET_1-1))
	require.False(t, bs.IsSet(SET_2-1))
	require.False(t, bs.IsSet(SET_2+1))
	require.False(t, bs.IsSet(SET_3-1))
	require.False(t, bs.IsSet(SET_3+1))

	for _, testCase := range []struct {
		idx string
		arg uint
		fn  func(*ekamath.BitSet, uint) (uint, bool)
		ok  bool
		val uint
	}{
		{"A1", SET_1, (*ekamath.BitSet).NextUp, true, SET_2},
		{"A2", SET_2, (*ekamath.BitSet).NextUp, true, SET_3},
		{"A3", SET_3, (*ekamath.BitSet).NextUp, false, SET_3},

		{"B1", SET_3, (*ekamath.BitSet).PrevUp, true, SET_2},
		{"B2", SET_2, (*ekamath.BitSet).PrevUp, true, SET_1},
		{"B3", SET_1, (*ekamath.BitSet).PrevUp, false, SET_1},

		{"C1", SET_1-1, (*ekamath.BitSet).NextDown, true, SET_1+1},
		{"C2", SET_1, (*ekamath.BitSet).NextDown, true, SET_1+1},
		{"C3", SET_2-1, (*ekamath.BitSet).NextDown, true, SET_2+1},
		{"C4", SET_2, (*ekamath.BitSet).NextDown, true, SET_2+1},
		{"C5", SET_3-1, (*ekamath.BitSet).NextDown, true, SET_3+1},
		{"C6", SET_3, (*ekamath.BitSet).NextDown, true, SET_3+1},

		{"D1", SET_3+1, (*ekamath.BitSet).PrevDown, true, SET_3-1},
		{"D2", SET_3, (*ekamath.BitSet).PrevDown, true, SET_3-1},
		{"D3", SET_2+1, (*ekamath.BitSet).PrevDown, true, SET_2-1},
		{"D4", SET_2, (*ekamath.BitSet).PrevDown, true, SET_2-1},
		{"D5", SET_2-1, (*ekamath.BitSet).PrevDown, true, SET_2-2},
		{"D6", SET_1, (*ekamath.BitSet).PrevDown, true, SET_1-1},
		{"D7", SET_1-1, (*ekamath.BitSet).PrevDown, false, SET_1-1},
	} {
		gotVal, gotOk := testCase.fn(bs, testCase.arg)
		require.True(t, gotOk == testCase.ok, "Unexpected OK for [%s] case", testCase.idx)
		require.EqualValues(t, int(testCase.val), int(gotVal), "Unexpected value for [%s] case", testCase.idx)
	}
}
