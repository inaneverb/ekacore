// Copyright © 2020-2021. All rights reserved.
// Author: Ilya Stroy.
// Contacts: iyuryevich@pm.me, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekarand

// Ruby original: https://github.com/usmanbashir/haikunator
// Go ver of Ruby original: https://github.com/yelinaung/go-haikunator

import (
	mrand "math/rand"

	"bytes"
	"strconv"

	"github.com/qioalice/ekago/ekastr/v4"
)

// Thanks to https://gist.github.com/hugsy/8910dc78d208e40de42deb29e62df913
// for english adjectives and nouns.

// Haikunate returns a randomized string with the following format:
// `<english_adjective>-<english_noun>-<4_digit>`.
func Haikunate() string {
	return HaikunateWithRange(0, 9999)
}

// HaikunateWithRange does the same as Haikunate() does,
// but the number will be in the ['from'..'to'] range.
func HaikunateWithRange(from, to uint) string {

	if from > to {
		from, to = to, from
	}

	var n = ekastr.RequiredForI64(int64(to)) // bytes required for max number
	var rn = mrand.Uint64() % uint64(to-from)
	var rnn = ekastr.RequiredForI64(int64(rn)) // bytes required for generated number
	var b bytes.Buffer

	b.Grow(n + 32)
	b.WriteString(haikuAdjectives[mrand.Int31n(int32(len(haikuAdjectives)))])
	b.WriteByte('-')
	b.WriteString(haikuNouns[mrand.Int31n(int32(len(haikuNouns)))])
	b.WriteByte('-')

	for i, n := 0, n-rnn; i < n; i++ {
		b.WriteByte('0')
	}

	return ekastr.B2S(strconv.AppendUint(b.Bytes(), rn, 10))
}
