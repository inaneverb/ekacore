package ekatime_test

import (
	"fmt"
	"math/rand"
	"sync/atomic"
	"testing"
	"time"

	"github.com/qioalice/ekago/v3/ekadeath"
	"github.com/qioalice/ekago/v3/ekatime"
)

func TestFoo(t *testing.T) {
	rand.Seed(time.Now().UnixNano())

	var (
		_1calls, _2calls, _3calls uint32
	)

	ekatime.OnceInMinute.Call(true, func(ts ekatime.Timestamp) {
		atomic.AddUint32(&_1calls, 1)
		st := time.Duration(rand.Intn(3000)) * time.Millisecond
		fmt.Println("Call 1", ts.String(), "sleep time", st)
		time.Sleep(st)
	})
	ekatime.OnceInMinute.Call(true, func(ts ekatime.Timestamp) {
		atomic.AddUint32(&_2calls, 1)
		st := time.Duration(rand.Intn(3000)) * time.Millisecond
		fmt.Println("Call 2", ts.String(), "sleep time", st)
		time.Sleep(st)
	})
	ekatime.OnceInMinute.Call(true, func(ts ekatime.Timestamp) {
		atomic.AddUint32(&_3calls, 1)
		st := time.Duration(rand.Intn(3000)) * time.Millisecond
		fmt.Println("Call 3", ts.String(), "sleep time", st)
		time.Sleep(st)
	})

	i := 0
	for range time.Tick(30 * time.Second) {
		i++
		fmt.Printf("Stat: [1]: %d, [2]: %d, [3]: %d\n",
			atomic.LoadUint32(&_1calls),
			atomic.LoadUint32(&_2calls),
			atomic.LoadUint32(&_3calls))

		if i == 10 {
			ekadeath.Exit()
		}
	}
}
