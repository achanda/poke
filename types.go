package poke

// Some types to be used everywhere

import (
	"fmt"
)

// Represents a range of ports
type PortRange struct {
	Start uint64
	End   uint64
}

func (pr *PortRange) String() string {
	return fmt.Sprintf("[%v,%v)", pr.Start, pr.End)
}

// Container for scan results from workers
type ScanResult struct {
	Port    uint64
	Success bool // True if the port is open
	Err     error
}
