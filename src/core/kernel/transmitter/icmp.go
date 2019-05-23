package transmitter

import (
	"fmt"
	"net"
	"time"
)

func sendICMPPacket(message *Message) {
	ping(message)
	message.Sequence += 1
}

func ping(message *Message) {
	if _, err := net.LookupHost(message.Address); err != nil { //ip不合法或域名不存在
		sendFeedback(message, "icmp", -1, -1)
	} else {
		const EchoRequestHeadLength = 8
		const EchoReplyHeadLength = 20
		id0, id1 := genIdentifier(message.Address)
		msg := make([]byte, message.Size+EchoRequestHeadLength)
		var timeout int64 = 2000

		msg[0] = 8                                      // echo,第一个字节表示报文类型，8表示回显请求
		msg[1] = 0                                      // code 0,ping的请求和应答，该code都为0
		msg[2] = 0                                      // checksum
		msg[3] = 0                                      // checksum
		msg[4], msg[5] = id0, id1                       //identifier[0] identifier[1], ID标识符 占2字节
		msg[6], msg[7] = genSequence(&message.Sequence) //sequence[0], sequence[1],序号占2字节
		length := message.Size + EchoRequestHeadLength
		check := checkSum(msg[0:length]) //计算检验和。
		msg[2] = byte(check >> 8)
		msg[3] = byte(check & 255)

		conn, err := net.Dial("ip4:icmp", message.Address)
		//conn, err = net.DialTimeout("ip:icmp", message.Address, time.Duration(timeout*1000*1000))
		if err != nil {
			fmt.Print(err)
			return
		}
		startTime := time.Now()
		_ = conn.SetDeadline(startTime.Add(700 * time.Millisecond)) //conn.SetReadDeadline可以在未收到数据的指定时间内停止Read等待，并返回错误err，然后判定请求超时
		message.NetworkMutex.Lock()
		_, err = conn.Write(msg[0:length]) //onn.Write方法执行之后也就发送了一条ICMP请求，同时进行计时和计次
		message.NetworkMutex.Unlock()
		//在使用Go语言的net.Dial函数时，发送echo request报文时，不用考虑i前20个字节的ip头；
		// 但是在接收到echo response消息时，前20字节是ip头。后面的内容才是icmp的内容，应该与echo request的内容一致
		receive := make([]byte, EchoReplyHeadLength+length)
		_, err = conn.Read(receive)
		endDuration := float64(time.Since(startTime)) / (1000 * 1000)
		//除了判断err!=nil，还有判断请求和应答的ID标识符，sequence序列码是否一致，以及ICMP是否超时（receive[ECHO_REPLY_HEAD_LEN] == 11，即ICMP报头的类型为11时表示ICMP超时）
		var ttl int
		if err != nil || receive[EchoReplyHeadLength+4] != msg[4] || receive[EchoReplyHeadLength+5] != msg[5] || receive[EchoReplyHeadLength+6] != msg[6] || receive[EchoReplyHeadLength+7] != msg[7] || endDuration >= float64(timeout) || receive[EchoReplyHeadLength] == 11 {
			ttl = -1
			endDuration = -1
		} else {
			ttl = int(receive[8])
		}
		_ = conn.Close()

		sendFeedback(message, "icmp", endDuration, ttl)
	}
}

func genIdentifier(host string) (byte, byte) {
	return host[0], host[1]
}
func genSequence(v *int16) (byte, byte) {
	ret1 := byte(*v >> 8)
	ret2 := byte(*v & 255)
	return ret1, ret2
}

func checkSum(msg []byte) uint16 {
	sum := 0
	length := len(msg)
	for i := 0; i < length-1; i += 2 {
		sum += int(msg[i])*256 + int(msg[i+1])
	}
	if length%2 == 1 {
		sum += int(msg[length-1]) * 256 // notice here, why *256?
	}
	sum = (sum >> 16) + (sum & 0xffff)
	answer := uint16(^sum)
	return answer
}
