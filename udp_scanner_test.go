package poke

import (
	"testing"
)

// Assumes we are running IPv6

func TestUdpScannerv4(t *testing.T) {
	for i := 1; i < 100; i++ {
		scr4 := NewUdpScanner("127.0.0.1", uint64(i), true)
		res := scr4.Scan()
		if *res == (ScanResult{}) {
			t.Fatalf("Got empty result while scanning %v", i)
		}
		if res.Err != nil {
			t.Fatalf("Got an error while scanning %v: %v", i, res.Err)
		}
		if res.Port != uint64(i) {
			t.Fatalf("Got back %v while scanning %v", res.Port, i)
		}
	}
}

func TestUdpScannerv6(t *testing.T) {
	for i := 1; i < 100; i++ {
		scr6 := NewUdpScanner("::1", uint64(i), false)
		res := scr6.Scan()
		if *res == (ScanResult{}) {
			t.Fatalf("Got empty result while scanning %v", i)
		}
		if res.Err != nil {
			t.Fatalf("Got an error while scanning %v: %v", i, res.Err)
		}
		if res.Port != uint64(i) {
			t.Fatalf("Got back %v while scanning %v", res.Port, i)
		}
	}
}
