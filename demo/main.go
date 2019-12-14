package main

import (
	"log"

	"github.com/caucy/batch_ping"
)

func main() {
	ipSlice := []string{}
	// ip list should not more than 65535

	ipSlice = append(ipSlice, "2400:da00:2::29") //support ipv6
	ipSlice = append(ipSlice, "baidu.com")

	bp, err := ping.NewBatchPinger(ipSlice, true) // true will need to be root, false may be permission denied

	if err != nil {
		log.Fatalf("new batch ping err %v", err)
	}
	bp.SetDebug(false) // debug == true will fmt debug log

	bp.SetSource("") // if have multi source ip, can use one ip
	bp.SetCount(10)

	bp.OnFinish = func(stMap map[string]*ping.Statistics) {
		for ip, st := range stMap {
			log.Printf("\n--- %s ping statistics ---\n", st.Addr)
			log.Printf("ip %s, %d packets transmitted, %d packets received, %v%% packet loss\n", ip,
				st.PacketsSent, st.PacketsRecv, st.PacketLoss)
			log.Printf("round-trip min/avg/max/stddev = %v/%v/%v/%v\n",
				st.MinRtt, st.AvgRtt, st.MaxRtt, st.StdDevRtt)
			log.Printf("rtts is %v \n", st.Rtts)
		}
	}

	err = bp.Run()
	if err != nil {
		log.Printf("run err %v \n", err)
	}
	bp.OnFinish(bp.Statistics())
}
