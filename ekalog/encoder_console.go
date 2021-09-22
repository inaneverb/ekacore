// Copyright © 2020-2021. All rights reserved.
// Author: Ilya Stroy.
// Contacts: iyuryevich@pm.me, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekalog

import (
	"io"
	"strings"

	"github.com/qioalice/ekago/v3/internal/ekaletter"
)

//noinspection GoSnakeCaseUsage
type (
	// CI_ConsoleEncoder is a type that built to be used as a part of CommonIntegrator
	// as an log Entry encoder to the TTY (console) with human readable format,
	// custom format supporting and ability to enable coloring.
	//
	// If you want to use CI_ConsoleEncoder, you need to instantiate object,
	// set log's text format string (if you dont want to use default one)
	// using SetFormat() and that is.
	// The last thing you need to do is to register CI_ConsoleEncoder with
	// CommonIntegrator using CommonIntegrator.WithEncoder().
	//
	// See SetFormat() docs to figure out how you can manage the log format.
	//
	// CI_ConsoleEncoder may look like heavy rock and it's only half true.
	// It's really flexible customizable text encoder, but the format string parsing
	// is performed only once. At the registration of CI_ConsoleEncoder.
	// Encode operations will be performed over split and parsed format string,
	// and it's as blazing fast as it's even possible.
	//
	// You MUST NOT to call EncodeEntry() method manually.
	// It is used by associated CommonIntegrator and it WILL lead to UB
	// if you will try to use it manually. May even panic.
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

		preEncodedFields        []byte
		preEncodedFieldsWritten int16
	}
)

var (
	// Make sure we won't break API.
	_ CI_Encoder = (*CI_ConsoleEncoder)(nil)
)

// SetFormat allows you to set format string that will be used as a template
// to generate and represent all log Entry objects.
//
// Calling this method many times will overwrite previous value of format string.
//
// The format string will be parsed and applied only after CI_ConsoleEncoder
// is registered with CommonIntegrator using CommonIntegrator.WithEncoder() method.
// After it is done, the format string cannot be changed and this method has no-op.
//
// -----
//
// Attached ekaerr.Error.
//
// You know that you may log ekaerr.Error objects using special log finishers,
// like Logger.Errore(), Logger.Warne(), etc.
// According the main rule and purpose of ekaerr packages,
// the attached fields are associated with its stackframe.
// So, fields of a some stackframe will be placed exactly near related stackframe.
// You can't write them all at the one place.
// Read more below about verbs to recognize how you can manage the output.
//
// -----
//
// Verbs. Introduction.
//
// The format string has its own verbs and verbs' parameters using which
// you can manage the place and view of Entry's parts.
//
// All verbs has the following format:
// "{{<verb_name>/<verb_parameter_1>/.../<verb_parameter_N>}}"
//
// You should build your own format string using verbs, separators, text, etc.
//
// For each verb there's verb name and its short aliases, so you can use any of that.
// Each verb has its own parameters and its defaults.
// You can overwrite one, two, N or all default verb parameters.
//
// The order of verb parameters doesn't matter.
// Letter case of verb names or its parameter names doesn't matter.
//
// Below you can see the verbs and its parameters.
// Keep in mind, that each verb's parameter has its own key,
// that determines what this verb parameter is.
//
// Verbs. List.
//
// 0. Empty verb / incorrect verb.
//    All text that is not verb or an incorrect verb will be used as is.
//    It will be a part of final log message.
//
// 1. Entry.Level verb.
//    Names: "level", "lvl", "l".
//
//    The verb will be replaced by Entry.Level string representation.
//    By default it writes a full Level's capitalized name.
//    The same as is returned by Level.String().
//
//    Parameters:
//    (Only first parameter is supported. 2nd and next will be ignored.
//    Incorrect parameter will be ignored.)
//
//    - "d", "D": Write Level's number instead of string.
//      From "0" for LEVEL_EMERGENCY down to "7" for LEVEL_DEBUG.
//    - "s": Write Level's short string variant.
//      The same as is returned by Level.String3().
//      Examples: "Emerg", "Ale", "Cri", "Err", "War", "Noe", "Inf", "Deb".
//    - "S": Write Level's short upper-cased string variant.
//      The same as is returned by Level.ToUpper3().
//      Examples: "EMERG", "ALE", "CRI", "ERR", "WAR", "NOE", "INF", "DEB".
//    - "ss": Write Level's full string variant.
//      This is default verb parameter. The same as is returned by Level.String().
//    - "SS": Write Level's full upper-cased string variant.
//      The same as is returned by Level.ToUpper().
//      Examples: "EMERGENCY", "ALERT", "CRITICAL", "WARNING", etc.
//
// 2. Entry.Time verb.
//    Names: "time", "t".
//
//    The verb will be replaced by Entry.Time string representation.
//    Default time format is: "Mon, Jan 02 15:04:05".
//
//    Parameters:
//    (Only first parameter is supported. 2nd and next will be ignored.
//    Incorrect parameter will be ignored.)
//
//    - "UNIX", "TIMESTAMP": Write time.Time as unix timestamp in seconds.
//      Example: 1619549194 for Tue Apr 27 2021 18:46:34 GMT+0000
//    - "ANSIC": Write time.Time as time.ANSIC format.
//      Example: "Mon Jan _2 15:04:05 2006".
//    - "UNIXDATE", "UNIX_DATE": Write time.Time as time.UnixDate format.
//      Example: "Mon Jan _2 15:04:05 MST 2006".
//    - "RUBYDATE", "RUBY_DATE": Write time.Time as time.RubyDate format.
//      Example: "Mon Jan 02 15:04:05 -0700 2006".
//    - "RFC822": Write time.Time as time.RFC822 format.
//      Example: "02 Jan 06 15:04 MST".
//    - "RFC822Z": Write time.Time as time.RFC822Z format.
//      Example: "02 Jan 06 15:04 -0700".
//    - "RFC850": Write time.Time as time.RFC850 format.
//      Example: "Monday, 02-Jan-06 15:04:05 MST".
//    - "RFC1123": Write time.Time as time.RFC1123 format.
//      Example: "Mon, 02 Jan 2006 15:04:05 MST".
//    - "RFC1123Z": Write time.Time as time.RFC1123Z format.
//      Example: "Mon, 02 Jan 2006 15:04:05 -0700".
//    - "RFC3339": Write time.Time as time.RFC3339 format.
//      Example: "2006-01-02T15:04:05Z07:00".
//    - "<your_own_time_format>: Uses string as time format.
//      time.Time.Format() will be called with that format string.
//
// 3. Entry's log body verb.
//    Names: "message", "body", "m", "b".
//
//    The verb will be replaced by Entry's log message.
//    You can include this verb only once.
//    2nd and next these verbs will be treated as just a text.
//
//    Parameters:
//    (The same key parameters will be overwritten by the last one.
//    The parsing of parameters will be stopped when an incorrect parameter found).
//
//    - "?^<text>": Places <text> before log's body if body is not empty.
//    - "?$<text>": Places <text> after log's body if body is not empty.
//
// 4. Entry's caller verb.
//    Names: "caller", "who", "w".
//
//    The verb will be replaced by caller's info.
//    It may include a package, function, line number of the caller.
//    The caller is the function that calls a log finisher.
//
//    Parameters:
//    (The same key parameters will be overwritten by the last one.
//    The parsing of parameters will be stopped when an incorrect parameter found).
//
//    - "0": Do not include caller's info.
//      It's meaningless by the first view, but the main purpose of that verb is
//      the fact that you can specify format using "f<format>" verb parameter,
//      that also will be used as a stacktrace's format.
//      So, if you want to specify format for stacktrace,
//      but don't want to include caller's info you need this verb's parameter.
//    - "f<format>", "F<format>": You can specify format of caller string,
//      that will be used as template to generate both of caller's info and stacktrace.
//      <format> string must be combined by your own way with:
//      - "w": Short function name. Only function, without package path.
//      - "W": Full function name. Includes package path.
//      - "f": Short filename. Only filename, without full path to that file.
//      - "F": Full filename. Includes full path to that file.
//      - "l", "L": File's line number.
//      - "p", "P": Full package path.
//      - <any_other>: Writes as is. Useful to split format's parts.
//
//    Example: "{{w/0/fF:L}}" is a caller's verb, that specifies format
//    "<full_file_name>:<line_number>" for both of caller and stacktrace,
//    but caller info won't be included to the log's output.
//
// 5. Entry's stacktrace verb.
//    Names: "stacktrace", "s".
//
//    The verb will be replaced by stacktrace, if it's presented.
//    Each stack frame may include its package, function, line number.
//
//    WARNING.
//    If you won't add this verb, the fields of your attached ekaerr.Error
//    won't be encoded, because they are heavily linked with stacktrace!
//
//    Parameters:
//    (The same key parameters will be overwritten by the last one.
//    The parsing of parameters will be stopped when an incorrect parameter found).
//
//    - "?^<text>": Places <text> before stacktrace if stacktrace is presented.
//    - "?$<text>": Places <text> after stacktrace if stacktrace is presented.
//
// 6. Entry's fields verb.
//    Names: "fields", "f".
//
//    The verb will be replaced by log Entry's fields.
//    Keep in mind, this verb has an affect only for log's fields.
//    If you have attached ekaerr.Error, its fields will be printed if
//    stacktrace's verb is included. But this verb has an affect of format
//    to that fields also.
//
//    Parameters:
//    (The same key parameters will be overwritten by the last one.
//    The parsing of parameters will be stopped when an incorrect parameter found).
//
//    - "?^<text>": Places <text> before FIRST field (in stackframe).
//    - "?$<text>": Places <text> after LAST field (in stackframe).
//    - "k<text>": Places <text> before EACH field's key.
//    - "v<text>": Places <text> before EACH field's value.
//    - "e<text>": Places <text> after EACH field's value (last field excluded).
//    - "l<text>": Places <text> at the each new line of log Entry's fields part.
//    - "le<text>": Places <text> at the each new line
//      of attached ekaerr.Error fields part.
//    - "*<number>": <number> is how much fields are placed at the one line.
//      (By default: 4. Use <= 0 value to place all fields at the one line).
//
// 7. TTY coloring verb.
//    Names: "color", "c".
//
//    This verb allows you to make some parts of final messages be colored.
//    Color in a TTY terms is just a special escape sequence.
//    Think about this verb as a HTML tag.
//
//    First of all, if you're placing this verb of starting coloring,
//    you need to place the verb of ending coloring.
//    Otherwise all next data will be colored too until next color verb is found.
//    It's not about CI_ConsoleEncoder but about how TTY coloring works.
//
//    Next thing you need to know is
//    There's "default" colors for each log's Level.
//    If you don't present a color using verb's parameters,
//    the default log's level color will be used.
//    You can overwrite default colors for level using SetColorFor() method.
//
//    How parameters works.
//    All parameters are split to its groups:
//    - Bold (enable/disable),
//    - Italic (enable/disable),
//    - Underline (enable/disable),
//    - Background color (overwrite/disable),
//    - Foreground color (overwrite/disable),
//    - Cancel any color.
//    You may specify more than 1 parameter to combine them.
//    If you specify more then 1 parameter for the same group, the last one will be used.
//    If you specify "cancel any color" parameter, all others parameters will be ignored.
//    If any parameter is invalid, all others will be ignored.
//
//    To specify color, you can use:
//    - ASCII colors: Use "ascii<color>" or "ascii(<color>)" syntax.
//      Allowable colors: [30..37]+[90..97] for foreground,
//      [40..47]+[100..107] for background.
//      Read more: https://en.wikipedia.org/wiki/ANSI_escape_code
//      It's UB to specify foreground color for background and vice-versa.
//      (Most likely ASCII affiliation of color will be used).
//    - X256 colors: Use "<color>" syntax. Allowable colors: [0..255].
//      Read more using the link above.
//    - HEX colors: Use "#<color>" format. Will be transformed to X256 colors.
//      Read more: https://en.wikipedia.org/wiki/Web_colors
//    - RGB colors: Use "rgb:<red>,<green>,<blue>" or "rgb(<red>,<green>,<blue>)"
//      or "rgb,<red>,<green>,<blue>" syntax.
//      All of <red>, <green>, <blue> must be in range [0..255].
//      Will be transformed to X256 colors.
//      Read more using link above.
//    - "-1": Disable coloring for desired type (background/foreground).
//
//    Parameters:
//
//    - No parameters: Uses log's Level default color.
//    - "0": Disables all TTY font effects: color, bold, italic, underline.
//    - "bold", "b": Use bold font.
//    - "nobold", "nob": Disable bold font, return back to normal.
//    - "italic", "i": Use italic font.
//    - "noitalic", "noi": Disable italic font, return back to normal.
//    - "underline", "u": Use underline font.
//    - "nounderline", "nou": Disable underline font, return back to normal.
//    - "fg:<color>: Set <color> for text's foreground. See above about colors.
//    - "bg:<color>": Set <color> for text's background. See above about colors.
//
//   Reminder.
//   Make sure your terminal supports X256 colors or use ASCII colors otherwise.
//   If your terminal doesn't support X256 colors and you will try to use it,
//   you may get an ugly escape sequences in your output.
//
//   Dropping colors for specific io.Writer.
//   You may want to disable coloring for specific io.Writer leaving it for another.
//   See CICE_DropColors() for more details.
//
// -----
//
// If you won't set any format string, the default one will be used.
// It is:
//
//   "{{c}}{{l}} {{t}}{{c/0}}\n" +      // include colored level, colored time
//   "{{m/?$\n}}" +                     // include message with \n if non-empty
//   "{{f/?$\n/v = /e, /l\t/le\t\t}}" + // include fields with " = " as key-value separator
//   "{{s/?$\n/e, }}" +                 // include stacktrace with \n if non-empty
//   "{{w/0/fd}}" +                     // omit caller, specify each stacktrace's frame format
//   "\n"
//
func (ce *CI_ConsoleEncoder) SetFormat(newFormat string) *CI_ConsoleEncoder {

	if s := strings.TrimSpace(newFormat); s == "" {
		// Just check whether 'newFormat' is non-empty string or a string
		// that contains not only whitespace characters. But do not ignore them.
		return ce
	}

	ce.format = newFormat
	return ce
}

// SetColorFor sets color what will be used as a replace for level-depended
// color verb from the 'format' string that is set by SetFormat() func
//
// Example: "c/fg:ascii:31/b", "c/fg:#123456/bg:rgb,50,50,50/i/u".
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

// PreEncodeField allows you to pre-encode some ekaletter.LetterField,
// that is must be used with EACH Entry that will be encoded using this CI_ConsoleEncoder.
//
// It's useful when you want some unchanged runtime data for each log message,
// like git hash commit, version, etc. Or if you want to create many loggers
// attach some different fields to them and log different logs using them.
//
// Unnamed fields are not allowed.
//
// WARNING!
// PreEncodeField() MUST BE USED ONLY IF CI_ConsoleEncoder HAS BEEN REGISTERED
// WITH SOME CommonIntegrator ALREADY. UB OTHERWISE, MAY PANIC!
func (ce *CI_ConsoleEncoder) PreEncodeField(f ekaletter.LetterField) {

	// Avoid calls of PreEncodeField() when CI_ConsoleEncoder has not built yet.
	if f.Key == "" || len(ce.formatParts) == 0 || f.IsInvalid() ||
		f.RemoveVary() && f.IsZero() {
		return
	}

	preEncodedFieldsLenBak := len(ce.preEncodedFields)
	ce.preEncodedFields = ce.encodeField(ce.preEncodedFields, f, false, ce.preEncodedFieldsWritten)

	if len(ce.preEncodedFields) != preEncodedFieldsLenBak {
		ce.preEncodedFieldsWritten++
	}
}

// EncodeEntry encodes passed Entry as text using provided (and parsed)
// format string that is set by SetFormat() method, returning a RAW encoded data.
//
// It may includes a special ASCII escape sequences to color output.
// It's enabled by default (if you didn't change format string).
//
// EncodeEntry is for internal purposes only and MUST NOT be called directly.
// UB otherwise, may panic.
func (ce *CI_ConsoleEncoder) EncodeEntry(e *Entry) []byte {

	// TODO: Reuse allocated buffers

	to := make([]byte, 0, ce.minimumBufferLen)

	// Use last ekaerr.Error's message as Entry's one if it's empty.
	if e.ErrLetter != nil {
		if l := len(e.ErrLetter.Messages); l > 0 && e.LogLetter.Messages[0].Body == "" {
			e.LogLetter.Messages[0].Body = e.ErrLetter.Messages[l-1].Body
			e.ErrLetter.Messages[l-1].Body = ""
		}
	}

	for _, part := range ce.formatParts {
		switch part.typ.Type() {

		case _CICE_FPT_VERB_JUST_TEXT:
			to = ce.encodeJustText(to, part)
		case _CICE_FPT_VERB_COLOR_CUSTOM:
			to = ce.encodeColor(to, part)
		case _CICE_FPT_VERB_COLOR_FOR_LEVEL:
			to = ce.encodeColorForLevel(to, e)
		case _CICE_FPT_VERB_BODY:
			to = ce.encodeBody(to, e)
		case _CICE_FPT_VERB_TIME:
			to = ce.encodeTime(e, part, to)
		case _CICE_FPT_VERB_LEVEL:
			to = ce.encodeLevel(to, part, e)
		case _CICE_FPT_VERB_STACKTRACE:
			to = ce.encodeStacktrace(to, e)
		case _CICE_FPT_VERB_CALLER:
			to = ce.encodeCaller(to, e)

		case _CICE_FPT_VERB_FIELDS:
			errLetterSystemFields := []ekaletter.LetterField(nil)
			if e.ErrLetter != nil {
				errLetterSystemFields = e.ErrLetter.SystemFields
			}
			to = ce.encodeFields(to, e.LogLetter.SystemFields, errLetterSystemFields, false, false)

			// Handle special case when ekaerr.Error's ekaletter.Letter has a fields
			// but has no stacktrace. It means that lightweight error has been created.
			lightweightErrorFields := []ekaletter.LetterField(nil)
			if e.ErrLetter != nil && len(e.ErrLetter.StackTrace) == 0 && len(e.ErrLetter.Fields) > 0 {
				lightweightErrorFields = e.ErrLetter.Fields
			}
			to = ce.encodeFields(to, e.LogLetter.Fields, lightweightErrorFields, false, true)
		}
	}

	// Restore ekaerr.Error's last message that was used as Entry's message.
	if e.ErrLetter != nil {
		if l := len(e.ErrLetter.Messages); l > 0 && e.ErrLetter.Messages[l-1].Body == "" {
			e.ErrLetter.Messages[l-1].Body = e.LogLetter.Messages[0].Body
			e.LogLetter.Messages[0].Body = ""
		}
	}

	return to
}

// CICE_DropColors returns an io.Writer, you can write a data with TTY colors to,
// that will rewrite everything but colors to the your dest io.Writer.
//
// You can pass returned io.Writer to the CommonIntegrator.WriteTo()
// method associating it with the CI_ConsoleEncoder that writes a log messages
// with a TTY colors. Returned io.Writer will rewrite message to your io.Writer
// but without these shell color sequences.
//
// It's useful when you want to write colored log messages to TTY but
// raw to the files. Or any other destination.
//
// Keep in mind, it's not "free" operation. It allocates RAM buffer,
// writes raw data to, parses it, droppings colors and rewrites cleared raw data.
func CICE_DropColors(dest io.Writer) io.Writer {
	return &_CICE_DropColors{dest: dest}
}
