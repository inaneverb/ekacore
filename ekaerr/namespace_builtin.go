// Copyright © 2020. All rights reserved.
// Author: Ilya Stroy.
// Contacts: iyuryevich@pm.me, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekaerr

var (
	// CommonErrors is a namespace for general purpose errors designed for universal usage.
	// These errors should typically be used in opaque manner, implying no handing in user code.
	// When handling is required, it is best to use custom error types.
	CommonErrors = newNamespace("Common", false)

	// REMINDER
	// DO NOT FORGOT USE PRIVATE NAMESPACE CONSTRUCTOR IF YOU WILL ADD A NEW
	// BUILTIN NAMESPACES!
	// AND USE CUSTOM == FALSE IN THAT CASE!
)
