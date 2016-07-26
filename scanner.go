package poke

// An interface for all the scanners
type Scanner interface {
	Scan() *ScanResult
}
