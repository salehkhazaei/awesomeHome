package broadcast

type BroadcastService struct {
	PacketMaxSize int
	BroadcastPort int

	subscribers []chan ReceivedPacket
}

func NewBroadcastService(packetMaxSize int, broadcastPort int) *BroadcastService {
	return &BroadcastService{
		PacketMaxSize: packetMaxSize,
		BroadcastPort: broadcastPort,
		subscribers:   make([]chan ReceivedPacket, 0),
	}
}

func (s *BroadcastService) Init() error {
	return s.initReceive()
}

func (s *BroadcastService) Subscribe(channel chan ReceivedPacket) {
	s.subscribers = append(s.subscribers, channel)
}

func (s *BroadcastService) Broadcast(data []byte) error {
	return s.sendPacket(data)
}
