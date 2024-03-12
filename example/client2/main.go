package main

import (
	"flag"
	"reflect"
	"sync"
	"sync/atomic"
	"time"

	"github.com/cat3306/gnetrpc"
	"github.com/cat3306/gnetrpc/protocol"
	"github.com/cat3306/gnetrpc/rpclog"
	"github.com/cat3306/gnetrpc/share"
)

type Builtin struct {
}

func (b *Builtin) Init(v ...interface{}) gnetrpc.IService {
	return b
}
func (b *Builtin) Alias() string {
	return ""
}
func (b *Builtin) Benchmark(ctx *protocol.Context) {

	atomic.AddUint64(&recvCnt, 1)
}

var (
	concurrency = flag.Int("c", 10, "concurrency")
	total       = flag.Int("n", 50000, "total requests for all clients")
	host        = flag.String("s", "127.0.0.1:7898", "server ip and port")
	rate        = flag.Int("r", 0, "throughputs")
	recvCnt     uint64
)

func main() {
	flag.Parse()
	wg := sync.WaitGroup{}
	totalT := time.Now().UnixNano()
	n := *concurrency
	m := *total / n
	wg.Add(*concurrency)
	for i := 0; i < *concurrency; i++ {

		go singleClient(&wg, m)

	}
	wg.Wait()
	for {
		if recvCnt == uint64(*total) {
			break
		}
	}
	totalT = time.Now().UnixNano() - totalT
	totalT = totalT / 1000000
	rpclog.Infof("took %d ms for %d requests", totalT, *total)
	rpclog.Infof("sent requests:%d\n", *total)
	rpclog.Infof("recv:%d\n", atomic.LoadUint64(&recvCnt))
	rpclog.Infof("throughput  (TPS)    : %d\n", int64(*total)*1000/totalT)
}

func singleClient(wg *sync.WaitGroup, m int) {
	defer wg.Done()
	client, err := gnetrpc.NewClient(*host, "tcp").
		Register(
			new(Builtin),
		).Run()
	if err != nil {
		panic(err)
	}
	for i := 0; i < m; i++ {
		err := ClientBenchmark(client)
		if err != nil {
			rpclog.Errorf("err:%s", err.Error())
		}
		//time.Sleep(time.Millisecond)
	}

}

func ClientBenchmark(client *gnetrpc.Client) error {
	args := prepareArgs()
	//rpclog.Infof(args.Field1)
	return client.Call("Builtin", "Benchmark", map[string]string{
		share.AuthKey: "鸳鸯擦，鸳鸯体，你爱我，我爱你",
	}, protocol.Json, args)
}

func prepareArgs() *BenchmarkMessage {
	b := true
	var i int32 = 100000
	var s = "许多往事在眼前一幕一幕，变的那麼模糊"

	var args BenchmarkMessage

	v := reflect.ValueOf(&args).Elem()
	num := v.NumField()
	for k := 0; k < num; k++ {
		field := v.Field(k)
		if field.Type().Kind() == reflect.Ptr {
			switch v.Field(k).Type().Elem().Kind() {
			case reflect.Int, reflect.Int32, reflect.Int64:
				field.Set(reflect.ValueOf(&i))
			case reflect.Bool:
				field.Set(reflect.ValueOf(&b))
			case reflect.String:
				field.Set(reflect.ValueOf(&s))
			}
		} else {
			switch field.Kind() {
			case reflect.Int, reflect.Int32, reflect.Int64:
				field.SetInt(100000)
			case reflect.Bool:
				field.SetBool(true)
			case reflect.String:
				field.SetString(s)
			}
		}

	}
	return &args
}
