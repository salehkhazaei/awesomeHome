package info

import (
	"encoding/json"
	"fmt"
	"ir.skhf/awesomeHome/utils"
)

var Major int = 1
var Minor int = 0

type AppInfo struct {
	Version string
	IP      string
}

func Info() (string, error) {
	appInfo := AppInfo{
		Version: Version(),
		IP:      utils.GetOutboundIP().String(),
	}

	appInfoJson, err := json.Marshal(appInfo)
	if err != nil {
		return "", err
	}

	return string(appInfoJson), nil
}

func Version() string {
	return fmt.Sprintf("%d.%d", Major, Minor)
}
