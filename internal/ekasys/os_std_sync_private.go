// Copyright © 2021. All rights reserved.
// Author: Ilya Stroy.
// Contacts: iyuryevich@pm.me, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekasys

import (
	"os"
	"sync"
)

type (
	stdSynced struct {
		f *os.File
		sync.Mutex
	}
)

var (
	Stdout *stdSynced
)

func (ss *stdSynced) Write(b []byte) (n int, err error) {
	ss.Lock()
	n, err = ss.f.Write(b)
	ss.Unlock()
	return n, err
}

func initStdoutSynced() {
	Stdout = new(stdSynced)
	Stdout.f = os.Stdout
}
