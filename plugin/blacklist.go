package plugin

import (
	"github.com/cat3306/gnetrpc"
	"github.com/cat3306/gnetrpc/rpclog"
	"github.com/panjf2000/gnet/v2"
	"net"
)

type BlacklistPlugin struct {
	blacklist     map[string]bool
	BlacklistMask []*net.IPNet // net.ParseCIDR("172.17.0.0/16") to get *net.IPNet
}

func (b *BlacklistPlugin) Add(ips ...string) *BlacklistPlugin {
	if b.blacklist == nil {
		b.blacklist = make(map[string]bool)
	}
	for _, ip := range ips {
		b.blacklist[ip] = true
	}
	return b
}
func (b *BlacklistPlugin) Type() gnetrpc.PluginType {
	return gnetrpc.PluginTypeOnOpen
}

func (b *BlacklistPlugin) OnDo(args ...interface{}) interface{} {
	conn := args[0].(gnet.Conn)
	ip, _, err := net.SplitHostPort(conn.RemoteAddr().String())
	if err != nil {
		return true
	}
	if b.blacklist[ip] {
		rpclog.Errorf("%s in black list refused service", conn.RemoteAddr().String())
		return false
	}

	remoteIP := net.ParseIP(ip)
	for _, mask := range b.BlacklistMask {
		if mask.Contains(remoteIP) {
			return false
		}
	}
	return true
}
