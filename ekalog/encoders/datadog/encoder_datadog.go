// Copyright © 2020. All rights reserved.
// Author: Ilya Stroy.
// Contacts: qioalice@gmail.com, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekalog_encoder_datadog

import (
	"github.com/qioalice/ekago/v2/ekalog"

	"github.com/json-iterator/go"
)

//noinspection GoSnakeCaseUsage
type (
	// CI_DatadogEncoder is like CI_JSONEncoder: it also encodes *ekalog.Entry object
	// to the JSON but with the following changes:
	//
	// 1. You can not set an indentation. It's always 0. No tabs, no new lines.
	//    Even at the end of data buffer, that contains JSON encoded log entry.
	//
	// 2. Log entry's timestamp has a different name: "timestamp_real".
	//    See encodeBase() method code comments to figure out why:
	//    https://github.com/qioalice/ekago/ekalog/encoders/datadog/encoder_datadog_private.go
	//
	// 3. All log's fields (not attached error's ones) started with "sys."
	//    are ignored, except those which names started with "sys.dd.".
	//    They will be added at the JSON root using "<tail>" as their names
	//    (from original name "sys.dd.<tail>").
	//
	//    WARNING!
	//    THEY WILL BE ADDED ONLY IF THEIR VALUE'S KIND IS STRING
	//    AND IF THEY NAMES ARE NOT JUST "sys.dd." BUT SOMETHING ELSE.
	//
	//    WARNING!
	//    KEEP IN MIND! ONLY LOG'S FIELDS, NOT ATTACHED ERROR'S ONES.
	//    FIELDS OF ATTACHED ERROR WILL BE PLACED AS IS (almost).
	//
	//    Note.
	//    It's used to allow you to add some special DataDog meta fields like:
	//        - "service" (use "sys.dd.service" as name),
	//        - "hostname": (use "sys.dd.hostname" as name),
	//        - "ddtags": (use "sys.dd.ddtags" as name),
	//        - "ddsource": (use "sys.dd.ddsource" as name),
	//        etc.
	//
	// 4. Only attached error's ID will be added, if attached error presented.
	//    No error's class, no error's class ID, no error's public message.
	//
	// 5. All log's fields (except for those indicated in p3),
	//    and all attached error's fields (each stack frame's fields)
	//    is encoded as JSON objects to the JSON's root.
	//
	// 6. Stacktrace (log's or attached error's) is encoded
	//    as JSON array of strings to the JSON's root
	//    (because DataDog does not supports arrays of objects).
	//    Each string will represent stack frame (caller) in the following format:
	//
	//        "(<stack_index>): <func_name_with_fullpath>(<short_filename>:<file_line>)"
	//
	// 7. Attached error's messages (each stack frame's message) are encoded
	//    as JSON array of strings to the JSON's root.
	//    (because DataDog does not supports arrays of objects).
	//    Each string will represent stack frame's message in the following format:
	//
	//        "(<stack_index>): <message>"
	//
	// 8. Marked attached error's stack frames are the same as unmarked.
	//    They are encoded w/o changes.
	//
	// -----
	//
	// You may see an examples of encoding by calling TestExampleLog() test func:
	// https://github.com/qioalice/ekago/ekalog/encoders/datadog/encoder_datadog_test.go
	//
	CI_DatadogEncoder struct {

		// api is jsoniter's API object.
		// Created at the first FreezeAndGetEncoder() call for object.
		// Won't be called twice. Only one.
		//
		// See FreezeAndGetEncoder() and doBuild() methods for more info.
		api jsoniter.API
	}
)

var (
	// Make sure we won't break API.
	_ ekalog.CI_Encoder = (*CI_DatadogEncoder)(nil).encode
)

// FreezeAndGetEncoder builds current CI_DatadogEncoder if it has not built yet
// returning a function (has an alias CI_Encoder) that can be used at the
// CommonIntegrator.WithEncoder() call while initializing.
//
// WARNING!
// DO NOT CALL THIS METHOD WHEN CI_DatadogEncoder == nil!
// YOU WILL GET THE ENCODER FUNCTION THAT WILL PANIC WHEN CALLED!
//
func (de *CI_DatadogEncoder) FreezeAndGetEncoder() ekalog.CI_Encoder {
	return de.doBuild().encode
}

