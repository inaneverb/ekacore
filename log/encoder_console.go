// Copyright © 2019. All rights reserved.
// Author: Ilya Yuryevich.
// Contacts: qioalice@gmail.com, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package log

import (
	"unsafe"

	"github.com/qioalice/gext/dangerous"
)

//
type ConsoleEncoder struct {
}

var (
	// Make sure we won't break API.
	_ CommonIntegratorEncoder = (*ConsoleEncoder)(nil).encode

	// Package's console encoder
	consoleEncoder     CommonIntegratorEncoder
	consoleEncoderAddr unsafe.Pointer
)

func init() {
	consoleEncoder = (&ConsoleEncoder{}).FreezeAndGetEncoder()
	consoleEncoderAddr = dangerous.TakeRealAddr(consoleEncoder)
}

//
func (ce *ConsoleEncoder) FreezeAndGetEncoder() CommonIntegratorEncoder {
	return ce.encode
}

//
func (ce *ConsoleEncoder) encode(e *Entry) []byte {
	return nil
}
