// Copyright © 2020. All rights reserved.
// Author: Ilya Stroy.
// Contacts: qioalice@gmail.com, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekagen

import (
	crand "crypto/rand"
	mrand "math/rand"

	"github.com/qioalice/ekago/v2/ekatime"
)

var (
	r *mrand.Rand
)

func init() {

	var (
		d = make([]byte, 32)
		seed = int64(99)
	)

	if _, legacyErr := crand.Read(d); legacyErr == nil {
		for _, b := range d {
			seed += int64(b)
		}
	}

	r = mrand.New(mrand.New(mrand.NewSource(seed)))
	r.Seed(ekatime.Now().I64())
}
