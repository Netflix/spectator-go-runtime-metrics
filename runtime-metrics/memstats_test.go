package runtime_metrics

import (
	"github.com/Netflix/spectator-go/spectator"
	"github.com/Netflix/spectator-go/spectator/writer"
	"runtime"
	"strconv"
	"testing"
	"time"
)

func TestUpdateMemStats(t *testing.T) {
	var clock ManualClock
	config := &spectator.Config{
		Location: "memory",
		CommonTags: map[string]string{
			"nf.app":     "test",
			"nf.cluster": "test-main",
			"nf.asg":     "test-main-v001",
			"nf.region":  "us-west-1",
		},
	}

	registry, err := spectator.NewRegistry(config)
	if err != nil {
		t.Error(err)
	}

	var mem memStatsCollector

	initializeMemStatsCollector(registry, &clock, &mem)
	clock.SetFromDuration(1 * time.Minute)

	var memStats runtime.MemStats
	memStats.Alloc = 100
	memStats.TotalAlloc = 200
	memStats.Sys = 300
	memStats.Mallocs = 10
	memStats.Frees = 5
	memStats.LastGC = uint64(30 * time.Second)
	memStats.NumGC = 2
	memStats.NumForcedGC = 1
	memStats.GCCPUFraction = .5
	memStats.PauseTotalNs = uint64(5 * time.Millisecond)
	updateMemStats(&mem, &memStats)

	expectedMeasurements := map[string]float64{
		"mem.numLiveObjects":     5,
		"mem.heapBytesAllocated": 100,
		"mem.maxHeapBytes":       300,
		"gc.timeSinceLastGC":     30,
		"gc.cpuPercentage":       50,
		"gc.pauseTime":           0.005,
	}
	memoryWriter := registry.GetWriter().(*writer.MemoryWriter)

	// Validate measurements
	measurements := memoryWriter.Lines
	validateMeasurements(t, measurements, expectedMeasurements)

	// reset memory writer
	memoryWriter.Lines = []string{}

	// Update metrics
	clock.SetFromDuration(2 * time.Minute)

	memStats.Alloc = 200
	memStats.TotalAlloc = 400
	memStats.Sys = 600
	memStats.Mallocs = 20
	memStats.Frees = 10
	memStats.LastGC = uint64(30 * time.Second)
	memStats.NumGC = 5
	memStats.NumForcedGC = 2
	memStats.GCCPUFraction = .4
	memStats.PauseTotalNs = uint64(15 * time.Millisecond)

	updateMemStats(&mem, &memStats)

	expectedMeasurements = map[string]float64{
		"mem.numLiveObjects":     10,
		"mem.heapBytesAllocated": 200,
		"mem.maxHeapBytes":       600,
		"mem.objectsAllocated":   20,
		"mem.objectsFreed":       10,
		"mem.allocationRate":     400,
		"gc.timeSinceLastGC":     90,
		"gc.cpuPercentage":       40,
		"gc.count":               5,
		"gc.forcedCount":         2,
		"gc.pauseTime":           0.01,
	}

	// Validate measurements
	measurements = memoryWriter.Lines
	validateMeasurements(t, measurements, expectedMeasurements)
}

func validateMeasurements(t *testing.T, lines []string, expectedMeasurements map[string]float64) {
	for _, line := range lines {
		// split line by ":" and get the first and third items
		_, metricId, value, err := spectator.ParseProtocolLine(line)
		if err != nil {
			t.Error(err)
		}

		actualValue, err := strconv.ParseFloat(value, 64)
		if err != nil {
			t.Error(err)
		}

		for metricName, expectedValue := range expectedMeasurements {
			if metricId.Name() == metricName {
				if expectedValue != actualValue {
					t.Errorf("Expected %f for %s but got %f", expectedValue, metricName, actualValue)
				}
			}
		}
	}
}
