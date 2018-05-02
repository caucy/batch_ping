package main 

import (
	"fmt"
	"net"
)

func main(){
	ip, err := net.ResolveIPAddr("ip", "www.baidu.com")
	fmt.Println(ip.String(), err)
}