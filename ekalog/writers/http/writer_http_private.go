// Copyright © 2020. All rights reserved.
// Author: Ilya Stroy.
// Contacts: qioalice@gmail.com, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekalog_writer_http

import (
	"bytes"
	"context"
	"sync/atomic"
	"time"

	"github.com/qioalice/ekago/v2/ekadeath"
	"github.com/qioalice/ekago/v2/ekaerr"

	"github.com/valyala/fasthttp"
)

//noinspection GoSnakeCaseUsage
const (
	// The values of CI_WriterHttp's 'casInitStatus' field.
	//
	// The string "-> <status>" (in comments) means
	// "To which status the current status can be changed to".
	// ---------

	// CI_WriterHttp object created, not started. Workers aren't spawned yet.
	// Channels aren't created yet. Buffers aren't allocated yet.
	//
	//   -> _DDWRITER_CAS_STATUS_INITIALIZING.
	//
	_DDWRITER_CAS_STATUS_NOT_INITIALIZED = int32(0)

	// CI_WriterHttp object under initializing right now by some goroutine.
	// If one goroutine tries to take a responsibility of initialization
	// but seeing that current status is "initializing", it will lock yourself
	// and waits until other goroutine initializes.
	//
	//   -> _DDWRITER_CAS_STATUS_READY
	//   -> _DDWRITER_CAS_STATUS_FINALLY_DISABLED
	//
	_DDWRITER_CAS_STATUS_INITIALIZING = int32(1)

	// CI_WriterHttp successfully initialized, worked. All workers spawned.
	//
	//   -> _DDWRITER_CAS_STATUS_TEMPORARY_DISABLED
	//
	_DDWRITER_CAS_STATUS_READY = int32(10)

	// Two cases:
	//
	// 1. There was a some network problem and CI_WriterHttp "paused" at this moment.
	//    (Only one worker is running at this moment and tries to restore connection).
	//
	// 2. Final dying is requested by destructor.
	//    Will be changed to "finally disabled" soon (almost instantly).
	//
	//   -> _DDWRITER_CAS_STATUS_READY
	//   -> _DDWRITER_CAS_STATUS_FINALLY_DISABLED
	//
	_DDWRITER_CAS_STATUS_TEMPORARY_DISABLED = int32(-2)

	// The CI_WriterHttp has been completely stop and will NEVER run again.
	// This status CAN NOT be changed.
	_DDWRITER_CAS_STATUS_FINALLY_DISABLED = int32(-4)
)

//noinspection GoSnakeCaseUsage
const (
	// Default values for _DDWRITER_CAS_STATUS_INITIALIZING's fields that are not set,
	// or had an incorrect values.

	_DEFAULT_DDWRITER_ENTRIES_TOTAL_BUF_SIZE = 1024
	_DEFAULT_DDWRITER_WORKER_NUM = 2
	_DEFAULT_DDWRITER_ENTRIES_PER_WORKER_BUF_SIZE = 32
	_DEFAULT_DDWRITER_WORKER_FLUSH_DELAY = 10 * time.Second
	_DEFAULT_DDWRITER_WORKER_FLUSH_DEFERRED_PER_ITER = 5
)

// configure is a private part of public configuration methods.
// Calls 'cb' passing 'dw' assuming that 'cb' will update some field in the 'dw'.
// Does it only if CI_WriterHttp has not been started (initialized) yet.
//
// Because it's private method, it guarantees that 'cb' != nil.
// Nil safe.
func (dw *CI_WriterHttp) configure(cb func(dw *CI_WriterHttp)) *CI_WriterHttp {

	if dw != nil {
		dw.slowInit.Lock()
		defer dw.slowInit.Unlock()

		if atomic.LoadInt32(&dw.casInitStatus) == _DDWRITER_CAS_STATUS_NOT_INITIALIZED {
			dw.beenPinged = false
			cb(dw)
		}
	}
	return dw
}

// canWrite reports whether Write() method can add a new encoded log entry
// as []byte to the 'entries' channel. If CI_WriterHttp is not initialized yet,
// it does an initialization and starts workers.
func (dw *CI_WriterHttp) canWrite() bool {

	switch atomic.LoadInt32(&dw.casInitStatus) {
	case _DDWRITER_CAS_STATUS_READY:
		return true
	case _DDWRITER_CAS_STATUS_FINALLY_DISABLED:
		return false
	}

	// Take a responsibility about CI_WriterHttp initializing.
	// We acquiring mutex, and we will initialize a writer.
	//
	// (Or we could acquire it only when someone else finishes an initialization:
	// in that case it's OK for us, we'll check it again and true is just returned then).
	dw.slowInit.Lock()

	// It's guarantees that after we getting a Lock, the status might be:
	// - Not initializing: We'll do it;
	// - Ready: Another goroutine already do it;
	// Another statuses are impossible.

	// Assume, the current status is "not initializing" and try to start initializing.
	continueInitialization := atomic.CompareAndSwapInt32(&dw.casInitStatus,
		_DDWRITER_CAS_STATUS_NOT_INITIALIZED, _DDWRITER_CAS_STATUS_INITIALIZING)

	if continueInitialization {
		err := dw.performInitialization()
		if err.IsNil() {
			atomic.StoreInt32(&dw.casInitStatus, _DDWRITER_CAS_STATUS_READY)
		} else {
			atomic.StoreInt32(&dw.casInitStatus, _DDWRITER_CAS_STATUS_FINALLY_DISABLED)
		}

		dw.slowInit.Unlock()

		err.LogAsError()
		return err.IsNil()
	}

	// Well, looks like another goroutine already did initialize it.
	// Thanks them, do nothing.
	dw.slowInit.Unlock()
	return dw.canWrite()
}

// performInitialization initializes a CI_WriterHttp. Returns error if no provider
// is selected or ping was failed.
// Tries to ping only if the user didn't do it himself.
func (dw *CI_WriterHttp) performInitialization() *ekaerr.Error {

	// At this code point, dw.slowInit mutex is acquired (locked).
	// Overwrite default or incorrect values.

	if dw.providerInitializer == nil {
		return ekaerr.IllegalArgument.
			New("CI_WriterHttp: Provider is not presented. " +
				"Call UseProvider<provider>() method or UseProviderManual().").
			Throw()
	}

	// Try to ping.
	if !dw.beenPinged {
		if err := dw.ping(false, nil); err.IsNotNil() {
			return err.
				AddMessage("CI_WriterHttp: Ping failed. " +
					"You may avoid internal ping call if you call Ping() explicitly.").
				Throw()
		}
	}

	// OK, run workers.

	dw.workerTickers = make([]*time.Ticker, dw.workerNum)
	dw.entries = make(chan []byte, dw.entriesBufferLen)
	dw.entriesPackDeferred = make(chan *bytes.Buffer, dw.deferredEntriesBufferLen)

	dw.workersWg.Add(int(dw.workerNum))
	dw.ctx, dw.cancelFunc = context.WithCancel(context.Background())

	for i := uint16(0); i < dw.workerNum; i++ {
		dw.workerTickers[i] = time.NewTicker(dw.workerFlushDelay)
		go dw.worker(i == 0, dw.entries, dw.workerTickers[i].C)
	}

	// OK, workers ran, register destructor
	// (we need to flush all changes before app will be closed).
	ekadeath.Reg(func() {
		lostEntries := atomic.LoadUint64(&dw.entriesCompletelyLostCounter)
		if lostEntries > 0 {
			ekaerr.RejectedOperation.
				New("CI_WriterHttp: Some log entries are lost and will never be logged.").
				AddFields("ci_writer_http_min_lost_entries_num", lostEntries).
				LogAsWarn()
		}
		dw.disable(false)
	})

	return nil
}

// initOverwriteZeroValues overwrites CI_WriterHttp's fields that are set to the
// incorrect values by setters or has not been set at all.
func (dw *CI_WriterHttp) initOverwriteZeroValues() {

	if dw.entriesBufferLen <= 0 {
		dw.entriesBufferLen = _DEFAULT_DDWRITER_ENTRIES_TOTAL_BUF_SIZE
	}

	if dw.workerNum <= 0 {
		dw.workerNum = _DEFAULT_DDWRITER_WORKER_NUM
	}

	if dw.workerEntriesBufferLen <= 0 {
		dw.workerEntriesBufferLen = _DEFAULT_DDWRITER_ENTRIES_PER_WORKER_BUF_SIZE
	}

	if dw.workerFlushDelay <= 0 {
		dw.workerFlushDelay = _DEFAULT_DDWRITER_WORKER_FLUSH_DELAY
	}

	if dw.workerFlushDeferredPerIter <= 0 {
		dw.workerFlushDeferredPerIter = _DEFAULT_DDWRITER_WORKER_FLUSH_DEFERRED_PER_ITER
	}
}

// disable disables temporary or finally the CI_WriterHttp object.
// Temporary disabling performs only when some HTTP request has been failed,
// but final disabling performs only at the stopping the whole app.
//
// If it's the final disabling, stops all internal timers, tickers, workers, etc.
// Flushes all pending log entries, closes connections, prepares self to disable.
//
// It guarantees that disable() will be locked until performInitialization()
// is processed. Thus, there won't be an "initialization <-> dying" data race.
func (dw *CI_WriterHttp) disable(temporary bool) {

	// Do not change the order or CAS and mutex acquiring.
	// It's made for preventing data race between disable() and performInitialization().
	dw.slowInit.Lock()

	atomic.CompareAndSwapInt32(&dw.casInitStatus,
		_DDWRITER_CAS_STATUS_READY, _DDWRITER_CAS_STATUS_TEMPORARY_DISABLED)

	if temporary {
		dw.slowInit.Unlock()
		return
	}

	atomic.StoreInt32(&dw.casInitStatus, _DDWRITER_CAS_STATUS_FINALLY_DISABLED)
	dw.cancelFunc()

	//close(dw.entries)
	//
	// WARNING! PANIC CAUTION! DO NOT UNCOMMENT!
	// 'entries' MUST NOT BE CLOSED, because of the fact that Write() may receive
	// a 'true' from canWrite() before 'entries' has been closed, and then
	// there will be an attempt to write to closed channel. And it's a panic.

	for i := uint16(0); i < dw.workerNum; i++ {
		dw.workerTickers[i].Stop()
	}

	// DO NOT CHANGE THE ORDER!
	dw.slowInit.Unlock()
	dw.workersWg.Wait()

	close(dw.entriesPackDeferred)
}

// ping tries to perform a dummy HTTP request to the log service provider
// explicitly by user (calling Ping() method) or implicitly (at the initialization).
func (dw *CI_WriterHttp) ping(

	doLock bool, // shall CI_WriterHttp's mutex must be locked while pinging
	cbs []func(req *fasthttp.Request), // additional ping HTTP request initializers

) *ekaerr.Error {

	if doLock {
		dw.slowInit.Lock()
		defer dw.slowInit.Unlock()
	}

	dw.beenPinged = true
	dw.initOverwriteZeroValues()

	buf := bytes.NewBuffer(nil)
	_, _ = buf.WriteString("[]")

	return dw.sendRequest(buf, cbs)
}

// worker is a CI_WriterHttp's worker that runs in the separate goroutine.
// Work goroutines spawns only one at the initialization.
func (dw *CI_WriterHttp) worker(

	masterWorker bool, // indicates whether this worker is master (not slave)
	encodedEntries <-chan []byte, // channel, log entries being processed are coming from
	ticker <-chan time.Time, // ticker "when accumulated log entries must be flushed"
) {
	defer dw.workersWg.Done()

	// Internal entries buffer. Reusable.
	// At this moment allocates 1 MB RAM. TODO: Probably add a setter?
	buf := bytes.NewBuffer(make([]byte, 0, 2 << 20))
	_, _ = buf.Write(dw.dataBefore)

	// ProcessAndSendBuf is a helper function, that removes last 'dateBetween'
	// from 'buf', writes 'dateAfter' instead ans send it to the log service.
	ProcessAndSendBuf := func(masterWorker bool, dw *CI_WriterHttp, buf *bytes.Buffer) {
		buf.Truncate(buf.Len()- len(dw.dataBetween))
		_, _ = buf.Write(dw.dataAfter)
		dw.processEntriesBuffer(masterWorker, buf)
		buf.Reset()
		_, _ = buf.Write(dw.dataBefore)
	}

	doneChan := dw.ctx.Done()

	// i is workerBuffer's index.
	for i := uint16(0) ;; {
		select {

		case <-doneChan:
			if i > 0 {
				// There saved unprocessed entries. Send them.
				ProcessAndSendBuf(masterWorker, dw, buf)
			}
			return

		case encodedEntry := <-encodedEntries:
			// Add received encoded entry to the internal worker pool.
			_, _ = buf.Write(encodedEntry) // always returns nil as error
			_, _ = buf.Write(dw.dataBetween)
			i++

			// If the pool is full, flush it using its HTTP/S API,
			// reuse the pool after flushing by setting its index = 0.
			if i == dw.workerEntriesBufferLen {
				ProcessAndSendBuf(masterWorker, dw, buf)
				i = 0
			}

		case <-ticker: // never been closed, even if Stop() is called
			// Oops, it's time for scheduled flush. It doesn't matter whether
			// pool is full or not yet. Flush anyway, if pool contains something.
			if i > 0 {
				ProcessAndSendBuf(masterWorker, dw, buf)
				i = 0
			}
		}
	}
}

// processEntriesBuffer tries to perform an HTTP request using 'buf' as HTTP POST
// request's body, changing the CI_WriterHttp's status:
//
// -> Temporary disabled if request has been failed, deferring failed request
//    to being processed later (when connection will be restored),
//
// -> Ready if CI_WriterHttp has been temporary disabled,
//    but connection had been successfully restored.
//
// While CI_WriterHttp is temporary disabled, tries to send HTTP request
// only if 'masterWorker' is true.
func (dw *CI_WriterHttp) processEntriesBuffer(

	masterWorker bool, // indicates whether this worker is master (not slave)
	buf *bytes.Buffer, // an HTTP POST request's body
) {
	switch status := atomic.LoadInt32(&dw.casInitStatus); {

	case status == _DDWRITER_CAS_STATUS_TEMPORARY_DISABLED && !masterWorker:

		// We won't send these entries at this moment, but after returning from
		// this method, 'buf' will be reused. We have to copy that.
		bufCopy := bytes.NewBuffer(make([]byte, 0, buf.Len()))
		_, _ = bufCopy.Write(buf.Bytes())

		// Try to send made copy later, when connection will be recovered.
		select {
		case dw.entriesPackDeferred <- bufCopy:
		default:
			// The buffer of encoded deferred log entries is full.
			// We can't do something with that.
			atomic.AddUint64(&dw.entriesCompletelyLostCounter, 1)
		}

		// This is not master worker.
		// Only master worker can check whether connections is established again.
		return
	}

	if err := dw.sendRequest(buf, nil); err.IsNotNil() {
		dw.disable(true)
		err.LogAsErrorw("Failed to log to DataDog") // TODO
		return
	}

	// Request finished w/ no error.
	// Maybe CI_WriterHttp was temporary disabled and now it's time to recover?
	atomic.CompareAndSwapInt32(&dw.casInitStatus,
		_DDWRITER_CAS_STATUS_TEMPORARY_DISABLED, _DDWRITER_CAS_STATUS_READY)

	deferredEntriesPackNum := uint16(len(dw.entriesPackDeferred))
	if deferredEntriesPackNum == 0 {
		return
	} else if deferredEntriesPackNum > dw.workerFlushDeferredPerIter {
		// We have to limit how much buffers will be processed again but only
		// if it's not the call in the destructor (the latest pushing attempt).
		if atomic.LoadInt32(&dw.casInitStatus) != _DDWRITER_CAS_STATUS_FINALLY_DISABLED {
			deferredEntriesPackNum = dw.workerFlushDeferredPerIter
		}
	}

	for i := uint16(0); i < deferredEntriesPackNum; i++ {
		select {
		case deferredEntriesPack := <- dw.entriesPackDeferred:
			deferredEntriesPackBak := bytes.NewBuffer(deferredEntriesPack.Bytes())

			if err := dw.sendRequest(deferredEntriesPack, nil); err.IsNotNil() {

				// Oops, failed again.
				dw.disable(true)
				err.LogAsErrorw("Failed to log to DataDog") // TODO

				// Try to defer entries pack again.
				select {
				case dw.entriesPackDeferred <- deferredEntriesPackBak:
				default:
					atomic.AddUint64(&dw.entriesCompletelyLostCounter, 1)
				}

				// Request failed. Next time will be better (hope).
				return
			}

		default:
			// There is less data than we expecting.
			// Guess another goroutine already did the job. There is nothing left to do.
			return
		}
	}
}

// sendRequest sends an HTTP POST request to the remote log provider using fasthttp,
// applying stored provider callback at the initialization to the fasthttp.Request
// object and then applying each callback from 'cbs' one by one.
//
// If HTTP request was failed (returned non 200, 202 HTTP codes),
// an error object will be returned.
func (dw *CI_WriterHttp) sendRequest(

	buf *bytes.Buffer, // The bytes that will be attached to HTTP request as POST body
	cbs []func(req *fasthttp.Request), // Callbacks will be called before sending request

) *ekaerr.Error {

	req := fasthttp.AcquireRequest(); defer fasthttp.ReleaseRequest(req)
	resp := fasthttp.AcquireResponse(); defer fasthttp.ReleaseResponse(resp)

	req.Header.SetMethod(fasthttp.MethodPost)
	req.SetBodyStream(buf, buf.Len())

	dw.providerInitializer(req)

	for i, n := 0, len(cbs); i < n; i++ {
		if cbs[i] != nil {
			cbs[i](req)
		}
	}

	if legacyErr := dw.c.DoRedirects(req, resp, 5); legacyErr != nil {
		return ekaerr.ExternalError.
			Wrap(legacyErr, "CI_WriterHttp: Failed to perform HTTP request.").
			AddFields("ci_writer_http_url", string(req.RequestURI())).
			Throw()
	}

	switch status := resp.StatusCode(); status {
	case fasthttp.StatusOK, fasthttp.StatusAccepted:
	default:
		return ekaerr.ExternalError.
			New("CI_WriterHttp: Unexpected HTTP status code.").
			AddFields("ci_writer_http_status_code", status,
				"ci_writer_http_url", string(req.RequestURI())).
			Throw()
	}

	return nil
}