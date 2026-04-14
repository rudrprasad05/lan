package storage

import (
	"crypto/ed25519"
	"encoding/json"
	"fmt"
	"os"
	"os/user"
	"runtime"

	"github.com/google/uuid"
)

type DeviceIdentity struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	DeviceType string `json:"device_type"`
	OS         string `json:"os"`
	OSVersion  string `json:"os_version"`
	Arch       string `json:"arch"`
	Hostname   string `json:"hostname"`
	PublicKey  string `json:"public_key"`
	PrivateKey string `json:"private_key"`
}

func LoadOrCreateIdentity() (*DeviceIdentity, error) {
	filePath := "device_identity.json"

	// Check if file exists
	if _, err := os.Stat(filePath); err == nil {
		// Load existing
		data, err := os.ReadFile(filePath)
		if err != nil {
			return nil, err
		}

		var identity DeviceIdentity
		if err := json.Unmarshal(data, &identity); err != nil {
			return nil, err
		}

		fmt.Println("Loaded existing identity")
		return &identity, nil
	}

	// Create new identity
	fmt.Println("Creating new device identity...")

	id := uuid.NewString()

	hostname, _ := os.Hostname()
	currentUser, _ := user.Current()

	osName := runtime.GOOS
	arch := runtime.GOARCH

	// simple OS version (we improve later)
	osVersion := "unknown"

	// generate keypair
	pub, priv, _ := ed25519.GenerateKey(nil)

	name := fmt.Sprintf("%s's device", currentUser.Username)

	identity := DeviceIdentity{
		ID:         id,
		Name:       name,
		DeviceType: "desktop",
		OS:         osName,
		OSVersion:  osVersion,
		Arch:       arch,
		Hostname:   hostname,
		PublicKey:  string(pub),
		PrivateKey: string(priv),
	}

	// Save to file
	data, err := json.MarshalIndent(identity, "", "  ")
	if err != nil {
		return nil, err
	}

	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return nil, err
	}

	fmt.Println("New identity created and saved")

	return &identity, nil
}
