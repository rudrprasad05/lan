package discovery

type DeviceMessage struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"` // laptop, desktop, etc
	OS   string `json:"os"`   // linux, windows, mac
	Arch string `json:"arch"`
	Port int    `json:"port"` // gRPC port later
	Key  string `json:"key"`
}
