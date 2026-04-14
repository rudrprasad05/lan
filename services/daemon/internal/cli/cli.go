package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"lan-share/daemon/internal/device"
	"lan-share/daemon/internal/discovery"
	"lan-share/daemon/internal/storage"
)

type CLI struct {
	Registry *discovery.Registry
	Service  *device.Service
}

func NewCLI(reg *discovery.Registry, svc *device.Service) *CLI {
	return &CLI{
		Registry: reg,
		Service:  svc,
	}
}

func (c *CLI) Start() {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("LAN CLI ready. Commands: list, trust <id>, reject <id>")

	for {
		fmt.Print("> ")
		scanner.Scan()

		input := strings.TrimSpace(scanner.Text())
		parts := strings.Split(input, " ")

		switch parts[0] {

		case "list":
			c.list()

		case "trust":
			if len(parts) < 2 {
				fmt.Println("usage: trust <id>")
				continue
			}

			c.Service.Trust(parts[1])
			fmt.Println("trusted:", parts[1])

		case "reject":
			if len(parts) < 2 {
				fmt.Println("usage: reject <id>")
				continue
			}

			c.Service.Reject(parts[1])
			fmt.Println("rejected:", parts[1])

		default:
			fmt.Println("unknown command")
		}
	}
}

func (c *CLI) list() {
	devs, err := storage.GetDevices()
	if err != nil {
		fmt.Println("db error:", err)
		return
	}

	for _, d := range devs {
		fmt.Printf("%s | %s | %s | %s\n",
			d.ID,
			d.Name,
			d.State,
			time.Unix(d.LastSeen, 0).Format("15:04:05"),
		)
	}
}
