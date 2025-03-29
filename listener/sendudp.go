package listener

import (
	"log"
	"net"

	"github.com/gsergey418/akademi/pb"
	"google.golang.org/protobuf/proto"
)

// Sends bytes to remoteAddr over UDP.
func (u *UDPListener) sendUDPBytes(remoteAddr *net.UDPAddr, buf []byte) error {
	log.Print("Writing response to ", remoteAddr, ": ", len(buf), " bytes.")
	_, err := u.udpConn.WriteTo(buf, remoteAddr)
	return err
}

// Sends pb.BaseMessage to remoteAddr.
func (u *UDPListener) sendUDPMessage(remoteAddr *net.UDPAddr, res, req *pb.BaseMessage) error {
	u.populateDefaultResponse(res, req)
	resBytes, err := proto.Marshal(res)
	if err != nil {
		return err
	}
	err = u.sendUDPBytes(remoteAddr, resBytes)
	return err
}
