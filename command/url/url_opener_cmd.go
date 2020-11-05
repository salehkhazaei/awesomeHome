package url

import (
	"ir.skhf/awesomeHome/process"
	"net/url"
)

type UrlOpenerCmd struct {
	Url string
}

func (cmd *UrlOpenerCmd) Run() error {
	return process.Exec(cmd.Url)
}

func (cmd *UrlOpenerCmd) Detect(commandStr string) bool {
	_, err := url.ParseRequestURI(commandStr)
	if err != nil {
		return false
	}
	cmd.Url = commandStr
	return true
}
