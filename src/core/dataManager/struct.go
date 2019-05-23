package dataManager

type cell struct {
	Endpoint    string `json:"endpoint"`
	Metric      string `json:"metric"`
	Timestamp   int64  `json:"timestamp"`
	Step        int    `json:"step"`
	Value       float64 `json:"value"`
	CounterType string `json:"counterType"`
	Tags        string `json:"tags"`
}