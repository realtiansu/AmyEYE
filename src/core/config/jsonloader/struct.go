package jsonloader

type Config struct {
	Data   Data    `json:"settings"`
	Points []Point `json:"points"`
}

type Data struct {
	Dir string `json:"direct"`
	Hostname string `json:"hostname"`
	RTCPP int  `json:"RTCPP"`
	RUDPP int  `json:"RUDPP"`
	DataShards int `json:"dataShards"`
	ParShards  int `json:"parShards"`
}

type Point struct {
	Address string   `json:"address"`
	Alias   string   `json:"alias"`
	Type    []string `json:"type"`
	Size    int      `json:"size"`
	Amount  int16    `json:"amount"`
	Crontab string   `json:"crontab"`
}
