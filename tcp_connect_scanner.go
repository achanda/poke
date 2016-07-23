package poke

import (
	"fmt"
	"net"
)

type TcpConnectScanner struct {
	Host string
	Port uint64
}

func NewTcpConnectScanner(host string, port uint64) Scanner {
	return TcpConnectScanner{host, port}
}

func (tcpcs TcpConnectScanner) Scan() *ScanResult {
	conn, err := net.Dial("tcp", fmt.Sprintf("%v:%v", tcpcs.Host, tcpcs.Port))
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
