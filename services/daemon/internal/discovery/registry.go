package discovery

import (
	"sync"
	"time"
)

type Device struct {
	DeviceMessage
	IP       string
	LastSeen time.Time
	State    DeviceState
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

func (r *Registry) Upsert(msg DeviceMessage, ip string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	d, exists := r.devices[msg.ID]

	if exists {
		d.LastSeen = time.Now()
		d.IP = ip
		return
	}

	r.devices[msg.ID] = &Device{
		DeviceMessage: msg,
		IP:            ip,
		LastSeen:      time.Now(),
		State:         Seen,
	}
}

func (r *Registry) GetAll() []Device {
	r.mu.Lock()
	defer r.mu.Unlock()

	out := make([]Device, 0, len(r.devices))
	for _, d := range r.devices {
		out = append(out, *d)
	}
	return out
}

func (r *Registry) SetState(id string, state DeviceState) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if d, ok := r.devices[id]; ok {
		d.State = state
	}
}
