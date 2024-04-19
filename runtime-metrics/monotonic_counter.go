package runtime_metrics

import (
	"github.com/Netflix/spectator-go"
	"sync"
	"sync/atomic"
)

// monotonicCounter is used track a monotonically increasing counter.
//
// You can find more about this type by viewing the relevant Java Spectator documentation here:
//
// https://netflix.github.io/spectator/en/latest/intro/gauge/#monotonic-counters
type monotonicCounter struct {
	value int64
	// Pointers need to be after counters to ensure 64-bit alignment. See
	// note in atomicnum.go
	registry    *spectator.Registry
	id          *spectator.Id
	counter     *spectator.Counter
	counterOnce sync.Once
}

// newMonotonicCounter generates a new monotonic counter, taking the registry so
// that it can lazy-load the underlying counter once `Set` is called the first
// time. It generates a new meter identifier from the name and tags.
func newMonotonicCounter(registry *spectator.Registry, name string, tags map[string]string) *monotonicCounter {
	return newMonotonicCounterWithId(registry, spectator.NewId(name, tags))
}

// newMonotonicCounterWithId generates a new monotonic counter, using the
// provided meter identifier.
func newMonotonicCounterWithId(registry *spectator.Registry, id *spectator.Id) *monotonicCounter {
	return &monotonicCounter{
		registry: registry,
		id:       id,
	}
}

// Set adds amount to the current counter.
func (c *monotonicCounter) Set(amount int64) {
	var uninitialized bool
	c.counterOnce.Do(func() {
		c.counter = c.registry.CounterWithId(c.id)
		uninitialized = true
	})

	if !uninitialized {
		prev := atomic.LoadInt64(&c.value)
		delta := amount - prev
		if delta >= 0 {
			c.counter.Add(delta)
		}
	}

	atomic.StoreInt64(&c.value, amount)
}

// Count returns the current counter value.
func (c *monotonicCounter) Count() int64 {
	return atomic.LoadInt64(&c.value)
}
