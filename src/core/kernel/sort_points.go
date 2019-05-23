package kernel

import (
	"core/common/errors"
	"core/config/jsonloader"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func pointsFormater(points *[]jsonloader.Point, manual bool) (*[]jsonloader.Point, error) {
	var pointsNew []jsonloader.Point
	for _, subConfig := range *points {
		if strings.HasPrefix(subConfig.Address, "list:") {
			point := subConfig
			alias := pointList(strings.TrimPrefix(subConfig.Address, "list:"))
			if points == nil {
				return nil, errors.New("can't open ip list file " + subConfig.Address)
			}
			for addr, alia := range alias {
				point.Address = addr
				point.Alias = alia
				pointsNew = append(pointsNew, point)
			}
		} else {
			pointsNew = append(pointsNew, subConfig)
		}
	}

	if manual {
		for _, pointName := range pointsNew {
			fmt.Println(pointName)
		}
		fmt.Println("--------------------------")
	}

	return &pointsNew, nil
}

func pointList(pointsFile string) map[string]string {
	alias := make(map[string]string)

	if isExist(pointsFile) {
		data, _ := ioutil.ReadFile(pointsFile)
		points := string(data)
		points = strings.TrimSuffix(points, "\n")
		points = strings.Replace(points, "\r", "", -1)
		line := strings.Split(points, "\n")

		for _, content := range line{
			ksp := strings.Split(content, "=")
			if len(ksp) == 2 {
				ali := strings.Trim(ksp[0], " ")
				addr := strings.Trim(ksp[1], " ")
				alias[addr] = ali
			} else if len(ksp) == 1 {
				alias[ksp[0]] = ""
			}
		}

		return alias
	}
	return nil
}

func isExist(f string) bool {
	_, err := os.Stat(f)
	return err == nil || os.IsExist(err)
}
