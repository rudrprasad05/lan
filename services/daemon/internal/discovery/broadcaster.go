package discovery

import (
	"encoding/json"
	"log"
	"net"
	"time"
)

func StartBroadcaster(msg DeviceMessage) {
	addr, _ := net.ResolveUDPAddr("udp4", "255.255.255.255:9999")

	conn, err := net.DialUDP("udp4", nil, addr)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	ticker := time.NewTicker(2 * time.Second)

	for range ticker.C {
		data, _ := json.Marshal(msg)

		_, err := conn.Write(data)
		if err != nil {
			log.Println("broadcast error:", err)
		}
	}
}
