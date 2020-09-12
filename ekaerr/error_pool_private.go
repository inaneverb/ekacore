package ekaerr

import (
	"runtime"
	"sync"
	"sync/atomic"
	"unsafe"

	"github.com/qioalice/ekago/v2/internal/ekafield"
	"github.com/qioalice/ekago/v2/internal/ekaletter"
)

type (
	errorPoolStat struct {
		AllocCalls   uint64
		NewCalls     uint64
		ReleaseCalls uint64
	}
)

//noinspection GoSnakeCaseUsage
const (
	// _LETTER_REUSE_MAX_LETTER_ITEMS is how much LetterItems chunks must be saved
	// into *Letter
	_LETTER_REUSE_MAX_LETTER_ITEMS = int16(4)

	// _LETTER_ITEM_ALLOC_CHUNK_SIZE is how much *LetterItem objects must be allocated
	// at a time and one by one at the RAM (as array) to decrease RAM fragmentation.
	_LETTER_ITEM_ALLOC_CHUNK_SIZE = int16(4)

	// _ERROR_POOL_INIT_COUNT is how much *Error (with *Letter) objects
	// errorPool pool will contain at the start.
	_ERROR_POOL_INIT_COUNT = 128
)

var (
	// errorPool is the pool of *Error (with *Letter) objects for being reused.
	errorPool sync.Pool

	eps_ errorPoolStat
)

//
//noinspection GoExportedFuncWithUnexportedType
func EPS() (eps errorPoolStat) {
	eps.AllocCalls = atomic.LoadUint64(&eps_.AllocCalls)
	eps.NewCalls = atomic.LoadUint64(&eps_.NewCalls)
	eps.ReleaseCalls = atomic.LoadUint64(&eps_.ReleaseCalls)
	return
}

// allocError creates a new *Error object, creates a new *Letter object inside Error,
// performs base initialization and returns it.
func allocError() interface{} {

	e := new(Error)
	e.letter = new(ekaletter.Letter)

	runtime.SetFinalizer(e.letter, releaseErrorForFinalizer)
	e.needSetFinalizer = false

	tail := (*ekaletter.LetterItem)(nil) // last allocated item from linked list

	e.letter.Items, tail = allocLetterItemsChunk()
	ekaletter.SetLastItem(e.letter, e.letter.Items)

	for i := int16(0); i < _LETTER_REUSE_MAX_LETTER_ITEMS-1; i++ {
		newHead, newTail := allocLetterItemsChunk()
		ekaletter.SetNextItem(tail, newHead)
		tail = newTail
	}

	// SystemFields is used for saving Entry's meta data.
	// https://github.com/qioalice/ekago/internal/letter/letter.go

	e.letter.SystemFields = make([]ekafield.Field, 4)

	e.letter.SystemFields[_ERR_SYS_FIELD_IDX_CLASS_ID].Key = "class_id"
	e.letter.SystemFields[_ERR_SYS_FIELD_IDX_CLASS_ID].Kind |=
		ekafield.KIND_FLAG_SYSTEM | ekafield.KIND_SYS_TYPE_EKAERR_CLASS_ID

	e.letter.SystemFields[_ERR_SYS_FIELD_IDX_CLASS_NAME].Key = "class_name"
	e.letter.SystemFields[_ERR_SYS_FIELD_IDX_CLASS_NAME].Kind |=
		ekafield.KIND_FLAG_SYSTEM | ekafield.KIND_SYS_TYPE_EKAERR_CLASS_NAME

	e.letter.SystemFields[_ERR_SYS_FIELD_IDX_PUBLIC_MESSAGE].Key = "public_message"
	e.letter.SystemFields[_ERR_SYS_FIELD_IDX_PUBLIC_MESSAGE].Kind |=
		ekafield.KIND_FLAG_SYSTEM | ekafield.KIND_SYS_TYPE_EKAERR_PUBLIC_MESSAGE

	e.letter.SystemFields[_ERR_SYS_FIELD_IDX_ERROR_ID].Key = "error_id"
	e.letter.SystemFields[_ERR_SYS_FIELD_IDX_ERROR_ID].Kind |=
		ekafield.KIND_FLAG_SYSTEM | ekafield.KIND_SYS_TYPE_EKAERR_UUID

	// We saving current e's ptr as *Letter's something for able to get an *Error
	// addr using its *Letter (used at the releaseErrorForGate()).
	ekaletter.SetSomething(e.letter, unsafe.Pointer(e))

	atomic.AddUint64(&eps_.AllocCalls, 1)
	return e
}

// initErrorPool initializes errorPool creating and storing
// exactly _ERROR_POOL_INIT_COUNT *Error objects to that pool.
func initErrorPool() {
	errorPool.New = allocError
	for i := 0; i < _ERROR_POOL_INIT_COUNT; i++ {
		errorPool.Put(allocError())
	}
}

// acquireError returns a new *Error object from the Error's pool or newly instantiated.
func acquireError() *Error {
	atomic.AddUint64(&eps_.NewCalls, 1)
	return errorPool.Get().(*Error).prepare()
}

// releaseError returns 'e' to the Error's pool for being reused in the future
// and that Error could be obtained later using acquireError().
func releaseError(e *Error) {
	atomic.AddUint64(&eps_.ReleaseCalls, 1)
	errorPool.Put(e.cleanup())
}

// releaseErrorForGate is just the same as releaseError but it tries to extract
// *Error addr assuming that 'errLetter' is the Error's one.
func releaseErrorForGate(errLetter *ekaletter.Letter) {
	if errLetter != nil {
		if e := (*Error)(ekaletter.GetSomething(errLetter)); e != nil {
			e.letter = errLetter
			releaseError(e)
		}
	}
}

//
func releaseErrorForFinalizer(errLetter *ekaletter.Letter) {
	if errLetter != nil {
		if e := (*Error)(ekaletter.GetSomething(errLetter)); e != nil {
			e.letter = errLetter
			e.needSetFinalizer = true
			releaseError(e)
		}
	}
}

// allocLetterItems allocates *LetterItem linked list that contains exactly
// _LETTER_ITEM_ALLOC_CHUNK_SIZE preallocated *LetterItem items.
// Returns the first and the last item of that list.
func allocLetterItemsChunk() (head, tail *ekaletter.LetterItem) {

	// preallocate as array to avoid RAM fragmentation
	pool := make([]ekaletter.LetterItem, _LETTER_ITEM_ALLOC_CHUNK_SIZE)
	num := _LETTER_ITEM_ALLOC_CHUNK_SIZE - 1
	ret := (*ekaletter.LetterItem)(nil)

	for ; num >= 0; num-- {
		ret = ekaletter.SetNextItem(&pool[num], ret)
		ekaletter.SetStackFrameIdx(ret, -1)
	}

	return ret, &pool[_LETTER_ITEM_ALLOC_CHUNK_SIZE-1]
}

// pruneLetterItemsChunk recursively removes a link in *LetterItem linked list
// starting from 'head', also frees all allocated resources.
func pruneLetterItemsChunk(head *ekaletter.LetterItem) {

	for head != nil {

		for i, n := 0, len(head.Fields); i < n; i++ {
			ekafield.Reset(&head.Fields[i])
		}

		head.Fields = nil
		head.Message = ""

		next := ekaletter.GetNextItem(head)
		ekaletter.SetNextItem(head, nil)
		head = next
	}
}
