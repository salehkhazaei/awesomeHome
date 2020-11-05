package google

import "testing"

func TestGoogleItCmdDetect(t *testing.T) {
	cmd := GoogleItCmd{}
	if cmd.Detect("google salam saleh") != true {
		t.Error("google it didnt detect well")
	}
}

func TestGoogleItCmdDontDetect(t *testing.T) {
	cmd := GoogleItCmd{}
	if cmd.Detect("hey google salam saleh") != false {
		t.Error("google it detect when it shouldnt")
	}
}

func TestFillDataOnDetect(t *testing.T) {
	cmd := GoogleItCmd{}
	cmd.Detect("google salam")
	if len(cmd.Query) <= 0 {
		t.Error("detect didnt fill data")
	}
}
