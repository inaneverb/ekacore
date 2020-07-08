// Copyright © 2020. All rights reserved.
// Author: Ilya Stroy.
// Contacts: qioalice@gmail.com, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekalog

import (
	"image/color"
	"math"
	"strconv"
	"strings"
	"time"
	"unsafe"

	"github.com/qioalice/ekago/v2/ekadanger"
	"github.com/qioalice/ekago/v2/internal/field"
	"github.com/qioalice/ekago/v2/internal/letter"
	"github.com/qioalice/ekago/v2/internal/xtermcolor"
)

// TODO: Update doc, comments

// "complex verb" - is any other verb than "just text" verb.

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

//
type colorBuilder struct {

	// https://en.wikipedia.org/wiki/ANSI_escape_code#SGR_parameters

	// [0..255] - xterm256 ANSI SGR color code
	// -1 if 'do cleanup color to terminal default' ( "\033[39m" or "\033[49m" )
	// -2 if 'not set, use those one that was used' (not included to SGR)
	bg, fg int16

	// 0 - 'not set, use those one that was used' (not included to SGR)
	// 1 - enable (included to SGR (01/03/04))
	// -1 - disable (included to SGR (22/23/24))
	bold, italic, underline int8

	// Yes, there is no support blanking text.
	// I think it's disgusting. It will never be added.
}

// [c]onsole [e]ncoder [f]ormat [p]art [t]ype (cefpt) predefined constants.
const (
	cefptMaskType formatPartType = 0x00FF
	cefptMaskData formatPartType = 0xFF00

	// TODO: formatPartType.setd,getd -> setter/getter of data with offset
	cefptOffsetData int = 8

	cefptVerbJustText formatPartType = 0x01
	cefptVerbColor    formatPartType = 0x02

	cefptVerbBody    formatPartType = 0x0A
	cefptVerbTime    formatPartType = 0x0B
	cefptVerbLevelD  formatPartType = 0x0C
	cefptVerbLevelS  formatPartType = 0x0D
	cefptVerbLevelSS formatPartType = 0x0E

	cefptVerbStack  formatPartType = 0x1A
	cefptVerbFields formatPartType = 0x2A
	cefptVerbCaller formatPartType = 0x3A
)

//
const (
	ceVerbStartIndicator rune = '{'
	ceVerbEndIndicator   rune = '}'
	ceVerbSeparator      byte = '/'
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
	consoleEncoderAddr = ekadanger.TakeRealAddr(consoleEncoder)
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

// doBuild builds the current console encoder only if it has not built yet.
// There is no-op if format string is empty.
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

	return ce.uniteJustTextVerbs()
}

// uniteJustTextVerbs just unites "just text" verbs in 'ce.formatParts'
// that follows each other. It could be occur when something with bad verbs
// were included in 'ce.format'.
func (ce *ConsoleEncoder) uniteJustTextVerbs() *ConsoleEncoder {

	var (
		// idx of "just text" verb next "just text" verbs will be united with
		justTextVerbIdx = -1
		// new out slice of verbs
		newFormatParts = make([]formatPart, 0, len(ce.formatParts))
	)

	for idx, verb := range ce.formatParts {
		switch verbTyp := verb.typ.Type(); {

		case justTextVerbIdx == -1 && verbTyp == cefptVerbJustText:
			justTextVerbIdx = idx

		case justTextVerbIdx != -1 && verbTyp == cefptVerbJustText:
			ce.formatParts[justTextVerbIdx].value += verb.value

		case justTextVerbIdx != -1 && verbTyp != cefptVerbJustText:
			newFormatParts = append(newFormatParts, ce.formatParts[justTextVerbIdx])
			justTextVerbIdx = -1
			fallthrough

		case justTextVerbIdx == -1 && verbTyp != cefptVerbJustText:
			newFormatParts = append(newFormatParts, verb)
		}
	}

	// it was definitely the last "just text" verb we should add
	if justTextVerbIdx != -1 {
		newFormatParts = append(newFormatParts, ce.formatParts[justTextVerbIdx])
	}

	ce.formatParts = newFormatParts[:]
	return ce
}

// parseFirstVerb parses 'format', extracts first verb (even if it's "just text"
// verb), saves it to ce.formatParts and then returns the rest of 'format' string.
func (ce *ConsoleEncoder) parseFirstVerb(format string) (rest string) {

	var (
		i   = 0
		pc  rune // prev char
		wv  bool // true if current parsing verb is complex verb (not "just text")
		wve bool // true if complex verb has been closed correctly
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

		case c == ceVerbStartIndicator && pc == ceVerbStartIndicator && i > 1:
			// > 1 (not > 0) because if string started with "{{", after first "{"
			// i already == 1.
			//
			// was "just text", found complex verb start
			i--

		case c == ceVerbEndIndicator && pc == ceVerbEndIndicator && wv:
			// ending of complex verb
			wve = true
			i++

		case c == ceVerbStartIndicator && pc == ceVerbStartIndicator:
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

	// what kind of verb we did parse and did verb has been closed correctly?
	if wv && wve {
		ce.minimumBufferLen += ce.rv(format[:i])
	} else {
		ce.minimumBufferLen += ce.rvJustText(format[:i])
	}

	return format[i:]
}

// rv (resolve verb) tries to determine what kind of complex verb 'verb' is,
// creates related 'formatPart', fills it and adds to 'ce.formatParts'.
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

	// it guarantees that "verb" starts from "{{",
	// so we can remove leading "{{" and trailing "}}"
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
			typ = cefptVerbLevelSS
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
func (cb *colorBuilder) init() {

	cb.bg, cb.fg = -2, -2
	cb.bold, cb.italic, cb.underline = 0, 0, 0
}

//
func (cb *colorBuilder) parseEntity(verbPart string) (parsed bool) {

	switch verbPart = strings.ToUpper(strings.TrimSpace(verbPart)); verbPart {
	// --- REMINDER! 1ST ARGUMENT IS ALWAYS UPPER CASED! ---

	case "":
		return true

	case "BOLD", "B":
		cb.bold = 1
		return true

	case "NOBOLD", "NOB":
		cb.bold = -1
		return true

	case "ITALIC", "I":
		cb.italic = 1
		return true

	case "NOITALIC", "NOI":
		cb.italic = -1
		return true

	case "UNDERLINE", "U":
		cb.underline = 1
		return true

	case "NOUNDERLINE", "NOU":
		cb.underline = -1
		return true
	}

	// okay, it's color, but which one? what format? RGB? HEX? RGBA? (lol, wtf?)
	// TODO: Add supporting of color's literals like "red", "pink", "blue", etc.

	// what's kind of color? default is fg
	var colorDestination *int16
	switch {

	case strings.HasPrefix(verbPart, "BG"):
		colorDestination = &cb.bg
		verbPart = strings.TrimSpace(verbPart[2:])

	case strings.HasPrefix(verbPart, "FG"): // already defaulted
		colorDestination = &cb.fg
		verbPart = strings.TrimSpace(verbPart[2:])

	default:
		colorDestination = &cb.fg
	}

	// handle special easy cases cases
	switch {
	case len(verbPart) == 0:
		// rare case, wasn't more chars after "bg" or "fg"
		return false

	case verbPart == "-1":
		// color cleanup (to default in term)
		*colorDestination = -1
		return

	case verbPart[0] == '#':
		// easy case if it's explicit hex
		return cb.parseHexTo(verbPart[1:], colorDestination)
	}

	// okay, maybe easy rgb/rgba?
	switch {
	case strings.HasPrefix(verbPart, "RGB"):
		return cb.parseRgbTo(verbPart[3:], colorDestination)

	case strings.HasPrefix(verbPart, "RGBA"):
		return cb.parseRgbTo(verbPart[4:], colorDestination)

	case strings.HasPrefix(verbPart, "RGB(") && verbPart[len(verbPart)-1] == ')':
		return cb.parseRgbTo(verbPart[4:len(verbPart)-1], colorDestination)

	case strings.HasPrefix(verbPart, "RGBA(") && verbPart[len(verbPart)-1] == ')':
		return cb.parseRgbTo(verbPart[4:len(verbPart)-1], colorDestination)
	}

	// okay maybe rgb by comma?
	if commas := strings.Count(verbPart, ","); commas >= 3 && commas <= 4 {
		return cb.parseRgbTo(verbPart, colorDestination)
	}

	// believe it's hex
	return cb.parseHexTo(verbPart, colorDestination)
}

//
func (cb *colorBuilder) parseHexTo(verbPart string, destination *int16) (parsed bool) {
	// --- REMINDER! 1ST ARGUMENT IS ALWAYS UPPER CASED! ---

	switch verbPart = strings.TrimSpace(verbPart); len(verbPart) {
	case 4:
		// short case with alpha, ignore alpha, extend to 6
		verbPart = verbPart[:3]
		fallthrough
	case 3:
		// short case, extend to 6
		var hexParts [6]uint8
		hexParts[0], hexParts[1] = verbPart[0], verbPart[0]
		hexParts[2], hexParts[3] = verbPart[1], verbPart[1]
		hexParts[4], hexParts[5] = verbPart[2], verbPart[2]
		verbPart = string(hexParts[:])
	case 8:
		// with alpha, ignore it
		verbPart = verbPart[:6]
	case 6:
		// default HEX case, handle later
	default:
		return false
	}

	termColor, err := xtermcolor.FromHexStr(verbPart)
	*destination = int16(termColor)

	return err == nil
}

//
func (cb *colorBuilder) parseRgbTo(verbPart string, destination *int16) (parsed bool) {
	// --- REMINDER! 1ST ARGUMENT IS ALWAYS UPPER CASED! ---

	rgbParts := strings.Split(strings.TrimSpace(verbPart), ",")
	if l := len(rgbParts); l < 3 && l > 4 {
		return false
	}

	var (
		r, g, b          int
		err1, err2, err3 error
	)

	r, err1 = strconv.Atoi(strings.TrimSpace(rgbParts[0]))
	g, err2 = strconv.Atoi(strings.TrimSpace(rgbParts[1]))
	b, err3 = strconv.Atoi(strings.TrimSpace(rgbParts[2]))

	if err1 != nil || err2 != nil || err3 != nil ||
		r < 0 || r > 255 || g < 0 || g > 255 || b < 0 || b > 255 {
		return false
	}

	rgb := color.RGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: 255}
	*destination = int16(xtermcolor.FromColor(rgb))

	return true
}

//
func (cb *colorBuilder) encode() string {

	// TODO: Here is too much Golang string mem reallocations
	//  maybe use []byte instead of string with only one allocation ?

	// https://en.wikipedia.org/wiki/ANSI_escape_code#SGR_parameters
	out := "\033["

	switch cb.bold {
	case 1:
		out += "01;" // enable bold
	case -1:
		out += "22;" // disable bold
	}

	switch cb.italic {
	case 1:
		out += "03;" // enable italic
	case -1:
		out += "23;" // disable italic
	}

	switch cb.underline {
	case 1:
		out += "04;" // enable underline
	case -1:
		out += "24;" // disable underline
	}

	switch cb.fg {
	case -2:
		// do nothing, use those one that was used
	case -1:
		// set foreground to term default
		out += "39;"
	default:
		out += "38;5;" + strconv.Itoa(int(cb.fg)) + ";"
	}

	switch cb.fg {
	case -2:
		// do nothing, use those one that was used
	case -1:
		// set background to term default
		out += "49;"
	default:
		out += "48;5;" + strconv.Itoa(int(cb.bg)) + ";"
	}

	if out[len(out)-1] != ';' {
		return "" // all values are default and ignored
	}

	out = out[:len(out)-1] + "m"
	return out
}

//
// ""
func (ce *ConsoleEncoder) rvColor(verb string) (predictedLen int) {

	var (
		cb           colorBuilder
		encodedColor string
		verbBak      = verb
	)

	cb.init()

	// skip verb identifier (this is color verb)
	if idx := strings.IndexByte(verb, ceVerbSeparator); idx != -1 {
		verb = verb[idx+1:]

		for idx = strings.IndexByte(verb, ceVerbSeparator); idx != -1; {
			if parsingFailed := !cb.parseEntity(verb[:idx]); parsingFailed {
				goto UNSUPPORTED_COLOR_ENTITY
			}
			verb = verb[idx+1:]
		}

		// parse last entity
		if parsingFailed := !cb.parseEntity(verb); parsingFailed {
			goto UNSUPPORTED_COLOR_ENTITY
		}
	}

	if encodedColor = cb.encode(); encodedColor == "" {
		goto UNSUPPORTED_COLOR_ENTITY
	}

	ce.formatParts = append(ce.formatParts, formatPart{
		typ:   cefptVerbColor,
		value: encodedColor,
	})

	return len(encodedColor)

UNSUPPORTED_COLOR_ENTITY: // label and goto only for make code more clean and clear
	return ce.rvJustText(verbBak)
}

//

// rvBody is a part of "resolve verb" functions.
//
func (ce *ConsoleEncoder) rvBody(verb string) (predictedLen int) {

	ce.formatParts = append(ce.formatParts, formatPart{
		typ: cefptVerbBody,
	})
	return 256
}

//
func (ce *ConsoleEncoder) rvCaller(verb string) (predictedLen int) {

	// TODO: Implement fields format

	ce.formatParts = append(ce.formatParts, formatPart{
		typ: cefptVerbCaller,
	})

	return 256
}

//
func (ce *ConsoleEncoder) rvFields(verb string) (predictedLen int) {

	// TODO: Implement fields format

	ce.formatParts = append(ce.formatParts, formatPart{
		typ: cefptVerbFields,
	})

	return 512
}

//
func (ce *ConsoleEncoder) rvStacktrace(verb string) (predictedLen int) {

	// TODO: Implement stacktrace format

	ce.formatParts = append(ce.formatParts, formatPart{
		typ: cefptVerbStack,
	})

	return 2048
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

	// TODO: Reuse allocated buffers

	buf := make([]byte, 0, ce.minimumBufferLen)

	for _, part := range ce.formatParts {
		switch part.typ.Type() {

		case cefptVerbJustText:
			buf = ce.encJustText(part, buf)

		case cefptVerbColor:
			buf = ce.encColor(part, buf)

		case cefptVerbBody:
			buf = ce.encBody(e, buf)

		case cefptVerbTime:
			buf = ce.encTime(e, part, buf)

		case cefptVerbLevelD, cefptVerbLevelS, cefptVerbLevelSS:
			buf = ce.encLevel(e, part, buf)

		case cefptVerbStack:
			//buf = ce.encStacktrace(e, part, buf)

		case cefptVerbFields:
			buf = ce.encFields(e.LogLetter.Items.Fields, part, buf)

		case cefptVerbCaller:
			//buf = ce.encCaller(e, part, buf)
		}
	}

	return buf
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
//func (ce *ConsoleEncoder) encCaller(e *Entry, fp formatPart, to []byte) []byte {
//
//	// TODO: Implement caller format
//
//	if e.Caller.PC == 0 {
//		return to
//	}
//
//	return bufw(to, e.Caller.DoFormat())
//}

//
//func (ce *ConsoleEncoder) encStacktrace(e *Entry, fp formatPart, to []byte) []byte {
//
//	// TODO: Implement stacktrace format
//	// Reminder: e.StackTrace.ExcludeInternal already called
//
//	if l := len(e.StackTrace); l > 0 {
//		for _, stackFrame := range e.StackTrace[:l-1] {
//			to = ce.encStackFrame(stackFrame, fp, to)
//			to = bufw(to, "\n")
//		}
//		to = ce.encStackFrame(e.StackTrace[l-1], fp, to)
//	}
//
//	return to
//}

//
//func (ce *ConsoleEncoder) encStackFrame(frame syse.StackFrame, fp formatPart, to []byte) []byte {
//
//	_, file := filepath.Split(frame.Function)
//	s := frame.Function + " (" + file + ":" + strconv.Itoa(frame.Line) + ")"
//
//	return bufw(to, s)
//}

//
func (ce *ConsoleEncoder) encFields(fields []field.Field, fp formatPart, to []byte) []byte {

	// TODO: Implement fields format

	lFields := len(fields)

	if lFields == 0 {
		return to
	}

	unnamedFieldIdx := 1

	for _, field_ := range fields[:lFields-1] {
		to = ce.encField(field_, to, &unnamedFieldIdx)
		to = bufw(to, ", ")
	}
	to = ce.encField(fields[lFields-1], to, &unnamedFieldIdx)

	return to
}

//
func (ce *ConsoleEncoder) encField(f field.Field, to []byte, unnamedFieldIdx *int) []byte {

	field_ := f.Key
	if field_ == "" {
		field_ = letter.UnnamedAsStr(*unnamedFieldIdx)
		*unnamedFieldIdx++
	}

	field_ += " = "

	switch f.Kind {
	case field.KIND_TYPE_BOOL:
		if f.IValue != 0 {
			field_ += "true"
		} else {
			field_ += "false"
		}

	case field.KIND_TYPE_INT, field.KIND_TYPE_INT_8, field.KIND_TYPE_INT_16, field.KIND_TYPE_INT_32, field.KIND_TYPE_INT_64:
		field_ += strconv.FormatInt(f.IValue, 10)

	case field.KIND_TYPE_UINT, field.KIND_TYPE_UINT_8, field.KIND_TYPE_UINT_16, field.KIND_TYPE_UINT_32, field.KIND_TYPE_UINT_64:
		field_ += strconv.FormatUint(uint64(f.IValue), 10)

	case field.KIND_TYPE_FLOAT_32:
		field_ += strconv.FormatFloat(float64(math.Float32frombits(uint32(f.IValue))), 'f', 2, 32)

	case field.KIND_TYPE_FLOAT_64:
		field_ += strconv.FormatFloat(float64(math.Float32frombits(uint32(f.IValue))), 'f', 2, 64)

	case field.KIND_TYPE_STRING:
		field_ += f.SValue

	default:
		return to
	}

	return bufw(to, field_)
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
	return bufw(to, e.LogLetter.Items.Message)
}
