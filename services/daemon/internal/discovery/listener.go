package discovery

import (
	"encoding/json"
	"lan-share/daemon/internal/node"
	"lan-share/daemon/internal/storage"
	"log"
	"net"
	"time"
)

func StartListener(ctx *node.NodeContext, reg *Registry) {
	identity, idErr := storage.LoadOrCreateIdentity()
	if idErr != nil {
		log.Fatal(idErr)
	}

	addr, err := net.ResolveUDPAddr("udp4", ":9999")
	if err != nil {
		log.Fatal(err)
	}

	conn, err := net.ListenUDP("udp4", addr)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	buf := make([]byte, 2048)

	for {
		n, remote, err := conn.ReadFromUDP(buf)
		if err != nil {
			log.Println(err)
			continue
		}

		var msg DeviceMessage
		if err := json.Unmarshal(buf[:n], &msg); err != nil {
			continue
		}

		reg.Upsert(msg, remote.IP.String())

		if identity.PublicKey == msg.Key {
			reg.SetState(msg.ID, Me)

			if err := storage.UpdateState(msg.ID, "me"); err != nil {
				log.Println("db sync error (state):", err)
			}

			if err := storage.UpdateLastSeen(msg.ID, time.Now().Unix()); err != nil {
				log.Println("db sync error (last_seen):", err)
			}

			continue
		}

		err = storage.UpsertDevice(storage.StoredDevice{
			ID:        msg.ID,
			Name:      msg.Name,
			PublicKey: msg.Key, // fill later in pairing phase
			State:     "seen",
			LastSeen:  time.Now().Unix(),
			TrustedAt: 0,
		})

		if err != nil {
			log.Println("db sync error:", err)
		}
	}
}
