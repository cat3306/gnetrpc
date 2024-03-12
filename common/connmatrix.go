package common

import (
	"sync"
	"sync/atomic"

	"github.com/cat3306/gnetrpc/rpclog"
	"github.com/valyala/bytebufferpool"

	"github.com/panjf2000/gnet/v2"
)

type ConnMatrix struct {
	async    bool //async if true use sync.RWMutex
	connMap  map[string]gnet.Conn
	asyncMap sync.Map
	cnt      int64
}

// false:sync true:async
func NewConnMatrix(async bool) *ConnMatrix {
	matrix := ConnMatrix{
		connMap: make(map[string]gnet.Conn),
		async:   async,
	}
	if !async {
		matrix.connMap = make(map[string]gnet.Conn)
	}
	return &matrix
}
func (c *ConnMatrix) SetAsync(async bool) {
	c.async = async
}

func (c *ConnMatrix) Add(conn gnet.Conn) {
	if c.async {
		c.asyncMap.Store(conn.Id(), conn)
		atomic.AddInt64(&c.cnt, 1)
	} else {
		c.connMap[conn.Id()] = conn
	}

}
func (c *ConnMatrix) Remove(id string) {
	if c.async {
		atomic.AddInt64(&c.cnt, -1)
		c.asyncMap.Delete(id)
	} else {
		delete(c.connMap, id)
	}

}

//	func (c *ConnMatrix) RemoveAll(msg string) {
//		if c.sync {
//			c.locker.Lock()
//			defer c.locker.Unlock()
//		}
//		if c.async{
//			c.asyncMap.
//		}
//		for k, v := range c.connMap {
//			v.Close(msg)
//			delete(c.connMap, k)
//		}
//	}
func (c *ConnMatrix) Broadcast(buffer *bytebufferpool.ByteBuffer) {
	var err error
	if c.async {
		c.asyncMap.Range(func(key, value any) bool {
			conn := value.(gnet.Conn)
			_, err = conn.Write(buffer.Bytes())
			if err != nil {
				rpclog.Errorf("Broadcast err:%s", err.Error())
			}
			return true
		})
	} else {
		for _, conn := range c.connMap {
			_, err = conn.Write(buffer.Bytes())
			if err != nil {
				rpclog.Errorf("Broadcast err:%s", err.Error())
			}
		}
	}
	bytebufferpool.Put(buffer)
}
func (c *ConnMatrix) Len() int {
	if c.async {
		return int(atomic.LoadInt64(&c.cnt))
	} else {
		return len(c.connMap)
	}
}
func (c *ConnMatrix) BroadcastExceptOne(buffer *bytebufferpool.ByteBuffer, id string) {
	var err error
	if c.async {
		c.asyncMap.Range(func(key, value any) bool {
			conn := value.(gnet.Conn)
			idKey := key.(string)
			if idKey == id {
				return true
			}
			_, err = conn.Write(buffer.Bytes())
			if err != nil {
				rpclog.Errorf("Broadcast err:%s", err.Error())
			}
			return true
		})
	} else {
		for k, v := range c.connMap {
			if k == id {
				continue
			}
			_, err = v.Write(buffer.Bytes())
			if err != nil {
				rpclog.Errorf("Broadcast err:%s", err.Error())
			}
		}
	}
	bytebufferpool.Put(buffer)
}
func (c *ConnMatrix) SendToConn(buffer *bytebufferpool.ByteBuffer, conn gnet.Conn) {
	err := conn.AsyncWrite(buffer.Bytes(), func(c gnet.Conn, err error) error {
		bytebufferpool.Put(buffer)
		return nil
	})
	if err != nil {
		rpclog.Errorf("conn.Write err:%s", err.Error())
	}
}
func (c *ConnMatrix) SendToOne(buffer *bytebufferpool.ByteBuffer, id string) {

	if c.async {
		v, ok := c.asyncMap.Load(id)
		if !ok {
			rpclog.Warnf("SendToOne not found conn id:%s", id)
			return
		}
		conn := v.(gnet.Conn)
		c.SendToConn(buffer, conn)
	} else {
		conn, ok := c.connMap[id]
		if !ok {
			rpclog.Warnf("SendToOne not found conn id:%s", id)
			return
		}
		c.SendToConn(buffer, conn)
	}
}

func (c *ConnMatrix) BroadcastSomeone(buffer *bytebufferpool.ByteBuffer, ids []string) {
	if c.async {
		for _, id := range ids {
			v, ok := c.asyncMap.Load(id)
			if ok {
				rpclog.Warnf("BroadcastSomeone not found conn id:%s", id)
				continue
			}
			conn := v.(gnet.Conn)
			c.SendToConn(buffer, conn)
		}
	} else {
		for _, id := range ids {
			conn, ok := c.connMap[id]
			if ok {
				rpclog.Warnf("BroadcastSomeone not found conn id:%s", id)
				continue
			}
			c.SendToConn(buffer, conn)
		}
	}
}
