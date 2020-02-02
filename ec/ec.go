// Copyright © 2019. All rights reserved.
// Author: Ilya Yuryevich.
// Contacts: qioalice@gmail.com, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ec

import (
	"github.com/satori/go.uuid"
)

//
type EC int

//
type ECXT struct {
	EC
	UUID uuid.UUID
}

//
const (
	EOK EC = iota

	ENotFound
	EAccessDenied

	EInternalError

	EInvalidLimit
	EInvalidOffset

	EInvalidArg

	EAlreadyHandled = 99

	ECustom = 100
)

//
func RegMsg(ec EC, msg string) {

}

var ecm = map[EC]string{

	EOK:       `No errors`,
	ENotFound: `Not found`,

	EInternalError: `Internal error: DB or Unknown error`,
	EInvalidLimit:  `Invalid SQL DB query limit value <= 0 or > 100`,
	EInvalidOffset: `Invalid SQL DB query offset value < 0`,

	EInvalidArg: `Invalid arg (most likely is zero)`,
}

//
func ECStr(ec EC) (ecs string) {
	if ecs = ecm[ec]; ecs == "" {
		ecs = `Unknown error`
	}
	return
}

//
func (ec EC) ECXT() (ret ECXT) {

	if ret.EC = ec; ret.EC != EOK {
		ret.UUID = uuid.Must(uuid.NewV4())
	}
	return
}

//
func (ec EC) ECXTNil() ECXT {
	return ECXT{EC: ec}
}

//
func (ec EC) ECXTForce() ECXT {
	return ECXT{EC: ec, UUID: uuid.Must(uuid.NewV4())}
}

//
func (ecxt ECXT) EOK() bool {
	return ecxt.EC == EOK
}
