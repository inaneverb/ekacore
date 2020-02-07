// Copyright © 2019. All rights reserved.
// Author: Ilya Yuryevich.
// Contacts: qioalice@gmail.com, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package log

import (
	"path/filepath"
	"strconv"

	"github.com/qioalice/gext/sys"
)

// printfHowMuchVerbs reports how much printf verbs 'format' has.
//
// Also fixes 'format' if it's has some incorrect things like e.g.
// no-escaped last percent, etc.
func printfHowMuchVerbs(format *string) (verbsCount int) {

	if format == nil || *format == "" {
		return
	}

	// because Golang uses UTF8, log messages can be written not using just ASCII.
	// yes, yes, I agree, it's some piece of shit, but who knows?
	// because of that this loop is so ugly.

	// Golang promises, that for-range loop of string splits it to runes
	// (and runes could be UTF8 characters).

	prevWasPercent := false
	for _, char := range *format {

		switch {
		case char != '%' && prevWasPercent:
			// prev char was a percent but current's one isn't
			// looks like printf verb
			verbsCount++
			prevWasPercent = false

		case char == '%' && prevWasPercent:
			// prev char was a percent but current's one also too
			// looks like percent escaping
			prevWasPercent = false

		case char == '%' && !prevWasPercent:
			// prev char wasn't a percent but current's one is
			// it could be a printf verb, but we don't know exactly at this moment
			prevWasPercent = true

		case char != '%' && !prevWasPercent:
			// just a common regular character, do nothing
		}
	}

	// fix format string if last char was a percent (and there is EOL)
	if prevWasPercent {
		*format += "%"
	}

	return
}

// implicitUnnamedFieldName returns string "unnamed_xx" where xx is idx.
// There is generated code to avoid 'strconv.Itoa' call for [0..32] 'idx' values.
func implicitUnnamedFieldName(idx int) string {

	if idx < 0 || idx > 32 {
		return "unnamed_" + strconv.Itoa(idx)
	}

	// TODO: Rewrite with binary search? Need more micro optimizations!

	switch idx {
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

// formCaller2 forms and returns a string by stack frame that contains caller's info.
func formCaller2(frame sys.StackFrame) string {

	_, fn := filepath.Split(frame.Function)
	_, file := filepath.Split(frame.File)

	return fn + " (" + file + ":" + strconv.Itoa(frame.Line) + ")"
}
