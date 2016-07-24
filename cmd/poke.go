package main

import (
	"flag"
	"fmt"
	"github.com/achanda/go-services"
	"github.com/achanda/poke"
	"github.com/achanda/poke/utils"
	"log"
	"os"
	"strconv"
	"strings"
)

// Need worker pool because running 1 goroutine per port exhausts file descriptors
const MAX_WORKERS = 100

// Run the port scanner
func main() {
	var host, port_range_arg, scanner_type string
	var ipver bool
	flag.StringVar(&host, "host", "", "host to scan")
	flag.StringVar(&port_range_arg, "ports", "", "ports to scan")
	flag.StringVar(&scanner_type, "scanner", "", "scanner type to use")
	flag.BoolVar(&ipver, "v4", true, "will we use IPv4")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
	}

	flag.Parse()

	scanner_type = strings.ToLower(scanner_type)
	if scanner_type == "" {
		fmt.Printf("unknown scanner type, defaulting to connect scan\n")
		scanner_type = "c"
	}

	if host == "" || port_range_arg == "" {
		flag.Usage()
		os.Exit(1)
	}

	hosts, err := utils.ParseHost(host)
	if err != nil {
		fmt.Printf("Could not parse the host param %v", host)
		os.Exit(2)
	}

	if len(hosts) == 1 {
		ipver = utils.IsIPv4(host)
	}

	prs, err := parsePorts(port_range_arg)
	if err != nil {
		log.Fatal(err)
	}

	portmap, err := services.GetServices()
	if err != nil {
		fmt.Printf("did not find a services file")
	}
	// Format results
	results := ScanPorts(host, prs, scanner_type, ipVer)
	proto := ""
	if scanner_type == "c" || scanner_type == "s" {
		proto = "tcp"
	} else {
		proto = "udp"
	}
	for _, host := range hosts {
		results := ScanPorts(host, prs, scanner_type, ipver)
		for port, success := range results {
			if success {
				if portmap != nil {
					servname := portmap[uint16(port)].Name
					if servname != "" {
						fmt.Printf("%v/%v open %v\n", port, proto, portmap[uint16(port)].Name)
					} else {
						fmt.Printf("%v/%v: open\n", port, proto)
					}
				}
			}
		}
	}
}

func parsePorts(ranges_str string) (*poke.PortRange, error) {
	parts := strings.Split(ranges_str, ":")
	if len(parts) != 2 {
		fmt.Printf("Please specify port range in the form start:end\n")
	}
	start, err := strconv.ParseUint(parts[0], 10, 0)
	if err != nil {
		fmt.Printf("Failed to convert %v to an int", parts[0])
		return nil, err
	}
	end, err := strconv.ParseUint(parts[1], 10, 0)
	if err != nil {
		fmt.Printf("Failed to convert %v to an int", parts[1])
		return nil, err
	}
	return &poke.PortRange{Start: start, End: end}, nil
}

// Run the scan with a worker pool; memory usage grows in proportion
// with number of ports scanned to prevent deadlock from blocking channels
func ScanPorts(host string, pr *poke.PortRange, scanner_type string, ipVer bool) map[uint64]bool {
	num_ports := pr.End - pr.Start + 1
	results := make(map[uint64]bool)
	jobpipe := make(chan uint64, num_ports)
	respipe := make(chan *poke.ScanResult, num_ports)

	fmt.Printf("Scanning %v\n", host)
	// Start workers
	for worker := 0; worker < MAX_WORKERS; worker++ {
		go scanWorker(host, jobpipe, respipe, scanner_type, ipVer)
	}

	// Seed w/ jobs
	for port := pr.Start; port < pr.End+1; port++ {
		jobpipe <- port
	}

	// Receive results
	received := uint64(0)
	for received < pr.End-pr.Start {
		res := <-respipe
		results[res.Port] = res.Success
		received += 1
	}
	return results
}

// Worker function; pull from job queue forever and return results on result
// queue
func scanWorker(host string, jobpipe chan uint64, respipe chan *poke.ScanResult, scanner_type string, ipVer bool) {
	for job := <-jobpipe; ; job = <-jobpipe {
		var sr poke.Scanner
		switch scanner_type {
		case "s":
			sr = poke.NewTcpSynScanner(host, job, ipVer)
		case "c":
			sr = poke.NewTcpConnectScanner(host, job, ipVer)
		}
		respipe <- sr.Scan()
	}
}
