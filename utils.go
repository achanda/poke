package poke

import (
	//"github.com/google/gopacket/routing"
	"math/rand"
	//"net"
	"time"
)

func random(min, max int) int {
	rand.Seed(time.Now().Unix())
	return rand.Intn(max-min) + min
}
