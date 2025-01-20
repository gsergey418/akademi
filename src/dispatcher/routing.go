package dispatcher

import (
	"github.com/gsergey418alt/akademi/core"
	"github.com/gsergey418alt/akademi/pb"
)

// Parser for the RoutingHeader
func (u *UDPDispatcher) parseRoutingHeader(msg *pb.BaseMessage) core.RoutingHeader {
	return core.RoutingHeader{
		NodeID:     core.BaseID(msg.NodeID),
		ListenPort: core.IPPort(msg.ListenPort),
	}
}

// Filler for the RoutingHeader
func (u *UDPDispatcher) writeRoutingHeader(msg *pb.BaseMessage) {
	msg.ListenPort = uint32(u.RoutingHeader.ListenPort)
	rID := RandomRequestID()
	msg.RequestID = rID[:]
	msg.NodeID = u.RoutingHeader.NodeID[:]
}
