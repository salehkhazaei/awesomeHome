package broadcast

import "testing"

func TestBroadcast(t *testing.T) {
	// define services
	broadcastService := NewBroadcastService(10240, 60504)

	// init services
	err := broadcastService.Init()
	if err != nil {
		panic(err)
	}

	// run forever
	channel := make(chan ReceivedPacket)
	broadcastService.Subscribe(channel)
	err = broadcastService.Broadcast([]byte("salam"))
	if err != nil {
		panic(err)
	}

	res := <-channel
	if string(res.Packet) != "salam" {
		t.Error(
			"Expected", "salam",
			"got", string(res.Packet),
		)
	}
}
