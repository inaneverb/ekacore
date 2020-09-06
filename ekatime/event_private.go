// Copyright © 2020. All rights reserved.
// Author: Ilya Stroy.
// Contacts: qioalice@gmail.com, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekatime

//noinspection GoSnakeCaseUsage
const (
	_EVENT_OFFSET_DATE       uint8 = 0
	_EVENT_OFFSET_ID         uint8 = _EVENT_OFFSET_DATE + _DATE_OFFSET_UNUSED
	_EVENT_OFFSET_IS_WORKDAY uint8 = _EVENT_OFFSET_ID + 15

	_EVENT_MASK_IS_WORKDAY   uint8 = 0x01
	_EVENT_MASK_ID           uint16 = 0x7FFF
)
