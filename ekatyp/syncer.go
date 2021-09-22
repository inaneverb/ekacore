// Copyright © 2019. All rights reserved.
// Author: Ilya Yuryevich.
// Contacts: iyuryevich@pm.me, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekatyp

import "io"

type Syncer interface {
	Sync() error
}

type WriteSyncer interface {
	io.Writer
	Syncer
}

type WriteSyncCloser interface {
	io.WriteCloser
	Syncer
}
