package api

import (
	"encoding/json"
	"net/http"

	"lan-share/daemon/api/res"
)

type HealthResponse struct {
	res.BaseResponse
	Status string `json:"status"`
}

func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	res := HealthResponse{
		BaseResponse: res.NewBaseResponse(),
		Status:       "ok",
	}

	_ = json.NewEncoder(w).Encode(res)
}
