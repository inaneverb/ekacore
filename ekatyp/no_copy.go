// Copyright © 2020. All rights reserved.
// Author: Ilya Stroy.
// Contacts: iyuryevich@pm.me, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekatyp

type NoCopy struct{}

func (*NoCopy) Lock()   {}
func (*NoCopy) Unlock() {}
