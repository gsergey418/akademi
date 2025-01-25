package listener

import (
	"bytes"
	"fmt"
	"log"
	"net"

	"github.com/gsergey418alt/akademi/core"
	"github.com/gsergey418alt/akademi/pb"
	"google.golang.org/protobuf/proto"
)

// Handle a slice of bytes as a UDP message.
func (u *UDPListener) handleUDPMessage(host *net.UDPAddr, buf []byte) error {
	log.Print("Message from ", host, ": ", len(buf), " bytes.")
	if len(buf) == 0 {
		return fmt.Errorf("empty packet, dropping")
	}
	req := &pb.BaseMessage{}
	err := proto.Unmarshal(buf, req)
	if err != nil {
		return err
	}

	if bytes.Equal(req.NodeID, u.AkademiNode.NodeID[:]) {
		return fmt.Errorf("request from self, dropping")
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
