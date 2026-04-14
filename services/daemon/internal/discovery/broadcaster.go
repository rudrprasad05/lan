package discovery

import (
	"encoding/json"
	"lan-share/daemon/internal/node"
	"log"
	"net"
	"runtime"
	"time"
)

func StartBroadcaster(ctx *node.NodeContext) {
	addr, _ := net.ResolveUDPAddr("udp4", "255.255.255.255:9999")

	conn, err := net.DialUDP("udp4", nil, addr)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	ticker := time.NewTicker(2 * time.Second)

	msg := DeviceMessage{
		ID:   ctx.Identity.ID,
		Name: ctx.Identity.Name,
		Type: ctx.Identity.DeviceType,
		OS:   runtime.GOOS,
		Arch: runtime.GOARCH,
		Port: 50052,
		Key:  ctx.Identity.PublicKey,
	}

	for range ticker.C {
		data, _ := json.Marshal(msg)

		_, err := conn.Write(data)
		if err != nil {
			log.Println("broadcast error:", err)
		}
	}
}
