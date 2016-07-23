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
	Host string
	Port uint64
	Conn net.PacketConn
}

func NewTcpSynScanner(host string, port uint64) Scanner {
	conn, err := net.ListenPacket("ip4:tcp", "0.0.0.0")
	if err != nil {
		panic(err)
	}
	return TcpSynScanner{host, port, conn}
}

func (tcpcs TcpSynScanner) Scan() *ScanResult {
	dip := net.ParseIP(tcpcs.Host)
	saddr, sport, err := utils.GetLocalIP(tcpcs.Host)
	if err != nil {
		panic(err)
	}
	ip := &layers.IPv4{
		SrcIP:    saddr,
		DstIP:    dip,
		Protocol: layers.IPProtocolTCP,
	}
	tcp := &layers.TCP{
		SrcPort: layers.TCPPort(sport),
		DstPort: layers.TCPPort(tcpcs.Port),
		SYN:     true,
		Seq:     1105024978,
		Window:  14600,
	}
	tcp.SetNetworkLayerForChecksum(ip)
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
