// Copyright © 2020. All rights reserved.
// Author: Ilya Yuryevich.
// Contacts: qioalice@gmail.com, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekasys

import (
	"os"
)

var (
	posixCachedUid = uint32(os.Getuid())
	posixCachedGid = uint32(os.Getgid())
)

func PosixCachedUid() uint32 { return posixCachedUid }
func PosixCachedGid() uint32 { return posixCachedGid }
