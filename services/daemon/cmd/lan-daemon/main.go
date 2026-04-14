package main

import (
	"lan-share/daemon/internal/cli"
	"lan-share/daemon/internal/device"
	"lan-share/daemon/internal/discovery"
	"lan-share/daemon/internal/node"
	"lan-share/daemon/internal/storage"
	"log"
)

func main() {

	storage.InitDB()

	identity, err := storage.LoadOrCreateIdentity()
	if err != nil {
		log.Fatal(err)
		return
	}

	ctx := &node.NodeContext{
		Identity: identity,
	}
	reg := discovery.NewRegistry()
	svc := device.NewService()

	go discovery.StartBroadcaster(ctx)
	go discovery.StartListener(ctx, reg)

	go cli.NewCLI(reg, svc).Start()

	log.Println("System running")

	select {}
}
