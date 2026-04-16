package api

import (
	"context"
	"lan-share/daemon/api/middleware"
	"lan-share/daemon/internal/discovery"
	"log"
	"net/http"
	"time"
)

type Server struct {
	addr   string
	reg    *discovery.Registry
	selfID string
	server *http.Server
}

func NewServer(addr string, reg *discovery.Registry, selfID string) *Server {
	mux := http.NewServeMux()

	s := &Server{
		addr:   addr,
		reg:    reg,
		selfID: selfID,
	}

	mux.HandleFunc("/api/health", s.healthHandler)
	mux.HandleFunc("/api/devices", s.devicesHandler)

	handler := middleware.CORSMiddleware(middleware.LoggingMiddleware(mux))

	s.server = &http.Server{
		Addr:              addr,
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
	}

	return s
}

func (s *Server) Start() error {
	log.Println("API listening on", s.addr)
	return s.server.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
