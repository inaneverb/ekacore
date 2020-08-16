// Copyright © 2020. All rights reserved.
// Author: Ilya Stroy.
// Contacts: qioalice@gmail.com, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekalog_writer_http

import (
	"bytes"
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/qioalice/ekago/v2/ekaerr"

	"github.com/valyala/fasthttp"
)

//noinspection GoSnakeCaseUsage
type (
	// CI_WriterHttp is a type that implements an io.Writer - legacy Golang interface,
	// doing write encoded log's entry as []byte to the some HTTP service using
	// desired API.
	//
	// Features:
	// -----------
	//
	// 1. Async transport.
	//    When you calling Write() it just pushes encoded entry (as []byte)
	//    to the worker and does not blocks the routine.
	//    Spawn as many workers (at the initialization) as you want.
	//    See SetWorkersNum() method.
	//
	// 2. Thread-safe.
	//    You may call CI_WriterHttp to as many goroutines as you want.
	//    The CommonIntegrator under the hood just calls Write() for all
	//    writers they holds on.
	//
	// 3. Accumulated bulk deferred requests.
	//    Has an internal buffer of encoded log entries (as []byte)
	//    (see SetBufferCap() method), each worker pulling data from
	//    and putting to its own internal buffer (see SetWorkerBufferCap() method).
	//
	//    When worker buffer is full or when it's time to flush the accumulated data
	//    (the stopping of app or flush time has come), sends accumulated data
	//    to the HTTP service.
	//
	// 4. Control how encoded log entries will be combined before sending.
	//    Each worker aggregates a single encoded log entries with others
	//    to the log entries pack (contains many encoded log entries) before send it.
	//
	//    You may set HOW they will be combined:
	//
	//        - Need to put something to buffer before first encoded log entry added?
	//          (see AddBefore() method).
	//
	//        - Need to put something to buffer after last encoded log entry added?
	//          (see AddAfter() method).
	//
	//        - Need to put something to buffer between encoded log entries?
	//          I mean between all of them, excluding before first and after last.
	//          (see AddBetween() method).
	//
	// 5. Flushing data each N time intervals.
	//    Each worker flushes accumulated data on demand (the internal buffer is full)
	//    or at the timeout you may set (see SetWorkerAutoFlushDelay() method).
	//
	// 6. Slow down when service unavailable or network error.
	//    If request sending at the some worker is failed
	//    (it's only network or log service issue if you configure writer well),
	//    the usage of HTTP API service will slow down, logging what happens
	//    and why and saving (not dropping) your logs to the internal buffer.
	//
	//    They will send when connection will be restored.
	//    You may set the deferred log entries pack buffer capacity
	//    (pack, because each worker builds a pack of log entries, remember? see p.3)
	//    (see SetDeferredBufferCap() method).
	//
	// 7. Graceful shutdown.
	//    Of course, if you're familiar of ekadeath package. If you're not yet,
	//    it's time to: https://github.com/qioalice/ekago/ekadeath .
	//
	//    When you calling ekadeath.Die(), ekadeath.Exit() or writing a log
	//    with the level that marked as fatal, you won't lost aggregated logs!
	//
	//    The package will sends the rest of logs for the last time for you.
	//    And you do not need to do something for that.
	//
	// 8. Auto-initialization:
	//    You do need to call methods like Start() or something like that.
	//    You may check whether you configuration is valid calling Ping() method,
	//    but it's not necessary.
	//
	//    Just call all configuration methods with the chaining style and pass
	//    CI_WriterHttp object to the CommonIntegrator's WriteTo() method
	//    (or to your own logging integrator) and there is!
	//    The CI_WriterHttp will be initialized at the first Write() call.
	//
	// 9. Configurable to use any service.
	//    It's a very customizable type using which you may stream your logs safely
	//    to the log aggregation services like:
	//         - DataDog: https://www.datadoghq.com/ : UseProviderDataDog(),
	//         - Rollbar: https://rollbar.com/       : UseProviderRollbar(),
	//         - GrayLog: https://www.graylog.org/   : UseProviderGrayLog(),
	//         - Sentry:  https://sentry.io/         : UseProviderSentry(),
	//         etc.
	//
	//     For those services, CI_WriterHttp has a methods (3rd column)
	//     that makes it easy to configure it for, but you may configure it manually.
	//
	//     If you want to manually set HTTP service use UseProviderManual() method,
	//     providing a fasthttp's Request initializer for your service.
	//
	// 10. Fast.
	//     Uses fasthttp ( https://github.com/valyala/fasthttp ) under the hood,
	//     as http client. Pools, reusing, caching, optimizations. All you need.
	//
	// --------
	//
	// WARNING!
	// DO NOT CALL Write() or Ping() METHODS UNTIL YOU FINISH ALL PREPARATIONS!
	// DO NOT PASS WRITER TO THE CommonIntegrator's WriteTo() METHOD UNTIL
	// YOU FINISH ALL PREPARATIONS!
	// IF YOU DO, THE CHANGES WILL NOT BE SAVED!
	//
	// WARNING! PANIC CAUTION!
	// YOU MUST SET THE LOG SERVICE YOU WANT TO WRITE LOG ENTRIES TO.
	// CHOOSE A PREDEFINED OR USE YOUR OWN.
	// IF YOU DO NOT DO THAT, THE INITIALIZATION WILL PANIC!
	//
	CI_WriterHttp struct {

		// Has getter or/and setter

		providerInitializer func(req *fasthttp.Request)

		entriesBufferLen uint16
		deferredEntriesBufferLen uint16

		workerNum uint16
		workerEntriesBufferLen uint16
		workerFlushDelay time.Duration

		dataBefore []byte
		dataAfter []byte
		dataBetween []byte

		workerFlushDeferredPerIter uint16

		// Internal parts

		casInitStatus int32
		slowInit sync.Mutex

		beenPinged bool

		ctx context.Context
		cancelFunc context.CancelFunc

		workersWg sync.WaitGroup

		workerTickers []*time.Ticker

		// This channel will never be closed.
		entries chan []byte
		entriesPackDeferred chan *bytes.Buffer

		entriesCompletelyLostCounter uint64

		c fasthttp.Client
	}
)

//noinspection GoSnakeCaseUsage
const (
	DATADOG_ADDR_US = "https://http-intake.logs.datadoghq.com/v1/input"
	DATADOG_ADDR_EU = "https://http-intake.logs.datadoghq.eu/v1/input"
)

var (
	ErrWriterIsNil = fmt.Errorf("CI_WriterHttp: writer is nil (not initialized)")
	ErrWriterDisabled = fmt.Errorf("CI_WriterHttp: writer is disabled (stopped)")
	ErrWriterBufferFull = fmt.Errorf("CI_WriterHttp: writer's buffer is full")
)

// UseProviderManual is a log service provider manual configurator.
// You MUST specify a callback that will set-up an HTTP request for desired provider.
//
// Nil safe. There is no-op if CI_WriterHttp already initialized.
func (dw *CI_WriterHttp) UseProviderManual(cb func(req *fasthttp.Request)) *CI_WriterHttp {
	return dw.configure(func(dw *CI_WriterHttp) {
		dw.providerInitializer = cb
	})
}

// UseProviderDataDog setups CI_WriterHttp for DataDog log service provider
// ( https://www.datadoghq.com/ ).
//
// You MUST specify 'addr' as desired DataDog's HTTP addr (you may use predefined
// constants DATADOG_ADDR_US, DATADOG_ADDR_EU) or use your own and DataDog
// service's token as 'token'.
//
// Nil safe. There is no-op if CI_WriterHttp already initialized.
func (dw *CI_WriterHttp) UseProviderDataDog(addr, token string) *CI_WriterHttp {
	return dw.UseProviderManual(func(req *fasthttp.Request) {
		req.Header.SetContentType("application/json")
		req.SetRequestURI(addr)
		req.Header.Set("DD-API-KEY", token)
	})
}

//
func (dw *CI_WriterHttp) SetBufferCap(cap uint16) *CI_WriterHttp {
	return dw.configure(func(dw *CI_WriterHttp) {
		dw.entriesBufferLen = cap
	})
}

//
func (dw *CI_WriterHttp) SetWorkerBufferCap(cap uint16) *CI_WriterHttp {
	return dw.configure(func(dw *CI_WriterHttp) {
		dw.workerEntriesBufferLen = cap
	})
}

//
func (dw *CI_WriterHttp) SetDeferredBufferCap(cap uint16) *CI_WriterHttp {
	return dw.configure(func(dw *CI_WriterHttp) {
		dw.deferredEntriesBufferLen = cap
	})
}

//
func (dw *CI_WriterHttp) SetWorkersNum(num uint16) *CI_WriterHttp {
	return dw.configure(func(dw *CI_WriterHttp) {
		dw.workerNum = num
	})
}

//
func (dw *CI_WriterHttp) SetWorkerAutoFlushDelay(delay time.Duration) *CI_WriterHttp {
	return dw.configure(func(dw *CI_WriterHttp) {
		dw.workerFlushDelay = delay
	})
}

// AddBefore sets the data that will be added to the encoded entries pack's buffer
// before the first encoded entry is added.
//
// Take a look:
// If your log service provider accepts many records as JSON array of objects,
// 'data' is "[" (JSON beginning of array char).
//
// Nil safe. There is no-op if CI_WriterHttp already initialized.
func (dw *CI_WriterHttp) AddBefore(data []byte) *CI_WriterHttp {
	return dw.configure(func(dw *CI_WriterHttp) {
		dw.dataBefore = data
	})
}

// AddAfter sets the data that will be added to the encoded entries pack's buffer
// after the last encoded entry is added.
//
// Take a look:
// If your log service provider accepts many records as JSON array of objects,
// 'data' is "]" (JSON ending of array char).
//
// Nil safe. There is no-op if CI_WriterHttp already initialized.
func (dw *CI_WriterHttp) AddAfter(data []byte) *CI_WriterHttp {
	return dw.configure(func(dw *CI_WriterHttp) {
		dw.dataAfter = data
	})
}

// AddBetween sets the data that will be added to the encoded entries pack's buffer
// between encoded log entries (but neither before first nor after last).
//
// Take a look:
// If your log service provider accepts many records as JSON array of objects,
// 'data' is "," (JSON separator of objects inside an array).
//
// Nil safe. There is no-op if CI_WriterHttp already initialized.
func (dw *CI_WriterHttp) AddBetween(data []byte) *CI_WriterHttp {
	return dw.configure(func(dw *CI_WriterHttp) {
		dw.dataBetween = data
	})
}

// Ping checks whether provider settings are correct and connection can be established.
//
// Keep in mind, that if you do not do ping by yourself, the CI_WriterHttp will try
// to ping when Write() will be called first time (initialization).
// But if you call method Ping(), there will be no internal Ping() call
// at the initialization.
// You may prefer explicit Ping() call, because here you may specify some callbacks
// 'cb' that are applied to the HTTP request before it will be send.
//
// WARNING!
// If you call Ping() by yourself and it returns an error, you may change something
// (provider settings, etc), and then try to ping again.
// But if Ping() is called at the initialization, the CI_WriterHttp can not be recover.
func (dw *CI_WriterHttp) Ping(cb ...func(req *fasthttp.Request)) *ekaerr.Error {
	switch {

	case dw == nil:
		return ekaerr.IllegalState.
			New("CI_WriterHttp: writer is nil (not initialized)").
			Throw()
	}

	return dw.ping(true, cb)
}

// Write sends 'p' to the internal entries being processed buffer and returns
// len(p) and nil if 'p' has been successfully queued.
//
// Initializes CI_WriterHttp object if it's not. If initialization once failed,
// the CI_WriterHttp can not be used anymore.
//
// Returned errors:
// - nil: OK, 'p' has been queued.
// - ErrWriterIsNil: CI_WriterHttp receiver is nil.
// - ErrWriterDisabled: CI_WriterHttp is stopped and will never start again.
// - ErrWriterBufferFull: Internal CI_WriterHttp's buffer of processed entries
//   is full. Next time set bigger buffer's length using SetBufferCap().
func (dw *CI_WriterHttp) Write(p []byte) (n int, err error) {
	switch {

	case dw == nil:
		return -1, ErrWriterIsNil

	case len(p) == 0:
		return 0, nil

	case !dw.canWrite():
		return -1, ErrWriterDisabled
	}

	select {

	case dw.entries <- p:
		return len(p), nil

	default:
		atomic.AddUint64(&dw.entriesCompletelyLostCounter, 1)
		return -1, ErrWriterBufferFull
	}
}