package tests

import (
	"fmt"
	"net"
	"testing"

	"github.com/gsergey418alt/akademi/core"
	"github.com/gsergey418alt/akademi/pb"
	"google.golang.org/protobuf/proto"
)

func TestProtobuf(t *testing.T) {
	listenAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:3865")
	if err != nil {
		panic(err)
	}
	conn, err := net.DialUDP("udp", nil, listenAddr)
	if err != nil {
		panic(err)
	}
	req := &pb.BaseMessage{}
	req.Message = &pb.BaseMessage_PingRequest{}
	fmt.Println(req)
	data, err := proto.Marshal(req)
	if err != nil {
		panic(err)
	}
	_, err = conn.Write(data)
	if err != nil {
		panic(err)
	}
	var buf [65535]byte
	l, err := conn.Read(buf[:])
	res := &pb.BaseMessage{}
	err = proto.Unmarshal(buf[:l], res)
	if err != nil {
		panic(err)
	}
	nodeID := core.BaseID(res.NodeID)
	fmt.Println("Received ping response. NodeID: ", nodeID.Base64Str())
}
