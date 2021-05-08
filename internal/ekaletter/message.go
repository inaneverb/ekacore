// Copyright © 2021. All rights reserved.
// Author: Ilya Stroy.
// Contacts: qioalice@gmail.com, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekaletter

type (
	// LetterMessage is a special type that represents a pair:
	// some message and its stack frame idx in the relevant ekasys.StackTrace slice.
	//
	// This struct is designed to be a part of Letter.
	// When Letter object is instantiated and initialized, it has its own stacktrace.
	// Using methods you can add a some message to each stacktrace's frame.
	// That's why you need a LetterMessage struct.
	//
	// It guarantees that StackFrameIdx < related ekasys.StackTrace's len,
	// and Body cannot be an empty string.
	LetterMessage struct {
		Body          string
		StackFrameIdx int16
	}
)
