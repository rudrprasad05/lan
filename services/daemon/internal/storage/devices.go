package storage

import (
	"context"

	db "lan-share/daemon/internal/storage/db"
)

type StoredDevice struct {
	ID        string
	Name      string
	PublicKey string
	State     string
	LastSeen  int64
	TrustedAt int64
}

func UpdateLastSeen(id string, t int64) error {
	return Queries.UpdateLastSeen(context.Background(), db.UpdateLastSeenParams{
		LastSeen: t,
		ID:       id,
	})
}
func UpdateState(id string, state string) error {
	return Queries.UpdateState(context.Background(), db.UpdateStateParams{
		State: state,
		ID:    id,
	})
}

func UpsertDevice(d StoredDevice) error {
	return Queries.UpsertDevice(context.Background(), db.UpsertDeviceParams{
		ID:        d.ID,
		Name:      d.Name,
		PublicKey: d.PublicKey,
		State:     d.State,
		LastSeen:  d.LastSeen,
		TrustedAt: d.TrustedAt,
	})
}

func GetDevices() ([]StoredDevice, error) {

	rows, err := Queries.GetDevices(context.Background())
	if err != nil {
		return nil, err
	}

	list := make([]StoredDevice, 0, len(rows))

	for _, r := range rows {
		list = append(list, StoredDevice{
			ID:        r.ID,
			Name:      r.Name,
			PublicKey: r.PublicKey,
			State:     r.State,
			LastSeen:  r.LastSeen,
			TrustedAt: r.TrustedAt,
		})
	}

	return list, nil
}
