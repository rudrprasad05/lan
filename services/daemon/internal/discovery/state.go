package discovery

type DeviceState string

const (
	Unknown  DeviceState = "unknown"
	Seen     DeviceState = "seen"
	Trusted  DeviceState = "trusted"
	Rejected DeviceState = "rejected"
)
