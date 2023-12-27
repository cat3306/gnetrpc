package main

import (
	"flag"
	"fmt"
	"github.com/cat3306/gnetrpc"
	"github.com/cat3306/gnetrpc/example/gameserver/conf"
	"github.com/cat3306/gnetrpc/example/gameserver/service"
	"github.com/cat3306/gnetrpc/example/gameserver/thirdmodule"
	"github.com/cat3306/gnetrpc/plugin"
)

var configFile = flag.String("f", "conf/conf.json", "the config file")

func main() {

	err := conf.Init(*configFile)
	if err != nil {
		panic(err)
	}
	thirdmodule.Init()
	s := gnetrpc.NewServer(
		gnetrpc.WithMulticore(true),
		gnetrpc.WithPrintRegisteredMethod(),
		gnetrpc.WithDefaultService(),
		gnetrpc.WithMainGoroutineChannelCap(10000),
	)
	//s.UseAuthFunc(func(ctx *protocol.Context, token string) error {
	//	if token != "鸳鸯擦，鸳鸯体，你爱我，我爱你" {
	//		return errors.New("你不爱我 !")
	//	}
	//	return nil
	//})
	s.AddPlugin(
		new(plugin.ConnectPlugin),
		new(plugin.ClosePlugin),
	)
	s.Register(
		new(service.Account),
		new(service.RoomMgr).Init(),
	)
	err = s.Run(gnetrpc.TcpNetwork, ":7898")
	fmt.Println(err)
}
