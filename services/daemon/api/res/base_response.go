package res

import (
	"runtime"
	"time"
)

type BaseResponse struct {
	Service   string `json:"service"`
	Timestamp int64  `json:"timestamp"`
	Version   string `json:"version"`
	OS        string `json:"os"`
	Arch      string `json:"arch"`
}

func NewBaseResponse() BaseResponse {
	return BaseResponse{
		Service:   "lan-daemon",
		Timestamp: time.Now().Unix(),
		Version:   "dev",
		OS:        runtime.GOOS,
		Arch:      runtime.GOARCH,
	}
}
