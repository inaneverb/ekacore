// Copyright © 2019-2023. All rights reserved.
// Author: Ilya Yuryevich.
// Contacts: iyuryevich@pm.me, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekaio

type Syncer interface {
	Sync() error
}
