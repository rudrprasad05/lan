package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"lan-share/daemon/internal/discovery"
)

type CLI struct {
	Registry *discovery.Registry
}

func NewCLI(reg *discovery.Registry) *CLI {
	return &CLI{Registry: reg}
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
			c.Registry.SetState(parts[1], discovery.Trusted)
			fmt.Println("trusted:", parts[1])

		case "reject":
			if len(parts) < 2 {
				fmt.Println("usage: reject <id>")
				continue
			}
			c.Registry.SetState(parts[1], discovery.Rejected)
			fmt.Println("rejected:", parts[1])

		default:
			fmt.Println("unknown command")
		}
	}
}

func (c *CLI) list() {
	devices := c.Registry.GetAll()

	for _, d := range devices {
		fmt.Printf("%s | %s | %s | %s | %s\n",
			d.ID,
			d.Name,
			d.IP,
			d.OS,
			d.State,
		)
	}
}
