package main

import (
	"ir.skhf/awesomeHome/components/broadcast"
	"ir.skhf/awesomeHome/components/command"
	"ir.skhf/awesomeHome/components/info"
	"ir.skhf/awesomeHome/components/process"
	"ir.skhf/awesomeHome/components/update"
	"ir.skhf/awesomeHome/components/webcam"
	"ir.skhf/awesomeHome/config"
	"ir.skhf/awesomeHome/server"
)

func main() {
	conf := config.NewAwesomeHomeConfig("", "")
	processService := process.NewProcessService()
	broadcastService := broadcast.NewBroadcastService(conf.BroadcastPacketMaxSize, conf.BroadcastPort)
	appInfoService := info.NewAppInfoService(broadcastService, conf.BroadcastSendTime)
	commandService := command.NewCommandService(processService)
	_ = update.NewSelfUpdateService()
	webcamService := webcam.NewWebcamService()

	broadcastService.Init()
	appInfoService.Init()
	webcamService.Init(nil, nil, nil, nil)

	httpServer := server.NewHttpServerService(conf.HttpServerPort)

	httpServer.Register("/", appInfoService.HandleHttp)
	httpServer.Register("/command", commandService.HandleHttp)
	httpServer.Register("/webcam", webcamService.HandleHttp)

	httpServer.Start()
}
