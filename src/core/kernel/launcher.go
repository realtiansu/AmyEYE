package kernel

import (
	"core/common/cron"
	"core/common/errors"
	"core/config/jsonloader"
	"core/kernel/transmitter"
	"strconv"
	"strings"
	"sync"
)

type Kernel struct {
	Running bool
	Crontab *cron.Cron
	Mutex   sync.Mutex
	Channel chan string
}

func AddProtoToCron(pointAddress string, alias string, pointSize int, protocol string, cronSec string, crontab *cron.Cron, networkMutex *sync.Mutex, channel *chan string, ss []int, amount int16) error {
	var message = &transmitter.Message{
		Address:       pointAddress,
		Size:          pointSize,
		Sequence:      1,
		Protocol:      protocol,
		NetworkMutex:  networkMutex,
		Channel:       channel,
		Alias:         alias,
		SpecialSupply: ss,
		Amount:        amount,
	}
	err2 := crontab.AddJob(cronSec, message)
	if err2 != nil {
		return errors.New("unable to add point %s setting in %s", pointAddress, cronSec)
	}

	return nil
}

func AddPointToCron(point jsonloader.Point, crontab *cron.Cron, data *jsonloader.Data, networkMutex *sync.Mutex, channel *chan string) error {
	var addr string
	var ss []int
	for _, protocol := range point.Type {
		protocol = strings.ToLower(protocol)
		if protocol == "udp" ||  protocol == "kcp" || protocol == "quic"  {
			addr = point.Address + ":" + strconv.Itoa(data.RUDPP)
			if  protocol == "kcp" {
				ss = append(ss, data.DataShards)
				ss = append(ss, data.ParShards)
			}
		} else {
			addr = point.Address
		}

		err := AddProtoToCron(addr, point.Alias, point.Size, protocol, point.Crontab, crontab, networkMutex, channel, ss, point.Amount)
		if err != nil {
			return err
		}
	}

	return nil
}

func addPointHandlers(points *[]jsonloader.Point, crontab *cron.Cron, data *jsonloader.Data, networkMutex *sync.Mutex, channel *chan string) error {
	for _, point := range *points {
		err := AddPointToCron(point, crontab, data, networkMutex, channel)
		if err != nil {
			return errors.New("create instance fail")
		}
	}

	return nil
}

func New(pointsConfig *[]jsonloader.Point, data *jsonloader.Data, manual bool) (*Kernel, error) {
	points, err := pointsFormater(pointsConfig, manual)
	if err != nil {
		return nil, err
	}

	crontab := cron.New()
	networkMutex := sync.Mutex{}
	channel := make(chan string)
	var kernel = &Kernel{
		Running: false,
		Crontab: crontab,
		Mutex:   networkMutex,
		Channel: channel,
	}
	err = addPointHandlers(points, crontab, data, &networkMutex, &channel)

	return kernel, err
}

func (kernel *Kernel) Start() error {
	if kernel.Running {
		return errors.New("already running")
	}
	kernel.Running = true
	kernel.Crontab.Start()

	return nil
}

func (kernel *Kernel) Stop() {
	if kernel.Running {
		kernel.Running = false
		kernel.Crontab.Stop()
	}
}

func (kernel *Kernel) Run() {
	if kernel.Running {
		return
	}
	kernel.Running = true
	kernel.Crontab.Run()
}
