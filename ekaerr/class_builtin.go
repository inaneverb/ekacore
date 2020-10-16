// Copyright © 2020. All rights reserved.
// Author: Ilya Stroy.
// Contacts: qioalice@gmail.com, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekaerr

var (
	// NotFound is a class for not found error
	NotFound = CommonErrors.NewClass("NotFound")

	// AlreadyExist is a class for an entity already exist error
	AlreadyExist = CommonErrors.NewClass("AlreadyExist")

	// IllegalArgument is a class for invalid argument error
	IllegalArgument = CommonErrors.NewClass("IllegalArgument")

	// IllegalState is a class for invalid state error
	IllegalState = CommonErrors.NewClass("IllegalState")

	// IllegalFormat is a class for invalid format error
	IllegalFormat = CommonErrors.NewClass("IllegalFormat")

	// InitializationFailed is a class for initialization error
	InitializationFailed = CommonErrors.NewClass("InitializationFailed")

	// DataUnavailable is a class for unavailable data error
	DataUnavailable = CommonErrors.NewClass("DataUnavailable")

	// ServiceUnavailable is a class for unavailable service error
	ServiceUnavailable = CommonErrors.NewClass("ServiceUnavailable")

	// UnsupportedOperation is a class for unsupported operation error
	UnsupportedOperation = CommonErrors.NewClass("UnsupportedOperation")

	// RejectedOperation is a class for rejected operation error
	RejectedOperation = CommonErrors.NewClass("RejectedOperation")

	// Interrupted is a class for interruption error
	Interrupted = CommonErrors.NewClass("Interrupted")

	// AssertionFailed is a class for assertion error
	AssertionFailed = CommonErrors.NewClass("AssertionFailed")

	// InternalError is a class for internal error
	InternalError = CommonErrors.NewClass("InternalError")

	// ExternalError is a class for external error
	ExternalError = CommonErrors.NewClass("ExternalError")

	// ConcurrentUpdate is a class for concurrent update error
	ConcurrentUpdate = CommonErrors.NewClass("ConcurrentUpdate")

	// TimeoutElapsed is a class for timeout error
	TimeoutElapsed = CommonErrors.NewClass("Timeout")

	// NotImplemented is an error class for lacking implementation
	NotImplemented = UnsupportedOperation.NewSubClass("NotImplemented")

	// UnsupportedVersion is a class for unsupported version error
	UnsupportedVersion = UnsupportedOperation.NewSubClass("UnsupportedVersion")
)
