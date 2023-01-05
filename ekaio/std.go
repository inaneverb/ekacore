// Copyright © 2021-2023. All rights reserved.
// Author: Ilya Stroy.
// Contacts: iyuryevich@pm.me, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekaio

import (
	"io"
	"os"
)

var (
	syncedStdin  = NewSyncedReader(os.Stdin)
	syncedStdout = NewSyncedWriter(os.Stdout)
	syncedStderr = NewSyncedWriter(os.Stderr)
)

// GetSyncedStdin returns an STDIN as an io.Reader, all operation of which
// are protected with mutex. Thus, all read op are synchronous.
func GetSyncedStdin() io.Reader { return syncedStdin }

// GetSyncedStdout returns an STDOUT as an io.Writer, all operation of which
// are protected with mutex. Thus, all write op are synchronous.
func GetSyncedStdout() io.Writer { return syncedStdout }

// GetSyncedStderr returns an STDERR as an io.Writer, all operation of which
// are protected with mutex. Thus, all write op are synchronous.
func GetSyncedStderr() io.Writer { return syncedStderr }
