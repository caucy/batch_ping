package main

import (
	"batch_ping/ping"
	"time"
	"fmt"
	"golang.org/x/net/icmp"
)

func main (){
	ipSlice := []string{}
	ipSlice = append(ipSlice, "122.228.74.183")
	ipSlice = append(ipSlice, "wwww.baidu.com")
	 ipSlice = append(ipSlice, "baidu.com")
	ipSlice = append(ipSlice, "121.42.9.142")
	ipSlice = append(ipSlice, "121.42.9.141")
	ipSlice = append(ipSlice, "121.42.9.144")
	ipSlice = append(ipSlice, "121.42.9.145")
	ipSlice = append(ipSlice, "121.42.9.146")
	ipSlice = append(ipSlice, "121.42.9.147")
	ipSlice = append(ipSlice, "121.42.9.148")
	ipSlice = append(ipSlice, "121.42.9.149")
	ipSlice = append(ipSlice, "121.42.9.150")


	bp, err := ping.NewBatchPinger(ipSlice, 4, time.Second*1, time.Second*10, true)

	if err != nil {
		fmt.Println(err)
	}

	bp.OnRecv = func(pkt *icmp.Echo, srcAddr string) {
		fmt.Printf("recv icmp_id=%d, icmp_seq=%d, srcAddr %v\n",
			pkt.ID, pkt.Seq, srcAddr)
	}

	bp.OnFinish = func(stMap map[string]*ping.Statistics) {
		for ip, st := range stMap{
			fmt.Printf("\n--- %s ping statistics ---\n", st.Addr)
			fmt.Printf("ip %s, %d packets transmitted, %d packets received, %v%% packet loss\n",ip,
				st.PacketsSent, st.PacketsRecv, st.PacketLoss)
			fmt.Printf("round-trip min/avg/max/stddev = %v/%v/%v/%v\n",
				st.MinRtt, st.AvgRtt, st.MaxRtt, st.StdDevRtt)
		}

	}

	bp.Run()
}