package command

import (
	"errors"
	"ir.skhf/awesomeHome/command/google"
	"ir.skhf/awesomeHome/command/url"
	"ir.skhf/awesomeHome/process"
)

type Command interface {
	Detect(commandStr string) bool
	Run(processService *process.ProcessService) error
}

type CommandService struct {
	ProcessService *process.ProcessService
}

func NewCommandService(processService *process.ProcessService) *CommandService {
	return &CommandService{
		ProcessService: processService,
	}
}

func (s *CommandService) Run(commandStr string) error {
	for _, cmdT := range s.GetCommands() {
		if cmdT.Detect(commandStr) {
			return cmdT.Run(s.ProcessService)
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
