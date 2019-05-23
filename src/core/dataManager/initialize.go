package dataManager

import (
	"core/common/util"
	"core/config/jsonloader"
	"time"
)

type dataMan struct {
	Running bool
	data    map[string][]float64
	channel chan string
	fileDir string
	hostname string
}

func New(channel chan string, config *jsonloader.Data) (*dataMan, error) {
	hostname := config.Hostname
	if hostname == "" {
		hostname = util.GenHostname()
	}

	var dataman = &dataMan{
		Running: false,
		channel: channel,
		data: make(map[string][]float64),
		fileDir: config.Dir,
		hostname:hostname,
	}

	return dataman, nil
}

func (dataMan *dataMan) Start(manual bool) {
	dataMan.Running = true
	//go listenCmd(dataMan.Cmd, &dataMan.fileDir)
	go listenChannel(dataMan, manual)
}

func (dataMan *dataMan) Stop(runningTime int) {
	dataMan.Running = false
	time.Sleep(250 * time.Millisecond)
	//统计下
	_ = sortStatistic(&dataMan.data, dataMan.fileDir, dataMan.hostname, runningTime)
}
