// Copyright © 2020. All rights reserved.
// Author: Ilya Stroy.
// Contacts: qioalice@gmail.com, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekasys

import (
	"io"
)

type (
	IStdSynced interface {
		io.Writer
	}
)

func Stdout() IStdSynced {
	return stdout
}
