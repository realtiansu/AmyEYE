package statistics

import (
	"sort"
)

func RttToAllStatistics(rtts []float64) []float64 {
	sort.Float64s(rtts)
	preRtts := preRtts(&rtts)

	rtt_avg := average(&preRtts)
	rtt_variance := variance(&preRtts, rtt_avg)
	rtt_quantile25, rtt_quantile50, rtt_quantile75 := quantile(&preRtts)
	rtt_min := preRtts[0]
	rtt_max := preRtts[len(preRtts)-1]

	jitters := jitter(&preRtts, rtt_avg)
	jitter_avg := average(&jitters)
	jitter_variance := variance(&jitters, jitter_avg)
	jitter_quantile25, jitter_quantile50, jitter_quantile75 := quantile(&jitters)
	jitter_min := jitters[0]
	jitter_max := jitters[len(preRtts)-1]

	loss_rate := lossRate(&rtts)

	//fmt.Println("rtt: ", rtts)
	//fmt.Println(rtt_avg, rtt_variance, rtt_min, rtt_max)
	//fmt.Println("rtt_quantile", rtt_quantile25, rtt_quantile50, rtt_quantile75)
	//fmt.Println("jitter", jitters)
	//fmt.Println(jitter_avg, jitter_variance, jitter_min, jitter_max)
	//fmt.Println("jitter_quantile", jitter_quantile25, jitter_quantile50, jitter_quantile75)
	//fmt.Println("loss rate", loss_rate)

	//					1			2			3		4			5				6				7
	return []float64{rtt_avg, rtt_variance, rtt_min, rtt_max, rtt_quantile25, rtt_quantile50, rtt_quantile75,
		//		8				9			10			11				12					13				14
		jitter_avg, jitter_variance, jitter_min, jitter_max, jitter_quantile25, jitter_quantile50, jitter_quantile75,
		//		15
		loss_rate}
}

func preRtts(rtts1 *[]float64) []float64 {
	var rtts2 []float64
	for _, rtt := range *rtts1 {
		if rtt > 0 {
			rtts2 = append(rtts2, rtt)
		}
	}
	if len(rtts2) == 0 {
		rtts2 = []float64{-1}
	}

	return rtts2
}