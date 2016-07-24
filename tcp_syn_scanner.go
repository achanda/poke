package poke

import (
	"fmt"
	"github.com/achanda/poke/utils"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"net"
	"time"
)

type TcpSynScanner struct {
	Host   string
	Port   uint64
	IsIPv4 bool
	Conn   net.PacketConn
}

func NewTcpSynScanner(host string, port uint64, ipver bool) Scanner {
	proto := "ip6:tcp"
	addr := net.IPv6zero.String()
	if ipver {
		proto = "ip4:tcp"
		addr = net.IPv4zero.String()
	}
	//fmt.Printf("%v %v\n", proto, addr)
	conn, err := net.ListenPacket(proto, addr)
	if err != nil {
		panic(err)
	}
	return TcpSynScanner{host, port, ipver, conn}
}

func (tcpcs TcpSynScanner) Scan() *ScanResult {
	dip := net.ParseIP(tcpcs.Host)
	saddr, sport, err := utils.GetLocalIP(tcpcs.Host)
	if err != nil {
		panic(err)
	}

	tcp := &layers.TCP{
		SrcPort: layers.TCPPort(sport),
		DstPort: layers.TCPPort(tcpcs.Port),
		SYN:     true,
		Seq:     1105024978,
		Window:  14600,
	}
	if tcpcs.IsIPv4 {
		ip4 := &layers.IPv4{
			SrcIP:    saddr,
			DstIP:    dip,
			Protocol: layers.IPProtocolTCP,
		}
		tcp.SetNetworkLayerForChecksum(ip4)
	} else {
		ip6 := &layers.IPv6{
			SrcIP:      saddr,
			DstIP:      dip,
			NextHeader: layers.IPProtocolTCP,
		}
		tcp.SetNetworkLayerForChecksum(ip6)
	}
	buf := gopacket.NewSerializeBuffer()
	opts := gopacket.SerializeOptions{
		ComputeChecksums: true,
		FixLengths:       true,
	}
	if err := gopacket.SerializeLayers(buf, opts, tcp); err != nil {
		fmt.Printf("%v", err)
	}
	if _, err := tcpcs.Conn.WriteTo(buf.Bytes(), &net.IPAddr{IP: dip}); err != nil {
		panic(err)
	}

	if err := tcpcs.Conn.SetDeadline(time.Now().Add(10 * time.Second)); err != nil {
		panic(err)
	}

	for {
		buf := make([]byte, 4096)
		n, addr, err := tcpcs.Conn.ReadFrom(buf)
		if err != nil {
			return &ScanResult{}
		}
		if addr.String() == dip.String() {
			packet := gopacket.NewPacket(buf[:n], layers.LayerTypeTCP, gopacket.Default)
			if tcpLayer := packet.Layer(layers.LayerTypeTCP); tcpLayer != nil {
				tcp, _ := tcpLayer.(*layers.TCP)
				if tcp.DstPort == layers.TCPPort(sport) {
					return &ScanResult{Port: tcpcs.Port, Success: !tcp.RST, Err: nil}
				}
			}
		}
	}
}
