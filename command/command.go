package command

import (
	"errors"
	"ir.skhf/awesomeHome/command/google"
	"ir.skhf/awesomeHome/command/url"
)

type Command interface {
	Detect(commandStr string) bool
	Run() error
}

type CommandService struct {
}

func NewCommandService() *CommandService {
	return &CommandService{}
}

func (s *CommandService) Run(commandStr string) error {
	for _, cmdT := range s.GetCommands() {
		if cmdT.Detect(commandStr) {
			return cmdT.Run()
		}
	}
	return errors.New("command not found")
}

func (s *CommandService) GetCommands() []Command {
	return []Command{
		&google.GoogleItCmd{},
		&url.UrlOpenerCmd{},
	}
}
