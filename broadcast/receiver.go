package broadcast

import (
	"fmt"
	"net"
)

type ReceivedPacket struct {
	Address net.Addr
	Packet  []byte
}

func (s *BroadcastService) initReceive() error {
	pc, err := net.ListenPacket("udp4", fmt.Sprintf(":%d", s.BroadcastPort))
	if err != nil {
		return err
	}

	go s.receiveLoop(pc)

	return nil
}

func (s *BroadcastService) receiveLoop(pc net.PacketConn) {
	defer pc.Close()

	for {
		buf := make([]byte, s.PacketMaxSize)
		n, addr, err := pc.ReadFrom(buf)
		if err != nil {
			return
		}

		packet := ReceivedPacket{
			Address: addr,
			Packet:  buf[:n],
		}

		for _, subscriber := range s.subscribers {
			subscriber <- packet
		}
	}
}
