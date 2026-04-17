package api

import (
	"context"
	"lan-share/daemon/api/middleware"
	"lan-share/daemon/internal/discovery"
	"lan-share/daemon/internal/filetransfer"
	"log"
	"net/http"
	"time"
)

type Server struct {
	addr      string
	reg       *discovery.Registry
	selfID    string
	fileStore *filetransfer.Service
	server    *http.Server
}

func NewServer(addr string, reg *discovery.Registry, selfID string, fileStore *filetransfer.Service) *Server {
	mux := http.NewServeMux()

	s := &Server{
		addr:      addr,
		reg:       reg,
		selfID:    selfID,
		fileStore: fileStore,
	}

	mux.HandleFunc("/api/health", s.healthHandler)
	mux.HandleFunc("/api/devices", s.devicesHandler)
	mux.HandleFunc("/api/files", s.filesHandler)
	mux.HandleFunc("/api/files/", s.fileDownloadHandler)

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
