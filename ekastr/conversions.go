// Copyright © 2020-2023. All rights reserved.
// Author: Ilya Stroy.
// Contacts: iyuryevich@pm.me, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekastr

import (
	"reflect"
	"unsafe"
)

// B2S converts byte slice to a string without memory allocation.
func B2S(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

// S2B converts string to a byte slice without memory allocation.
//
// WARNING! PANIC CAUTION!
// It will lead to panic if you're going to modify a byte slice
// that is received from the string literal. You should modify only strings
// that are dynamically allocated (in heap), not in RO memory. Example:
//
//		var b = S2B("string")
//		b[0] = 0         // <-- panic here
//	 b = append(b, 0) // but this will work, because of copy
func S2B(s string) (b []byte) {

	var sh = (*reflect.StringHeader)(unsafe.Pointer(&s))
	var bh = (*reflect.SliceHeader)(unsafe.Pointer(&b))

	bh.Data = sh.Data
	bh.Len = len(s) // We need to ensure s is still alive
	bh.Cap = len(s) // We need to ensure s is still alive

	// https://groups.google.com/g/golang-nuts/c/Zsfk-VMd_fU

	return b
}
