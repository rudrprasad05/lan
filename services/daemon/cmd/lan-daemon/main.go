package main

import (
	"lan-share/daemon/internal/discovery"
	"lan-share/daemon/internal/storage"
	"log"
)

func main() {
	storage.InitDB()

	log.Println("Daemon started")

	identity, err := storage.LoadOrCreateIdentity()
	if err != nil {
		log.Fatal(err)
	}

	reg := discovery.NewRegistry()

	msg := discovery.DeviceMessage{
		ID:   identity.ID,
		Name: identity.Name,
		Type: identity.DeviceType,
		OS:   identity.OS,
		Arch: identity.Arch,
		Port: 50052,
	}

	go discovery.StartBroadcaster(msg)
	go discovery.StartListener(reg)

	log.Println("Discovery running...")

	select {} // block forever
}
