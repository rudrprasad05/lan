package api

import (
	"context"
	"lan-share/daemon/api/middleware"
	"lan-share/daemon/internal/discovery"
	"lan-share/daemon/internal/filetransfer"
	"log"
	"net"
	"net/http"
	"time"
)

type Server struct {
	addr       string
	reg        *discovery.Registry
	selfID     string
	fileStore  *filetransfer.Service
	httpClient *http.Client
	server     *http.Server
}

func NewServer(addr string, reg *discovery.Registry, selfID string, fileStore *filetransfer.Service) *Server {
	mux := http.NewServeMux()

	s := &Server{
		addr:      addr,
		reg:       reg,
		selfID:    selfID,
		fileStore: fileStore,
		httpClient: &http.Client{
			Timeout: 0,
		},
	}

	mux.HandleFunc("/api/health", s.healthHandler)
	mux.HandleFunc("/api/devices", s.devicesHandler)
	mux.HandleFunc("/api/files", s.filesHandler)
	mux.HandleFunc("/api/files/receive", s.receiveFileHandler)
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

func (s *Server) apiURL(ip, endpoint string) string {
	host := net.JoinHostPort(ip, "43821")
	return "http://" + host + endpoint
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
