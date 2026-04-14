package node

import "lan-share/daemon/internal/storage"

type NodeContext struct {
	Identity *storage.DeviceIdentity
}
