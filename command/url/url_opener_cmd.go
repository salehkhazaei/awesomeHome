package url

import (
	"ir.skhf/awesomeHome/process"
)

type UrlOpenerCmd struct {
	Url string
}

func (cmd *UrlOpenerCmd) Run() error {
	return process.Exec(cmd.Url)
}
