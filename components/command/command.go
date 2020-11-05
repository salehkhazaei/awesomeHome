package command

import (
	"errors"
	"fmt"
	"io"
	"ir.skhf/awesomeHome/components/command/google"
	"ir.skhf/awesomeHome/components/command/url"
	"ir.skhf/awesomeHome/components/process"
	"net/http"
	"strings"
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

func (s *CommandService) HandleHttp(w http.ResponseWriter, r *http.Request) {
	buf := new(strings.Builder)
	_, err := io.Copy(buf, r.Body)
	if err != nil {
		fmt.Printf("error occured during reading request: %v\n", err)
		return
	}

	err = s.Run(buf.String())
	if err != nil {
		fmt.Printf("error occured during handling http request %v, error: %v\n", r, err)
		return
	}

	_, _ = w.Write([]byte("done"))
}
