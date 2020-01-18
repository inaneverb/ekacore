// Copyright © 2019. All rights reserved.
// Author: Ilya Yuryevich.
// Contacts: qioalice@gmail.com, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package log

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
