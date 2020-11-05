package main

import (
	"ir.skhf/awesomeHome/components/broadcast"
	"ir.skhf/awesomeHome/components/command"
	"ir.skhf/awesomeHome/components/info"
	"ir.skhf/awesomeHome/components/process"
	"ir.skhf/awesomeHome/components/update"
	"ir.skhf/awesomeHome/config"
	"ir.skhf/awesomeHome/server"
)

func main() {
	conf := config.NewAwesomeHomeConfig("", "")
	processService := process.NewProcessService()
	broadcastService := broadcast.NewBroadcastService(conf.BroadcastPacketMaxSize, conf.BroadcastPort)
	appInfoService := info.NewAppInfoService(broadcastService, conf.BroadcastSendTime)
	_ = command.NewCommandService(processService)
	_ = update.NewSelfUpdateService()

	broadcastService.Init()
	appInfoService.Init()

	httpServer := server.NewHttpServerService(conf.HttpServerPort)

	httpServer.Register("/", appInfoService.HandleHttp)

	httpServer.Start()
}
