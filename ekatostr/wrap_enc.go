// Copyright © 2023. All rights reserved.
// Author: Ilya Stroy.
// Contacts: inaneverb@pm.me, https://github.com/inaneverb
// License: https://opensource.org/licenses/MIT

package ekatostr

import (
	"bytes"
	"unsafe"

	"github.com/inaneverb/ekacore/ekaunsafe/v4"
)

// wrapEncWithPrefix generates and returns an _Encoder, that writes
// given 'prefix' to output and then gives control to the output to the 'enc'.
// It checks, whether 'enc' wrote something to the output buffer,
// and if it's not, then it clears the prefix from the output.
func wrapEncWithPrefix(prefix string, enc _Encoder) _Encoder {

	// Kinda hack. We're writing prefix anyway, but if 'enc' didn't
	// write something, then just shrink buffer at the size of prefix
	// from the right.

	return func(to []byte, v unsafe.Pointer, bh uint8) []byte {
		to = append(to, prefix...)

		var toBak = to
		if to = enc(to, v, bh); len(to) == len(toBak) {
			to = to[:len(toBak)-len(prefix)]
		}

		return to
	}
}

// wrapEncBytes generates and returns an _Encoder that shall serve []byte
// through bytes.Buffer to be able to consume using given 'readerEnc' _Encoder,
// that shall work with io.Reader or *bytes.Buffer directly.
func wrapEncBytes(readerEnc _Encoder) _Encoder {
	return func(to []byte, v unsafe.Pointer, bh uint8) []byte {
		var r = bytes.NewReader(*(*[]byte)(v))
		return readerEnc(to, ekaunsafe.UnpackInterface(r).Word, bh)
	}
}
