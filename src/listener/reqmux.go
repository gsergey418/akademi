package listener

import (
	"github.com/gsergey418alt/akademi/core"
	"github.com/gsergey418alt/akademi/pb"
)

// Multiplexer for the BaseMessage type.
func (u *UDPListener) reqMux(req *pb.BaseMessage) (*pb.BaseMessage, error) {
	res := &pb.BaseMessage{}
	switch {
	case req.GetPingRequest() != nil:
		res.Message = &pb.BaseMessage_PingResponse{}
	case req.GetFindNodeRequest() != nil:
		msg := &pb.FindNodeResponse{}
		nodes, err := u.AkademiNode.GetClosestNodes(core.BaseID(req.GetFindNodeRequest().NodeID), core.BucketSize)
		if err != nil {
			return nil, err
		}
		for _, v := range nodes {
			msg.RoutingEntry = append(msg.RoutingEntry, &pb.RoutingEntry{
				Address: string(v.Host),
				NodeID:  v.NodeID[:],
			})
		}
		res.Message = &pb.BaseMessage_FindNodeResponse{FindNodeResponse: msg}
	case req.GetFindKeyRequest() != nil:
		msg := &pb.FindKeyResponse{}
		data := u.AkademiNode.Get(core.BaseID(req.GetFindKeyRequest().KeyID))
		if data != nil {
			msg.Data = data
		} else {
			nodes, err := u.AkademiNode.GetClosestNodes(core.BaseID(req.GetFindKeyRequest().KeyID), core.BucketSize)
			if err != nil {
				return nil, err
			}
			for _, v := range nodes {
				msg.RoutingEntry = append(msg.RoutingEntry, &pb.RoutingEntry{
					Address: string(v.Host),
					NodeID:  v.NodeID[:],
				})
			}
		}
		res.Message = &pb.BaseMessage_FindKeyResponse{FindKeyResponse: msg}
	case req.GetStoreRequest() != nil:
		err := u.AkademiNode.Set(req.GetStoreRequest().Data)
		if err != nil {
			return nil, err
		}
		res.Message = &pb.BaseMessage_StoreResponse{}
	}
	return res, nil
}

// Populates the default response protobuf.
func (u *UDPListener) populateDefaultResponse(res, req *pb.BaseMessage) {
	res.RequestID = req.RequestID
	res.ListenPort = uint32(u.ListenAddr.Port)
	res.NodeID = u.AkademiNode.NodeID[:]
}
