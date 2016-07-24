package utils

import (
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
