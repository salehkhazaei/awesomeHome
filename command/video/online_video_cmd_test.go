package video

import "testing"

func TestOnlineVideoCmdDetect(t *testing.T) {
	cmd := OnlineVideoCmd{}
	if cmd.Detect("https://file-examples-com.github.io/uploads/2017/04/file_example_MP4_480_1_5MG.mp4") != true {
		t.Error("url opener didnt detect well")
	}
}

func TestOnlineVideoCmdDontDetect(t *testing.T) {
	cmd := OnlineVideoCmd{}
	if cmd.Detect("g https://google.com") != false {
		t.Error("url opener detect when it shouldnt")
	}
}

func TestFillDataOnDetect(t *testing.T) {
	cmd := OnlineVideoCmd{}
	cmd.Detect("https://file-examples-com.github.io/uploads/2017/04/file_example_MP4_480_1_5MG.mp4")
	if len(cmd.Url) <= 0 {
		t.Error("detect didnt fill data")
	}
}
