// Copyright © 2019. All rights reserved.
// Author: Ilya Stroy.
// Contacts: iyuryevich@pm.me, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekaletter

import (
	"strconv"
)

// KeyOrUnnamed returns LetterField's 'Key' field if it's not empty or string
// "unnamed_<unnamedIdx>". Has generated code for 'unnamedIdx' <= 32.
// Increments 'unnamedIdx' before use it (if it have to be used).
//
// Returns an empty string if 'unnamedIdx' == nil and it must be used.
func (f LetterField) KeyOrUnnamed(unnamedIdx *int16) string {

	if f.Key != "" {
		return f.Key
	}

	if unnamedIdx == nil {
		return ""
	}

	*unnamedIdx++

	if *unnamedIdx < 0 || *unnamedIdx > 32 {
		return "unnamed_" + strconv.Itoa(int(*unnamedIdx))
	}

	switch *unnamedIdx {
	case 0:
		return "unnamed_00"
	case 1:
		return "unnamed_01"
	case 2:
		return "unnamed_02"
	case 3:
		return "unnamed_03"
	case 4:
		return "unnamed_04"
	case 5:
		return "unnamed_05"
	case 6:
		return "unnamed_06"
	case 7:
		return "unnamed_07"
	case 8:
		return "unnamed_08"
	case 9:
		return "unnamed_09"
	case 10:
		return "unnamed_10"
	case 11:
		return "unnamed_11"
	case 12:
		return "unnamed_12"
	case 13:
		return "unnamed_13"
	case 14:
		return "unnamed_14"
	case 15:
		return "unnamed_15"
	case 16:
		return "unnamed_16"
	case 17:
		return "unnamed_17"
	case 18:
		return "unnamed_18"
	case 19:
		return "unnamed_19"
	case 20:
		return "unnamed_20"
	case 21:
		return "unnamed_21"
	case 22:
		return "unnamed_22"
	case 23:
		return "unnamed_23"
	case 24:
		return "unnamed_24"
	case 25:
		return "unnamed_25"
	case 26:
		return "unnamed_26"
	case 27:
		return "unnamed_27"
	case 28:
		return "unnamed_28"
	case 29:
		return "unnamed_29"
	case 30:
		return "unnamed_30"
	case 31:
		return "unnamed_31"
	case 32:
		return "unnamed_32"
	default:
		return ""
	}
}

// RemoveVary removes "sign of variety" for the current LetterField.
// Returns true if it was and removed now. False otherwise.
func (f *LetterField) RemoveVary() bool {
	if f.Key != "" && f.Key[len(f.Key)-1] == '?' {
		f.Key = f.Key[:len(f.Key)-1]
		return true
	}
	return false
}
