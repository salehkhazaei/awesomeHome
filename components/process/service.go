package process

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"syscall"
)

type ProcessService struct{}

func NewProcessService() *ProcessService {
	return &ProcessService{}
}

func (s *ProcessService) Exec(cmd string, args ...string) error {
	binary, err := exec.LookPath(cmd)
	if err != nil {
		return err
	}

	env := os.Environ()
	err = syscall.Exec(binary, args, env)
	if err != nil {
		return err
	}

	return nil
}

func (s *ProcessService) VLC(filename string) error {
	return s.Exec("vlc", filename)
}

func (s *ProcessService) OpenBrowser(url string) error {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}

	return err
}
