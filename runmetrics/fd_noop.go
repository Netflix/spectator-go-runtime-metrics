//go:build !linux
// +build !linux

package runmetrics

func fdStats(s *sysStatsCollector) {
	// do nothing
}
