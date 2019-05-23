package transmitter

import (
	"fmt"
	"math/rand"
	"time"
)

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func Allocator(message *Message) {
	if message.Sequence <= message.Amount {

		protocol := message.Protocol
		if protocol == "icmp" {
			sendICMPPacket(message)
		}
		if protocol == "udp" {
			sendUDPPacket(message)
		}
		//if protocol == "tcp" {
		//	fmt.Printf("protocol:%s\n", protocol)
		//}
		if protocol == "kcp" {
			sendKcpPacket(message)
		}
		if protocol == "quic" {
			sendQuicPacket(message)
		}

	}

}

func sendFeedback(message *Message, proto string, duration float64, ttl int) {
	var name string
	if message.Alias != ""{
		name = message.Alias
	} else {
		name = message.Address
	}
	content := fmt.Sprintf("%d,%s,%s,%d,%.2f,%d,%d", message.Sequence, name, proto, message.Size, duration, ttl, time.Now().Unix())
	//fmt.Println(content)
	*message.Channel <- content
}

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
