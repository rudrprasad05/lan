package storage

import (
	"database/sql"
	"fmt"
	crypto "lan-share/daemon/internal/cypto"
	"os"
	"os/user"
	"runtime"
	"time"

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
	// 1. Try load existing
	row := DB.QueryRow(`
		SELECT id, name, device_type, os, os_version, arch, hostname, public_key, private_key, created_at
		FROM device_identity
		LIMIT 1
	`)

	var identity DeviceIdentity

	err := row.Scan(
		&identity.ID,
		&identity.Name,
		&identity.DeviceType,
		&identity.OS,
		&identity.OSVersion,
		&identity.Arch,
		&identity.Hostname,
		&identity.PublicKey,
		&identity.PrivateKey,
		&identity.CreatedAt,
	)

	if err == nil {
		fmt.Println("Loaded existing identity from DB")
		return &identity, nil
	}

	if err != sql.ErrNoRows {
		return nil, err
	}

	// 2. Create new identity
	fmt.Println("Creating new device identity...")

	id := uuid.NewString()

	hostname, _ := os.Hostname()
	currentUser, _ := user.Current()

	osName := runtime.GOOS
	arch := runtime.GOARCH
	osVersion := "unknown"

	pub, priv, err := crypto.GenerateKeyPair()
	if err != nil {
		return nil, err
	}

	name := fmt.Sprintf("%s's device", currentUser.Username)

	identity = DeviceIdentity{
		ID:         id,
		Name:       name,
		DeviceType: "desktop",
		OS:         osName,
		OSVersion:  osVersion,
		Arch:       arch,
		Hostname:   hostname,
		PublicKey:  string(pub),
		PrivateKey: string(priv),
		CreatedAt:  time.Now().Unix(),
	}

	// 3. Insert into DB
	_, err = DB.Exec(`
		INSERT INTO device_identity 
		(id, name, device_type, os, os_version, arch, hostname, public_key, private_key, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`,
		identity.ID,
		identity.Name,
		identity.DeviceType,
		identity.OS,
		identity.OSVersion,
		identity.Arch,
		identity.Hostname,
		identity.PublicKey,
		identity.PrivateKey,
		identity.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	fmt.Println("New identity created and stored in DB")

	return &identity, nil
}
