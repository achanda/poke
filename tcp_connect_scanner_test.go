package poke

import (
	"testing"
)

func TestConnectScannerv4(t *testing.T) {
	for i := 1; i < 100; i++ {
		scr4 := NewTcpConnectScanner("0.0.0.0", uint64(i), true)
		res := scr4.Scan()
		// We should get a connection refused from all ports except 22
		// since we are ssh'd in. This does not take into account that
		// SSH can be running on another port
		if i != 22 && *res == (ScanResult{}) {
			t.Fatalf("Expected to get connection refused while scanning %v", i)
		}
		if i == 22 && res.Success == false {
			t.Fatalf("Expected port 22 to be open")
		}
		// For all ports, result should have the given port number
		if res.Port != uint64(i) {
			t.Fatalf("Got back %v while scanning %v", res.Port, i)
		}
	}
}

func TestConnectScannerv6(t *testing.T) {
	for i := 1; i < 100; i++ {
		scr6 := NewTcpConnectScanner("::1", uint64(i), false)
		res := scr6.Scan()
		if i != 22 && *res == (ScanResult{}) {
			t.Fatalf("Expected to error out while scanning %v", res.Err)
		}
		// Assumes SSH is running over IPv6
		if i == 22 && res.Success == false {
			t.Fatalf("Expected port 22 to be open")
		}
		if res.Port != uint64(i) {
			t.Fatalf("Got back %v while scanning %v", res.Port, i)
		}
	}
}
