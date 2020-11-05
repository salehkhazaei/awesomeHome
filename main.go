package main

import (
	"ir.skhf/awesomeHome/broadcast"
	"ir.skhf/awesomeHome/command"
	"ir.skhf/awesomeHome/config"
	"ir.skhf/awesomeHome/info"
	"ir.skhf/awesomeHome/process"
	"ir.skhf/awesomeHome/update"
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
}
