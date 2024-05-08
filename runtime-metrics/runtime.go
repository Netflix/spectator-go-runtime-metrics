package runtime_metrics

import (
	"github.com/Netflix/spectator-go/spectator"
	"github.com/Netflix/spectator-go/spectator/meter"
	"runtime"
	"time"
)

type sysStatsCollector struct {
	registry      *spectator.Registry
	curOpen       *meter.Gauge
	maxOpen       *meter.Gauge
	numGoroutines *meter.Gauge
}

func goRuntimeStats(s *sysStatsCollector) {
	s.numGoroutines.Set(float64(runtime.NumGoroutine()))
}

// CollectSysStats collects system stats: current/max file handles, number of
// goroutines
func CollectSysStats(registry *spectator.Registry) {
	var s sysStatsCollector
	s.registry = registry
	s.maxOpen = registry.Gauge("fh.max", nil)
	s.curOpen = registry.Gauge("fh.allocated", nil)
	s.numGoroutines = registry.Gauge("go.numGoroutines", nil)

	ticker := time.NewTicker(30 * time.Second)
	go func() {
		log := registry.GetLogger()
		for range ticker.C {
			log.Debugf("Collecting system stats")
			fdStats(&s)
			goRuntimeStats(&s)
		}
	}()
}

// CollectRuntimeMetrics starts the collection of memory and file handle metrics
//
//goland:noinspection GoUnusedExportedFunction
func CollectRuntimeMetrics(registry *spectator.Registry) {
	CollectMemStats(registry)
	CollectSysStats(registry)
}
