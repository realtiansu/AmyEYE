package dataManager

import (
	"core/common/statistics"
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"
	"time"
)


func sortStatistic(data *map[string][]float64, fileDir string, hostname string, runningTime int) error {
	timeStr:=time.Now().Format("2006-01-02 15:04:05")
	sortedData := make(map[string][]float64)
	for option, cell := range *data {
		result := statistics.RttToAllStatistics(cell)
		sortedData[option] = result

		//写文件
		content := statistics.FloatSliceToString(&result)
		content = fmt.Sprintf("%s,%s", timeStr, content)
		filePath := filepath.Join(fileDir, fmt.Sprintf("%sbytes.txt", option))
		_ = record(content, filePath)
	}

	_ = genJson(&sortedData, runningTime, hostname)

	return nil
}

func genJson(sortedData *map[string][]float64, runningTime int, hostname string) error {
	mapp := mapping()
	timestamp := time.Now().Unix()
	var cells []cell

	//for option, _ := range *sortedData {
	//	fmt.Println(option)
	//}

	for option, data := range *sortedData {
		fileNameDetail := strings.Split(option, "-")
		for i, val := range data {
			var cell = cell{
				Endpoint: hostname,
				Metric: fileNameDetail[1] + "_to_" + fileNameDetail[0],
				Timestamp: timestamp,
				Step: runningTime + 5,
				Value: val,
				CounterType: "GAUGE",
				Tags: "detail=" + mapp[i] + ",packet_size=" + fileNameDetail[2] + ",",
			}
			cells = append(cells, cell)
		}
	}

	jsonStr, err := json.Marshal(cells)
	if err != nil {
		return err
	}
	fmt.Println(string(jsonStr))
	return nil
}

func mapping() []string {
	return []string{
		"rtt_avg",
		"rtt_var",
		"rtt_min",
		"rtt_max",
		"rtt_quantile25",
		"rtt_quantile50",
		"rtt_quantile75",
		"jitter_avg",
		"jitter_var",
		"jitter_min",
		"jitter_max",
		"jitter_quantile25",
		"jitter_quantile50",
		"jitter_quantile75",
		"loss_rate",
	}
}