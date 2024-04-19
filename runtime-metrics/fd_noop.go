//go:build !linux
// +build !linux

package runtime_metrics

func fdStats(s *sysStatsCollector) {
	// do nothing
}
