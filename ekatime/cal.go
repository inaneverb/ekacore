// Copyright © 2020. All rights reserved.
// Author: Ilya Stroy.
// Contacts: qioalice@gmail.com, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekatime

import (
	"time"
	"unsafe"

	"github.com/qioalice/ekago/v2/ekaerr"
)

type (
	// Calendar
	//
	// WARNING!
	// NOT THREAD SAFETY FOR MIXING INITIALIZING AND ACCESSING!
	// DATA RACE ATTENTION! PANIC ATTENTION!
	//
	//   The Calendar's object is not thread safety at all.
	//   It means that you need to initialize it before using. E.g.:
	//   You must call WhenNewDay(), RegJsonEncoder(), RegYourOwnEncoder() BEFORE
	//   first call of Today() and Run(). Undefined behaviour and data race otherwise.
	Calendar struct {

		// TODO: Make Calendar type being thread-safe (at least partially).

		events []Event
		today unsafe.Pointer

		newDayTimer *time.Timer
		newDayCallback func(*Today)

		todayEncoderJson func(*Today) []byte
		todayEncoderUsers1 func(*Today) []byte
	}
)

//
func (c *Calendar) EventAdd(event Event) *Calendar {
	if c != nil && c.eventIdx(event) == -1 {
		c.events = append(c.events, event)
	}
	return c
}

//
func (c *Calendar) EventRemove(event Event) *Calendar {
	if c != nil {
		if idx := c.eventIdx(event); idx != -1 {
			c.events[idx] = c.events[len(c.events)-1]
			c.events = c.events[:len(c.events)-1]
		}
	}
	return c
}

//
func (c *Calendar) WhenNewDay(cb func(*Today)) *Calendar {
	if c != nil {
		c.newDayCallback = cb
	}
	return nil
}

//
func (c *Calendar) RegJsonEncoder(encoder func(*Today) []byte) *Calendar {
	if c != nil {
		c.todayEncoderJson = encoder
	}
	return c
}

//
func (c *Calendar) RegYourOwnEncoder(number int, encoder func(*Today) []byte) *Calendar {
	if c != nil {
		_ = number // for the future and because I dont wanna break API
		c.todayEncoderUsers1 = encoder
	}
	return c
}

//
func (c *Calendar) Today() *Today {
	return c.getToday()
}

//
func (c *Calendar) Run() *ekaerr.Error {

	// TODO: c == nil check

	// First call updateToday() at the same goroutine to prevent call Today()
	// before it will be initialized.
	c.updateToday()

	// Now, we must call updateToday() each time when new day has come.
	c.newDayTimer = time.AfterFunc(TillNextMidnight(), c.newDayHasCome)

	return nil
}