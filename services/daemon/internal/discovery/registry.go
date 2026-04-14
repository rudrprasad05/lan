package discovery

import (
	"sync"
	"time"
)

type Device struct {
	DeviceMessage
	IP       string
	LastSeen time.Time
}

type Registry struct {
	mu      sync.Mutex
	devices map[string]*Device
}

func NewRegistry() *Registry {
	return &Registry{
		devices: make(map[string]*Device),
	}
}

func (r *Registry) Upsert(d DeviceMessage, ip string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.devices[d.ID] = &Device{
		DeviceMessage: d,
		IP:            ip,
		LastSeen:      time.Now(),
	}
}

func (r *Registry) GetAll() []Device {
	r.mu.Lock()
	defer r.mu.Unlock()

	var list []Device
	for _, d := range r.devices {
		list = append(list, *d)
	}
	return list
}

func (r *Registry) Cleanup(timeout time.Duration) {
	r.mu.Lock()
	defer r.mu.Unlock()

	for id, d := range r.devices {
		if time.Since(d.LastSeen) > timeout {
			delete(r.devices, id)
		}
	}
}
