package transmitter

import "sync"

type Message struct {
	Address       string
	Size          int
	Sequence      int16
	Protocol      string
	NetworkMutex  *sync.Mutex
	Channel       *chan string
	Alias         string
	Rport         int
	SpecialSupply []int
	Amount        int16
}
