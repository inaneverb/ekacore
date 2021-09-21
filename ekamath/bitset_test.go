// Copyright © 2021. All rights reserved.
// Author: Ilya Stroy.
// Contacts: qioalice@gmail.com, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekamath_test

import (
	"runtime"
	"testing"

	"github.com/qioalice/ekago/v3/ekamath"

	"github.com/stretchr/testify/require"
)

func TestBitSet(t *testing.T) {

	bs := ekamath.NewBitSet(32)
	t.Cleanup(func() {
		bs.DebugFullDump()
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

		{"C1", SET_1 - 1, (*ekamath.BitSet).NextDown, true, SET_1 + 1},
		{"C2", SET_1, (*ekamath.BitSet).NextDown, true, SET_1 + 1},
		{"C3", SET_2 - 1, (*ekamath.BitSet).NextDown, true, SET_2 + 1},
		{"C4", SET_2, (*ekamath.BitSet).NextDown, true, SET_2 + 1},
		{"C5", SET_3 - 1, (*ekamath.BitSet).NextDown, true, SET_3 + 1},
		{"C6", SET_3, (*ekamath.BitSet).NextDown, true, SET_3 + 1},

		{"D1", SET_3 + 1, (*ekamath.BitSet).PrevDown, true, SET_3 - 1},
		{"D2", SET_3, (*ekamath.BitSet).PrevDown, true, SET_3 - 1},
		{"D3", SET_2 + 1, (*ekamath.BitSet).PrevDown, true, SET_2 - 1},
		{"D4", SET_2, (*ekamath.BitSet).PrevDown, true, SET_2 - 1},
		{"D5", SET_2 - 1, (*ekamath.BitSet).PrevDown, true, SET_2 - 2},
		{"D6", SET_1, (*ekamath.BitSet).PrevDown, true, SET_1 - 1},
		{"D7", SET_1 - 1, (*ekamath.BitSet).PrevDown, false, SET_1 - 1},
	} {
		gotVal, gotOk := testCase.fn(bs, testCase.arg)
		require.True(t, gotOk == testCase.ok, "Unexpected OK for [%s] case", testCase.idx)
		require.EqualValues(t, int(testCase.val), int(gotVal), "Unexpected value for [%s] case", testCase.idx)
	}
}

func TestBitSet_Operations(t *testing.T) {

	bs1 := ekamath.NewBitSet(32)
	bs2 := ekamath.NewBitSet(13)

	set1 := []uint{1, 3, 6, 7, 3, 6, 10, 14, 30}
	set2 := []uint{1, 2, 3, 3, 1, 4, 10, 31}

	for _, set1Elem := range set1 {
		bs1.Up(set1Elem)
	}

	for _, set2Elem := range set2 {
		bs2.Up(set2Elem)
	}

	require.EqualValues(t, []uint{1, 3, 6, 7, 10, 14, 30}, bs1.DebugOnesAsSlice(32))
	require.EqualValues(t, []uint{1, 2, 3, 4, 10, 31}, bs2.DebugOnesAsSlice(32))

	bs3 := bs1.Clone()
	bs3.Union(bs2)

	require.EqualValues(t, []uint{1, 2, 3, 4, 6, 7, 10, 14, 30, 31}, bs3.DebugOnesAsSlice(32))

	bs3 = bs1.Clone()
	bs3.Intersection(bs2)

	require.EqualValues(t, []uint{1, 3, 10}, bs3.DebugOnesAsSlice(32))

	bs3 = bs1.Clone()
	bs3.Difference(bs2)

	require.EqualValues(t, []uint{6, 7, 14, 30}, bs3.DebugOnesAsSlice(32))

	bs3 = bs1.Clone()
	bs3.SymmetricDifference(bs2)

	require.EqualValues(t, []uint{2, 4, 6, 7, 14, 30, 31}, bs3.DebugOnesAsSlice(32))

	bs3 = bs1.Clone()
	bs3.Complement()

	zeroes := make([]uint, 0, 64)
	for v, e := bs3.NextDown(0); e; v, e = bs3.NextDown(v) {
		zeroes = append(zeroes, v)
	}

	require.EqualValues(t, []uint{1, 3, 6, 7, 10, 14, 30}, zeroes)

	bs3.ShrinkUpTo(18)

	require.EqualValues(t, []uint{2, 4, 5, 8, 9, 11, 12, 13, 15, 16, 17, 18}, bs3.DebugOnesAsSlice(32))
}

func TestBitSet_CountBetween(t *testing.T) {

	bs1 := ekamath.NewBitSet(32)

	set1 := []uint{1, 2, 4, 5, 7, 10, 11, 17, 18, 19, 23, 25, 28, 29, 30, 31, 32}
	for _, set1Elem := range set1 {
		bs1.Up(set1Elem)
	}

	require.EqualValues(t, set1, bs1.DebugOnesAsSlice(32))

	c := bs1.CountBetween(1, 32)
	require.EqualValues(t, len(set1), int(c))

	c = bs1.CountBetween(3, 20)
	require.EqualValues(t, 8, int(c))

	c = bs1.CountBetween(12, 16)
	require.EqualValues(t, 0, int(c))

	c = bs1.CountBetween(1, 2)
	require.EqualValues(t, 2, int(c))
}

func TestBitSet_CountBetween2(t *testing.T) {

	bs2 := ekamath.NewBitSet(256)

	set2 := []uint{
		/*   1..64  */ 3, 4, 6, 10, 15, 16, 33, 34, 36, 63, 64,
		/*  65..128 */ 65, 67, 128,
		/* 129..192 */ 129, 142, 145, 146,
		/* 193..256 */ 200, 209, 210, 250,
	}
	for _, set2Elem := range set2 {
		bs2.Up(set2Elem)
	}

	require.EqualValues(t, set2, bs2.DebugOnesAsSlice(256))

	c := bs2.CountBetween(1, 256)
	require.EqualValues(t, len(set2), c)

	c = bs2.CountBetween(10, 66)
	require.EqualValues(t, 9, c)

	c = bs2.CountBetween(10, 144)
	require.EqualValues(t, 13, c)

	c = bs2.CountBetween(10, 36)
	require.EqualValues(t, 6, c)

	c = bs2.CountBetween(10, 200)
	require.EqualValues(t, 16, c)

	c = bs2.CountBetween(40, 210)
	require.EqualValues(t, 12, c)

	c = bs2.CountBetween(5, 205)
	require.EqualValues(t, 17, c)
}

func TestBitSet_EncodeDecode(t *testing.T) {

	const MAX = uint(256)
	belongs := func(data []uint, elem uint) bool {
		for _, dataElem := range data {
			if elem == dataElem {
				return true
			}
		}
		return false
	}

	bs2 := ekamath.NewBitSet(MAX)

	set2 := []uint{
		/*   1..64  */ 3, 4, 6, 10, 15, 16, 33, 34, 36, 63, 64,
		/*  65..128 */ 65, 67, 128,
		/* 129..192 */ 129, 142, 145, 146,
		/* 193..256 */ 200, 209, 210, 250,
	}
	for _, set2Elem := range set2 {
		bs2.Up(set2Elem)
	}

	encodedBinary, err := bs2.MarshalBinary()
	require.NoError(t, err)

	runtime.GC()

	bs2 = ekamath.NewBitSet(MAX)
	err = bs2.UnmarshalBinary(encodedBinary)
	require.NoError(t, err)

	runtime.GC()

	for i := uint(1); i <= MAX; i++ {
		have := bs2.IsSet(i)
		must := belongs(set2, i)
		require.True(t, have == must, "Have: %t, Must: %t, Elem: %v", have, must, i)
	}

	runtime.GC()

	encodedText, err := bs2.MarshalText()
	require.NoError(t, err)

	runtime.GC()

	bs2 = ekamath.NewBitSet(MAX)
	err = bs2.UnmarshalText(encodedText)
	require.NoError(t, err)

	runtime.GC()

	for i := uint(1); i <= MAX; i++ {
		have := bs2.IsSet(i)
		must := belongs(set2, i)
		require.True(t, have == must, "Have: %t, Must: %t, Elem: %v", have, must, i)
	}
}
