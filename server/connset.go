package server

import (
	"github.com/panjf2000/gnet/v2"
	"sync"
)

type ConnSet struct {
	connections map[int]gnet.Conn
	locker      sync.RWMutex
}
