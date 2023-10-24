package util

import (
	"github.com/cat3306/gnetrpc/rpclog"
	"math"
	"reflect"
	"runtime"
	"runtime/debug"
	"time"
)

type PanicRepeatRunArgs struct {
	Sleep time.Duration
	Try   int
}

func runPanicLess(f func()) (panicLess bool) {
	defer func() {
		err := recover()
		panicLess = err == nil
		if err != nil {
			name := runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
			s := debug.Stack()
			rpclog.Errorf("%s err:%v,%s", name, err, s)
		}
	}()
	f()
	return
}

func PanicRepeatRun(f func(), args ...PanicRepeatRunArgs) {
	param := PanicRepeatRunArgs{
		Sleep: 0,
		Try:   6,
	}
	if len(args) != 0 {
		param = args[0]
	}
	if param.Try == 0 {
		param.Try = math.MaxInt64
	}
	total := param.Try
	for !runPanicLess(f) && param.Try >= 1 {
		if param.Sleep != 0 {
			time.Sleep(param.Sleep)
		}
		param.Try--
	}
	if param.Try == 0 {
		name := runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
		rpclog.Errorf("%s:finally failed,total:%d", name, total)
	}
}
