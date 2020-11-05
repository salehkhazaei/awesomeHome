package url

import "testing"

func TestUrlOpenerCmdDetect(t *testing.T) {
	cmd := UrlOpenerCmd{}
	if cmd.Detect("https://google.com") != true {
		t.Error("url opener didnt detect well")
	}
}

func TestUrlOpenerCmdDontDetect(t *testing.T) {
	cmd := UrlOpenerCmd{}
	if cmd.Detect("g https://google.com") != true {
		t.Error("url opener detect when it shouldnt")
	}
}

func TestFillDataOnDetect(t *testing.T) {
	cmd := UrlOpenerCmd{}
	cmd.Detect("https://google.com")
	if len(cmd.Url) <= 0 {
		t.Error("detect didnt fill data")
	}
}
