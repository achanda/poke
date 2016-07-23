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
	//ipa, err := net.ResolveIPAddr("ip", tcpcs.Host)
	//if err != nil {
	//	return nil
	//return &ScanResult{Port: 1234, Success: true, Err: nil}
	//}
	//fmt.Println(ipa)
	//conn, err := net.ListenIP("ip:tcp", ipa)
	//if err != nil {
	//	panic(err)
	//return nil
	//return &ScanResult{Port: 1234, Success: true, Err: nil}
	//}
	//defer conn.Close()

	//sp := random(1024, 2024)
	dip := net.ParseIP(tcpcs.Host)
	//sip, err := getIP(dip)
	//fmt.Printf("%v", sip)
	//if err != nil {
	//return nil
	//	return &ScanResult{Port: 1234, Success: true, Err: nil}
	//}
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
	//fmt.Printf("Dst: %v\n", tcpcs.Port)
	tcp.SetNetworkLayerForChecksum(ip)
	buf := gopacket.NewSerializeBuffer()
	opts := gopacket.SerializeOptions{
		ComputeChecksums: true,
		FixLengths:       true,
	}
	if err := gopacket.SerializeLayers(buf, opts, tcp); err != nil {
		fmt.Printf("%v", err)
		//return &ScanResult{Port: 1234, Success: true, Err: nil}
	}
	if _, err := tcpcs.Conn.WriteTo(buf.Bytes(), &net.IPAddr{IP: dip}); err != nil {
		//panic(err)
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
					//fmt.Printf("Remote port %v ", tcp.DstPort)
					//fmt.Printf("SYN %v ", tcp.SYN)
					//fmt.Printf("ACK %v ", tcp.ACK)
					//fmt.Printf("RST %v \n", tcp.RST)
					return &ScanResult{Port: tcpcs.Port, Success: !tcp.RST, Err: nil}
				}
			}
		}
	}

	/*
		var rtcp layers.TCP
		var payload gopacket.Payload
		data := make([]byte, 4096)
		parser := gopacket.NewDecodingLayerParser(layers.LayerTypeTCP, &rtcp, &payload)
		_, _, err1 := tcpcs.Conn.ReadFromIP(data)
		if err1 != nil {
			fmt.Printf("%v", err1)
		}
		decoded := []gopacket.LayerType{}
		if err := parser.DecodeLayers(data, &decoded); err != nil || len(decoded) == 0 {
			fmt.Printf("Failed to code packet %v\n", err)
		}
		for _, layerType := range decoded {
			fmt.Printf("%v", layerType)
			//if layerType == layers.LayerTypeTCP {
			//fmt.Printf("Src %v ", rtcp.SrcPort)
			//fmt.Printf("Dst %v ", rtcp.DstPort)
			fmt.Printf("Reset %v ", rtcp.RST)
			fmt.Printf("Syn %v ", rtcp.SYN)
			fmt.Printf("Ack %v\n", rtcp.ACK)
			return &ScanResult{Port: uint64(tcpcs.Port), Success: (rtcp.ACK || rtcp.RST), Err: nil}
			//}
		}
		//return &ScanResult{Port: 1234, Success: true, Err: nil}
		return nil
	*/
}
