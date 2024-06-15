package runtime_metrics

import (
	"runtime"
	"time"

	"github.com/Netflix/spectator-go/v2/spectator"
	"github.com/Netflix/spectator-go/v2/spectator/meter"
)

type memStatsCollector struct {
	clock            Clock
	registry         spectator.Registry
	bytesAlloc       *meter.Gauge
	allocationRate   *meter.MonotonicCounterUint
	totalBytesSystem *meter.Gauge
	numLiveObjects   *meter.Gauge
	objectsAllocated *meter.MonotonicCounterUint
	objectsFreed     *meter.MonotonicCounterUint

	gcLastPauseTimeValue uint64
	gcPauseTime          *meter.Timer
	gcAge                *meter.Gauge
	gcCount              *meter.MonotonicCounter
	forcedGcCount        *meter.MonotonicCounter
	gcPercCpu            *meter.Gauge
}

func memStats(m *memStatsCollector) {
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)
	updateMemStats(m, &mem)
}

func updateMemStats(m *memStatsCollector, mem *runtime.MemStats) {
	m.bytesAlloc.Set(float64(mem.Alloc))
	m.allocationRate.Set(mem.TotalAlloc)
	m.totalBytesSystem.Set(float64(mem.Sys))
	m.numLiveObjects.Set(float64(mem.Mallocs - mem.Frees))
	m.objectsAllocated.Set(mem.Mallocs)
	m.objectsFreed.Set(mem.Frees)

	nanosPause := mem.PauseTotalNs - m.gcLastPauseTimeValue
	m.gcPauseTime.Record(time.Duration(nanosPause))
	m.gcLastPauseTimeValue = mem.PauseTotalNs

	nanos := m.clock.Nanos()
	timeSinceLastGC := nanos - int64(mem.LastGC)
	secondsSinceLastGC := float64(timeSinceLastGC) / 1e9
	m.gcAge.Set(secondsSinceLastGC)

	m.gcCount.Set(float64(mem.NumGC))
	m.forcedGcCount.Set(float64(mem.NumForcedGC))
	m.gcPercCpu.Set(mem.GCCPUFraction * 100)
}

func initializeMemStatsCollector(registry spectator.Registry, clock Clock, mem *memStatsCollector) {
	mem.clock = clock
	mem.registry = registry
	mem.bytesAlloc = registry.Gauge("mem.heapBytesAllocated", nil)
	mem.allocationRate = registry.MonotonicCounterUint("mem.allocationRate", nil)
	mem.totalBytesSystem = registry.Gauge("mem.maxHeapBytes", nil)
	mem.numLiveObjects = registry.Gauge("mem.numLiveObjects", nil)
	mem.objectsAllocated = registry.MonotonicCounterUint("mem.objectsAllocated", nil)
	mem.objectsFreed = registry.MonotonicCounterUint("mem.objectsFreed", nil)
	mem.gcPauseTime = registry.Timer("gc.pauseTime", nil)
	mem.gcAge = registry.Gauge("gc.timeSinceLastGC", nil)
	mem.gcCount = registry.MonotonicCounter("gc.count", nil)
	mem.forcedGcCount = registry.MonotonicCounter("gc.forcedCount", nil)
	mem.gcPercCpu = registry.Gauge("gc.cpuPercentage", nil)
}

func CollectMemStatsWithClock(registry spectator.Registry, clock Clock) {
	var mem memStatsCollector
	initializeMemStatsCollector(registry, clock, &mem)

	ticker := time.NewTicker(30 * time.Second)
	go func() {
		log := registry.GetLogger()
		for range ticker.C {
			log.Debugf("Collecting memory stats")
			memStats(&mem)
		}
	}()
}

// CollectMemStats collects memory stats
//
// See: https://golang.org/pkg/runtime/#MemStats
func CollectMemStats(registry spectator.Registry) {
	CollectMemStatsWithClock(registry, &SystemClock{})
}
