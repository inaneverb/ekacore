// Copyright © 2019. All rights reserved.
// Author: Ilya Yuryevich.
// Contacts: qioalice@gmail.com, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package log

import (
	"strconv"
	"strings"
	"time"
	"unsafe"

	"github.com/qioalice/gext/dangerous"
)

//
type ConsoleEncoder struct {

	// -----
	// RAW FORMAT
	// This is the ConsoleEncoder's parts that are considered raw.
	// These parts are set directly by some methods or with minimum aggregation.
	// They will be parsed, converted and casted to internal structures
	// that are more convenient to be used as source of log message format
	// generator.
	// -----
	format string

	// -----
	// BUILT FORMAT
	// This part represents parsed raw format string.
	// First of all RAW format string parsed to the two group of entities:
	// - Just string, not verb (writes as is),
	// - Format verbs (will be substituted to the log's parts, such as
	//   log's message, timestamp, log's level, log's fields, stacktrace, etc).
	// These parts (w/o reallocation) will be stored with the same sequence
	// as they represented in 'format' but there will be specified verb's types also.
	// Moreover their common length will be calculated and stored to decrease
	// destination []byte buffer reallocations (just allocate big buffer,
	// at least as more as 'formatParts' required).
	// -----
	formatParts []formatPart
	// Sum of: len of just text parts + predicted len of log's parts.
	minimumBufferLen int
}

// formatPart represents built format's part (parsed RAW format string).
// - for 'typ' == 'cefptVerbJustText', 'value' is original RAW format string's part;
// - for 'typ' == 'cefptVerbColor', 'value' is calculated bash escape sequence
// (like "\033[01;03;38;05;144m");
// - for other 'typ' variants, 'value' is empty, 'cause it's log's verbs
// and they're runtime calculated entities.
type formatPart struct {
	typ   formatPartType
	value string
}

//
type formatPartType uint16

// [c]onsole [e]ncoder [f]ormat [p]art [t]ype (cefpt) predefined constants.
const (
	cefptMaskType formatPartType = 0x00FF
	cefptMaskData formatPartType = 0xFF00

	// TODO: formatPartType.setd,getd -> setter/getter of data with offset
	cefptOffsetData int = 8

	cefptVerbJustText formatPartType = 0x01
	cefptVerbColor    formatPartType = 0x02
	cefptVerbBody     formatPartType = 0x0A
	cefptVerbTime     formatPartType = 0x0B
	cefptVerbLevelD   formatPartType = 0x0C
	cefptVerbLevelS   formatPartType = 0x0D
	cefptVerbLevelSS  formatPartType = 0x0E
)

//
const (
	ceVerbStartIndicator rune = '{'
	ceVerbEndIndicator   rune = '}'
	ceVerbSeparator      byte = '|'
)

// [c]onsole [e]ncoder [v]erb [t]ype values, that are used in format string
// to determine what kind of verb must be used.
var (
	cevtCaller     = []string{"caller", "who", "w"}
	cevtColor      = []string{"color", "c"}
	cevtLevel      = []string{"level", "lvl", "l"}
	cevtTime       = []string{"time", "t"}
	cevtMessage    = []string{"message", "body", "m", "b"}
	cevtFields     = []string{"fields", "f"}
	cevtStacktrace = []string{"stacktrace", "s"}
)

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

func (fpt formatPartType) Type() formatPartType {
	return fpt & cefptMaskType
}

func (fpt formatPartType) Data() formatPartType {
	return fpt & cefptMaskData
}

//
func (ce *ConsoleEncoder) FreezeAndGetEncoder() CommonIntegratorEncoder {

	return ce.doBuild().encode
}

//
func (ce *ConsoleEncoder) SetFormat() {
	// TODO: ce.minimumBufferLen = 0
}

//
func (ce *ConsoleEncoder) doBuild() *ConsoleEncoder {

	switch {
	case ce == nil:
		return nil

	case len(ce.formatParts) > 0:
		// do not build if it's so already
		return ce
	}

	// start parsing ce.format
	// all parsing loops are for-range based (because there is UTF-8 support)
	// (yes, you can use not only ASCII parts in your format string,
	// and yes if you do it, you are mad. stop it!).
	for rest := ce.format; rest != ""; rest = ce.parseFirstVerb(rest) {
	}

	// TODO: Optimization: Unite paired "just text" verbs
	return ce
}

// parseFirstVerb parses 'format', extracts first verb (even if it's "just text"
// verb), saves it to ce.formatParts and then returns the rest of 'format' string.
func (ce *ConsoleEncoder) parseFirstVerb(format string) (rest string) {

	// "complex verb" - is any other verb than "just text" verb.

	var (
		i  = 0
		pc rune // prev char
		wv bool // true if currently parsing complex verb
	)

	// loop is for-range based (because there is UTF-8 support)
	// (yes, you can use not only ASCII parts in your format string,
	// and yes if you do it, you are mad. stop it!).
	for _, c := range format {
		switch {
		case c == ceVerbStartIndicator && pc == ceVerbStartIndicator && wv:
			// unexpected "{{" inside complex verb, treat all prev as "just text",
			// try to treat as starting complex verb
			wv = false
			i--

		case c == ceVerbStartIndicator && pc == ceVerbStartIndicator && i > 0:
			// was "just text", found complex verb start
			i--

		case c == ceVerbEndIndicator && pc == ceVerbEndIndicator && wv:
			// ending of complex verb
			i++

		case c == ceVerbStartIndicator && pc == ceVerbStartIndicator:
			// this is the beginning of 'format' and of complex verb
			wv = true
			continue

		default:
			pc = c
			continue
		}
		break
	}

	// what kind of verb we have?
	if wv {
		ce.minimumBufferLen += ce.rv(format[:i])
	} else {
		ce.minimumBufferLen += ce.rvJustText(format[:i])
	}

	return format[i:]
}

// rv (resolve verb) tries to determine what kind of complex verb 'verb' is,
// creates related 'formatPart', fills it and adds to ce.formatParts.
// Returned predicted minimum length of bytes that required in buffer to store
// formatted verb.
func (ce *ConsoleEncoder) rv(verb string) (predictedLen int) {

	// hpm is "has prefix many"
	// just like strings.HasPrefix, but you can check many prefixes at the same time.
	hpm := func(verb string, prefixes []string) bool {
		for _, prefix := range prefixes {
			if strings.HasPrefix(verb, prefix) {
				return true
			}
		}
		return false
	}

	// it guarantees that "verb" starts from "{{"
	switch verb = verb[2 : len(verb)-2]; {

	case hpm(verb, cevtCaller):
		return ce.rvCaller(verb)

	case hpm(verb, cevtColor):
		return ce.rvColor(verb)

	case hpm(verb, cevtLevel):
		return ce.rvLevel(verb)

	case hpm(verb, cevtTime):
		return ce.rvTime(verb)

	case hpm(verb, cevtMessage):
		return ce.rvBody(verb)

	case hpm(verb, cevtFields):
		return ce.rvFields(verb)

	case hpm(verb, cevtStacktrace):
		return ce.rvStacktrace(verb)

	default:
		// incorrect verb, treat it as "just text" verb
		return ce.rvJustText(verb)
	}
}

// rvLevel is a part of "resolve verb" functions.
//
func (ce *ConsoleEncoder) rvLevel(verb string) (predictedLen int) {

	var typ = cefptVerbLevelSS // by default print full string level

	if idx := strings.IndexByte(verb, ceVerbSeparator); idx != -1 {
		switch verb[idx+1:] {
		case "d":
			typ = cefptVerbLevelD
			predictedLen = 3
		case "s":
			typ = cefptVerbLevelS
			predictedLen = 3
		case "ss":
			typ = cefptVerbLevelS
			// guess it's enough to store any logger.Level full string
			// representation
			predictedLen = 16
		}
	}

	ce.formatParts = append(ce.formatParts, formatPart{
		typ: typ,
	})

	return predictedLen
}

// rvTime is a part of "resolve verb" functions.
//
func (ce *ConsoleEncoder) rvTime(verb string) (predictedLen int) {

	// TODO: Add possibility to encode time as unix timestamp.
	var format string

	if idx := strings.IndexByte(verb, ceVerbSeparator); idx != -1 {
		if format = strings.TrimSpace(verb[idx+1:]); format == "" {
			format = time.RFC1123
		}
	} else {
		format = time.RFC1123
	}

	ce.formatParts = append(ce.formatParts, formatPart{
		typ:   cefptVerbTime,
		value: format,
	})

	return len(format) + 10 // stock for some weekdays
}

//
func (ce *ConsoleEncoder) rvColor(verb string) (predictedLen int) {

}

//
func (ce *ConsoleEncoder) rvCaller(verb string) (predictedLen int) {

}

// rvBody is a part of "resolve verb" functions.
//
func (ce *ConsoleEncoder) rvBody(verb string) (predictedLen int) {

	ce.formatParts = append(ce.formatParts, formatPart{
		typ: cefptVerbBody,
	})
	return 256
}

//
func (ce *ConsoleEncoder) rvFields(verb string) (predictedLen int) {

}

//
func (ce *ConsoleEncoder) rvStacktrace(verb string) (predictedLen int) {

}

// rvJustText is a part of "resolve verb" functions.
//
func (ce *ConsoleEncoder) rvJustText(text string) (predictedLen int) {

	if text != "" {
		ce.formatParts = append(ce.formatParts, formatPart{
			typ:   cefptVerbJustText,
			value: text,
		})
	}

	return len(text)
}

//
func (ce *ConsoleEncoder) encode(e *Entry) []byte {

}

//
func (ce *ConsoleEncoder) encLevel(e *Entry, fp formatPart, to []byte) []byte {

	formattedLevel := ""
	switch fp.typ {

	case cefptVerbLevelD:
		formattedLevel = strconv.Itoa(int(e.Level))

	case cefptVerbLevelS:
		formattedLevel = e.Level.String()[:3]

	case cefptVerbLevelSS:
		formattedLevel = e.Level.String()
	}

	return append(bufgr(to, len(formattedLevel)), formattedLevel...)
}

//
func (ce *ConsoleEncoder) encTime(e *Entry, fp formatPart, to []byte) []byte {

	formattedTime := e.Time.Format(fp.value)
	if formattedTime == "" {
		formattedTime = e.Time.Format(time.RFC1123)
	}

	return append(bufgr(to, len(formattedTime)), formattedTime...)
}

//
func (ce *ConsoleEncoder) encCaller(e *Entry, fp formatPart, to []byte) []byte {

}

//
func (ce *ConsoleEncoder) encFields(e *Entry, fp formatPart, to []byte) []byte {

}

//
func (ce *ConsoleEncoder) encStacktrace(e *Entry, fp formatPart, to []byte) []byte {

}

// easy case because ASCII sequence already generated at the rvColor method.
func (ce *ConsoleEncoder) encColor(fp formatPart, to []byte) []byte {
	return bufw(to, fp.value)
}

// easy case because fp.value is the text we should add.
func (ce *ConsoleEncoder) encJustText(fp formatPart, to []byte) []byte {
	return bufw(to, fp.value)
}

// easy case because e.Message is the text we should add.
func (ce *ConsoleEncoder) encBody(e *Entry, to []byte) []byte {
	return bufw(to, e.Message)
}
