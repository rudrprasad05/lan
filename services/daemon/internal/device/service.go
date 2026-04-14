package device

import (
	"context"
	"log"
	"time"

	"lan-share/daemon/internal/storage"
	db "lan-share/daemon/internal/storage/db"
)

type Service struct{}

func NewService() *Service {
	return &Service{}
}
func (s *Service) OnSeen(id, name, ip string) {
	ctx := context.Background()
	now := time.Now().Unix()

	existing, err := storage.Queries.GetDeviceByID(ctx, id)

	finalName := name
	state := "seen"
	trustedAt := int64(0)

	// ONLY check error, NOT nil
	if err == nil {
		if existing.Name != "" {
			finalName = existing.Name
		}
		state = existing.State
		trustedAt = existing.TrustedAt
	} else {
		// ignore "not found"
		if err.Error() != "sql: no rows in result set" {
			log.Println("db error:", err)
			return
		}
	}

	err = storage.Queries.UpsertDevice(ctx, db.UpsertDeviceParams{
		ID:        id,
		Name:      finalName,
		PublicKey: "",
		State:     state,
		LastSeen:  now,
		TrustedAt: trustedAt,
	})

	if err != nil {
		log.Println("db sync error:", err)
	}
}

func (s *Service) Trust(id string) {
	ctx := context.Background()
	now := time.Now().Unix()

	err := storage.Queries.UpdateDeviceStateAndTrust(ctx, db.UpdateDeviceStateAndTrustParams{
		State:     "trusted",
		TrustedAt: now,
		ID:        id,
	})

	if err != nil {
		log.Println("trust error:", err)
	}
}

func (s *Service) Reject(id string) {
	ctx := context.Background()

	err := storage.Queries.UpdateDeviceStateAndTrust(ctx, db.UpdateDeviceStateAndTrustParams{
		State:     "rejected",
		TrustedAt: 0,
		ID:        id,
	})

	if err != nil {
		log.Println("reject error:", err)
	}
}
