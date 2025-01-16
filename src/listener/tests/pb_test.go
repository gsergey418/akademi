package tests

import (
	"net"
	"testing"

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
	msg := &pb.BaseMessage{}
	msg.Content = "Hello, World!"
	data, err := proto.Marshal(msg)
	if err != nil {
		panic(err)
	}
	_, err = conn.Write(data)
	if err != nil {
		panic(err)
	}
}
