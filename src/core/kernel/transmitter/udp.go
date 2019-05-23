package transmitter

import (
	"bytes"
	"net"
	"time"
)

func sendUDPPacket(message *Message) {
	_ = udp(message)
	message.Sequence += 1
}

func udp(message *Message) error {
	conn, err := net.DialTimeout("udp4", message.Address, time.Second)
	if err != nil {
		return err
	}
	defer conn.Close()

	toSend := []byte( randSeq(message.Size) )
	startTime := time.Now()
	_ = conn.SetDeadline(startTime.Add(750 * time.Millisecond))
	message.NetworkMutex.Lock()
	_, err = conn.Write(toSend)
	message.NetworkMutex.Unlock()
	if err != nil {
		return err
	}

	data := make([]byte, 1024)
	_, err = conn.Read(data)
	data = data[:len(toSend)]
	endDuration := float64(time.Since(startTime)) / (1000 * 1000)
	if err != nil {
		return err
	}

	if bytes.Equal(data, toSend) {
		sendFeedback(message, "udp", endDuration, -1)
	} else {
		sendFeedback(message, "udp",-1, -1)
	}

	return nil
}
