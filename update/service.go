package update

import (
	"net/http"

	"github.com/inconshreveable/go-update"
)

type SelfUpdateService struct {
}

func NewSelfUpdateService() *SelfUpdateService {
	return &SelfUpdateService{}
}

func (s *SelfUpdateService) doUpdate(url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	err = update.Apply(resp.Body, update.Options{})
	if err != nil {
		return err
	}
	return err
}
