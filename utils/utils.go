package utils

// This package provides some utility functions

import (
	"math/rand"
	"net"
	"time"
)

// Generates a random int between max and min
func Random(min, max int) int {
	rand.Seed(time.Now().Unix())
	return rand.Intn(max-min) + min
}

// Gets the local IP and port for a given destination address
// local IP type matches that of the destination
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

// Tries to check if the given string is a valid IPv4 address
func IsIPv4(host string) bool {
	ip := net.ParseIP(host)
	if ip != nil {
		return ip.To4() != nil
	}
	return false
}

// Given a CIDR, gets all IP addresses in it including net and broadcast
func getIPs(cidr string) ([]string, error) {
	var hosts []string
	ip, ipnet, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil, err
	}
	for ip := ip.Mask(ipnet.Mask); ipnet.Contains(ip); increment(ip) {
		hosts = append(hosts, ip.String())
	}
	return hosts, nil
}

// Parses the host param. If it's a name, tries to resolve it
// If it's a CIDR, returns all IP addresses in it
func ParseHost(host string) ([]string, error) {
	// Check if we got a CIDR
	ips, err := getIPs(host)
	if err == nil {
		return ips, err
	}

	var hosts []string
	addr := net.ParseIP(host)
	if addr == nil {
		ipa, err := net.ResolveIPAddr("ip", host)
		if err != nil {
			hosts = append(hosts, net.IP{}.String())
			return hosts, err
		}
		hosts = append(hosts, ipa.IP.String())
		return hosts, nil
	}
	hosts = append(hosts, addr.String())
	return hosts, nil
}

// Given an IP, gets the next one
func increment(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}

// Tries to check if the given host is up
func IsHostUp(host string) bool {
	_, err := net.DialTimeout("ip", host, time.Duration(10))
	if err == nil {
		return true
	}
	return false
}
