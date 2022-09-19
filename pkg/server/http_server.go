package server

import (
	"context"
	"log"
	"net/http"
	"time"
)

type HTTPServer struct {
	server          *http.Server
	shutdownTimeout time.Duration
}

func NewHTTPServer(addr string, router http.Handler, shutdownTimeout time.Duration) *HTTPServer {
	return &HTTPServer{
		server: &http.Server{
			Addr:    addr,
			Handler: router,
		},
		shutdownTimeout: shutdownTimeout,
	}
}

func (s *HTTPServer) Start() {
	go func() {
		if err := s.server.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatal("http server handler")
		}

	}()
}

func (s *HTTPServer) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
	defer cancel()

	if err := s.server.Shutdown(ctx); err != nil {
		log.Fatal("http shutdown error")
	}

}
