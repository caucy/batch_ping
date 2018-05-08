package ping

import (
	"net"
	"time"
	"fmt"
	"sync"
	"os"
	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
	"golang.org/x/net/ipv6"
)

var IsFinished bool = false

type BatchPinger struct {
	// pingers []*Pinger

	// ip 反查输入域名与ip
	domainIpMap map[string]string

	MapPinger map[string]*Pinger

	// Interval is the wait time between each packet send. Default is 1s.
	Interval time.Duration

	// Timeout specifies a timeout before ping exits, regardless of how many
	// packets have been received.
	Timeout time.Duration

	// Count tells pinger to stop after sending (and receiving) Count echo
	// packets. If this option is not specified, pinger will operate until
	// interrupted.

	// count 是单个ip ping 的次数
	Count int

	// totalcount 是 n* count ，这里注意去重
	TotalCount int

	// record send count
	SendCount int

	ReceivCount int

	// Debug runs in debug mode
	Debug bool

	// OnRecv is called when Pinger receives and processes a packet
	OnRecv func(*icmp.Echo, string)

	// OnFinish is called when Pinger exits
	OnFinish func(map[string]*Statistics)

	// stop chan bool
	ipv4 bool
	source  string
	network string
	Id 	int
}


// NewPinger returns a new Pinger struct pointer
// interval in secs
func NewBatchPinger(ipSlice []string, count int, interval time.Duration,
	timeout time.Duration, ipv4 bool) (*BatchPinger, error) {
	var batachPinger = BatchPinger{
		Interval: interval,
		Timeout:  timeout,
		Count:  count,
		network:  "ip",
		ipv4: 	  ipv4,
		Id: getId(),
		MapPinger: make(map[string]*Pinger),
		domainIpMap: make(map[string]string),
	}
	id :=getId()

	for _, ipDomain := range ipSlice {
		//检查域名和ip 是否合法
		ipaddr, err := net.ResolveIPAddr("ip", ipDomain)
		if err != nil {
			return nil, err
		}

		if isIPv4(ipaddr.IP) {
			ipv4 = true
		} else if isIPv6(ipaddr.IP) {
			ipv4 = false
		}
		pinger := newPinger(ipaddr, id, ipv4)
		//利用map 去重 
		ipaddrStr := ipaddr.String()
		batachPinger.MapPinger[ipaddrStr] = pinger
		batachPinger.domainIpMap[ipaddrStr] = ipDomain
	}

	batachPinger.TotalCount = batachPinger.Count * len(batachPinger.MapPinger)

	return &batachPinger, nil
}

// id 取进程 id
func getId() int {
	return os.Getpid()
}

func (bp *BatchPinger) Run() {
	var conn *icmp.PacketConn
	if bp.ipv4 {
		fmt.Printf("source:%v, network:%v\n", bp.source, bp.network)
		if conn = bp.Listen(ipv4Proto[bp.network], bp.source); conn == nil {
			return
		}
	} else {
		if conn = bp.Listen(ipv6Proto[bp.network], bp.source); conn == nil {
			return
		}
	}
	defer conn.Close()
	defer bp.finish()

	var wg sync.WaitGroup

	wg.Add(2)
	go bp.RecvICMP(conn, &wg, )
	go bp.SendICMP(conn, &wg)

	// 发送icmp, 传递整体seq, 整体计数加一
	wg.Wait()
	fmt.Println("it is time to exit")
}

func (bp *BatchPinger) Listen(netProto string, source string) *icmp.PacketConn {
	conn, err := icmp.ListenPacket(netProto, source)
	if err != nil {
		fmt.Printf("Error listening for ICMP packets: %s\n", err.Error())
		return nil
	}
	return conn
}

func (bp *BatchPinger) RecvICMP(conn *icmp.PacketConn, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		if ! IsFinished{
			bytes := make([]byte, 512)
			conn.SetReadDeadline(time.Now().Add(time.Millisecond * 100))
			n, srcAdd, err := conn.ReadFrom(bytes)
			if err != nil {
				if neterr, ok := err.(*net.OpError); ok {
					if neterr.Timeout() {
						// Read timeout
						continue
					} 
				}
			}
			recvPkg := &packet{bytes: bytes, nbytes: n, addr: srcAdd.String()}
			bp.ProcessPacket(recvPkg)

		}else {
			break
		}
			
	}
}

func (bp *BatchPinger) SendICMP(conn *icmp.PacketConn, wg *sync.WaitGroup) {
	defer wg.Done()

	err := bp.BatchSendICMP(conn)
	if err != nil {
		fmt.Println(err.Error())
	}
	timeout := time.NewTicker(bp.Timeout)
	interval := time.NewTicker(bp.Interval)

	// 每次更新判断是否结束，超时判断整体是否结束

	for {
		if ! IsFinished {
			select {
			case <-timeout.C:
				fmt.Println("tick timeout")
				IsFinished = true
	
			case <-interval.C:
				fmt.Println("tick interval")
				err = bp.BatchSendICMP(conn)
				if err != nil {
					fmt.Println("FATAL: ", err.Error())
				}
			}
		}else{
			break
		}
		
	}
	
}

func (bp *BatchPinger) BatchSendICMP(conn *icmp.PacketConn) error {
	//检查是否全部发送完
	if bp.SendCount >= bp.TotalCount{
		return nil
	}

	for ip, pinger := range bp.MapPinger {
		err := pinger.SendICMP(conn, bp.SendCount)
		fmt.Printf("ip %v, sended\n", ip)
		bp.SendCount += 1
		if err != nil {
			fmt.Printf("icmp send err:%v\n", err)
			return err
		}
	}
	return nil
}

func (bp *BatchPinger) ProcessPacket(recv *packet) error {
	
	var bytes []byte
	var proto int
	if bp.ipv4 {
		if bp.network == "ip" {
			bytes = ipv4Payload(recv.bytes)
		} else {
			bytes = recv.bytes
		}
		proto = protocolICMP
	} else {
		bytes = recv.bytes
		proto = protocolIPv6ICMP
	}
	srcAddr := recv.addr

	var m *icmp.Message
	var err error
	if m, err = icmp.ParseMessage(proto, bytes[:recv.nbytes]); err != nil {
		fmt.Println("Error parsing icmp message")
		return nil
	}

	if m.Type != ipv4.ICMPTypeEchoReply && m.Type != ipv6.ICMPTypeEchoReply {
		// Not an echo reply, ignore it
		fmt.Println("not reply")
		return nil
	}
	body := m.Body.(*icmp.Echo)
	
	if body.ID != bp.Id{
		fmt.Println("not the process id")
		return nil
	}


	rtt := time.Since(bytesToTime(body.Data[:timeSliceLength]))
	pinger := bp.MapPinger[srcAddr]

	pinger.rtts = append(pinger.rtts, rtt)
	pinger.PacketsRecv += 1
	bp.ReceivCount += 1

	handler := bp.OnRecv
	if handler != nil {
		handler(body, srcAddr)
	}

	// 退出条件是， 发送数据包 == 待发送所有包
	if bp.ReceivCount == bp.TotalCount{
		IsFinished = true
	}
	return nil
	
}

func (bp *BatchPinger) finish() {
	handler := bp.OnFinish
	if handler != nil {
		s := bp.Statistics()
		handler(s)
	}
}

func (bp *BatchPinger) Statistics() map[string]*Statistics {
	stMap := map[string]*Statistics{}
	for ip, pinger := range bp.MapPinger{
		x := pinger.Statistics()
		ipDomain := bp.domainIpMap[ip]
		stMap[ipDomain] = x
	}
	return stMap
}

