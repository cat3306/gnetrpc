package server

import (
	"github.com/cat3306/gnetrpc"
	"testing"
)

func TestServer(t *testing.T) {
	s := NewServer(WithMulticore(true))
	err := s.Run(gnetrpc.TcpNetwork, ":7898")
	t.Fatalf(err.Error())
}
