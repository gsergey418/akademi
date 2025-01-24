package listener

import (
	"bytes"
	"fmt"
	"log"
	"net"

	"github.com/gsergey418alt/akademi/core"
	"github.com/gsergey418alt/akademi/node"
	"github.com/gsergey418alt/akademi/pb"
	"google.golang.org/protobuf/proto"
)

// The UDPListener struct receives protocol buffers
// encoded data via a UDP socket.
type UDPListener struct {
	ListenAddr  *net.UDPAddr
	AkademiNode *node.AkademiNode

	udpReadBuffer [65535]byte
	udpConn       *net.UDPConn
}

// Parse listenAddrString and set AkademiNode.
func (u *UDPListener) Initialize(listenAddrStr string, a *node.AkademiNode) error {
	listenAddr, err := net.ResolveUDPAddr("udp", listenAddrStr)
	if err != nil {
		return err
	}
	u.ListenAddr = listenAddr
	u.AkademiNode = a
	return nil
}

// Sturcture that contains the source host and message bytes.
type UDPBytes struct {
	host *net.UDPAddr
	data []byte
}

// Main listener goroutine.
func (u *UDPListener) ListenUDP(msgChan chan UDPBytes, errChan chan error) {
	for {
		l, host, err := u.udpConn.ReadFromUDP(u.udpReadBuffer[:])
		if err != nil {
			log.Print(err)
		}
		msgChan <- UDPBytes{host, u.udpReadBuffer[:l]}
	}
}

// Handler goroutine.
func (u *UDPListener) UDPWorker(msgChan chan UDPBytes, errChan chan error) {
	for {
		udpBytes := <-msgChan
		err := u.handleUDPMessage(udpBytes.host, udpBytes.data)
		if err != nil {
			log.Print(err)
		}
	}
}

// Opens a UDP socket on UDPListener.ListenAddr.
func (u *UDPListener) Listen() error {
	conn, err := net.ListenUDP("udp", u.ListenAddr)
	u.udpConn = conn
	if err != nil {
		return err
	}
	defer conn.Close()
	log.Print("UDP listener started on address ", u.ListenAddr, ".")

	errChan := make(chan error)
	msgChan := make(chan UDPBytes, 1024)

	go u.ListenUDP(msgChan, errChan)
	go u.UDPWorker(msgChan, errChan)

	log.Fatal(<-errChan)
	return nil
}

// Populates the default response protobuf.
func (u *UDPListener) populateDefaultResponse(res, req *pb.BaseMessage) {
	res.RequestID = req.RequestID
	res.ListenPort = uint32(u.ListenAddr.Port)
	res.NodeID = u.AkademiNode.NodeID[:]
}

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

// Handle a slice of bytes as a UDP message.
func (u *UDPListener) handleUDPMessage(host *net.UDPAddr, buf []byte) error {
	log.Print("Message from ", host, ": ", len(buf), " bytes.")
	req := &pb.BaseMessage{}
	err := proto.Unmarshal(buf, req)
	if err != nil {
		return err
	}

	if bytes.Equal(req.NodeID, u.AkademiNode.NodeID[:]) {
		return fmt.Errorf("Request from self, dropping.")
	}

	r := core.RoutingEntry{
		Host:   core.Host(fmt.Sprintf("%s:%d", host.IP, req.ListenPort)),
		NodeID: core.BaseID(req.NodeID),
	}
	err = u.AkademiNode.UpdateRoutingTable(r)
	if err != nil {
		return err
	}

	res, err := u.reqMux(req)
	if err != nil {
		res = &pb.BaseMessage{}
		msg := &pb.ErrorMessage{Text: err.Error()}
		res.Message = &pb.BaseMessage_ErrorMessage{ErrorMessage: msg}
	}
	err = u.sendUDPMessage(host, res, req)

	return err
}
