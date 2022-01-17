// Copyright © 2020. All rights reserved.
// Author: Ilya Stroy.
// Contacts: iyuryevich@pm.me, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekalog

import (
	"strings"
)

// hpm is "has prefix many" just like strings.HasPrefix,
// but you can check many prefixes at the same time.
func hpm(verb string, prefixes []string) bool {
	for _, prefix := range prefixes {
		if strings.HasPrefix(verb, prefix) {
			return true
		}
	}
	return false
}

// bufgr is buffer grow - a function that takes some buffer 'buf',
// and checks whether it has at least 'required' free bytes. Returns 'buf' if it so.
// Otherwise creates a new buffer, with X as a new capacity, where:
//
// 		X ~= 'cap(buf) * 1.5 + required', but can be more,
//
// and copies all data from 'buf' to a new buffer, then returns it.
func bufgr(buf []byte, required int) []byte {

	if cap(buf)-len(buf) >= required {
		return buf
	}
	// https://github.com/golang/go/wiki/SliceTricks#extend
	return append(buf, make([]byte, required+cap(buf)/2)...)[:len(buf)]

	// TODO: Maybe strict and guarantee that new cap === required + cap(buf) * 1.5?
	//  ATM Golang may reserve a few bytes for internal prediction purpose I guess.
}

// bufw writes 'text' to 'buf', growing 'buf' if it's need to write 'text'.
// Returns grown 'buf' (if it has been grown) or originally 'buf'.
// So, it's recommend to use it func like Golang's one 'append'.
func bufw(buf []byte, text string) []byte {
	return append(bufgr(buf, len(text)), text...)
}

// bufw2 writes 'raw' to 'buf', growing 'buf' if it's need to write 'raw'.
// Returns grown 'buf' (if it has been grown) or originally 'buf'.
// So, it's recommend to use it func like Golang's one 'append'.
func bufw2(buf, raw []byte) []byte {
	return append(bufgr(buf, len(raw)), raw...)
}

// bufwc adds 'c' to 'buf', growing 'buf' if it's required to write that byte.
// Returns grown 'buf' (if it's so) or originally 'buf'.
// Use this function as Golang's append().
//
// It's a new function. In the old code a bufw or a bufw2 still may be used
// to write a single byte. Replace them by this function.
func bufwc(buf []byte, c byte) []byte {
	return append(bufgr(buf, 1), c)
}
