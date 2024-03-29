package util

import (
	"sync"

	"github.com/panjf2000/gnet/v2"
)

func GetSyncMapFromConn(conn gnet.Conn) *sync.Map {
	ctx := conn.Context()
	if ctx == nil {
		m := &sync.Map{}
		conn.SetContext(m)
		return m
	}
	m := ctx.(*sync.Map)
	return m
}
