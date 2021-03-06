package info

import (
	"encoding/json"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestInfoJson(t *testing.T) {
	appInfoService := NewAppInfoService(nil, time.Millisecond)

	js, err := appInfoService.Info()
	if err != nil {
		panic(err)
	}

	appInfo := AppInfo{}
	err = json.Unmarshal([]byte(js), &appInfo)
	if err != nil {
		panic(err)
	}

	if len(appInfo.Version) == 0 {
		t.Error("Invalid version")
	}
}

func TestVersion(t *testing.T) {
	appInfoService := NewAppInfoService(nil, time.Millisecond)
	arr := strings.Split(appInfoService.Version(), ".")
	if len(arr) != 2 {
		t.Error("invalid version format")
	}

	major, err := strconv.Atoi(arr[0])
	if err != nil {
		panic(err)
	}

	minor, err := strconv.Atoi(arr[1])
	if err != nil {
		panic(err)
	}

	if major != Major {
		t.Error("invalid major")
	}

	if minor != Minor {
		t.Error("invalid minor")
	}
}
