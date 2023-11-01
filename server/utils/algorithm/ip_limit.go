package algorithm

import (
	"fmt"
	"math/big"
	"net"
)

//IpLimitChoose 判断是否在某个IP区间
func IpLimitChoose(ip, startIP, endIP string) bool {
	ipN := InetAtoN(ip)
	ipStartN := InetAtoN(startIP)
	ipEndN := InetAtoN(endIP)
	if ipN >= ipStartN && ipN <= ipEndN {
		return true
	}
	return false
}

//InetNtoA 整型转IP
func InetNtoA(ip int64) string {
	return fmt.Sprintf("%d.%d.%d.%d",
		byte(ip>>24), byte(ip>>16), byte(ip>>8), byte(ip))
}

//InetAtoN IP转整型
func InetAtoN(ip string) int64 {
	ret := big.NewInt(0)
	ret.SetBytes(net.ParseIP(ip).To4())
	return ret.Int64()
}
