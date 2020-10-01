// Copyright © 2020. All rights reserved.
// Author: Ilya Stroy.
// Contacts: qioalice@gmail.com, https://github.com/qioalice
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
	stdout *stdSynced
)

func (ss *stdSynced) Write(b []byte) (n int, err error) {
	ss.Lock()
	n, err = ss.f.Write(b)
	ss.Unlock()
	return n, err
}

func initStdoutSynced() {
	stdout = new(stdSynced)
	stdout.f = os.Stdout
}
