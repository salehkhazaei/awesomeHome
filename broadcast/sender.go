package broadcast

import (
	"fmt"
	"ir.skhf/awesomeHome/utils"
	"net"
)

func (s *BroadcastService) sendPacket(data []byte) error {
	address := fmt.Sprintf("%s:%d", s.currentBroadcastAddress(), s.BroadcastPort)
	remote, err := net.ResolveUDPAddr("udp4", address)
	if err != nil {
		return err
	}
	list, err := net.DialUDP("udp4", nil, remote)
	if err != nil {
		return err
	}
	defer list.Close()

	_, err = list.Write(data)
	if err != nil {
		return err
	}

	return nil
}

func (s *BroadcastService) currentBroadcastAddress() string {
	mask := net.CIDRMask(24, 32)
	ip := utils.GetOutboundIP()

	broadcast := net.IP(make([]byte, 4))
	for i := range ip {
		broadcast[i] = ip[i] | ^mask[i]
	}

	return broadcast.String()
}
