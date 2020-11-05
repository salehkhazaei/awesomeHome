package server

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type HttpServerService struct {
	Port int
}

func NewHttpServerService(port int) *HttpServerService {
	return &HttpServerService{
		Port: port,
	}
}

func (s *HttpServerService) Register(path string, handlerFunc func(w http.ResponseWriter, r *http.Request)) {
	http.HandleFunc(path, handlerFunc)
}

func (s *HttpServerService) Start() {
	addr := fmt.Sprintf(":%d", s.Port)
	log.Fatal(http.ListenAndServe(addr, nil))
}
