// Copyright © 2020. All rights reserved.
// Author: Ilya Yuryevich.
// Contacts: iyuryevich@pm.me, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekasys

import (
	"os"
)

var (
	posixCachedUid = uint32(os.Getuid())
	posixCachedGid = uint32(os.Getgid())
)

// PosixCachedUid returns the cached value of os.Getuid() call.
func PosixCachedUid() uint32 { return posixCachedUid }

// PosixCachedGid returns the cached value of os.Getgid() call.
func PosixCachedGid() uint32 { return posixCachedGid }
