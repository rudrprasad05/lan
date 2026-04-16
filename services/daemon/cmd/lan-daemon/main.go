package main

import (
	"lan-share/daemon/api"
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
	apiServer := api.NewServer("127.0.0.1:43821", reg, identity.ID)
	go func() {
		if err := apiServer.Start(); err != nil {
			log.Println("api server stopped:", err)
		}
	}()

	go discovery.StartBroadcaster(ctx)
	go discovery.StartListener(ctx, reg)

	// go cli.NewCLI(reg, svc).Start()

	log.Println("System running")

	select {}
}
