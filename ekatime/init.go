// Copyright © 2020. All rights reserved.
// Author: Ilya Stroy.
// Contacts: iyuryevich@pm.me, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekatime

func init() {

	// WARNING!
	// CHANGE THE ORDER CAREFULLY!

	initTable5() // must be first always!

	// ---------

	initDateNumStr()
	initTimeNumStr()

	initWeekday()
}
