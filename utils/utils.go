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

func GetLocalIP(dst string) (net.IP, int, error) {
	serverAddr, err := net.ResolveUDPAddr("udp", net.JoinHostPort(dst, "12345"))
	if err != nil {
		return net.IP{}, -1, err
	}

	var conn *net.UDPConn
	if conn, err = net.DialUDP("udp", nil, serverAddr); err == nil {
		if udpaddr, ok := conn.LocalAddr().(*net.UDPAddr); ok {
			return udpaddr.IP, udpaddr.Port, nil
		}
	}
	return net.IP{}, -1, err
}
