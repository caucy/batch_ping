# batch-ping


ICMP batch Ping library for Go, inspired by
[go-ping](https://github.com/sparrc/go-ping)

Here is a very simple example :


```go
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

	bp, err := ping.NewBatchPinger(ipSlice, false) // true will need to be root

	if err != nil {
		log.Fatalf("new batch ping err %v", err)
	}
	bp.SetDebug(false) // debug == true will fmt debug log

	bp.SetSource("") // if hava multi source ip, can use one isp
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



```

It sends ICMP packet(s) and waits for a response. If it receives a response,
it calls the "receive" callback. When it's finished, can call the "OnFinish"
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

## feature:
 
#### 1, bind source ip

 ```
bp.SetSource("") // if hava multi isp ip, can use one 
 ```
 

#### 2, support ipv6

can  support use ipv4 and ipv6 at the same time

#### 3, support ping multi ip 

NewBatchPinger can support multi ip ping

#### 4, support two model

can use unprivileged mode , need not to be root


## Attention:
ping can support ping many ip, id is pid，and the addr's seq is the same.

fix same bug of github.com/sparrc/go-ping, such as if the dst server use the iptable ban ping, go-ping will hang .



