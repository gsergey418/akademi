package tests

import (
	"fmt"
	"log"
	"net"

	"github.com/gsergey418alt/akademi/pb"
	"google.golang.org/protobuf/proto"
)

func main() {
	addr, err := net.ResolveUDPAddr("udp", "127.0.0.1:3865")
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	var udpReadBuffer [65535]byte
	for {
		l, host, err := conn.ReadFromUDP(udpReadBuffer[:])
		if err != nil {
			log.Print(err)
		}
		req := &pb.BaseMessage{}
		err = proto.Unmarshal(udpReadBuffer[:l], req)
		fmt.Println("Received request:", req)
		if err != nil {
			panic(err)
		}
		res := &pb.BaseMessage{}
		msg := &pb.ErrorMessage{Text: "Test error."}
		res.Message = &pb.BaseMessage_ErrorMessage{ErrorMessage: msg}
		fmt.Println("Sending response:", res)
		bytes, err := proto.Marshal(res)
		if err != nil {
			panic(err)
		}
		conn.WriteTo(bytes, host)
	}
}
