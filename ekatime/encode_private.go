// Copyright © 2020. All rights reserved.
// Author: Ilya Stroy.
// Contacts: qioalice@gmail.com, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekatime

// batoi is atoi from ASCII string with 2 characters b1, b2.
// It allows to write algorithms to decode Date, Time, Timestamp from string
// faster up to 5x than using strconv.Atoi().
func batoi(b1, b2 byte) (i int16, valid bool) {
	if b2 = b2 - '0'; b2 > 9 {
		return 0, false
	}
	i = int16(b2)
	if b1 == '-' {
		return -i, true
	}
	if b1 = b1 - '0'; b1 > 9 {
		return 0, false
	}
	return int16(b1) * 10 + i, true
}
