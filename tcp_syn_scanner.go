package poke

import (
	"fmt"
	"net"
)

type TcpSynScanner struct {
	Host string
	Port uint64
}

func newTcpSynScanner(host string, port uint64) Scanner {
	return TcpSynScanner{host, port}
}

func (tcpcs TcpSynScanner) Scan() *ScanResult {
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
