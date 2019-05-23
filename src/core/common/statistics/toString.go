package statistics

import (
	"fmt"
	"strings"
)

func FloatSliceToString(data *[]float64) string {
	str := ""
	for _, num := range *data {
		str = fmt.Sprintf("%s,%.2f",str, num)
	}

	return strings.TrimLeft(str, ",")
}