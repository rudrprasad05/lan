package storage

import (
	"context"
	"fmt"
	"log"
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
	DeviceType DeviceType `json:"deviceType"`
	OS         string
	OSVersion  string
	Arch       string
	Hostname   string
	PublicKey  string
	PrivateKey string
	CreatedAt  int64
}

type DeviceType string

const (
	DeviceDesktop DeviceType = "desktop"
	DeviceLaptop  DeviceType = "laptop"
	DeviceMobile  DeviceType = "mobile"
)

func DetectDeviceType() DeviceType {
	log.Println(runtime.GOOS)
	switch runtime.GOOS {

	case "android", "ios":
		log.Println("mb")
		return DeviceMobile

	case "darwin":
		log.Println("mac")
		return DeviceLaptop

	case "windows", "linux":
		// Try battery detection (laptop vs desktop)
		if hasBattery() {
			return DeviceLaptop
		}
		return DeviceDesktop

	default:
		return DeviceDesktop
	}
}

func hasBattery() bool {
	_, err := os.Stat("/sys/class/power_supply/BAT0")
	return err == nil
}

func ParseDeviceType(s string) DeviceType {
	switch DeviceType(s) {
	case DeviceDesktop, DeviceLaptop, DeviceMobile:
		return DeviceType(s)
	default:
		return DeviceDesktop // fallback
	}
}

func LoadOrCreateIdentity() (*DeviceIdentity, error) {

	ctx := context.Background()

	// 1. Try load existing
	row, err := Queries.GetIdentity(ctx)
	if err == nil {

		return &DeviceIdentity{
			ID:         row.ID,
			Name:       row.Name,
			DeviceType: DetectDeviceType(),
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
		DeviceType: DetectDeviceType(),
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
		DeviceType: string(identity.DeviceType),
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
