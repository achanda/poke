package utils

import (
	//"github.com/google/gopacket/routing"
	"math/rand"
	"net"
	"time"
)

func random(min, max int) int {
	rand.Seed(time.Now().Unix())
	return rand.Intn(max-min) + min
}

func GetLocalIP(dst string) (net.IPAddr, error) {
	serverAddr, err := net.ResolveUDPAddr("udp", net.JoinHostPort(dst, "12345"))
	if err != nil {
		return net.IPAddr{}, err
	}

	var conn *net.UDPConn
	if conn, err = net.DialUDP("udp", nil, serverAddr); err == nil {
		if udpaddr, ok := conn.LocalAddr().(*net.UDPAddr); ok {
			laddr, err := net.ResolveIPAddr("ip", udpaddr.IP.String())
			if err != nil {
				return net.IPAddr{}, err
			}
			return *laddr, nil
		}
	}
	return net.IPAddr{}, err
}
