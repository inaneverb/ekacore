// Copyright © 2020-2021. All rights reserved.
// Author: Ilya Stroy.
// Contacts: iyuryevich@pm.me, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekalog

import (
	"bytes"
	"io"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/qioalice/ekago/v3/ekasys"
	"github.com/qioalice/ekago/v3/internal/ekaletter"

	"github.com/json-iterator/go"
)

//noinspection GoSnakeCaseUsage
type (
	// _CICE_FormatPart represents built format's part (parsed RAW format string).
	// - for 'typ' == '_CICE_FPT_VERB_JUST_TEXT', 'value' is original RAW format string's part;
	// - for 'typ' == '_CICE_FPT_VERB_COLOR_CUSTOM', 'value' is calculated bash escape sequence
	//   (like "\033[01;03;38;05;144m");
	// - for other 'typ' variants, 'value' is empty, 'cause it's log's verbs
	// and they're runtime calculated entities.
	_CICE_FormatPart struct {
		typ   _CICE_FormatPartType
		value string
	}

	// _CICE_FormatPartType is a special type of _CICE_FormatPart's field 'typ' that contains
	// an info what kind of format verb current _CICE_FormatPart object is
	// and how exactly it will be converted to the text at the runtime.
	_CICE_FormatPartType uint16

	_CICE_FieldsFormat struct {
		isSet                bool
		beforeFields         string
		afterFields          string
		beforeKey            string
		afterKey             string
		afterValue           string
		afterNewLine         string
		afterNewLineForError string
		itemsPerLine         int16
	}

	_CICE_BodyFormat struct {
		isSet      bool
		beforeBody string
		afterBody  string
	}

	_CICE_CallerFormat struct {
		isSet     bool
		isDefault bool
		parts     [10]struct {
			typ int16
			val string
		}
	}

	_CICE_StacktraceFormat struct {
		isSet       bool
		beforeStack string
		afterStack  string
	}

	_CICE_ErrorFormat struct {
		isSet       bool
		beforeError string
		afterError  string
	}

	_CICE_DropColors struct {
		buf  bytes.Buffer
		dest io.Writer
	}
)

//noinspection GoSnakeCaseUsage
const (
	// Common Integrator Console Encoder Format Part Type (CICE FPT)
	// predefined constants.

	_CICE_FPT_MASK_TYPE _CICE_FormatPartType = 0x00_FF
	_CICE_FPT_MASK_DATA _CICE_FormatPartType = 0xFF_00

	_CICE_FPT_VERB_JUST_TEXT       _CICE_FormatPartType = 0x01
	_CICE_FPT_VERB_COLOR_CUSTOM    _CICE_FormatPartType = 0x02
	_CICE_FPT_VERB_COLOR_FOR_LEVEL _CICE_FormatPartType = 0x03
	_CICE_FPT_VERB_BODY            _CICE_FormatPartType = 0x0A
	_CICE_FPT_VERB_TIME            _CICE_FormatPartType = 0x0B
	_CICE_FPT_VERB_LEVEL           _CICE_FormatPartType = 0x0C
	_CICE_FPT_VERB_STACKTRACE      _CICE_FormatPartType = 0x1A
	_CICE_FPT_VERB_FIELDS          _CICE_FormatPartType = 0x2A
	_CICE_FPT_VERB_CALLER          _CICE_FormatPartType = 0x3A

	// Common Integrator Console Encoder Level Format (CICE LF)
	// type constants.

	_CICE_LF_NUMBER           _CICE_FormatPartType = 1
	_CICE_LF_SHORT_NORMAL     _CICE_FormatPartType = 2
	_CICE_LF_SHORT_UPPER_CASE _CICE_FormatPartType = 3
	_CICE_LF_FULL_NORMAL      _CICE_FormatPartType = 4
	_CICE_LF_FULL_UPPER_CASE  _CICE_FormatPartType = 5

	// Common Integrator Console Encoder Time Format (CICE TF)
	// type constants.

	_CICE_TF_TIMESTAMP _CICE_FormatPartType = 1
	_CICE_TF_ANSIC     _CICE_FormatPartType = 2
	_CICE_TF_UNIXDATE  _CICE_FormatPartType = 3
	_CICE_TF_RUBYDATE  _CICE_FormatPartType = 4
	_CICE_TF_RFC822    _CICE_FormatPartType = 5
	_CICE_TF_RFC822_Z  _CICE_FormatPartType = 6
	_CICE_TF_RFC850    _CICE_FormatPartType = 7
	_CICE_TF_RFC1123   _CICE_FormatPartType = 8
	_CICE_TF_RFC1123_Z _CICE_FormatPartType = 9
	_CICE_TF_RFC3339   _CICE_FormatPartType = 10

	// Common Integrator Console Encoder Caller Format (CICE CF)
	// type constants.

	_CICE_CF_TYPE_SEPARATOR  int16 = -1
	_CICE_CF_TYPE_FUNC_SHORT int16 = 1
	_CICE_CF_TYPE_FUNC_FULL  int16 = 2
	_CICE_CF_TYPE_FILE_SHORT int16 = 3
	_CICE_CF_TYPE_FILE_FULL  int16 = 4
	_CICE_CF_TYPE_LINE_NUM   int16 = 5
	_CICE_CF_TYPE_PKG_SHORT  int16 = 6 // unused
	_CICE_CF_TYPE_PKG_FULL   int16 = 7

	// Common Integrator Console Encoder (CICE) verb predefined constants.

	_CICE_VERB_START_INDICATOR rune = '{'
	_CICE_VERB_END_INDICATOR   rune = '}'
	_CICE_VERB_SEPARATOR       byte = '/'

	// Common Integrator Console Integrator Standard Colors (CICE SC)
	// predefined constants.

	_CICE_SC_DEBUG     string = "c/fg:#b2b2b2"
	_CICE_SC_INFO      string = "c/fg:#87d7ff"
	_CICE_SC_NOTICE    string = "c/fg:#00d700"
	_CICE_SC_WARNING   string = "c/fg:#ff8700"
	_CICE_SC_ERROR     string = "c/fg:#ff0000"
	_CICE_SC_CRITICAL  string = "c/fg:#ffaf00/b"
	_CICE_SC_ALERT     string = "c/fg:#ffaf00/b/u"
	_CICE_SC_EMERGENCY string = "c/fg:#ff00ff/b/u"

	// Common Integrator Console Encoder (CICE) defaults.

	_CICE_DEFAULT_FORMAT string = "{{c}}[{{l/SS}}] <{{t}}>:{{c/0}} " + // include colored level, colored time
		"{{m/?$\n}}" + // include message with \n if non-empty
		"{{f/?$\n/v = /e, /l\t/le\t\t}}" + // include fields with " = " as key-value separator
		"{{s/?$\n/e, }}" + // include stacktrace with \n if non-empty
		"{{w/0/fd}}" + // omit caller, specify each stacktrace's frame format
		"\n"

	_CICE_DEFAULT_TIME_FORMAT string = "Mon Jan 02 15:04:05"
)

// CommonIntegrator CI_ConsoleEncoder Verb Types (CICE VT)
// that are used in format string to determine what kind of verb must be used.
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
	// Make sure we won't break API by declaring package's console encoder.
	defaultConsoleEncoder CI_Encoder
)

// Type extracts verb's type from current _CICE_FormatPartType that can be compared
// with _CEFPT_VERB... constants.
func (fpt _CICE_FormatPartType) Type() _CICE_FormatPartType {
	return fpt & _CICE_FPT_MASK_TYPE
}

// Data extracts verb's type data from current _CICE_FormatPartType that is needed
// for some internal verb's encoders.
func (fpt _CICE_FormatPartType) Data() _CICE_FormatPartType {
	return (fpt & _CICE_FPT_MASK_DATA) >> 8
}

// doBuild builds the current CI_ConsoleEncoder only if it has not built yet.
// There is no-op if format string is empty or encoder already built.
func (ce *CI_ConsoleEncoder) doBuild() *CI_ConsoleEncoder {

	switch {
	case ce == nil:
		return nil

	case len(ce.formatParts) > 0:
		// do not build if it's so already
		return ce
	}

	ce.format = strings.TrimSpace(ce.format)
	if ce.format == "" {
		ce.format = _CICE_DEFAULT_FORMAT
	}

	// start parsing ce.format
	// all parsing loops are for-range based (because there is UTF-8 support)
	// (yes, you can use not only ASCII parts in your format string,
	// and yes if you do it, you are mad. stop it!).
	for rest := ce.format; rest != ""; rest = ce.parseFirstVerb(rest) {
	}

	ce.uniteJustTextVerbs()
	ce.setStandardParts()

	return ce
}

// uniteJustTextVerbs unites "just text" verbs in 'ce.formatParts'
// that follows each other. It may happen when something with bad verbs
// were included in 'ce.format'.
func (ce *CI_ConsoleEncoder) uniteJustTextVerbs() *CI_ConsoleEncoder {

	var (
		// idx of "just text" verb next "just text" verbs will be united with
		justTextVerbIdx = -1
		// new out slice of verbs
		newFormatParts = make([]_CICE_FormatPart, 0, len(ce.formatParts))
	)

	for idx, verb := range ce.formatParts {
		switch verbTyp := verb.typ.Type(); {

		case justTextVerbIdx == -1 && verbTyp == _CICE_FPT_VERB_JUST_TEXT:
			justTextVerbIdx = idx

		case justTextVerbIdx != -1 && verbTyp == _CICE_FPT_VERB_JUST_TEXT:
			ce.formatParts[justTextVerbIdx].value += verb.value

		case justTextVerbIdx != -1 && verbTyp != _CICE_FPT_VERB_JUST_TEXT:
			newFormatParts = append(newFormatParts, ce.formatParts[justTextVerbIdx])
			justTextVerbIdx = -1
			fallthrough

		case justTextVerbIdx == -1 && verbTyp != _CICE_FPT_VERB_JUST_TEXT:
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

// setStandardParts saves standard colors for standard log levels
// if they has not been set yet.
func (ce *CI_ConsoleEncoder) setStandardParts() *CI_ConsoleEncoder {

	if ce.colorMap == nil {
		ce.colorMap = make(map[Level]string)
	}

	if ce.colorMap[LEVEL_DEBUG] == "" {
		ce.colorMap[LEVEL_DEBUG] = ce.rvColorHelper(_CICE_SC_DEBUG)
	}
	if ce.colorMap[LEVEL_INFO] == "" {
		ce.colorMap[LEVEL_INFO] = ce.rvColorHelper(_CICE_SC_INFO)
	}
	if ce.colorMap[LEVEL_NOTICE] == "" {
		ce.colorMap[LEVEL_NOTICE] = ce.rvColorHelper(_CICE_SC_NOTICE)
	}
	if ce.colorMap[LEVEL_WARNING] == "" {
		ce.colorMap[LEVEL_WARNING] = ce.rvColorHelper(_CICE_SC_WARNING)
	}
	if ce.colorMap[LEVEL_ERROR] == "" {
		ce.colorMap[LEVEL_ERROR] = ce.rvColorHelper(_CICE_SC_ERROR)
	}
	if ce.colorMap[LEVEL_CRITICAL] == "" {
		ce.colorMap[LEVEL_CRITICAL] = ce.rvColorHelper(_CICE_SC_CRITICAL)
	}
	if ce.colorMap[LEVEL_ALERT] == "" {
		ce.colorMap[LEVEL_ALERT] = ce.rvColorHelper(_CICE_SC_ALERT)
	}
	if ce.colorMap[LEVEL_EMERGENCY] == "" {
		ce.colorMap[LEVEL_EMERGENCY] = ce.rvColorHelper(_CICE_SC_EMERGENCY)
	}

	if !ce.cf.isSet {
		ce.cf.isDefault = true
	}

	return ce
}

// parseFirstVerb parses 'format', extracts first verb (even if it's "just text"
// verb), saves it to ce.formatParts and then returns the rest of 'format' string.
func (ce *CI_ConsoleEncoder) parseFirstVerb(format string) (rest string) {

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
		case c == _CICE_VERB_START_INDICATOR && pc == _CICE_VERB_START_INDICATOR && wv:
			// unexpected "{{" inside complex verb, treat all prev as "just text",
			// try to treat as starting complex verb
			wv = false
			i--

		case c == _CICE_VERB_START_INDICATOR && pc == _CICE_VERB_START_INDICATOR && i > 1:
			// > 1 (not > 0) because if string started with "{{", after first "{"
			// i already == 1.
			//
			// was "just text", found complex verb start
			i--

		case c == _CICE_VERB_END_INDICATOR && pc == _CICE_VERB_END_INDICATOR && wv:
			// ending of complex verb
			wve = true
			i++

		case c == _CICE_VERB_START_INDICATOR && pc == _CICE_VERB_START_INDICATOR:
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

	// what kind of verb did we parse and whether verb has been closed correctly?
	if wv && wve {
		ce.minimumBufferLen += ce.rv(format[:i])
	} else {
		ce.minimumBufferLen += ce.rvJustText(format[:i])
	}

	return format[i:]
}

// rv (resolve verb) tries to determine what kind of complex verb 'verb' is,
// creates related '_CICE_FormatPart', fills it and adds to 'ce.formatParts'.
// Returned predicted minimum length of bytes that required in buffer to store
// formatted verb.
func (ce *CI_ConsoleEncoder) rv(verb string) (predictedLen int) {

	// applyOnce is a helper func to avoid many if-else statements in tht switch below.
	applyOnce := func(isSet *bool, fallback, applicator func(string) int, verb string) int {
		if !*isSet {
			*isSet = true
			return applicator(verb)
		} else {
			return fallback(verb)
		}
	}

	// it guarantees that "verb" starts from "{{",
	// so we can remove leading "{{" and trailing "}}"
	switch verb = verb[2 : len(verb)-2]; {

	case hpm(verb, cevtCaller):
		return applyOnce(&ce.cf.isSet, ce.rvJustText, ce.rvCaller, verb)

	case hpm(verb, cevtColor):
		return ce.rvColor(verb)

	case hpm(verb, cevtLevel):
		return ce.rvLevel(verb)

	case hpm(verb, cevtTime):
		return ce.rvTime(verb)

	case hpm(verb, cevtMessage):
		return applyOnce(&ce.bf.isSet, ce.rvJustText, ce.rvBody, verb)

	case hpm(verb, cevtFields):
		return applyOnce(&ce.ff.isSet, ce.rvJustText, ce.rvFields, verb)

	case hpm(verb, cevtStacktrace):
		return applyOnce(&ce.sf.isSet, ce.rvJustText, ce.rvStacktrace, verb)

	default:
		// incorrect verb, treat it as "just text" verb
		return ce.rvJustText(verb)
	}
}

// rvHelper is a part of "resolve verb" functions but moreover it's a helper.
// This function literally have no recognition algorithm but splits 'verb'
// to verb parts (ignoring first part, assuming that its verb's name) and then
// calls 'resolver' for each that part.
// Stops splitting and calling 'resolver' if it's return 'false' one time.
//
// It guarantees that if 'verbPart' is empty, 'resolver' won't be called for that
// and processing will stopped.
//
// Requirements:
// 'verb' != "", 'resolver' != nil. Otherwise no-op.
func (_ *CI_ConsoleEncoder) rvHelper(verb string, resolver func(verbPart string) (continue_ bool)) {

	if verb == "" || resolver == nil {
		return
	}

	// skip verb name
	if idx := strings.IndexByte(verb, _CICE_VERB_SEPARATOR); idx != -1 {
		verb = verb[idx+1:]
	} else {
		// there is no verb separator, it's simple verb. Do nothing
		return
	}

	for idx := strings.IndexByte(verb, _CICE_VERB_SEPARATOR); idx != -1; {
		if idx == 0 || !resolver(verb[:idx]) {
			goto STOP
		}
		verb = verb[idx+1:]
		idx = strings.IndexByte(verb, _CICE_VERB_SEPARATOR)
	}

	// parse last entity
	if !resolver(verb) {
		goto STOP
	}

STOP:
	return
}

func (ce *CI_ConsoleEncoder) rvJustText(text string) (predictedLen int) {

	if text != "" {
		ce.formatParts = append(ce.formatParts, _CICE_FormatPart{
			typ:   _CICE_FPT_VERB_JUST_TEXT,
			value: text,
		})
	}

	return len(text)
}

func (ce *CI_ConsoleEncoder) rvLevel(verb string) (predictedLen int) {

	formattedLevel := _CICE_LF_FULL_NORMAL
	predictedLen = 9

	if idx := strings.IndexByte(verb, _CICE_VERB_SEPARATOR); idx != -1 {
		switch verb[idx+1:] {
		case "d", "D":
			formattedLevel = _CICE_LF_NUMBER
			predictedLen = 1
		case "s":
			formattedLevel = _CICE_LF_SHORT_NORMAL
			predictedLen = 5
		case "S":
			formattedLevel = _CICE_LF_SHORT_UPPER_CASE
			predictedLen = 5
		case "ss":
			formattedLevel = _CICE_LF_FULL_NORMAL
			predictedLen = 9
		case "SS":
			formattedLevel = _CICE_LF_FULL_UPPER_CASE
			predictedLen = 9
		}
	}

	ce.formatParts = append(ce.formatParts, _CICE_FormatPart{
		typ: _CICE_FPT_VERB_LEVEL | (formattedLevel << 8),
	})

	return predictedLen
}

func (ce *CI_ConsoleEncoder) rvTime(verb string) (predictedLen int) {

	format := _CICE_DEFAULT_TIME_FORMAT
	formattedTime := _CICE_FormatPartType(0)

	(*CI_ConsoleEncoder)(nil).rvHelper(verb, func(verbPart string) (continue_ bool) {
		if verbPart = strings.TrimSpace(format); verbPart != "" {
			switch predefined := strings.ToUpper(verbPart); predefined {
			case "UNIX", "TIMESTAMP":
				formattedTime = _CICE_TF_TIMESTAMP
			case "ANSIC":
				formattedTime = _CICE_TF_ANSIC
			case "UNIXDATE", "UNIX_DATE":
				formattedTime = _CICE_TF_UNIXDATE
			case "RUBYDATE", "RUBY_DATE":
				formattedTime = _CICE_TF_RUBYDATE
			case "RFC822":
				formattedTime = _CICE_TF_RFC822
			case "RFC822Z":
				formattedTime = _CICE_TF_RFC822_Z
			case "RFC850":
				formattedTime = _CICE_TF_RFC850
			case "RFC1123":
				formattedTime = _CICE_TF_RFC1123
			case "RFC1123Z":
				formattedTime = _CICE_TF_RFC1123_Z
			case "RFC3339":
				formattedTime = _CICE_TF_RFC3339
			default:
				format = verbPart
			}
		}
		return false // only first time verb is allowed and will be parsed
	})

	ce.formatParts = append(ce.formatParts, _CICE_FormatPart{
		typ:   _CICE_FPT_VERB_TIME | (formattedTime << 8),
		value: format,
	})

	return len(time.Now().Format(format)) + 10 // stock for some weekdays
}

func (ce *CI_ConsoleEncoder) rvColor(verb string) (predictedLen int) {

	if idx := strings.IndexByte(verb, _CICE_VERB_SEPARATOR); idx == -1 {
		ce.formatParts = append(ce.formatParts, _CICE_FormatPart{
			typ: _CICE_FPT_VERB_COLOR_FOR_LEVEL,
		})
		return ce.colorMapMax
	}

	if encodedColor := ce.rvColorHelper(verb); encodedColor != "" {
		ce.formatParts = append(ce.formatParts, _CICE_FormatPart{
			typ:   _CICE_FPT_VERB_COLOR_CUSTOM,
			value: encodedColor,
		})
		return len(encodedColor)
	} else {
		return ce.rvJustText(verb)
	}
}

func (_ *CI_ConsoleEncoder) rvColorHelper(colorVerb string) string {

	cb := colorBuilder{}
	cb.init()

	(*CI_ConsoleEncoder)(nil).rvHelper(colorVerb, func(verbPart string) (continue_ bool) {
		return cb.parseEntity(verbPart)
	})

	return cb.encode()
}

// rvBody is a part of "resolve verb" functions.
// rvBody tries to parse 'verb' as logger Entry's body and anyway indicates that
// here will be stored Entry's body.
//
// Verb's arguments:
// - For non-empty body:
//   - "?^<text>": <text> will be prepended to the Entry's body at the runtime.
//   - "?$<text>": <text> will be appended to the Entry's body at the runtime.
func (ce *CI_ConsoleEncoder) rvBody(verb string) (predictedLen int) {

	(*CI_ConsoleEncoder)(nil).rvHelper(verb, func(verbPart string) (continue_ bool) {
		switch {
		case strings.HasPrefix(verbPart, "?^"):
			ce.bf.beforeBody = verbPart[2:]
		case strings.HasPrefix(verbPart, "?$"):
			ce.bf.afterBody = verbPart[2:]
		default:
			return false
		}
		return true
	})

	ce.formatParts = append(ce.formatParts, _CICE_FormatPart{
		typ: _CICE_FPT_VERB_BODY,
	})

	return 256 + len(ce.bf.beforeBody) + len(ce.bf.afterBody)
}

func (ce *CI_ConsoleEncoder) rvCaller(verb string) (predictedLen int) {

	isAdd := true
	formatPrefixes := []string{"f", "F"}

	(*CI_ConsoleEncoder)(nil).rvHelper(verb, func(verbPart string) (continue_ bool) {
		switch {
		case verbPart == "0":
			isAdd = false
		case hpm(verbPart, formatPrefixes):
			predictedLen += ce.rvCallerFormat(verbPart[1:])
		default:
			return false
		}
		return true
	})

	if isAdd {
		ce.formatParts = append(ce.formatParts, _CICE_FormatPart{
			typ: _CICE_FPT_VERB_CALLER,
		})
		return predictedLen
	} else {
		return 0
	}
}

func (ce *CI_ConsoleEncoder) rvCallerFormat(f string) (predictedLen int) {

	if f == "" || f == "d" || f == "D" {
		ce.cf.isDefault = true
		return 256
	}

	j := 0 // index of ce.cf.parts
	for _, fc := range f {

		t := _CICE_CF_TYPE_SEPARATOR // by default threat it as a separator
		switch fc {

		case 'w':
			t = _CICE_CF_TYPE_FUNC_SHORT
		case 'W':
			t = _CICE_CF_TYPE_FUNC_FULL
		case 'f':
			t = _CICE_CF_TYPE_FILE_SHORT
		case 'F':
			t = _CICE_CF_TYPE_FILE_FULL
		case 'l', 'L':
			t = _CICE_CF_TYPE_LINE_NUM
		case 'p', 'P':
			t = _CICE_CF_TYPE_PKG_FULL

		default:
			switch ce.cf.parts[j].typ {
			default:
				j++
				fallthrough
			case 0:
				ce.cf.parts[j].typ = _CICE_CF_TYPE_SEPARATOR
				fallthrough
			case -1:
				ce.cf.parts[j].val += string(fc)
			}
			continue
		}

		// maybe it's already existed? skip it if it so
		for i, n := 0, len(ce.cf.parts); i < n && ce.cf.parts[i].typ != 0; i++ {
			if ce.cf.parts[i].typ == t {
				t = _CICE_CF_TYPE_SEPARATOR
			}
		}

		if t != _CICE_CF_TYPE_SEPARATOR {
			if ce.cf.parts[j].typ != 0 {
				j++
			}
			ce.cf.parts[j].typ = t
		}
	}

	return 256
}

// rvFields is a part of "resolve verb" functions.
// rvFields tries to parse 'verb' as logger Entry's fields (not attached Error's)
// but anyway indicates that here will be stored Entry's body.
//
// Verb's arguments:
// - If at least one LetterField presented:
//   - "?^<text>": <text> will be write before any LetterField is written at the runtime.
//   - "?$<text>": <text> will be appended to the end of last LetterField at the runtime.
//   - "k<text>": <text> will be written before LetterField's keys is written.
//   - "v<text>": <text> will be written before LetterField's value is written.
//   - "e<text>": <text> will be written after LetterField's value excluding last.
//   - "l<text>": <text> will be written at the each new line of fields' part set.
//   - "*<int>": <int> is how much fields are placed at the one line
//     (by default: 4. Use <= 0 value to place all fields at the one line).
func (ce *CI_ConsoleEncoder) rvFields(verb string) (predictedLen int) {

	ce.ff.itemsPerLine = 4

	(*CI_ConsoleEncoder)(nil).rvHelper(verb, func(verbPart string) (continue_ bool) {
		switch upperCased := strings.ToUpper(verbPart); {

		case strings.HasPrefix(verbPart, "?^"):
			ce.ff.beforeFields = verbPart[2:]
		case strings.HasPrefix(verbPart, "?$"):
			ce.ff.afterFields = verbPart[2:]
		case strings.HasPrefix(upperCased, "LE"):
			ce.ff.afterNewLineForError = verbPart[2:]
		case upperCased[0] == 'L':
			ce.ff.afterNewLine = verbPart[1:]
		case upperCased[0] == 'K':
			ce.ff.beforeKey = verbPart[1:]
		case upperCased[0] == 'V':
			ce.ff.afterKey = verbPart[1:]
		case upperCased[0] == 'E':
			ce.ff.afterValue = verbPart[1:]

		case verbPart[0] == '*':
			if perLine_, err := strconv.Atoi(verbPart[1:]); err == nil {
				if perLine_ < 0 {
					ce.ff.itemsPerLine = 0
				} else {
					ce.ff.itemsPerLine = int16(perLine_)
				}
			}

		default:
			return false
		}
		return true
	})

	ce.formatParts = append(ce.formatParts, _CICE_FormatPart{
		typ: _CICE_FPT_VERB_FIELDS,
	})

	return 512 + len(ce.ff.beforeFields) + len(ce.ff.afterFields) + len(ce.ff.beforeKey) +
		len(ce.ff.afterKey) + len(ce.ff.afterValue) + len(ce.ff.afterNewLine)
}

func (ce *CI_ConsoleEncoder) rvStacktrace(verb string) (predictedLen int) {

	(*CI_ConsoleEncoder)(nil).rvHelper(verb, func(verbPart string) (continue_ bool) {
		switch {
		case strings.HasPrefix(verbPart, "?^"):
			ce.sf.beforeStack = verbPart[2:]
		case strings.HasPrefix(verbPart, "?$"):
			ce.sf.afterStack = verbPart[2:]
		default:
			return false
		}
		return true
	})

	ce.formatParts = append(ce.formatParts, _CICE_FormatPart{
		typ: _CICE_FPT_VERB_STACKTRACE,
	})

	return 2048
}

func (ce *CI_ConsoleEncoder) encodeJustText(to []byte, fp _CICE_FormatPart) []byte {
	return bufw(to, fp.value)
}

func (ce *CI_ConsoleEncoder) encodeLevel(to []byte, fp _CICE_FormatPart, e *Entry) []byte {

	formattedLevel := ""
	switch fp.typ.Data() {

	case _CICE_LF_NUMBER:
		formattedLevel = strconv.Itoa(int(e.Level))
	case _CICE_LF_SHORT_NORMAL:
		formattedLevel = e.Level.String3()
	case _CICE_LF_SHORT_UPPER_CASE:
		formattedLevel = e.Level.ToUpper3()
	case _CICE_LF_FULL_NORMAL:
		formattedLevel = e.Level.String()
	case _CICE_LF_FULL_UPPER_CASE:
		formattedLevel = strings.ToUpper(e.Level.String())
	}

	return bufw(to, formattedLevel)
}

func (ce *CI_ConsoleEncoder) encodeTime(e *Entry, fp _CICE_FormatPart, to []byte) []byte {

	formattedTime := ""

	switch fp.typ.Data() {
	case _CICE_TF_TIMESTAMP:
		formattedTime = strconv.FormatInt(e.Time.Unix(), 10)
	case _CICE_TF_ANSIC:
		formattedTime = e.Time.Format(time.ANSIC)
	case _CICE_TF_UNIXDATE:
		formattedTime = e.Time.Format(time.UnixDate)
	case _CICE_TF_RUBYDATE:
		formattedTime = e.Time.Format(time.RubyDate)
	case _CICE_TF_RFC822:
		formattedTime = e.Time.Format(time.RFC822)
	case _CICE_TF_RFC822_Z:
		formattedTime = e.Time.Format(time.RFC822Z)
	case _CICE_TF_RFC850:
		formattedTime = e.Time.Format(time.RFC850)
	case _CICE_TF_RFC1123:
		formattedTime = e.Time.Format(time.RFC1123)
	case _CICE_TF_RFC1123_Z:
		formattedTime = e.Time.Format(time.RFC1123Z)
	case _CICE_TF_RFC3339:
		formattedTime = e.Time.Format(time.RFC3339)
	default:
		formattedTime = e.Time.Format(fp.value)
	}

	return bufw(to, formattedTime)
}

func (ce *CI_ConsoleEncoder) encodeColor(to []byte, fp _CICE_FormatPart) []byte {
	return bufw(to, fp.value)
}

func (ce *CI_ConsoleEncoder) encodeColorForLevel(to []byte, e *Entry) []byte {
	if color := ce.colorMap[e.Level]; color != "" {
		return bufw(to, color)
	}
	return to
}

func (ce *CI_ConsoleEncoder) encodeBody(to []byte, e *Entry) []byte {

	body := e.LogLetter.Messages[0].Body
	if body == "" {
		return to
	}

	if ce.bf.beforeBody != "" {
		to = bufw(to, ce.bf.beforeBody)
	}

	to = bufw(to, body)

	if ce.bf.afterBody != "" {
		to = bufw(to, ce.bf.afterBody)
	}

	return to
}

func (ce *CI_ConsoleEncoder) encodeCaller(to []byte, e *Entry) []byte {

	var frame ekasys.StackFrame

	switch {
	case len(e.LogLetter.StackTrace) > 0:
		frame = e.LogLetter.StackTrace[0]

	case e.ErrLetter != nil:
		frame = e.ErrLetter.StackTrace[0]

	default:
		return to
	}

	return ce.encodeStackFrame(to, frame, nil, ekaletter.LetterMessage{})
}

func (ce *CI_ConsoleEncoder) encodeFields(to []byte, fs, addFs []ekaletter.LetterField, isErrors, addPreEncoded bool) []byte {

	if len(fs) == 0 {
		return to
	}

	if !isErrors && ce.ff.beforeFields != "" {
		to = bufw(to, ce.ff.beforeFields)
	}

	var (
		unnamedFieldIdx, writtenFields int16
	)

	addField := func(to []byte, f *ekaletter.LetterField, isErrors bool, unnamedFieldIex, writtenFields *int16) []byte {
		if strings.HasPrefix(f.Key, "sys.") {
			return to
		}

		keyBak := f.Key

		if f.Key == "" && !f.IsSystem() {
			f.Key = f.KeyOrUnnamed(&unnamedFieldIdx)
		}

		toLenBak := len(to)
		to = ce.encodeField(to, *f, isErrors, *writtenFields)
		if len(to) != toLenBak {
			*writtenFields++
		}

		f.Key = keyBak
		return to
	}

	for i, n := int16(0), int16(len(fs)); i < n; i++ {
		to = addField(to, &fs[i], isErrors, &unnamedFieldIdx, &writtenFields)
	}
	for i, n := int16(0), int16(len(addFs)); i < n; i++ {
		to = addField(to, &addFs[i], isErrors, &unnamedFieldIdx, &writtenFields)
	}

	if addPreEncoded && ce.preEncodedFieldsWritten > 0 {
		if ce.ff.afterValue != "" {
			to = to[:len(to)-len(ce.ff.afterValue)]
		}
		if l := len(to); to[l-1] != '\n' {
			to = bufw(to, "\n")
		}
		to = bufw2(to, ce.preEncodedFields)
	}

	if writtenFields > 0 {

		// remove last "after value" and write "after fields"
		if ce.ff.afterValue != "" {
			to = to[:len(to)-len(ce.ff.afterValue)]
		}

		if !isErrors && ce.ff.afterFields != "" {
			to = bufw(to, ce.ff.afterFields)
		}

	} else {

		// no fields were written
		// remove 'beforeFields' part
		to = to[:len(to)-len(ce.ff.beforeFields)]
	}

	return to
}

func (ce *CI_ConsoleEncoder) encodeField(to []byte, f ekaletter.LetterField, isErrors bool, fieldNum int16) []byte {

	// Maybe field wants to be started with new line?
	oldKey := f.Key
	if f.Key = strings.TrimSpace(f.Key); len(oldKey) != len(f.Key) {
		to = bufw(to, "\n")
	}

	// write new line and new line title
	if ce.ff.itemsPerLine > 0 && fieldNum != 0 && fieldNum%ce.ff.itemsPerLine == 0 {
		to = bufw(to, "\n")
	}

	if wasNewLine := to[len(to)-1] == '\n'; wasNewLine && !isErrors && len(ce.ff.afterNewLine) > 0 {
		to = bufw(to, ce.ff.afterNewLine)
	} else if wasNewLine && isErrors && len(ce.ff.afterNewLineForError) > 0 {
		to = bufw(to, ce.ff.afterNewLineForError)
	}

	// Write "before key", key, "after key", value and "after value"

	if ce.ff.beforeKey != "" {
		to = bufw(to, ce.ff.beforeKey)
	}
	to = bufw(to, f.Key)
	if ce.ff.afterKey != "" {
		to = bufw(to, ce.ff.afterKey)
	}
	to = ce.encodeFieldValue(to, f)
	if ce.ff.afterValue != "" {
		to = bufw(to, ce.ff.afterValue)
	}

	return to
}

func (ce *CI_ConsoleEncoder) encodeFieldValue(to []byte, f ekaletter.LetterField) []byte {

	if f.Kind.IsSystem() {
		switch f.Kind.BaseType() {

		case ekaletter.KIND_SYS_TYPE_EKAERR_UUID, ekaletter.KIND_SYS_TYPE_EKAERR_CLASS_NAME:
			to = bufw(to, `"`)
			to = bufw(to, f.SValue)
			to = bufw(to, `"`)

		case ekaletter.KIND_SYS_TYPE_EKAERR_CLASS_ID:
			to = strconv.AppendInt(to, f.IValue, 10)

		default:
			to = bufw(to, `"<unsupported system field>"`)
		}

	} else if f.Kind.IsNil() {
		to = bufw(to, "null")

	} else if f.Kind.IsInvalid() {
		to = bufw(to, "<invalid_field>")

	} else {
		switch f.Kind.BaseType() {

		case ekaletter.KIND_TYPE_BOOL:
			to = strconv.AppendBool(to, f.IValue != 0)

		case ekaletter.KIND_TYPE_INT,
			ekaletter.KIND_TYPE_INT_8, ekaletter.KIND_TYPE_INT_16,
			ekaletter.KIND_TYPE_INT_32, ekaletter.KIND_TYPE_INT_64:
			to = strconv.AppendInt(to, f.IValue, 10)

		case ekaletter.KIND_TYPE_UINT,
			ekaletter.KIND_TYPE_UINT_8, ekaletter.KIND_TYPE_UINT_16,
			ekaletter.KIND_TYPE_UINT_32, ekaletter.KIND_TYPE_UINT_64:
			to = strconv.AppendUint(to, uint64(f.IValue), 10)

		case ekaletter.KIND_TYPE_FLOAT_32:
			f := float64(math.Float32frombits(uint32(f.IValue)))
			to = strconv.AppendFloat(to, f, 'f', 2, 32)

		case ekaletter.KIND_TYPE_FLOAT_64:
			f := math.Float64frombits(uint64(f.IValue))
			to = strconv.AppendFloat(to, f, 'f', 2, 64)

		case ekaletter.KIND_TYPE_UINTPTR, ekaletter.KIND_TYPE_ADDR:
			to = bufw(to, "0x")
			to = strconv.AppendInt(to, f.IValue, 16)

		case ekaletter.KIND_TYPE_STRING:
			to = strconv.AppendQuote(to, f.SValue)

		case ekaletter.KIND_TYPE_COMPLEX_64:
			r := math.Float32frombits(uint32(f.IValue >> 32))
			i := math.Float32frombits(uint32(f.IValue))
			c := complex128(complex(r, i))
			// TODO: Use strconv.AppendComplex() when it will be released.
			to = bufw(to, strconv.FormatComplex(c, 'f', 2, 32))

		case ekaletter.KIND_TYPE_COMPLEX_128:
			c := f.Value.(complex128)
			// TODO: Use strconv.AppendComplex() when it will be released.
			to = bufw(to, strconv.FormatComplex(c, 'f', 2, 64))

		case ekaletter.KIND_TYPE_UNIX:
			to = bufw(to, time.Unix(f.IValue, 0).Format("Jan 2 15:04:05"))

		case ekaletter.KIND_TYPE_UNIX_NANO:
			to = bufw(to, time.Unix(0, f.IValue).Format("Jan 2 15:04:05.000000000"))

		case ekaletter.KIND_TYPE_DURATION:
			to = bufw(to, time.Duration(f.IValue).String())

		case ekaletter.KIND_TYPE_MAP, ekaletter.KIND_TYPE_EXTMAP:
			// TODO: Add support of extracted maps.
			if jsonedMap, legacyErr := jsoniter.Marshal(f.Value); legacyErr == nil {
				to = bufw2(to, jsonedMap)
			} else {
				to = bufw(to, "<unsupported_map>")
			}

		case ekaletter.KIND_TYPE_STRUCT:
			if jsonedStruct, legacyErr := jsoniter.Marshal(f.Value); legacyErr == nil {
				to = bufw2(to, jsonedStruct)
			} else {
				to = bufw(to, "<unsupported_struct>")
			}

		case ekaletter.KIND_TYPE_ARRAY:
			if jsonedArray, legacyErr := jsoniter.Marshal(f.Value); legacyErr == nil {
				to = bufw2(to, jsonedArray)
			} else {
				to = bufw(to, "<unsupported_array>")
			}

		default:
			to = bufw(to, "<unsupported_field>")
		}
	}

	return to
}

func (ce *CI_ConsoleEncoder) encodeStacktrace(to []byte, e *Entry) []byte {

	stacktrace := e.LogLetter.StackTrace
	if len(stacktrace) == 0 && e.ErrLetter != nil {
		stacktrace = e.ErrLetter.StackTrace
	}

	n := int16(len(stacktrace))
	if n == 0 {
		return to
	}

	if ce.sf.beforeStack != "" {
		to = bufw(to, ce.sf.beforeStack)
	}

	var (
		fi       = 0 // fi for fields' index
		mi       = 0 // mi for messages' index
		fields   []ekaletter.LetterField
		messages []ekaletter.LetterMessage
	)

	if e.ErrLetter != nil {
		fields = e.ErrLetter.Fields
		messages = e.ErrLetter.Messages
	}

	for i := int16(0); i < n; i++ {
		messageForStackFrame := ekaletter.LetterMessage{}
		fieldsForStackFrame := []ekaletter.LetterField(nil)
		fiEnd := 0

		//goland:noinspection GoNilness
		if mi < len(messages) && messages[mi].StackFrameIdx == i {
			messageForStackFrame = messages[mi]
			mi++
		}

		if fi < len(fields) && fields[fi].StackFrameIdx == i {
			fiEnd = fi + 1
			for fiEnd < len(fields) && fields[fiEnd].StackFrameIdx == i {
				fiEnd++
			}
		}

		if fiEnd != 0 {
			fieldsForStackFrame = fields[fi:fiEnd]
		}

		to = ce.encodeStackFrame(to, stacktrace[i], fieldsForStackFrame, messageForStackFrame)

		if i < n-1 {
			to = bufw(to, "\n")
		}
	}

	if ce.sf.afterStack != "" {
		to = bufw(to, ce.sf.afterStack)
	}

	return to
}

func (ce *CI_ConsoleEncoder) encodeStackFrame(

	to []byte,
	frame ekasys.StackFrame,
	fields []ekaletter.LetterField,
	message ekaletter.LetterMessage,

) []byte {

	if !ce.cf.isDefault {
		// Reminder: frame.DoFormat does once:
		// "<package>/<func> (<short_file>:<file_line>) <full_package_path>".

		for i, n := 0, len(ce.cf.parts); i < n && ce.cf.parts[i].typ != 0; i++ {
			switch ce.cf.parts[i].typ {

			case _CICE_CF_TYPE_SEPARATOR:
				to = bufw(to, ce.cf.parts[i].val)

			case _CICE_CF_TYPE_FUNC_SHORT:
				frame.DoFormat()
				to = bufw(to, frame.Format[:frame.FormatFileOffset-1])

			case _CICE_CF_TYPE_FUNC_FULL:
				to = bufw(to, frame.Function)

			case _CICE_CF_TYPE_FILE_SHORT:
				frame.DoFormat()
				i := strings.IndexByte(frame.Format, ':')
				to = bufw(to, frame.Format[frame.FormatFileOffset+1:i])

			case _CICE_CF_TYPE_FILE_FULL:
				to = bufw(to, frame.File)

			case _CICE_CF_TYPE_LINE_NUM:
				to = bufw(to, strconv.Itoa(frame.Line))

			case _CICE_CF_TYPE_PKG_FULL:
				to = bufw(to, frame.Format[frame.FormatFullPathOffset:])
			}
		}
	} else {
		to = bufw(to, frame.DoFormat())
	}

	if message.Body != "" || len(fields) > 0 {
		to = bufw(to, "\n")

		if ce.ff.afterNewLineForError != "" {
			to = bufw(to, ce.ff.afterNewLineForError)
		}

		if message.Body != "" {
			to = bufw(to, message.Body)
			to = bufw(to, "\n")
		}

		lToBak := len(to)
		to = ce.encodeFields(to, fields, nil, true, false)

		// ce.encodeFields may write no fields. Then we must clear last "\n"
		if len(to) == lToBak {
			to = to[:len(to)-1]
		}
	}

	return to
}

func (dc *_CICE_DropColors) Write(p []byte) (n int, err error) {
	if n, err = dc.buf.Write(p); err != nil && err != io.EOF {
		return n, err
	}
	b := dc.buf.Bytes()
	j := 0
	for i, n, write := 0, len(b), true; i < n; i++ {
		if write && b[i] == '\033' {
			write = false
			i += 2
		} else if !write && b[i] == 'm' {
			write = true
		} else if write {
			b[j] = b[i]
			j++
		}
	}
	return dc.dest.Write(b[:j])
}
