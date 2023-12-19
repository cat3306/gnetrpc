package gnetrpc

import (
	"fmt"
	"github.com/cat3306/gnetrpc/rpclog"
	"github.com/valyala/bytebufferpool"
	"sync"

	"github.com/lithammer/shortuuid/v4"
	"github.com/panjf2000/gnet/v2"
)

type ConnMatrix struct {
	locker  sync.RWMutex
	connMap map[string]gnet.Conn
}

func NewConnMatrix() *ConnMatrix {
	return &ConnMatrix{
		connMap: make(map[string]gnet.Conn),
	}
}
func (c *ConnMatrix) GenId(fd int) string {
	return fmt.Sprintf("%s@%d", shortuuid.New(), fd)
}
func (c *ConnMatrix) Add(conn gnet.Conn) {
	c.locker.Lock()
	defer c.locker.Unlock()
	conn.SetId(c.GenId(conn.Fd()))
	c.connMap[conn.Id()] = conn
}
func (c *ConnMatrix) Remove(id string) {
	c.locker.Lock()
	defer c.locker.Unlock()
	delete(c.connMap, id)
}

func (c *ConnMatrix) Broadcast(buffer *bytebufferpool.ByteBuffer) {
	c.locker.RLock()
	tmpList := make([]gnet.Conn, 0, len(c.connMap))
	for _, v := range c.connMap {
		tmpList = append(tmpList, v)
	}
	c.locker.RUnlock()
	wg := sync.WaitGroup{}
	wg.Add(len(tmpList))
	fmt.Println(len(tmpList))
	for _, v := range tmpList {
		err := v.AsyncWrite(buffer.Bytes(), func(c gnet.Conn, err error) error {
			wg.Done()
			return nil
		})
		if err != nil {
			rpclog.Errorf("Broadcast err:%s", err.Error())
		}
	}
	wg.Wait()
	bytebufferpool.Put(buffer)
}
func (c *ConnMatrix) BroadcastExceptOne(buffer *bytebufferpool.ByteBuffer, id string) {
	c.locker.RLock()
	tmpList := make([]gnet.Conn, 0, len(c.connMap))
	for _, v := range c.connMap {
		if v.Id() == id {
			continue
		}
		tmpList = append(tmpList, v)
	}
	c.locker.RUnlock()
	wg := sync.WaitGroup{}
	wg.Add(len(tmpList))
	for _, v := range tmpList {
		err := v.AsyncWrite(buffer.Bytes(), func(c gnet.Conn, err error) error {
			wg.Done()
			return nil
		})
		if err != nil {
			rpclog.Errorf("Broadcast err:%s", err.Error())
		}
	}
	wg.Wait()
	bytebufferpool.Put(buffer)
}

func (c *ConnMatrix) SendToOne(buffer *bytebufferpool.ByteBuffer, id string) {
	c.locker.RLock()
	conn, ok := c.connMap[id]
	c.locker.RUnlock()
	if !ok {
		rpclog.Errorf("not found conn,id:%d", id)
		return
	}
	err := conn.AsyncWrite(buffer.Bytes(), func(c gnet.Conn, err error) error {
		bytebufferpool.Put(buffer)
		return nil
	})
	if err != nil {
		rpclog.Errorf("conn.Write err:%s", err.Error())
	}
}

func (c *ConnMatrix) BroadcastSomeone(buffer *bytebufferpool.ByteBuffer, ids []string) {
	tmpList := make([]gnet.Conn, 0, len(ids))
	c.locker.RLock()
	for _, id := range ids {
		conn, ok := c.connMap[id]
		if ok {
			tmpList = append(tmpList, conn)
		} else {
			rpclog.Warnf("SendToSome not found id:%s", id)
		}
	}
	c.locker.RUnlock()
	wg := sync.WaitGroup{}
	wg.Add(len(tmpList))
	for _, conn := range tmpList {
		err := conn.AsyncWrite(buffer.Bytes(), func(c gnet.Conn, err error) error {
			wg.Done()
			return nil
		})
		if err != nil {
			rpclog.Errorf("SendToSome Write err:%s", err.Error())
		}
	}
	wg.Wait()
	bytebufferpool.Put(buffer)
}
