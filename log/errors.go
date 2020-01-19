// Copyright © 2019. All rights reserved.
// Author: Ilya Yuryevich.
// Contacts: qioalice@gmail.com, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package log

import (
	"github.com/qioalice/gext/errors"
)

var (
	errLoggerObjectInvalid = errors.IllegalArgument.
		New("Logger object is invalid")
)
