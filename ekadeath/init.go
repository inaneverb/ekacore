// Copyright © 2019-2021. All rights reserved.
// Author: Ilya Stroy.
// Contacts: iyuryevich@pm.me, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekadeath

import (
	"os"
	"os/signal"
	"syscall"
)

// Spawns goroutine which can handle SIGKILL, SIGTERM that leads to call Die(1).
func init() {
	var ch = make(chan os.Signal)
	signal.Notify(ch, os.Interrupt, syscall.SIGKILL, syscall.SIGTERM)
	go func() {
		_ = <-ch // blocks
		Die(1)
	}()
}
