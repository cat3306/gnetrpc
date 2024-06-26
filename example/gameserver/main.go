package main

import (
	"flag"

	"github.com/cat3306/gnetrpc"
	"github.com/cat3306/gnetrpc/example/gameserver/conf"
	"github.com/cat3306/gnetrpc/example/gameserver/service"
	"github.com/cat3306/gnetrpc/example/gameserver/thirdmodule"
	"github.com/cat3306/gnetrpc/rpclog"
)

var configFile string

func main() {
	flag.StringVar(&configFile, "f", "conf/conf.json", "the config file")
	flag.Parse()
	err := conf.Init(configFile)
	if err != nil {
		panic(err)
	}
	thirdmodule.Init()
	s := gnetrpc.NewServer(
		gnetrpc.WithPrintRegisteredMethod(),
		gnetrpc.WithDefaultService(),
		gnetrpc.WithMainGoroutineChannelCap(10000),
		//gnetrpc.WithReusePort(true),
		gnetrpc.WithNumEventLoop(512),
		gnetrpc.WithReadBufferCap(1024*1024),
		gnetrpc.WithWriteBufferCap(1024*1024),
	)
	// s.UseAuthFunc(func(ctx *protocol.Context, token string) error {
	// 	if token != "" {
	// 		return errors.New("")
	// 	}
	// 	return nil
	// })
	s.AddPlugin(
		new(service.ConnOpenPlugin),
		new(service.ConnClosePlugin),
		new(service.ShutdownPlugin),
	)
	s.Register(
		new(service.Account).Init(),
		new(service.RoomMgr).Init(),
		new(service.GameMgr).Init(),
	)
	err = s.Run(gnetrpc.UdpNetwork, "0.0.0.0:7898")
	if err != nil {
		rpclog.Infof("run err:%s", err.Error())
	}
}
