package main

import (
	"fmt"
	"log"

	"lan-share/daemon/internal/storage"
)

func main() {
	storage.InitDB()

	log.Println("Daemon started")

	identity, err := storage.LoadOrCreateIdentity()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Device Identity:")
	fmt.Println("ID:", identity.ID)
	fmt.Println("Name:", identity.Name)
	fmt.Println("OS:", identity.OS, identity.OSVersion)
	fmt.Println("Arch:", identity.Arch)
}
