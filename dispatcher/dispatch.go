package dispatcher

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/gsergey418/akademi/core"
	"github.com/gsergey418/akademi/pb"
	"google.golang.org/protobuf/proto"
)

// dispatchUDPBytes is a function that manages sending
// raw bytes over a UDP connection.
func (u *UDPDispatcher) dispatchUDPBytes(host core.Host, buf []byte) ([]byte, error) {
	conn, err := net.Dial("udp", string(host))
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	err = conn.SetDeadline(time.Now().Add(5 * time.Second))
	if err != nil {
		return nil, err
	}

	self := strings.Split(conn.LocalAddr().String(), ":")[0] + ":" + strconv.Itoa(int(u.RoutingHeader.ListenPort))
	if self == string(host) {
		return nil, fmt.Errorf("Refusing to dispatch request to self.")
	}

	log.Print("Dispatching request to ", host, ": ", len(buf), " bytes.")
	_, err = conn.Write(buf)
	if err != nil {
		return nil, err
	}
	var udpReadBuffer [65535]byte
	l, err := conn.Read(udpReadBuffer[:])
	if err != nil {
		return nil, err
	}
	log.Print("Response from ", host, ": ", l, " bytes.")
	return udpReadBuffer[:l], nil
}

// dispatchUDPMessage is a function that wraps around
// dispatchUDPBytes, operating on pb.BaseMessage structs.
func (u *UDPDispatcher) dispatchUDPMessage(host core.Host, req *pb.BaseMessage) (*pb.BaseMessage, error) {
	u.writeRoutingHeader(req)
	buf, err := proto.Marshal(req)
	if err != nil {
		return nil, err
	}
	resBytes, err := u.dispatchUDPBytes(host, buf)
	if err != nil {
		return nil, err
	}
	res := &pb.BaseMessage{}
	err = proto.Unmarshal(resBytes, res)
	if err != nil {
		return nil, err
	}
	if res.GetErrorMessage() != nil {
		return nil, fmt.Errorf("%s", res.GetErrorMessage().Text)
	}
	if bytes.Compare(res.RequestID, req.RequestID) != 0 {
		return nil, fmt.Errorf("SessionError: Request and response RequestIDs don't match!")
	}
	return res, nil
}
