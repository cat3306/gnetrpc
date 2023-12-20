package gnetrpc

import (
	"github.com/panjf2000/gnet/v2"
	"github.com/panjf2000/gnet/v2/pkg/logging"
	"time"
)

type OptionFn func(*serverOption)

func WithPrintRegisteredMethod() OptionFn {
	return func(s *serverOption) {
		s.printMethod = true
	}
}
func WithDefaultService() OptionFn {
	return func(s *serverOption) {
		s.defaultService = true
	}
}
func WithMulticore(multicore bool) OptionFn {
	return func(s *serverOption) {
		s.gnetOptions.Multicore = multicore
	}
}
func WithMainGoroutineChannelCap(cap int) OptionFn {
	return func(s *serverOption) {
		s.mainGoroutineChannelCap = cap
	}
}
func WithClientAsyncMode() OptionFn {
	return func(s *serverOption) {
		s.clientAsyncMode = true
	}
}
func WithLockOSThread(lockOSThread bool) OptionFn {
	return func(s *serverOption) {
		s.gnetOptions.LockOSThread = lockOSThread
	}
}

func WithReadBufferCap(readBufferCap int) OptionFn {
	return func(s *serverOption) {
		s.gnetOptions.ReadBufferCap = readBufferCap
	}
}

func WithWriteBufferCap(writeBufferCap int) OptionFn {
	return func(s *serverOption) {
		s.gnetOptions.WriteBufferCap = writeBufferCap
	}
}

func WithLoadBalancing(lb int) OptionFn {
	return func(s *serverOption) {
		s.gnetOptions.LB = gnet.LoadBalancing(lb)
	}
}

func WithNumEventLoop(numEventLoop int) OptionFn {
	return func(s *serverOption) {
		s.gnetOptions.NumEventLoop = numEventLoop
	}
}

func WithReusePort(reusePort bool) OptionFn {
	return func(s *serverOption) {
		s.gnetOptions.ReusePort = reusePort
	}
}

func WithReuseAddr(reuseAddr bool) OptionFn {
	return func(s *serverOption) {
		s.gnetOptions.ReuseAddr = reuseAddr
	}
}

func WithTCPKeepAlive(tcpKeepAlive time.Duration) OptionFn {
	return func(s *serverOption) {
		s.gnetOptions.TCPKeepAlive = tcpKeepAlive
	}
}

func WithTCPNoDelay(tcpNoDelay int) OptionFn {
	return func(s *serverOption) {
		s.gnetOptions.TCPNoDelay = gnet.TCPSocketOpt(tcpNoDelay)
	}
}

func WithSocketRecvBuffer(recvBuf int) OptionFn {
	return func(s *serverOption) {
		s.gnetOptions.SocketRecvBuffer = recvBuf
	}
}

func WithSocketSendBuffer(sendBuf int) OptionFn {
	return func(s *serverOption) {
		s.gnetOptions.SocketSendBuffer = sendBuf
	}
}

func WithTicker(ticker bool) OptionFn {
	return func(s *serverOption) {
		s.gnetOptions.Ticker = ticker
	}
}

func WithLogPath(fileName string) OptionFn {
	return func(s *serverOption) {
		s.gnetOptions.LogPath = fileName
	}
}

func WithLogLevel(lvl logging.Level) OptionFn {
	return func(s *serverOption) {
		s.gnetOptions.LogLevel = lvl
	}
}

func WithLogger(logger logging.Logger) OptionFn {
	return func(s *serverOption) {
		s.gnetOptions.Logger = logger
	}
}

func WithMulticastInterfaceIndex(idx int) OptionFn {
	return func(s *serverOption) {
		s.gnetOptions.MulticastInterfaceIndex = idx
	}
}

func WithAntExpiryDuration(expiryDuration time.Duration) OptionFn {
	return func(s *serverOption) {
		s.antOption.ExpiryDuration = expiryDuration
	}
}

func WithPreAlloc(preAlloc bool) OptionFn {
	return func(s *serverOption) {
		s.antOption.PreAlloc = preAlloc
	}
}

func WithMaxBlockingTasks(maxBlockingTasks int) OptionFn {
	return func(s *serverOption) {
		s.antOption.MaxBlockingTasks = maxBlockingTasks
	}
}

func WithNonblocking(nonblocking bool) OptionFn {
	return func(s *serverOption) {
		s.antOption.Nonblocking = nonblocking
	}
}

func WithPanicHandler(panicHandler func(interface{})) OptionFn {
	return func(s *serverOption) {
		s.antOption.PanicHandler = panicHandler
	}
}

//func WithAntLogger(logger Logger) Option {
//	return func(opts *Options) {
//		opts.Logger = logger
//	}
//}

func WithDisablePurge(disable bool) OptionFn {
	return func(s *serverOption) {
		s.antOption.DisablePurge = disable
	}
}
