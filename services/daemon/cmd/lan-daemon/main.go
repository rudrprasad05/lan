package main

import (
	"log"
	"runtime"

	"lan-share/daemon/internal/cli"
	"lan-share/daemon/internal/discovery"
	"lan-share/daemon/internal/storage"
)

func main() {

	storage.InitDB()

	identity, err := storage.LoadOrCreateIdentity()
	if err != nil {
		log.Fatal(err)
	}

	reg := discovery.NewRegistry()

	msg := discovery.DeviceMessage{
		ID:   identity.ID,
		Name: identity.Name,
		Type: identity.DeviceType,
		OS:   runtime.GOOS,
		Arch: runtime.GOARCH,
		Port: 50052,
	}

	go discovery.StartBroadcaster(msg)
	go discovery.StartListener(reg)

	go cli.NewCLI(reg).Start()

	log.Println("System running")

	select {}
}
