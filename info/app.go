package info

import (
	"encoding/json"
	"fmt"
)

var Major int = 1
var Minor int = 0

type AppInfo struct {
	Version string
}

func Info() (string, error) {
	appInfo := AppInfo{Version: Version()}

	appInfoJson, err := json.Marshal(appInfo)
	if err != nil {
		return "", err
	}

	return string(appInfoJson), nil
}

func Version() string {
	return fmt.Sprintf("%d.%d", Major, Minor)
}
