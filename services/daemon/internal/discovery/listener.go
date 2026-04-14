package discovery

import (
	"encoding/json"
	"log"
	"net"
)

func StartListener(reg *Registry) {
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

		log.Printf("discovered device: %s (%s)\n", msg.Name, remote.IP.String())
	}
}
