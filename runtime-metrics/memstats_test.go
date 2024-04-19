package runtime_metrics

import (
	"github.com/Netflix/spectator-go"
	"reflect"
	"runtime"
	"testing"
	"time"
)

// TODO this test will need to be rewritten once we migrate to spectator-go thin client because we'll no longer have access to the Meters() method
func TestUpdateMemStats(t *testing.T) {
	var clock ManualClock
	config := makeConfig("")
	registry := spectator.NewRegistry(config)
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

	ms := myMeters(registry)
	if len(ms) != 11 {
		t.Error("Expected 11 meters registered, got", len(ms))
	}

	expectedValues := map[string]float64{
		"mem.numLiveObjects":     5,
		"mem.heapBytesAllocated": 100,
		"mem.maxHeapBytes":       300,
		"gc.timeSinceLastGC":     float64(30),
		"gc.cpuPercentage":       50,
	}
	for _, m := range ms {
		name := m.MeterId().Name()
		if name == "gc.pauseTime" {
			assertTimer(t, m.(*spectator.Timer), 1, 5*1e6, 25*1e12, 5*1e6)
		} else {
			expected := expectedValues[name]
			measures := m.Measure()
			if expected > 0 {
				if len(measures) != 1 {
					t.Fatalf("Expected one value from %v: got %d", m.MeterId(), len(measures))
				}
				if v := measures[0].Value(); v != expected {
					t.Errorf("%v: expected %f. got %f", m.MeterId(), expected, v)
				}
			} else {
				if len(measures) != 0 {
					t.Errorf("Unexpected measurements from %v: got %d measurements", m.MeterId(), len(measures))
				}
			}
		}
	}

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
	ms = registry.Meters()
	expectedValues = map[string]float64{
		"mem.numLiveObjects":     10,
		"mem.heapBytesAllocated": 200,
		"mem.maxHeapBytes":       600,
		"mem.objectsAllocated":   10,
		"mem.objectsFreed":       5,
		"mem.allocationRate":     200,
		"gc.timeSinceLastGC":     float64(90),
		"gc.cpuPercentage":       40,
		"gc.count":               3,
		"gc.forcedCount":         1,
	}
	for _, m := range ms {
		name := m.MeterId().Name()
		switch name {
		case "gc.pauseTime":
			assertTimer(t, m.(*spectator.Timer), 1, 10*1e6, 100*1e12, 10*1e6)
		case "spectator.registrySize":
		default:
			expected := expectedValues[name]
			measures := m.Measure()
			if expected > 0 {
				if len(measures) != 1 {
					t.Errorf("Expected one value from %v: got %d", m.MeterId(), len(measures))
				}
				if v := measures[0].Value(); v != expected {
					t.Errorf("%v: expected %f. got %f", m.MeterId(), expected, v)
				}
			} else if len(measures) != 0 {
				t.Errorf("Unexpected measurements from %v: got %d measurements", m.MeterId(), len(measures))
			}
		}
	}
}

func assertTimer(t *testing.T, timer *spectator.Timer, count int64, total int64, totalSq float64, max int64) {
	ms := timer.Measure()
	if len(ms) != 4 {
		t.Error("Expected 4 measurements from a Timer, got ", len(ms))
	}

	expected := make(map[string]float64)
	expected[timer.MeterId().WithStat("count").MapKey()] = float64(count)
	expected[timer.MeterId().WithStat("totalTime").MapKey()] = float64(total) / 1e9
	expected[timer.MeterId().WithStat("totalOfSquares").MapKey()] = totalSq / 1e18
	expected[timer.MeterId().WithStat("max").MapKey()] = float64(max) / 1e9

	got := make(map[string]float64)
	for _, v := range ms {
		got[v.Id().MapKey()] = v.Value()
	}
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Expected measurements (count=%d, total=%d, totalSq=%.0f, max=%d)", count, total, totalSq, max)
		for _, m := range ms {
			t.Errorf("Got %s %v = %f", m.Id().Name(), m.Id().Tags(), m.Value())
		}
	}

	// ensure timer is reset after being measured
	if timer.Count() != 0 || timer.TotalTime() != 0 {
		t.Error("Timer should be reset after being measured")
	}
}
