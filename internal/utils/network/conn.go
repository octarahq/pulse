package network

import (
	gnet "github.com/shirou/gopsutil/v4/net"
)

func GetNetConnOpen() int {
	conns, err := gnet.Connections("tcp")
	if err != nil {
		return 0
	}
	return len(conns)
}

func GetNetConnTcp() int {
	conns, err := gnet.Connections("tcp")
	if err != nil {
		return 0
	}
	return len(conns)
}

func GetNetConnUdp() int {
	conns, err := gnet.Connections("udp")
	if err != nil {
		return 0
	}
	return len(conns)
}

func GetNetConnSsh() int {
	conns, err := gnet.Connections("tcp")
	if err != nil {
		return 0
	}
	count := 0
	for _, conn := range conns {
		if conn.Raddr.Port == 22 {
			count++
		}
	}
	return count
}

func GetNetConnEstablished() int {
	conns, err := gnet.Connections("tcp")
	if err != nil {
		return 0
	}
	count := 0
	for _, conn := range conns {
		if conn.Status == "ESTABLISHED" {
			count++
		}
	}
	return count
}
