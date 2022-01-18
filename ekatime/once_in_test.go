package ekatime_test

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/qioalice/ekago/v3/ekatime"
)

func TestFoo(t *testing.T) {
	rand.Seed(time.Now().UnixNano())

	//goland:noinspection GoSnakeCaseUsage
	const WAIT_MS_MAX = 500

	//goland:noinspection GoSnakeCaseUsage
	const (
		INV_NOW_1 = true
		INV_NOW_2 = true
		INV_NOW_3 = true
		INV_NOW_4 = false
		INV_NOW_5 = false
		INV_NOW_6 = true
		INV_NOW_7 = false
		INV_NOW_8 = false
	)

	onceIn := &ekatime.OnceInMinute

	onceIn.Call(INV_NOW_1, func(ts ekatime.Timestamp) {
		st := time.Duration(rand.Intn(WAIT_MS_MAX)) * time.Millisecond
		time.Sleep(st)
		fmt.Println("Call 1", ts.String(), "sleep time", st)
	})
	onceIn.Call(INV_NOW_2, func(ts ekatime.Timestamp) {
		st := time.Duration(rand.Intn(WAIT_MS_MAX)) * time.Millisecond
		time.Sleep(st)
		fmt.Println("Call 2", ts.String(), "sleep time", st)
	})
	onceIn.Call(INV_NOW_3, func(ts ekatime.Timestamp) {
		st := time.Duration(rand.Intn(WAIT_MS_MAX)) * time.Millisecond
		time.Sleep(st)
		fmt.Println("Call 3", ts.String(), "sleep time", st)
	})
	onceIn.Call(INV_NOW_4, func(ts ekatime.Timestamp) {
		st := time.Duration(rand.Intn(WAIT_MS_MAX)) * time.Millisecond
		time.Sleep(st)
		fmt.Println("Call 4", ts.String(), "sleep time", st)
	})
	onceIn.Call(INV_NOW_5, func(ts ekatime.Timestamp) {
		st := time.Duration(rand.Intn(WAIT_MS_MAX)) * time.Millisecond
		time.Sleep(st)
		fmt.Println("Call 5", ts.String(), "sleep time", st)
	})
	onceIn.Call(INV_NOW_6, func(ts ekatime.Timestamp) {
		st := time.Duration(rand.Intn(WAIT_MS_MAX)) * time.Millisecond
		time.Sleep(st)
		fmt.Println("Call 6", ts.String(), "sleep time", st)
	})
	onceIn.Call(INV_NOW_7, func(ts ekatime.Timestamp) {
		st := time.Duration(rand.Intn(WAIT_MS_MAX)) * time.Millisecond
		time.Sleep(st)
		fmt.Println("Call 7", ts.String(), "sleep time", st)
	})
	onceIn.Call(INV_NOW_8, func(ts ekatime.Timestamp) {
		st := time.Duration(rand.Intn(WAIT_MS_MAX)) * time.Millisecond
		time.Sleep(st)
		fmt.Println("Call FINAL", ts.String(), "sleep time", st)
	})

	select {}
}
