// Copyright © 2020-2021. All rights reserved.
// Author: Ilya Stroy.
// Contacts: qioalice@gmail.com, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekarand

import (
	crand "crypto/rand"
	mrand "math/rand"

	"encoding/binary"
	"time"
)

func init() {

	var crEntropy [8]byte
	_, _ = crand.Read(crEntropy[:])

	d := int64(binary.BigEndian.Uint64(crEntropy[:]))
	mrand.Seed(time.Now().UnixNano() + d)
}
