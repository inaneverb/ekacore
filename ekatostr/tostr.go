// Copyright © 2023. All rights reserved.
// Author: Ilya Stroy.
// Contacts: inaneverb@pm.me, https://github.com/inaneverb
// License: https://opensource.org/licenses/MIT

package ekatostr

import (
	"reflect"

	"github.com/inaneverb/ekacore/ekaunsafe/v4"
)

//goland:noinspection GoSnakeCaseUsage
const (
	// --- Public flags --- //

	BH_LOW_ON    uint8 = 0x01 // Include parsing arrays, structs, maps
	BH_LOW_JSON  uint8 = 0x02 // Use JSON for arrays, structs, maps
	BH_HIGH_ON   uint8 = 0x04 // Include functions, channels, interfaces
	BH_GO_PTR    uint8 = 0x08 // Dereference pointers before
	BH_SKIP_ZERO uint8 = 0x80 // Ignore zero values
)

// ToStrTo generates a string representation of 'v', adjusting behaviour
// by flags, specified in 'bh', and saving result into 'to',
// re-allocating space (if necessary), and returning original (or grown) buf.
//
// Tip:
// You may back up a length of buffer before calling this func,
// and after it compare original length with the length of returned buffer.
// If they're the same, it means that BH_SKIP_ZERO flag was led
// to skipping the zero value.
//
// WARNING!
// Specifying BH_GO_PTR allows to request follow the pointers,
// but it won't lead to stack overflow, since there's the deep limit,
// after which the pointer address will be used as is. You may expect,
// that border be a small pow of 2, like 4 or 8.
func ToStrTo(to []byte, v any, bh uint8) ([]byte, reflect.Kind) {

	if bh&BH_LOW_JSON != 0 {
		bh |= BH_LOW_ON // Fix 'bh' if it's unexpected state.
	}

	var toBak = to // Make a backup to check later what should we return
	var i = ekaunsafe.UnpackInterface(v)
	var rt = ekaunsafe.ReflectTypeOfRType(i.Type)

	to = getEnc(i.Type)(to, i.Word, bh)

	if len(to) == len(toBak) {
		to = nil
	}

	return to, rt.Kind()
}

func ToStr(v any) string {
	var disabled = BH_SKIP_ZERO | BH_LOW_JSON
	var b, _ = ToStrTo(nil, v, (^uint8(0))&^disabled)
	return ekaunsafe.BytesToString(b)
}
