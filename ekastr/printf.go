// Copyright © 2020. All rights reserved.
// Author: Ilya Stroy.
// Contacts: iyuryevich@pm.me, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekastr

// PrintfVerbsCount reports how much printf verbs 'format' has.
//
// Also fixes 'format' if it has some incorrect things like e.g.
// no-escaped last percent, etc.
func PrintfVerbsCount(format *string) int {

	if format == nil || *format == "" {
		return 0
	}

	// Because Golang uses UTF8, log messages may contain not only ASCII.
	// Yes, yes, I agree, it's some piece of shit, but who knows?
	// Because of that this loop is so ugly.

	// Golang guarantees, that for-range loop of string splits it to runes
	// (and runes could be UTF8 characters).

	var prevWasPercent = false
	var verbsCount = 0

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
			// it could be a printf verb, but we have to make sure
			prevWasPercent = true

		case char != '%' && !prevWasPercent:
			// just a common regular character, do nothing
		}
	}

	// Fix format string if last char was a percent (and there is EOL).

	if prevWasPercent {
		*format += "%"
	}

	return verbsCount
}
