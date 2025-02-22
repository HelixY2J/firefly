package registry

import "time"

type NodeInfo struct {
	NodeID   string
	Address  string
	LastSeen time.Time
	IsMaster bool
}
