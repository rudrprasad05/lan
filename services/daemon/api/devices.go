package api

import (
	"encoding/json"
	"net/http"

	"lan-share/daemon/api/res"
	"lan-share/daemon/internal/discovery"
)

type DeviceListResponse struct {
	Devices []DeviceResponse `json:"devices"`
}

type DeviceResponse struct {
	res.BaseResponse
	ID       string `json:"id"`
	Name     string `json:"name"`
	Type     string `json:"type"`
	OS       string `json:"os"`
	Arch     string `json:"arch"`
	IP       string `json:"ip"`
	State    string `json:"state"`
	LastSeen int64  `json:"lastSeen"`
}

func (s *Server) devicesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	devices := s.reg.GetAll()
	response := DeviceListResponse{
		Devices: make([]DeviceResponse, 0, len(devices)),
	}

	for _, device := range devices {
		if shouldSkipDevice(device, s.selfID) {
			continue
		}

		response.Devices = append(response.Devices, DeviceResponse{
			BaseResponse: res.NewBaseResponse(),
			ID:           device.ID,
			Name:         device.Name,
			Type:         device.Type,
			IP:           device.IP,
			State:        string(device.State),
			LastSeen:     device.LastSeen.Unix(),
		})
	}

	_ = json.NewEncoder(w).Encode(response)
}

func shouldSkipDevice(device discovery.Device, selfID string) bool {
	return device.ID == selfID || device.State == discovery.Me
}
