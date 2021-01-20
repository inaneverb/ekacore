// Copyright © 2020. All rights reserved.
// Author: Ilya Stroy.
// Contacts: qioalice@gmail.com, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekastr

/*

*/
func Interpolate(s string, cbVerbFound, cbTextFound func(v string)) {
	f1 := func(v []byte) { cbVerbFound(B2S(v)) }
	f2 := func(v []byte) { cbTextFound(B2S(v)) }
	Interpolateb(S2B(s), f1, f2)
}

/*

*/
func Interpolateb(p []byte, cbVerbFound, cbTextFound func(v []byte)) {

	var (
		part   []byte
		isVerb bool
	)

	if len(p) == 0 {
		return
	}

	for len(p) > 0 {
		part, p, isVerb = i9nGetNext(p)
		if isVerb {
			cbVerbFound(part)
		} else {
			cbTextFound(part)
		}
	}
}

/*

*/
// parseFirstVerb parses 'format', extracts first verb (even if it's "just text"
// verb), saves it to ce.formatParts and then returns the rest of 'format' string.
func i9nGetNext(p []byte) (part, nextP []byte, isVerb bool) {

	if len(p) == 0 {
		return nil, nil, false
	}

	//goland:noinspection GoSnakeCaseUsage
	const (
		VERB_START_INDICATOR byte = '{'
		VERB_END_INDICATOR   byte = '}'
	)

	var (
		i   = 0
		pc  byte // prev char
		wv  bool // true if current parsing verb is complex verb (not "just text")
		wve bool // true if complex verb has been closed correctly
	)

	for _, c := range p {
		switch {
		case c == VERB_START_INDICATOR && pc == VERB_START_INDICATOR && wv:
			// unexpected "{{" inside complex verb, treat all prev as "just text",
			// try to treat as starting complex verb
			wv = false
			i--

		case c == VERB_START_INDICATOR && pc == VERB_START_INDICATOR && i > 1:
			// > 1 (not > 0) because if string started with "{{", after first "{"
			// i already == 1.
			//
			// was "just text", found complex verb start
			i--

		case c == VERB_END_INDICATOR && pc == VERB_END_INDICATOR && wv:
			// ending of complex verb
			wve = true
			i++

		case c == VERB_START_INDICATOR && pc == VERB_START_INDICATOR:
			// this is the beginning of 'format' and of complex verb
			wv = true
			i++
			continue

		default:
			pc = c
			i++
			continue
		}
		break
	}

	return p[:i], p[i:], wv && wve
}
