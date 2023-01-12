// Copyright © 2020. All rights reserved.
// Author: Ilya Stroy.
// Contacts: iyuryevich@pm.me, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekastr

// Interpolate splits string into pieces, calling 'cbVerbFound' or 'cbTextFound'
// for each piece (depending on piece's kind) consequentially.
//
// There are only 2 pieces of piece's kind:
//   - Verb. Some text inside double brackets. Like: "{{verb}}".
//     So, if verb is found, a related callback is called ('cbVerbFound').
//   - Just text. Anything other but verb.
//
// Note.
// In the example above, "verb" string will be passed to the 'cbVerbFound',
// not "{{verb}}.
//
// So, the main reason is:
// You define your verbs, like "print", "call" or something like that,
// use them in string, like "blah-blah1 {{print}} blah-blah2 {{call}} blah3"
// and then define your 'cbVerbFound' with switch statement,
// where you determine behaviour depend on verb - "print" or "call".
//
// Note.
// It guarantees, that order of pieces of string is preserved.
// So, using example of previous note, the order of calls is:
// - 'cbTextFound', argument "blah-blah1 ",
// - 'cbVerbFound', argument "print",
// - 'cbTextFound', argument " blah-blah2 ",
// - 'cbVerbFound', argument "call",
// - 'cbTextFound', argument " blah3".
// Keep your eyes on the spaces in "just text" pieces.
func Interpolate(s string, cbVerbFound, cbTextFound func(v string)) {
	var f1 = func(v []byte) { cbVerbFound(FromBytes(v)) }
	var f2 = func(v []byte) { cbTextFound(FromBytes(v)) }
	InterpolateBytes(ToBytes(s), f1, f2)
}

// InterpolateBytes is the same as Interpolate, but works with []byte
// instead of string.
func InterpolateBytes(p []byte, cbVerbFound, cbTextFound func(v []byte)) {

	var part []byte
	var isVerb bool

	for len(p) > 0 {
		part, p, isVerb = getNextPart(p)
		if isVerb {
			cbVerbFound(part)
		} else {
			cbTextFound(part)
		}
	}
}

// getNextPart parses 'format', extracts next part of string.
// It could be "just text" or a verb. Returns found part, rest of the string
// and a flag whether returned part is a verb.
func getNextPart(p []byte) (part, tail []byte, isVerb bool) {

	if len(p) == 0 {
		return nil, nil, false
	}

	const (
		VerbStart byte = '{'
		VerbEnd   byte = '}'
	)

	var i = 0
	var pc byte  // prev char
	var wv bool  // true if current parsing part is a verb (not "just text")
	var wve bool // true if verb has been closed correctly

	for _, c := range p {
		switch {
		case c == VerbStart && pc == VerbStart && wv:
			// unexpected "{{" inside complex verb, treat all prev as "just text",
			// try to treat these chars as complex verb starting
			wv = false
			i--

		case c == VerbStart && pc == VerbStart && i > 1:
			// > 1 (not > 0) because if string started with "{{", after first "{"
			// i already == 1.
			//
			// was "just text", found complex verb start
			i--

		case c == VerbEnd && pc == VerbEnd && wv:
			// verb's ending
			wve = true
			i++

		case c == VerbStart && pc == VerbStart:
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
