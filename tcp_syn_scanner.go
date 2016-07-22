package poke

import (
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"net"
)

type TcpSynScanner struct {
	Host string
	Port uint64
}

func newTcpSynScanner(host string, port uint64) Scanner {
	return TcpSynScanner{host, port}
}

func (tcpcs TcpSynScanner) Scan() *ScanResult {
	ipa, err := net.ResolveIPAddr("ip", tcpcs.Host)
	if err != nil {
		return nil
	}
	conn, err := net.ListenIP("ip:tcp", ipa)
	if err != nil {
		return nil
	}
	defer conn.Close()

	sp := random(1024, 2024)
	dip := net.ParseIP(tcpcs.Host)
	//sip, err := getIP(dip)
	//fmt.Printf("%v", sip)
	if err != nil {
		return nil
	}
	ip := &layers.IPv4{
		//SrcIP:    sip,
		DstIP:    dip,
		Protocol: layers.IPProtocolTCP,
	}
	tcp := &layers.TCP{
		SrcPort: layers.TCPPort(sp),
		DstPort: layers.TCPPort(tcpcs.Port),
		SYN:     true,
	}
	fmt.Printf("Src: %v Dst: %v\n", sp, tcpcs.Port)
	tcp.SetNetworkLayerForChecksum(ip)
	buf := gopacket.NewSerializeBuffer()
	opts := gopacket.SerializeOptions{
		ComputeChecksums: true,
		FixLengths:       true,
	}
	if err := gopacket.SerializeLayers(buf, opts, tcp); err != nil {
		//fmt.Printf("%v", err)
	}
	conn.Write(buf.Bytes())

	var rtcp layers.TCP
	var payload gopacket.Payload
	data := make([]byte, 4096)
	parser := gopacket.NewDecodingLayerParser(layers.LayerTypeTCP, &rtcp, &payload)
	_, _, err1 := conn.ReadFromIP(data)
	if err1 != nil {
		fmt.Printf("%v", err1)
	}
	decoded := []gopacket.LayerType{}
	if err := parser.DecodeLayers(data, &decoded); err != nil || len(decoded) == 0 {
		fmt.Printf("Failed to code packet %v\n", err)
	}
	for _, layerType := range decoded {
		if layerType == layers.LayerTypeTCP {
			//if rtcp.SrcPort == 22 {
			//fmt.Printf("%v", rtcp.SYN && rtcp.ACK || rtcp.RST)
			fmt.Printf("Src %v ", rtcp.SrcPort)
			fmt.Printf("Dst %v ", rtcp.DstPort)
			fmt.Printf("Reset %v ", rtcp.RST)
			fmt.Printf("Syn %v ", rtcp.SYN)
			fmt.Printf("Ack %v\n", rtcp.ACK)
			//}
			return &ScanResult{Port: uint64(rtcp.SrcPort), Success: (rtcp.SYN && rtcp.ACK || rtcp.RST), Err: nil}
		}
	}
	return nil
}
