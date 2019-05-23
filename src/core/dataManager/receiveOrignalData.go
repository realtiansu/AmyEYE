package dataManager

import (
	"fmt"
	"strconv"
	"strings"
)

func listenChannel(man *dataMan, manual bool) {
	for true {
		if !man.Running {
			break
		}
		feedback := <-man.channel
		if manual {
			fmt.Println(feedback)
		}
		receiveOriginalData(feedback, man)
	}
}

func receiveOriginalData(feedback string, man *dataMan) {
	list := strings.Split(feedback, ",")
	option := fmt.Sprintf("%s-%s-%s", list[1], list[2], list[3])	//addr-proto-size
	duration, _ := strconv.ParseFloat(list[4], 64)
	man.data[option] = append(man.data[option], float64(duration))
}
