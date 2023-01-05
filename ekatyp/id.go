// Copyright © 2021. All rights reserved.
// Author: Ilya Stroy.
// Contacts: iyuryevich@pm.me, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekatyp

import (
	"database/sql"
	"database/sql/driver"
	"encoding"
	"encoding/json"
	"errors"
	"fmt"
)

type ID interface {
	encoding.BinaryMarshaler
	encoding.BinaryUnmarshaler
	encoding.TextMarshaler
	encoding.TextUnmarshaler
	json.Marshaler
	json.Unmarshaler
	fmt.Stringer
	sql.Scanner
	driver.Valuer

	Equal(other ID) bool
	IsNil() bool
	SetNil()
	Bytes() []byte
}

var (
	ErrIDNilDestination = errors.New("ID: Nil destination")
)
