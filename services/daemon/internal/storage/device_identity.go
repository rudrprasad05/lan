package storage

import (
	"context"
	"fmt"
	"os"
	"os/user"
	"runtime"
	"time"

	crypto "lan-share/daemon/internal/cypto"
	"lan-share/daemon/internal/storage/db"

	"github.com/google/uuid"
)

type DeviceIdentity struct {
	ID         string
	Name       string
	DeviceType string
	OS         string
	OSVersion  string
	Arch       string
	Hostname   string
	PublicKey  string
	PrivateKey string
	CreatedAt  int64
}

func LoadOrCreateIdentity() (*DeviceIdentity, error) {

	ctx := context.Background()

	// 1. Try load existing
	row, err := Queries.GetIdentity(ctx)
	if err == nil {

		return &DeviceIdentity{
			ID:         row.ID,
			Name:       row.Name,
			DeviceType: row.DeviceType,
			OS:         row.Os,
			OSVersion:  row.OsVersion,
			Arch:       row.Arch,
			Hostname:   row.Hostname,
			PublicKey:  row.PublicKey,
			PrivateKey: row.PrivateKey,
			CreatedAt:  row.CreatedAt,
		}, nil
	}

	// 2. Create new identity

	id := uuid.NewString()

	hostname, _ := os.Hostname()
	currentUser, _ := user.Current()

	pub, priv, err := crypto.GenerateKeyPair()
	if err != nil {
		return nil, err
	}

	identity := DeviceIdentity{
		ID:         id,
		Name:       fmt.Sprintf("%s's device", currentUser.Username),
		DeviceType: "desktop",
		OS:         runtime.GOOS,
		OSVersion:  "unknown",
		Arch:       runtime.GOARCH,
		Hostname:   hostname,
		PublicKey:  string(pub),
		PrivateKey: string(priv),
		CreatedAt:  time.Now().Unix(),
	}

	// 3. Insert via sqlc
	err = Queries.CreateIdentity(ctx, db.CreateIdentityParams{
		ID:         identity.ID,
		Name:       identity.Name,
		DeviceType: identity.DeviceType,
		Os:         identity.OS,
		OsVersion:  identity.OSVersion,
		Arch:       identity.Arch,
		Hostname:   identity.Hostname,
		PublicKey:  identity.PublicKey,
		PrivateKey: identity.PrivateKey,
		CreatedAt:  identity.CreatedAt,
	})

	if err != nil {
		return nil, err
	}

	return &identity, nil
}
