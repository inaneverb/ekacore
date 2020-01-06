// Copyright © 2019. All rights reserved.
// Author: Ilya Yuryevich.
// Contacts: qioalice@gmail.com, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package types

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
