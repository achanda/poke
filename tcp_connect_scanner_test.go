package poke

import (
	"testing"
)

func TestConnectScannerv4(t *testing.T) {
	for i := 1; i < 100; i++ {
		//scr4 := NewTcpConnectScanner("127.0.0.1", uint64(i), true)
		scr4 := NewTcpConnectScanner("0.0.0.0", uint64(i), true)
		res := scr4.Scan()
		if i != 22 && *res == (ScanResult{}) {
			t.Fatalf("Expected to get connection refused while scanning %v", i)
		}
		if i == 22 && res.Success == false {
			t.Fatalf("Expected port 22 to be open")
		}
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
		if i == 22 && res.Success == false {
			t.Fatalf("Expected port 22 to be open")
		}
		if res.Port != uint64(i) {
			t.Fatalf("Got back %v while scanning %v", res.Port, i)
		}
	}
}
