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
