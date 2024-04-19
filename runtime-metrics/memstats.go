package runtime_metrics

import (
	"github.com/Netflix/spectator-go"
	"runtime"
	"time"
)

type memStatsCollector struct {
	clock            Clock
	registry         *spectator.Registry
	bytesAlloc       *spectator.Gauge
	allocationRate   *monotonicCounter
	totalBytesSystem *spectator.Gauge
	numLiveObjects   *spectator.Gauge
	objectsAllocated *monotonicCounter
	objectsFreed     *monotonicCounter

	gcLastPauseTimeValue uint64
	gcPauseTime          *spectator.Timer
	gcAge                *spectator.Gauge
	gcCount              *monotonicCounter
	forcedGcCount        *monotonicCounter
	gcPercCpu            *spectator.Gauge
}

func memStats(m *memStatsCollector) {
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)
	updateMemStats(m, &mem)
}

func updateMemStats(m *memStatsCollector, mem *runtime.MemStats) {
	m.bytesAlloc.Set(float64(mem.Alloc))
	m.allocationRate.Set(int64(mem.TotalAlloc))
	m.totalBytesSystem.Set(float64(mem.Sys))
	m.numLiveObjects.Set(float64(mem.Mallocs - mem.Frees))
	m.objectsAllocated.Set(int64(mem.Mallocs))
	m.objectsFreed.Set(int64(mem.Frees))

	nanosPause := mem.PauseTotalNs - m.gcLastPauseTimeValue
	m.gcPauseTime.Record(time.Duration(nanosPause))
	m.gcLastPauseTimeValue = mem.PauseTotalNs

	nanos := m.clock.Nanos()
	timeSinceLastGC := nanos - int64(mem.LastGC)
	secondsSinceLastGC := float64(timeSinceLastGC) / 1e9
	m.gcAge.Set(secondsSinceLastGC)

	m.gcCount.Set(int64(mem.NumGC))
	m.forcedGcCount.Set(int64(mem.NumForcedGC))
	m.gcPercCpu.Set(mem.GCCPUFraction * 100)
}

func initializeMemStatsCollector(registry *spectator.Registry, clock Clock, mem *memStatsCollector) {
	mem.clock = clock
	mem.registry = registry
	mem.bytesAlloc = registry.Gauge("mem.heapBytesAllocated", nil)
	mem.allocationRate = newMonotonicCounter(registry, "mem.allocationRate", nil)
	mem.totalBytesSystem = registry.Gauge("mem.maxHeapBytes", nil)
	mem.numLiveObjects = registry.Gauge("mem.numLiveObjects", nil)
	mem.objectsAllocated = newMonotonicCounter(registry, "mem.objectsAllocated", nil)
	mem.objectsFreed = newMonotonicCounter(registry, "mem.objectsFreed", nil)
	mem.gcPauseTime = registry.Timer("gc.pauseTime", nil)
	mem.gcAge = registry.Gauge("gc.timeSinceLastGC", nil)
	mem.gcCount = newMonotonicCounter(registry, "gc.count", nil)
	mem.forcedGcCount = newMonotonicCounter(registry, "gc.forcedCount", nil)
	mem.gcPercCpu = registry.Gauge("gc.cpuPercentage", nil)
}

func CollectMemStatsWithClock(registry *spectator.Registry, clock Clock) {
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
func CollectMemStats(registry *spectator.Registry) {
	CollectMemStatsWithClock(registry, &SystemClock{})
}
