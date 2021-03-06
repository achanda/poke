package utils

import (
	"math"
	"net"
	"reflect"
	"testing"
)

func TestRandom(t *testing.T) {
	num := Random(1, 10)
	if num < 1 || num > 10 {
		t.Fatalf("Generated random number is outside range")
	}
}

func TestGetLocalIP(t *testing.T) {
	ip, port, err := GetLocalIP("127.0.0.1")
	if err != nil {
		t.Fatalf("%v", err)
	}
	if reflect.DeepEqual(ip, net.IP{}) {
		t.Fatalf("Did not expect empty IP: %v", ip)
	}
	if port == -1 {
		t.Fatalf("Did not expect port to be -1")
	}
}

func TestGetLocalIPInvalid(t *testing.T) {
	_, port, err := GetLocalIP("1234.2.3.4")
	if err == nil {
		t.Fatalf("Expected error")
	}
	if port != -1 {
		t.Fatalf("Expected to get back -1, got %v", port)
	}

	_, port1, err1 := GetLocalIP("blahblahblah")
	if err1 == nil {
		t.Fatalf("Expected error")
	}
	if port1 != -1 {
		t.Fatalf("Expected to get back -1, got %v", port1)
	}
}

func TestIsIPv4(t *testing.T) {
	ip1 := "106.10.138.240"
	ip2 := "2604:a880:1:20::9f9:9001"
	if !IsIPv4(ip1) {
		t.Fatalf("%v is a valid IPv4", ip1)
	}
	if IsIPv4(ip2) {
		t.Fatalf("%v is not a valid IPv4", ip2)
	}
}

func TestParseHost(t *testing.T) {
	host1 := "localhost"
	host2 := "106.10.138.240"
	host3 := "2604:a880:1:20::9f9:9001"
	host4 := "192.168.11.0/24"

	ip1, err := ParseHost(host1)
	if err != nil {
		t.Fatalf("%v", err)
	}
	if !reflect.DeepEqual(ip1[0], "127.0.0.1") {
		t.Fatalf("%v should resolve to a valid IPv4 address", ip1)
	}

	ip2, err := ParseHost(host2)
	if err != nil {
		t.Fatalf("%v", err)
	}
	if !reflect.DeepEqual(ip2[0], host2) {
		t.Fatalf("%v should be a valid IPv4 address", host2)
	}

	ip3, err := ParseHost(host3)
	if err != nil {
		t.Fatalf("%v", err)
	}
	if !reflect.DeepEqual(ip3[0], host3) {
		t.Fatalf("%v should be a valid IPv6 address", host3)
	}

	ip4, err := ParseHost(host4)
	expected_count := int(math.Pow(2, 32-24))
	if err != nil {
		t.Fatalf("%v", err)
	}
	if len(ip4) != expected_count {
		t.Fatalf("Expected to get %v, got %v", expected_count, len(ip4))
	}
}
