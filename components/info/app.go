package info

import (
	"encoding/json"
	"fmt"
	"ir.skhf/awesomeHome/components/broadcast"
	"ir.skhf/awesomeHome/utils"
	"net/http"
	"time"
)

type AppInfoService struct {
	BroadcastService *broadcast.BroadcastService
	SendTime         time.Duration
}

func NewAppInfoService(
	broadcastService *broadcast.BroadcastService,
	sendTime time.Duration,
) *AppInfoService {
	return &AppInfoService{
		BroadcastService: broadcastService,
		SendTime:         sendTime,
	}
}

var Major int = 1
var Minor int = 0

type AppInfo struct {
	Version string
	IP      string
}

func (s *AppInfoService) Info() (string, error) {
	appInfo := AppInfo{
		Version: s.Version(),
		IP:      utils.GetOutboundIP().String(),
	}

	appInfoJson, err := json.Marshal(appInfo)
	if err != nil {
		return "", err
	}

	return string(appInfoJson), nil
}

func (s *AppInfoService) Version() string {
	return fmt.Sprintf("%d.%d", Major, Minor)
}

func (s *AppInfoService) Init() {
	go s.SendLoop()
}

func (s *AppInfoService) SendLoop() {
	for {
		jsonStr, err := s.Info()
		if err != nil {
			fmt.Printf("failed to build broadcast json due to %v\n", err)
			continue
		}

		err = s.BroadcastService.Broadcast([]byte(jsonStr))
		if err != nil {
			fmt.Printf("failed to send broadcast due to %v\n", err)
			continue
		}

		time.Sleep(s.SendTime)
	}
}

func (s *AppInfoService) HandleHttp(w http.ResponseWriter, r *http.Request) {
	jsonStr, err := s.Info()
	if err != nil {
		fmt.Printf("failed to build broadcast json due to %v\n", err)
		return
	}

	if _, err := w.Write([]byte(jsonStr)); err != nil {
		fmt.Printf("error occured during writing http request %v, error: %v\n", r, err)
	}
}
