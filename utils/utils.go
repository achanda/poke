package utils

import (
	"math/rand"
	"net"
	"time"
)

func Random(min, max int) int {
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

func IsIPv4(host string) bool {
	ip := net.ParseIP(host)
	if ip != nil {
		return ip.To4() != nil
	}
	return false
}

func ParseHost(host string) (net.IP, error) {
	addr := net.ParseIP(host)
	if addr == nil {
		ipa, err := net.ResolveIPAddr("ip", host)
		if err != nil {
			return net.IP{}, err
		}
		return ipa.IP, nil
	}
	return addr, nil
}
