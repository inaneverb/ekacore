package errors

import (
	"sync/atomic"
)

var privateCounter uint64 = 1

//
func getPrivateID() uint64 {
	return atomic.AddUint64(&privateCounter, 1)
}

// generate full text error
func gfte(parentName, message, hidden string) string {

	isSkip := func(char byte) bool {
		return char <= 33 || // ASCII up to space and '!'
			char == '.' || char == ',' || char == ':' || char == ';'
	}

	msg := parentName
	if msg != "" {
		msg += ": "
	}

	if message != "" {
		msg += message
	}

	if hidden != "" {
		// skip last spaces
		i := len(msg) - 1
		for ; i >= 0 && isSkip(msg[i]); i-- {
		}

		msg = msg[:i+1]

		msg += ", cause " + hidden
	}

	return msg
}
