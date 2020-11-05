package process

import (
	"os"
	"os/exec"
	"syscall"
)

func Exec(cmd string, args ...string) error {
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
