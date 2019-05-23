package transmitter

import (
	"bytes"
	"core/common/reedsolomon"
	"fmt"
	"net"
	"time"
)

func sendKcpPacket(message *Message) {
	_ = kcp(message)
	message.Sequence += 1
}

func kcp(message *Message) error {
	conn, err := net.DialTimeout("udp4", message.Address, time.Second)
	if err != nil {
		return err
	}
	defer conn.Close()

	dataShards := message.SpecialSupply[0]
	parShards := message.SpecialSupply[1]
	//fmt.Println(dataShards, parShards)
	enc, err := reedsolomon.New(dataShards, parShards)
	if err != nil {
		fmt.Println(err)
		return err
	}
	mes := []byte( randSeq(message.Size) )
	shards, err := enc.Split(mes)

	if err != nil {
		fmt.Println(err)
		return err
	}
	err = enc.Encode(shards)
	if err != nil {
		fmt.Println(err)
		return err
	}
	var buffer bytes.Buffer
	for _, shard := range shards {
		buffer.Write(shard)
	}
	toSend := buffer.Bytes()

	startTime := time.Now()
	_ = conn.SetDeadline(startTime.Add(750 * time.Millisecond))
	message.NetworkMutex.Lock()
	_, err = conn.Write(toSend)								//发送
	message.NetworkMutex.Unlock()
	if err != nil {
		fmt.Println(err)
		return err
	}

	data := make([]byte, 1024)
	_, err = conn.Read(data)
	data = data[:len(toSend)]
	endDuration := float64(time.Since(startTime)) / (1000 * 1000)

	shards2 := deShard(data, len(shards))
	ok, err := enc.Verify(shards2)

	if ok {
		//fmt.Println("1")
		sendFeedback(message, "kcp", endDuration, -1)
	} else {
		err = enc.Reconstruct(shards2)
		if err != nil {
			fmt.Println(err)
			return err }
		ok, err = enc.Verify(shards2)
		if err != nil { fmt.Println(err);return err }
		if ok {
			//fmt.Println("2")
			sendFeedback(message, "kcp", endDuration, -1)
		} else {
			//fmt.Println("3")
			sendFeedback(message, "kcp",-1, -1)
		}
	}

	return nil
}

func deShard(data []byte, leng int) [][]byte {
	shards := make([][]byte, leng)
	step := len(data)/ leng

	for i := 0; i < leng; i++ {
		shards[i] = data[i*step:(i+1)*step]
	}

	return shards
}