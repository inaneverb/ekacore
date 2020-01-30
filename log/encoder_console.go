// Copyright © 2019. All rights reserved.
// Author: Ilya Yuryevich.
// Contacts: qioalice@gmail.com, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package log

//import (
//	"unsafe"
//
//	"github.com/qioalice/gext/dangerous"
//)
//
////
//type ConsoleEncoder struct {
//
//	// -----
//	// RAW FORMAT
//	// This is the ConsoleEncoder's parts that are considered raw.
//	// These parts are set directly by some methods or with minimum aggregation.
//	// They will be parsed, converted and casted to internal structures
//	// that are more convenient to be used as source of log message format
//	// generator.
//	// -----
//	format string
//
//	// -----
//	// BUILT FORMAT
//	// This part represents parsed raw format string.
//	// First of all RAW format string parsed to the two group of entities:
//	// - Just string, not verb (writes as is),
//	// - Format verbs (will be substituted to the log's parts, such as
//	//   log's message, timestamp, log's level, log's fields, stacktrace, etc).
//	// These parts (w/o reallocation) will be stored with the same sequence
//	// as they represented in 'format' but there will be specified verb's types also.
//	// Moreover their common length will be calculated and stored to decrease
//	// destination []byte buffer reallocations (just allocate big buffer,
//	// at least as more as 'formatParts' required).
//	// -----
//	formatParts []formatPart
//	// Sum of: len of just text parts + predicted len of log's parts.
//	minimumBufferLen int
//	// The same as 'minimumBufferLen' but more, with a margin to decrease
//	// reallocations, improve RAM performance.
//	optimalBufferLen int
//}
//
//// formatPart represents built format's part (parsed RAW format string).
//// - for 'typ' == 'cefptJustText', 'value' is original RAW format string's part
//// (slice golang addr offset);
//// - for 'typ' == 'cefptColor', 'value' is calculated bash escape sequence
//// (like "\033[01;03;38;05;144m");
//// - for other 'typ' variants, 'value' is empty, 'cause it's log's verbs
//// and they're runtime calculated entities.
//type formatPart struct {
//
//	typ formatPartType
//	value string
//}
//
////
//type formatPartType int
//
//// [c]onsole [e]ncoder [f]ormat [p]art [t]ype (cefpt) predefined constants.
//const (
//
//	cefptJustText = 1
//	cefptColor = 2
//)
//
//var (
//	// Make sure we won't break API.
//	_ CommonIntegratorEncoder = (*ConsoleEncoder)(nil).encode
//
//	// Package's console encoder
//	consoleEncoder     CommonIntegratorEncoder
//	consoleEncoderAddr unsafe.Pointer
//)
//
//func init() {
//	consoleEncoder = (&ConsoleEncoder{}).FreezeAndGetEncoder()
//	consoleEncoderAddr = dangerous.TakeRealAddr(consoleEncoder)
//}
//
////
//func (ce *ConsoleEncoder) FreezeAndGetEncoder() CommonIntegratorEncoder {
//
//	if ce.isBuilt() {
//		return ce.encode
//	} else {
//		return ce.prepareToBuild().doBuild().encode
//	}
//}
//
////
//func (ce *ConsoleEncoder) isBuilt() bool {
//
//}
//
////
//func (ce *ConsoleEncoder) prepareToBuild() *ConsoleEncoder {
//
//}
//
////
//func (ce *ConsoleEncoder) doBuild() *ConsoleEncoder {
//
//	// start parsing ce.format
//	// all parsing loops are for-range based (because there is UTF-8 support)
//	// (yes, you can use not only ASCII parts in your format string,
//	// and yes if you do it, you are mad. stop it!).
//	var (
//		i = 0
//		formatLen = len(ce.format)
//		format = ce.format
//		wasBrace = false
//		wasBraceEnd = false
//		wasVerbStart = false
//		wasVerbEnd = false
//		minimumBufferLen = 0
//	)
//	// handle first 'just text' part (if it exists)
//	for _, char := range format {
//		switch {
//
//		case char == '{' && wasBrace:
//			wasBrace = false
//			wasVerbStart = true
//
//		case char == '{':
//			wasBrace = true
//			i++
//
//		default:
//			i++
//		}
//	}
//
//	// there was 'just text'part
//	if i != 0 {
//		ce.formatParts = append(ce.formatParts, formatPart{
//			typ: cefptJustText,
//			value: format[:i],
//		})
//		format = format[i:]
//		minimumBufferLen += i
//	}
//
//	for idx, char := range format {
//
//	}
//}
//
////
//func (ce *ConsoleEncoder) resolveVerb(verb string) (predictedLen int) {
//
//}
//
////
//func (ce *ConsoleEncoder) encode(e *Entry) []byte {
//
//	//
//	return nil
//}
