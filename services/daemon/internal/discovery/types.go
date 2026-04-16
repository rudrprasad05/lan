package discovery

import "lan-share/daemon/internal/storage"

type DeviceMessage struct {
	ID         string             `json:"id"`
	Name       string             `json:"name"`
	DeviceType storage.DeviceType `json:"deviceType"` // laptop, desktop, etc
	OS         string             `json:"os"`         // linux, windows, mac
	Arch       string             `json:"arch"`
	Port       int                `json:"port"` // gRPC port later
	Key        string             `json:"key"`
}
