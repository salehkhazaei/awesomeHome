package google

import (
	"ir.skhf/awesomeHome/process"
	"net/url"
)

type GoogleItCmd struct {
	Query string
}

func (cmd *GoogleItCmd) Run() error {
	return process.Exec("https://google.com?q=" + url.QueryEscape(cmd.Query))
}

func (cmd *GoogleItCmd) Detect(commandStr string) bool {
	googleCmd := "google "
	if commandStr[:len(googleCmd)] == googleCmd {
		cmd.Query = commandStr[len(googleCmd):]
		return true
	}
	return false
}
