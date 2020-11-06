package main

import (
	"fmt"
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
	fmt.Println("starting awesome home")

	conf := config.NewAwesomeHomeConfig("", "")
	processService := process.NewProcessService()
	broadcastService := broadcast.NewBroadcastService(conf.BroadcastPacketMaxSize, conf.BroadcastPort)
	appInfoService := info.NewAppInfoService(broadcastService, conf.BroadcastSendTime)
	commandService := command.NewCommandService(processService)
	_ = update.NewSelfUpdateService()
	webcamService := webcam.NewWebcamService()

	fmt.Println("services created")

	err := broadcastService.Init()
	if err != nil {
		panic(err)
	}

	fmt.Println("broadcast service initialized")

	appInfoService.Init()

	fmt.Println("app info service initialized")

	err = webcamService.Init("/dev/video0", "", "", false)
	if err != nil {
		panic(err)
	}

	fmt.Println("webcam service initialized")

	httpServer := server.NewHttpServerService(conf.HttpServerPort)

	httpServer.Register("/", appInfoService.HandleHttp)
	httpServer.Register("/command", commandService.HandleHttp)
	httpServer.Register("/webcam", webcamService.HandleHttp)

	fmt.Println("starting http server")

	httpServer.Start()
}
