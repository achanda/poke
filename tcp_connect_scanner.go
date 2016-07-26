package poke

// This file implements a simple TCP connect scanner

import (
	"net"
	"strconv"
)

type TcpConnectScanner struct {
	Host string
	Port uint64
	Isv4 bool
}

func NewTcpConnectScanner(host string, port uint64, ipver bool) Scanner {
	return TcpConnectScanner{host, port, ipver}
}

func (tcpcs TcpConnectScanner) Scan() *ScanResult {
	proto := "tcp6"
	if tcpcs.Isv4 {
		proto = "tcp4"
	}
	conn, err := net.Dial(proto, net.JoinHostPort(tcpcs.Host, strconv.FormatUint(tcpcs.Port, 10)))
	result := ScanResult{
		Port:    tcpcs.Port,
		Success: err == nil,
		Err:     err,
	}
	if conn != nil {
		conn.Close()
	}
	return &result
}
