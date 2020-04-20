package pingUtil

import (
	"fmt"
	"math"
	"sort"
)

type Statistic struct {
	PackTrans int		// total number of packets transmitted (request sent)
	PackRecv int		// total number of packets received (reply received)
	Rtts []float64		// array of round-trip time
	RttAvg float64		// average of round-trip time
	RttStd float64		// standard deviation of round-trip time
	Dst string			// destination hostname/IP address
	RttsLen int			// total number of round-trip times (length of Rtts)
}

func (s Statistic)PrintStats() {

	// prepare statistics
	s.RttsLen = len(s.Rtts)
	s.SortRtts()
	s.SetRttAvg()
	s.SetRttStd()

	// Statistics print statements
	fmt.Printf("\n--- %s ping statistics ---\n", s.Dst)

	fmt.Printf("%d packets transmitted, %d packets received, %.1f%% packet loss\n",
		s.PackTrans, s.PackRecv, (1.0 - float64(s.PackRecv) / float64(s.PackTrans)) * 100.0)

	fmt.Printf("round-trip min/avg/max/stddev = %.4f/%.4f/%.4f/%.4f ms\n",
		s.GetRttMin()*1000, s.RttAvg*1000, s.GetRttMax()*1000, s.RttStd*1000)
}

func (s Statistic)GetRttMin() float64{
	/*
		Get minimum round-trip time
	 */

	return s.Rtts[0]
}

func (s Statistic)GetRttMax() float64{
	/*
		Get maximum round-trip time
	 */

	return s.Rtts[s.RttsLen-1]
}

func (s *Statistic)SetRttAvg() {
	/*
		Save average round-trip time
	 */

	var total float64

	for i:=0; i < s.RttsLen; i++ {
		total += s.Rtts[i]
	}
	s.RttAvg = total / float64(s.RttsLen)
}

func (s *Statistic)SetRttStd() {
	/*
		Save standard deviation of round-trip time
	 */

	var total float64
	for i:=0; i < s.RttsLen; i++ {
		total += math.Pow(s.Rtts[i] - s.RttAvg, 2)
	}
	total /= float64(s.RttsLen)
	s.RttStd = math.Pow(total, 0.5)
}

func (s *Statistic)SortRtts() {
	/*
		sort the array of round-trip time so that
		the minimum is at starting index and the
		maximum is at the last index
	 */

	sort.Float64s(s.Rtts)
}