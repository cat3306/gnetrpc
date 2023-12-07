package gnetrpc

import (
	"sync"

	"github.com/panjf2000/gnet/v2"
)

type connMatrix struct {
	locker  sync.RWMutex
	connMap map[int]gnet.Conn
}

func newConnMatrix() *connMatrix {
	return &connMatrix{
		connMap: make(map[int]gnet.Conn),
	}
}
func (c *connMatrix) Add(conn gnet.Conn) {
	c.locker.Lock()
	defer c.locker.Unlock()
	c.connMap[conn.Fd()] = conn
}
func (c *connMatrix) Remove(fd int) {
	c.locker.Lock()
	defer c.locker.Unlock()
	delete(c.connMap, fd)
}
