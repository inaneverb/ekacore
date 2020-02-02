package errors

import (
	"github.com/qioalice/gext/sys"

	"github.com/modern-go/reflect2"
)

type Error struct {
	parentId   uint64
	parentName string

	StackTrace sys.StackTrace

	Message string
	Hidden  string
}

var (
	// Used for internal type comparision.
	reflectTypeError    = reflect2.TypeOf(Error{})
	reflectTypeErrorPtr = reflect2.TypeOf((*Error)(nil))
)

//
func (e *Error) Error() string {
	return e.Full()
}

//
func (e *Error) Full() string {

	if e == nil || e.parentId == 0 {
		return ""
	}

	return gfte(e.parentName, e.Message, e.Hidden)
}

//
func Unwrap(e error) *Error {

	if reflect2.IsNil(e) {
		return nil
	}

	// TODO: use reflect2 for assigning and type checking
	if e, ok := e.(*Error); ok {
		return e
	}

	return nil
}

//
func Cast(e error) *Error {

	// TODO: Implement (cast any error to Error)
	return nil
}

//
func newError(parentId uint64, parentName string) *Error {

	if parentId == 0 {
		return nil
	}

	return &Error{
		parentId:   parentId,
		parentName: parentName,
	}
}
