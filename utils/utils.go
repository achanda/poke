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

func getIPs(cidr string) ([]net.IP, error) {
	var hosts []net.IP
	ip, ipnet, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil, err
	}
	for ip := ip.Mask(ipnet.Mask); ipnet.Contains(ip); increment(ip) {
		hosts = append(hosts, ip)
	}
	// Remove network and broadcast
	return hosts[1 : len(hosts)-1], nil
}

func ParseHost(host string) ([]net.IP, error) {
	// Check if we got a CIDR
	ips, err := getIPs(host)
	if err == nil {
		return ips, err
	}

	var hosts []net.IP
	addr := net.ParseIP(host)
	if addr == nil {
		ipa, err := net.ResolveIPAddr("ip", host)
		if err != nil {
			hosts = append(hosts, net.IP{})
			return hosts, err
		}
		hosts = append(hosts, ipa.IP)
		return hosts, nil
	}
	hosts = append(hosts, addr)
	return hosts, nil
}

func increment(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}
