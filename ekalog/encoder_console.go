// Copyright © 2020. All rights reserved.
// Author: Ilya Stroy.
// Contacts: qioalice@gmail.com, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekalog

import (
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/qioalice/ekago/v2/ekasys"
	"github.com/qioalice/ekago/v2/internal/field"
	"github.com/qioalice/ekago/v2/internal/letter"
)

//noinspection GoSnakeCaseUsage
type (
	// CI_ConsoleEncoder is a type that built to be used as a part of CommonIntegrator
	// as an log Entries encoder to the TTY (console) with human readable format,
	// custom format supporting and ability to enable coloring.
	//
	// If you want to use CI_ConsoleEncoder, you need to instantiate object,
	// set log's text format string using SetFormat() and then call
	// FreezeAndGetEncoder() method. By that you'll get the function that has
	// an alias CI_Encoder and you can add it as encoder by
	// CommonIntegrator.WithEncoder().
	//
	// See https://github.com/qioalice/ekago/ekalog/integrator.go ,
	// https://github.com/qioalice/ekago/ekalog/integrator_common.go for more info.
	CI_ConsoleEncoder struct {

		// RAW FORMAT
		// This is what user is set by SetFormat() method.
		//
		// It will be parsed, converted and casted to internal structures
		// that are more convenient to be used as source of log message format
		// generator.
		format string

		// BUILT FORMAT
		// This part represents parsed raw format string.
		//
		// First of all RAW format string parsed to the two group of entities:
		// - Just string, not verb (writes as is),
		// - Format verbs (will be substituted to the log's parts, such as
		//   log's message, timestamp, log's level, log's fields, stacktrace, etc).
		//
		// These parts (w/o reallocation) will be stored with the same sequence
		// as they represented in 'format' but there will be specified verb's types also.
		// Moreover their common length will be calculated and stored to decrease
		// destination []byte buffer reallocations (just allocate big buffer,
		// at least as more as 'formatParts' required).
		formatParts []_CICE_FormatPart

		colorMap    map[Level]string // map of default colors for each level
		colorMapMax int              // max used len of ASCII color encoded seq.

		// Sum of: len of just text parts + predicted len of log's parts.
		minimumBufferLen int

		ff _CICE_FieldsFormat     // parts of fields formatting
		bf _CICE_BodyFormat       // parts of body formatting
		cf _CICE_CallerFormat     // parts of caller, stacktrace's frames formatting
		sf _CICE_StacktraceFormat // parts of stacktrace formatting
		ef _CICE_ErrorFormat      // parts of attached error formatting
	}

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
)

//noinspection GoSnakeCaseUsage
const (
	// Common Integrator Console Encoder Format Part Type (CICE FPT)
	// predefined constants.

	_CICE_FPT_MASK_TYPE            _CICE_FormatPartType = 0x00_FF
	_CICE_FPT_MASK_DATA            _CICE_FormatPartType = 0xFF_00

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

	_CICE_LF_NUMBER                _CICE_FormatPartType = 1
	_CICE_LF_SHORT_NORMAL          _CICE_FormatPartType = 2
	_CICE_LF_SHORT_UPPER_CASE      _CICE_FormatPartType = 3
	_CICE_LF_FULL_NORMAL           _CICE_FormatPartType = 4
	_CICE_LF_FULL_UPPER_CASE       _CICE_FormatPartType = 5

	// Common Integrator Console Encoder Time Format (CICE TF)
	// type constants.

	_CICE_TF_TIMESTAMP             _CICE_FormatPartType = 1
	_CICE_TF_ANSIC                 _CICE_FormatPartType = 2
	_CICE_TF_UNIXDATE              _CICE_FormatPartType = 3
	_CICE_TF_RUBYDATE              _CICE_FormatPartType = 4
	_CICE_TF_RFC822                _CICE_FormatPartType = 5
	_CICE_TF_RFC822_Z              _CICE_FormatPartType = 6
	_CICE_TF_RFC850                _CICE_FormatPartType = 7
	_CICE_TF_RFC1123               _CICE_FormatPartType = 8
	_CICE_TF_RFC1123_Z             _CICE_FormatPartType = 9
	_CICE_TF_RFC3339               _CICE_FormatPartType = 10

	// Common Integrator Console Encoder Caller Format (CICE CF)
	// type constants.

	_CICE_CF_TYPE_SEPARATOR        int16 = -1
	_CICE_CF_TYPE_FUNC_SHORT       int16 = 1
	_CICE_CF_TYPE_FUNC_FULL        int16 = 2
	_CICE_CF_TYPE_FILE_SHORT       int16 = 3
	_CICE_CF_TYPE_FILE_FULL        int16 = 4
	_CICE_CF_TYPE_LINE_NUM         int16 = 5
	_CICE_CF_TYPE_PKG_SHORT        int16 = 6 // unused
	_CICE_CF_TYPE_PKG_FULL         int16 = 7

	// Common Integrator Console Encoder (CICE) verb predefined constants.

	_CICE_VERB_START_INDICATOR     rune = '{'
	_CICE_VERB_END_INDICATOR       rune = '}'
	_CICE_VERB_SEPARATOR           byte = '/'

	// Common Integrator Console Integrator Standard Colors (CICE SC)
	// predefined constants.

	_CICE_SC_DEBUG                 string = `c/fg:ascii:36`
	_CICE_SC_INFO                  string = `c/fg:ascii:32`
	_CICE_SC_WARNING               string = `c/fg:ascii:33/b`
	_CICE_SC_ERROR                 string = `c/fg:ascii:31/b`
	_CICE_SC_FATAL                 string = `c/fg:ascii:35/b`

	// Common Integrator Console Encoder (CICE) defaults.

	_CICE_DEFAULT_FORMAT           string =
		// include colored level, colored time
		"{{c}}{{l}} {{t}}{{c/0}}\n" +
		// include message with \n if non-empty
		"{{m/?$\n}}" +
		// include fields with " = " as key-value separator, colored key
		"{{f/?$\n/v = /e, /l\t/le\t\t}}" +
		// include stacktrace with \n if non-empty
		"{{s/?$\n/e, }}" +
		// omit caller, specify each stacktrace's frame format
		"{{w/0/fd}}" +
		//
		"\n"

	_CICE_DEFAULT_TIME_FORMAT string = "Mon, Jan 02 15:04:05"
)

// Common Integrator Console Encoder Verb Types (CICE VT)
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
	// Make sure we won't break API by declaring package's console encoder
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

// SetFormat
func (ce *CI_ConsoleEncoder) SetFormat(newFormat string) *CI_ConsoleEncoder {

	if s := strings.TrimSpace(newFormat); s == "" {
		// Just check whether 'newFormat' is non-empty string or a string
		// that contains not only whitespace characters. But do not ignore them.
		return ce
	}

	// TODO: calculate buf size

	ce.format = newFormat
	return ce
}

// SetColorFor sets color what will be used as a replace for level-depended
// color verb from the 'format' string that is set by SetFormat() func.
func (ce *CI_ConsoleEncoder) SetColorFor(level Level, color string) *CI_ConsoleEncoder {

	if ce.colorMap == nil {
		ce.colorMap = make(map[Level]string)
	}

	if encodedColor := ce.rvColorHelper(color); encodedColor != "" {
		ce.colorMap[level] = encodedColor
		if l := len(encodedColor); ce.colorMapMax < l {
			ce.colorMapMax = l
		}
	}

	return ce
}

// FreezeAndGetEncoder builds current CI_ConsoleEncoder if it has not built yet
// and if format string has been provided by SetFormat() returning a function
// (has an alias CI_Encoder) that can be used at the
// CommonIntegrator.WithEncoder() call while initializing.
func (ce *CI_ConsoleEncoder) FreezeAndGetEncoder() CI_Encoder {
	return ce.doBuild().encode
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

	if ce.colorMap[LEVEL_WARNING] == "" {
		ce.colorMap[LEVEL_WARNING] = ce.rvColorHelper(_CICE_SC_WARNING)
	}

	if ce.colorMap[LEVEL_ERROR] == "" {
		ce.colorMap[LEVEL_ERROR] = ce.rvColorHelper(_CICE_SC_ERROR)
	}

	if ce.colorMap[LEVEL_FATAL] == "" {
		ce.colorMap[LEVEL_FATAL] = ce.rvColorHelper(_CICE_SC_FATAL)
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

	// what kind of verb we did parse and did verb has been closed correctly?
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

// rvJustText is a part of "resolve verb" functions.
//
func (ce *CI_ConsoleEncoder) rvJustText(text string) (predictedLen int) {

	if text != "" {
		ce.formatParts = append(ce.formatParts, _CICE_FormatPart{
			typ:   _CICE_FPT_VERB_JUST_TEXT,
			value: text,
		})
	}

	return len(text)
}

// rvLevel is a part of "resolve verb" functions.
//
func (ce *CI_ConsoleEncoder) rvLevel(verb string) (predictedLen int) {

	formattedLevel := _CICE_LF_FULL_NORMAL
	if idx := strings.IndexByte(verb, _CICE_VERB_SEPARATOR); idx != -1 {
		switch verb[idx+1:] {
		case "d", "D": formattedLevel = _CICE_LF_NUMBER
		case "s":      formattedLevel = _CICE_LF_SHORT_NORMAL
		case "S":      formattedLevel = _CICE_LF_SHORT_UPPER_CASE
		case "ss":     formattedLevel = _CICE_LF_FULL_NORMAL
		case "SS":     formattedLevel = _CICE_LF_FULL_UPPER_CASE
		}
	}

	ce.formatParts = append(ce.formatParts, _CICE_FormatPart{
		typ: _CICE_FPT_VERB_LEVEL | (formattedLevel << 8),
	})

	return predictedLen
}

// rvTime is a part of "resolve verb" functions.
//
func (ce *CI_ConsoleEncoder) rvTime(verb string) (predictedLen int) {

	format := _CICE_DEFAULT_TIME_FORMAT
	formattedTime := _CICE_FormatPartType(0)

	(*CI_ConsoleEncoder)(nil).rvHelper(verb, func(verbPart string) (continue_ bool) {
		if verbPart = strings.TrimSpace(format); verbPart != "" {
			switch predefined := strings.ToUpper(verbPart); predefined {
			case "UNIX", "TIMESTAMP":     formattedTime = _CICE_TF_TIMESTAMP
			case "ANSIC":                 formattedTime = _CICE_TF_ANSIC
			case "UNIXDATE", "UNIX_DATE": formattedTime = _CICE_TF_UNIXDATE
			case "RUBYDATE", "RUBY_DATE": formattedTime = _CICE_TF_RUBYDATE
			case "RFC822":                formattedTime = _CICE_TF_RFC822
			case "RFC822Z":               formattedTime = _CICE_TF_RFC822_Z
			case "RFC850":                formattedTime = _CICE_TF_RFC850
			case "RFC1123":               formattedTime = _CICE_TF_RFC1123
			case "RFC1123Z":              formattedTime = _CICE_TF_RFC1123_Z
			case "RFC3339":               formattedTime = _CICE_TF_RFC3339
			default:                      format = verbPart
			}
		}
		return false // only first time verb is allowed and will be parsed
	})

	ce.formatParts = append(ce.formatParts, _CICE_FormatPart{
		typ:   _CICE_FPT_VERB_TIME | (formattedTime << 8),
		value: format,
	})

	return len(format) + 10 // stock for some weekdays
}

//
// ""
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

//
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
		case strings.HasPrefix(verbPart, "?^"): ce.bf.beforeBody = verbPart[2:]
		case strings.HasPrefix(verbPart, "?$"): ce.bf.afterBody = verbPart[2:]
		default:                                return false
		}
		return true
	})

	ce.formatParts = append(ce.formatParts, _CICE_FormatPart{
		typ: _CICE_FPT_VERB_BODY,
	})

	return 256 + len(ce.bf.beforeBody) + len(ce.bf.afterBody)
}

//
func (ce *CI_ConsoleEncoder) rvCaller(verb string) (predictedLen int) {

	isAdd := true
	formatPrefixes := []string{"f", "F"}

	(*CI_ConsoleEncoder)(nil).rvHelper(verb, func(verbPart string) (continue_ bool) {
		switch {
		case verbPart == "0":               isAdd = false
		case hpm(verbPart, formatPrefixes): predictedLen += ce.rvCallerFormat(verbPart[1:])
		default:                            return false
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

//
func (ce *CI_ConsoleEncoder) rvCallerFormat(f string) (predictedLen int) {

	if f == "" || f == "d" || f == "D" {
		ce.cf.isDefault = true
		return 256
	}

	j := 0 // index of ce.cf.parts
	for _, fc := range f {

		t := _CICE_CF_TYPE_SEPARATOR // by default threat it as a separator
		switch fc {

		case 'w':      t = _CICE_CF_TYPE_FUNC_SHORT
		case 'W':      t = _CICE_CF_TYPE_FUNC_FULL
		case 'f':      t = _CICE_CF_TYPE_FILE_SHORT
		case 'F':      t = _CICE_CF_TYPE_FILE_FULL
		case 'l', 'L': t = _CICE_CF_TYPE_LINE_NUM
		case 'p', 'P': t = _CICE_CF_TYPE_PKG_FULL

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
// - If at least one Field presented:
//   - "?^<text>": <text> will be write before any Field is written at the runtime.
//   - "?$<text>": <text> will be appended to the end of last Field at the runtime.
//   - "k<text>": <text> will be written before Field's keys is written.
//   - "v<text>": <text> will be written before Field's value is written.
//   - "e<text>": <text> will be written after Field's value excluding last.
//   - "l<text>": <text> will be written at the each new line of fields' part set.
//   - "*<int>": <int> is how much fields are placed at the one line
//     (by default: 4. Use <= 0 value to place all fields at the one line).
func (ce *CI_ConsoleEncoder) rvFields(verb string) (predictedLen int) {

	ce.ff.itemsPerLine = 4

	(*CI_ConsoleEncoder)(nil).rvHelper(verb, func(verbPart string) (continue_ bool) {
		switch upperCased := strings.ToUpper(verbPart); {

		case strings.HasPrefix(verbPart, "?^"):   ce.ff.beforeFields = verbPart[2:]
		case strings.HasPrefix(verbPart, "?$"):   ce.ff.afterFields = verbPart[2:]
		case strings.HasPrefix(upperCased, "LE"): ce.ff.afterNewLineForError = verbPart[2:]
		case upperCased[0] == 'L':                ce.ff.afterNewLine = verbPart[1:]
		case upperCased[0] == 'K':                ce.ff.beforeKey = verbPart[1:]
		case upperCased[0] == 'V':                ce.ff.afterKey = verbPart[1:]
		case upperCased[0] == 'E':                ce.ff.afterValue = verbPart[1:]

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

// rvStacktrace is a part of "resolve verb" functions.
func (ce *CI_ConsoleEncoder) rvStacktrace(verb string) (predictedLen int) {

	(*CI_ConsoleEncoder)(nil).rvHelper(verb, func(verbPart string) (continue_ bool) {
		switch {
		case strings.HasPrefix(verbPart, "?^"): ce.sf.beforeStack = verbPart[2:]
		case strings.HasPrefix(verbPart, "?$"): ce.sf.afterStack = verbPart[2:]
		default:                                return false
		}
		return true
	})

	ce.formatParts = append(ce.formatParts, _CICE_FormatPart{
		typ: _CICE_FPT_VERB_STACKTRACE,
	})

	return 2048
}

//
func (ce *CI_ConsoleEncoder) encode(e *Entry) []byte {

	// TODO: Reuse allocated buffers

	buf := make([]byte, 0, ce.minimumBufferLen)
	allowEmpty := e.LogLetter.Items.Flags.TestAll(FLAG_INTEGRATOR_IGNORE_EMPTY_PARTS)

	for _, part := range ce.formatParts {
		switch part.typ.Type() {

		case _CICE_FPT_VERB_JUST_TEXT:       buf = ce.encodeJustText(buf, part)
		case _CICE_FPT_VERB_COLOR_CUSTOM:    buf = ce.encodeColor(buf, part)
		case _CICE_FPT_VERB_COLOR_FOR_LEVEL: buf = ce.encodeColorForLevel(buf, e)
		case _CICE_FPT_VERB_BODY:            buf = ce.encodeBody(buf, e)
		case _CICE_FPT_VERB_TIME:            buf = ce.encodeTime(e, part, buf)
		case _CICE_FPT_VERB_LEVEL:           buf = ce.encodeLevel(buf, part, e)
		case _CICE_FPT_VERB_STACKTRACE:      buf = ce.encodeStacktrace(buf, e, allowEmpty)
		case _CICE_FPT_VERB_CALLER:          buf = ce.encodeCaller(buf, e)

		case _CICE_FPT_VERB_FIELDS:
			buf = ce.encodeFields(buf, e.LogLetter.SystemFields, allowEmpty, false)
			if e.ErrLetter != nil {
				buf = ce.encodeFields(buf, e.ErrLetter.SystemFields, allowEmpty, false)
			}
			buf = ce.encodeFields(buf, e.LogLetter.Items.Fields, allowEmpty, false)
		}
	}

	return buf
}

// easy case because fp.value is the text we should add.
func (ce *CI_ConsoleEncoder) encodeJustText(to []byte, fp _CICE_FormatPart) []byte {
	return bufw(to, fp.value)
}

//
func (ce *CI_ConsoleEncoder) encodeLevel(to []byte, fp _CICE_FormatPart, e *Entry) []byte {

	formattedLevel := ""
	switch fp.typ.Data() {

	case _CICE_LF_NUMBER:           formattedLevel = strconv.Itoa(int(e.Level))
	case _CICE_LF_SHORT_NORMAL:     formattedLevel = e.Level.String()[:3]
	case _CICE_LF_SHORT_UPPER_CASE: formattedLevel = strings.ToUpper(e.Level.String()[:3])
	case _CICE_LF_FULL_NORMAL:      formattedLevel = e.Level.String()
	case _CICE_LF_FULL_UPPER_CASE:  formattedLevel = strings.ToUpper(e.Level.String())
	}

	return bufw(to, formattedLevel)
}

//
func (ce *CI_ConsoleEncoder) encodeTime(e *Entry, fp _CICE_FormatPart, to []byte) []byte {

	formattedTime := ""

	switch fp.typ.Data() {
	case _CICE_TF_TIMESTAMP: formattedTime = strconv.FormatInt(e.Time.Unix(), 10)
	case _CICE_TF_ANSIC:     formattedTime = e.Time.Format(time.ANSIC)
	case _CICE_TF_UNIXDATE:  formattedTime = e.Time.Format(time.UnixDate)
	case _CICE_TF_RUBYDATE:  formattedTime = e.Time.Format(time.RubyDate)
	case _CICE_TF_RFC822:    formattedTime = e.Time.Format(time.RFC822)
	case _CICE_TF_RFC822_Z:  formattedTime = e.Time.Format(time.RFC822Z)
	case _CICE_TF_RFC850:    formattedTime = e.Time.Format(time.RFC850)
	case _CICE_TF_RFC1123:   formattedTime = e.Time.Format(time.RFC1123)
	case _CICE_TF_RFC1123_Z: formattedTime = e.Time.Format(time.RFC1123Z)
	case _CICE_TF_RFC3339:   formattedTime = e.Time.Format(time.RFC3339)
	default:                 formattedTime = e.Time.Format(fp.value)
	}

	return bufw(to, formattedTime)
}

// easy case because ASCII sequence already generated at the rvColor method.
func (ce *CI_ConsoleEncoder) encodeColor(to []byte, fp _CICE_FormatPart) []byte {
	return bufw(to, fp.value)
}

//
func (ce *CI_ConsoleEncoder) encodeColorForLevel(to []byte, e *Entry) []byte {

	if color := ce.colorMap[e.Level]; color != "" {
		return bufw(to, color)
	}

	// TODO: mocked (find the closes log level and use that's color)
	return to
}

// easy case because e.Message is the text we should add.
func (ce *CI_ConsoleEncoder) encodeBody(to []byte, e *Entry) []byte {

	body := e.LogLetter.Items.Message
	//if body == "" && e.ErrLetter != nil {
	//	// TODO: It's guaranteed that first err letter's item has non-empty
	//	//  and marked message (it's what error has been constructed w/)
	//	//  but anyway the better way is to find first marked letter item w/ message
	//	body = e.ErrLetter.Items.Message
	//}

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

	return ce.encodeStackFrame(to, frame, nil, false)
}

//
func (ce *CI_ConsoleEncoder) encodeFields(

	to []byte,
	fields []field.Field,
	allowEmpty,
	isErrors bool,

) []byte {

	if len(fields) == 0 {
		return to
	}

	if !isErrors && ce.ff.beforeFields != "" {
		to = bufw(to, ce.ff.beforeFields)
	}

	newLine := ce.ff.afterNewLine
	if isErrors {
		newLine = ce.ff.afterNewLineForError
	}
	if newLine != "" {
		to = bufw(to, newLine)
	}

	unnamedFieldIdx := 1
	writtenFieldIdx := int16(0)
	for i, n := int16(0), int16(len(fields)); i < n; i++ {

		// write new line and new line title
		if ce.ff.itemsPerLine > 0 &&
			writtenFieldIdx != 0 &&
			writtenFieldIdx%ce.ff.itemsPerLine == 0 {

			to = bufw(to, "\n")
			newLine := ce.ff.afterNewLine
			if isErrors {
				newLine = ce.ff.afterNewLineForError
			}
			if newLine != "" {
				to = bufw(to, newLine)
			}
		}

		// Despite the fact, that letter.ParseTo() already contains IsZero() call
		// it's only vary fields filtering (fields with the "?" at the end of name).
		//
		// So, there may be a non-vary zero fields that may be not allowed.
		// Some user's system fields (started with "sys." as their names) also
		// must not be added.
		//
		// allowEmpty must be first at the SCE, because we don't need IsZero() call
		// (it's kinda heavy) if empty fields are allowed.
		// https://en.wikipedia.org/wiki/Short-circuit_evaluation
		if (!allowEmpty && fields[i].IsZero()) || strings.HasPrefix(fields[i].Key, "sys.") {
			continue
		} else {
			writtenFieldIdx++
		}

		// before key (key identifier)
		if ce.ff.beforeKey != "" {
			to = bufw(to, ce.ff.beforeKey)
		}

		// write key
		if fields[i].Key != "" {
			to = bufw(to, fields[i].Key)
		} else {
			to = bufw(to, letter.UnnamedAsStr(unnamedFieldIdx))
			unnamedFieldIdx++
		}

		// before value (value identifier)
		if ce.ff.afterKey != "" {
			to = bufw(to, ce.ff.afterKey)
		}

		// write value

		// ----- SYSTEM FIELDS -----

		if fields[i].Kind.IsSystem() {
			switch fields[i].Kind.BaseType() {

			case field.KIND_SYS_TYPE_EKAERR_UUID, field.KIND_SYS_TYPE_EKAERR_CLASS_NAME,
			field.KIND_SYS_TYPE_EKAERR_PUBLIC_MESSAGE:
				to = bufw(to, `"`)
				to = bufw(to, fields[i].SValue)
				to = bufw(to, `"`)

			case field.KIND_SYS_TYPE_EKAERR_CLASS_ID:
				to = bufw(to, strconv.FormatInt(fields[i].IValue, 10))

			default:
				to = bufw(to, `"<unsupported system field>"`)
			}
			goto END_FIELD_PROCESSING
		}

		// ----- NIL FIELDS -----

		if fields[i].Kind.IsNil() {
			to = bufw(to, "null")
			goto END_FIELD_PROCESSING
		}

		// ----- ARRAY FIELDS -----
		// todo

		// ----- BASE TYPE FIELDS -----

		switch fields[i].Kind.BaseType() {

		case field.KIND_TYPE_BOOL:
			if fields[i].IValue != 0 {
				to = bufw(to, "true")
			} else {
				to = bufw(to, "false")
			}

		case field.KIND_TYPE_INT,
		field.KIND_TYPE_INT_8, field.KIND_TYPE_INT_16,
		field.KIND_TYPE_INT_32, field.KIND_TYPE_INT_64:
			to = bufw(to, strconv.FormatInt(fields[i].IValue, 10))

		case field.KIND_TYPE_UINT,
		field.KIND_TYPE_UINT_8, field.KIND_TYPE_UINT_16,
		field.KIND_TYPE_UINT_32, field.KIND_TYPE_UINT_64:
			to = bufw(to, strconv.FormatUint(uint64(fields[i].IValue), 10))

		case field.KIND_TYPE_FLOAT_32:
			f := float64(math.Float32frombits(uint32(fields[i].IValue)))
			to = bufw(to, strconv.FormatFloat(f, 'f', 2, 32))

		case field.KIND_TYPE_FLOAT_64:
			f := float64(math.Float32frombits(uint32(fields[i].IValue)))
			to = bufw(to, strconv.FormatFloat(f, 'f', 2, 64))

		case field.KIND_TYPE_STRING:
			to = bufw(to, `"`)
			to = bufw(to, fields[i].SValue)
			to = bufw(to, `"`)

		default:
		}

	END_FIELD_PROCESSING:

		// write after value
		if ce.ff.afterValue != "" {
			to = bufw(to, ce.ff.afterValue)
		}

	} // end loop of fields

	// remove last after value
	if ce.ff.afterValue != "" {
		to = to[:len(to)-len(ce.ff.afterValue)]
	}

	if !isErrors && ce.ff.afterFields != "" {
		to = bufw(to, ce.ff.afterFields)
	}

	return to
}

// encodeStacktrace
func (ce *CI_ConsoleEncoder) encodeStacktrace(to []byte, e *Entry, allowEmpty bool) []byte {

	stacktrace := e.LogLetter.StackTrace
	if len(stacktrace) == 0 && e.ErrLetter != nil {
		stacktrace = e.ErrLetter.StackTrace
	}

	lStacktrace := int16(len(stacktrace))
	if lStacktrace == 0 {
		return to
	}

	if ce.sf.beforeStack != "" {
		to = bufw(to, ce.sf.beforeStack)
	}

	letterItem := (*letter.LetterItem)(nil)
	letterItemIdx := int16(0)
	if e.ErrLetter != nil {
		letterItem = e.ErrLetter.Items
		letterItemIdx = letterItem.StackFrameIdx()
	}

	for i := int16(0); i < lStacktrace; i++ {
		letterItemPassed := (*letter.LetterItem)(nil)
		if letterItem != nil && letterItemIdx == i {
			letterItemPassed = letterItem
			letterItem = letterItem.Next()
			letterItemIdx = letterItem.StackFrameIdx()
		}
		to = ce.encodeStackFrame(to, stacktrace[i], letterItemPassed, allowEmpty)
		if i < lStacktrace-1 {
			to = bufw(to, "\n")
		}
	}

	if ce.sf.afterStack != "" {
		to = bufw(to, ce.sf.afterStack)
	}

	return to
}

//
func (ce *CI_ConsoleEncoder) encodeStackFrame(

	to []byte,
	frame ekasys.StackFrame,
	letterItem *letter.LetterItem,
	allowEmpty bool,

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

	if letterItem != nil {
		to = bufw(to, "\n")

		if ce.ff.afterNewLineForError != "" {
			to = bufw(to, ce.ff.afterNewLineForError)
		}

		if letterItem.Flags.TestAll(letter.FLAG_MARKED_LETTER_ITEM) {
			to = bufw(to, `(*) `)
		}

		to = bufw(to, letterItem.Message)
		to = bufw(to, "\n")

		lToBak := len(to)
		to = ce.encodeFields(to, letterItem.Fields, allowEmpty, true)

		// ce.encodeFields may write no fields. Then we must clear last "\n"
		if len(to) == lToBak {
			to = to[:len(to)-1]
		}
	}

	return to
}
