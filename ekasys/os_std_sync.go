// Copyright © 2020. All rights reserved.
// Author: Ilya Stroy.
// Contacts: iyuryevich@pm.me, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekasys

import (
	"io"

	"github.com/qioalice/ekago/v3/internal/ekasys"
)

type (
	IStdSynced interface {
		io.Writer
	}
)

func Stdout() IStdSynced {
	return ekasys.Stdout
}
