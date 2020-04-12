package log

import (
	"testing"

	"github.com/qioalice/gext/ec"
	"github.com/qioalice/gext/errors"
)

func TestFoo1(t *testing.T) {
	arg := String("arg3", "val3")
	Warnec(errors.InitializationFailed.New("hi"), ec.EInvalidArg,
		"arg1", "val1", "arg2", "val2", &arg)
}

func TestFoo2(t *testing.T) {
	Warnec(errors.InitializationFailed.New("hi"), ec.EInvalidArg,
		"arg1", "val1", "arg2", "val2", String("arg3", "val3"))
}
