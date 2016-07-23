package poke

import (
	"fmt"
	"github.com/achanda/poke/utils"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"net"
)

type TcpSynScanner struct {
	Host string
	Port uint64
	Conn net.IPConn
}

func newTcpSynScanner(host string, port uint64, conn net.IPConn) Scanner {
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
	saddr, err := utils.GetLocalIP(tcpcs.Host)
	if err != nil {
		panic(err)
	}
	ip := &layers.IPv4{
		SrcIP:    saddr.IP,
		DstIP:    dip,
		Protocol: layers.IPProtocolTCP,
	}
	tcp := &layers.TCP{
		//SrcPort: layers.TCPPort(sp),
		DstPort: layers.TCPPort(tcpcs.Port),
		SYN:     true,
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
	tcpcs.Conn.Write(buf.Bytes())

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
}
