// Copyright © 2020-2021. All rights reserved.
// Author: Ilya Stroy.
// Contacts: iyuryevich@pm.me, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekaerr

import (
	"github.com/modern-go/reflect2"

	"github.com/qioalice/ekago/v4/internal/ekaletter"
)

func init() {
	errorPool.New = allocError

	ekaletter.BridgeErrorGetLetter = bridgeGetLetter

	ekaletter.BridgeErrorGetStackIdx = bridgeGetStackIdx
	ekaletter.BridgeErrorSetStackIdx = bridgeSetStackIdx

	// It's prohibited to use some types as Error's fields.
	ignoredTypes := []uintptr{
		reflect2.RTypeOf(Class{}), reflect2.RTypeOf((*Class)(nil)),
		reflect2.RTypeOf(Namespace{}), reflect2.RTypeOf((*Namespace)(nil)),
		reflect2.RTypeOf(Error{}), reflect2.RTypeOf((*Error)(nil)),
	}
	ekaletter.RTypesBeingIgnoredForParsing =
		append(ekaletter.RTypesBeingIgnoredForParsing, ignoredTypes...)
}
