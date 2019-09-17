# batch-ping


ICMP batch Ping library for Go, inspired by
[go-ping](https://github.com/sparrc/go-ping)

Here is a very simple example :


```go
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

```

It sends ICMP packet(s) and waits for a response. If it receives a response,
it calls the "receive" callback. When it's finished, it calls the "finish"
callback.

## Installation:

```
go get github.com/caucy/batch_ping
```


## Note on linux support :

This library attempts to send an
"unprivileged" ping via UDP. On linux, this must be enabled by setting

```
sudo sysctl -w net.ipv4.ping_group_range="0   2147483647"
```

## To do:
 1, bind source ip

 2, support ipv6

## Attention:
1, ping can support ping many ip, id is pid，and seq will grows

2, ip cannot be dereplic, this should be fix.

3, can not support use ipv4 and ipv6 at the same time

4, tick need close 


