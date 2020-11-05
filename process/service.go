package process

import (
	"os"
	"os/exec"
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
