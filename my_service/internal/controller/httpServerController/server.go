package httpServerController

import (
	"net/http"
	"time"
)

type Server struct {
	server http.Server
}

func (s *Server) RunServer(port string, router http.Handler) error {
	s.server = http.Server{
		Addr:           ":" + port,
		Handler:        router,
		MaxHeaderBytes: 1 << 20,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
	}
	return s.server.ListenAndServe()
}
