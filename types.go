package poke

import (
	"fmt"
)

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
	Success bool
	Err     error
}
