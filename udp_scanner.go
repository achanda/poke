package poke

import (
	"fmt"
	"github.com/achanda/poke/utils"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"net"
	"time"
)

type UdpScanner struct {
	Host   string
	Port   uint64
	IsIPv4 bool
	Conn   net.PacketConn
}

func NewUdpScanner(host string, port uint64, ipver bool) Scanner {
	proto := "ip6:tcp"
	addr := net.IPv6zero.String()
	if ipver {
		proto = "ip4:tcp"
		addr = net.IPv4zero.String()
	}
	conn, err := net.ListenPacket(proto, addr)
	if err != nil {
		panic(err)
	}
	return UdpScanner{host, port, ipver, conn}
}

func (us UdpScanner) Scan() *ScanResult {
	defer us.Conn.Close()
	dip := net.ParseIP(us.Host)
	saddr, sport, err := utils.GetLocalIP(us.Host)
	if err != nil {
		panic(err)
	}
	udp := &layers.UDP{
		SrcPort: layers.UDPPort(sport),
		DstPort: layers.UDPPort(us.Port),
	}
	if us.IsIPv4 {
		ip4 := &layers.IPv4{
			SrcIP:    saddr,
			DstIP:    dip,
			Protocol: layers.IPProtocolTCP,
		}
		udp.SetNetworkLayerForChecksum(ip4)
	} else {
		ip6 := &layers.IPv6{
			SrcIP:      saddr,
			DstIP:      dip,
			NextHeader: layers.IPProtocolTCP,
		}
		udp.SetNetworkLayerForChecksum(ip6)
	}
	buf := gopacket.NewSerializeBuffer()
	opts := gopacket.SerializeOptions{
		ComputeChecksums: true,
		FixLengths:       true,
	}
	if err := gopacket.SerializeLayers(buf, opts, udp); err != nil {
		fmt.Printf("%v", err)
	}
	if _, err := us.Conn.WriteTo(buf.Bytes(), &net.IPAddr{IP: dip}); err != nil {
		panic(err)
	}

	if err := us.Conn.SetDeadline(time.Now().Add(10 * time.Second)); err != nil {
		panic(err)
	}

	for {
		buf2 := make([]byte, 4096)
		n, addr, err := us.Conn.ReadFrom(buf2)
		if err != nil {
			return &ScanResult{}
		}
		if addr.String() == dip.String() {
			packet := gopacket.NewPacket(buf2[:n], layers.LayerTypeUDP, gopacket.Default)
			if udp := packet.Layer(layers.LayerTypeUDP); udp != nil {
				return &ScanResult{Port: us.Port, Success: true, Err: nil}
			}
		}
	}
}
