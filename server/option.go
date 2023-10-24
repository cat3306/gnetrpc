package server

import (
	"github.com/panjf2000/gnet/v2"
	"github.com/panjf2000/gnet/v2/pkg/logging"
	"time"
)

type OptionFn func(*Server)

func WithPrintRegisteredMethod() OptionFn {
	return func(s *Server) {
		s.option.printMethod = true
	}
}

func WithMulticore(multicore bool) OptionFn {
	return func(s *Server) {
		s.option.gnetOptions.Multicore = multicore
	}
}

func WithLockOSThread(lockOSThread bool) OptionFn {
	return func(s *Server) {
		s.option.gnetOptions.LockOSThread = lockOSThread
	}
}

func WithReadBufferCap(readBufferCap int) OptionFn {
	return func(s *Server) {
		s.option.gnetOptions.ReadBufferCap = readBufferCap
	}
}

func WithWriteBufferCap(writeBufferCap int) OptionFn {
	return func(s *Server) {
		s.option.gnetOptions.WriteBufferCap = writeBufferCap
	}
}

func WithLoadBalancing(lb int) OptionFn {
	return func(s *Server) {
		s.option.gnetOptions.LB = gnet.LoadBalancing(lb)
	}
}

func WithNumEventLoop(numEventLoop int) OptionFn {
	return func(s *Server) {
		s.option.gnetOptions.NumEventLoop = numEventLoop
	}
}

func WithReusePort(reusePort bool) OptionFn {
	return func(s *Server) {
		s.option.gnetOptions.ReusePort = reusePort
	}
}

func WithReuseAddr(reuseAddr bool) OptionFn {
	return func(s *Server) {
		s.option.gnetOptions.ReuseAddr = reuseAddr
	}
}

func WithTCPKeepAlive(tcpKeepAlive time.Duration) OptionFn {
	return func(s *Server) {
		s.option.gnetOptions.TCPKeepAlive = tcpKeepAlive
	}
}

func WithTCPNoDelay(tcpNoDelay int) OptionFn {
	return func(s *Server) {
		s.option.gnetOptions.TCPNoDelay = gnet.TCPSocketOpt(tcpNoDelay)
	}
}

func WithSocketRecvBuffer(recvBuf int) OptionFn {
	return func(s *Server) {
		s.option.gnetOptions.SocketRecvBuffer = recvBuf
	}
}

func WithSocketSendBuffer(sendBuf int) OptionFn {
	return func(s *Server) {
		s.option.gnetOptions.SocketSendBuffer = sendBuf
	}
}

func WithTicker(ticker bool) OptionFn {
	return func(s *Server) {
		s.option.gnetOptions.Ticker = ticker
	}
}

func WithLogPath(fileName string) OptionFn {
	return func(s *Server) {
		s.option.gnetOptions.LogPath = fileName
	}
}

func WithLogLevel(lvl logging.Level) OptionFn {
	return func(s *Server) {
		s.option.gnetOptions.LogLevel = lvl
	}
}

func WithLogger(logger logging.Logger) OptionFn {
	return func(s *Server) {
		s.option.gnetOptions.Logger = logger
	}
}

func WithMulticastInterfaceIndex(idx int) OptionFn {
	return func(s *Server) {
		s.option.gnetOptions.MulticastInterfaceIndex = idx
	}
}
