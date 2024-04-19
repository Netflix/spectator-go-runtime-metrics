package runtime_metrics

import (
	"github.com/Netflix/spectator-go"
	"strings"
	"time"
)

func makeConfig(uri string) *spectator.Config {
	return &spectator.Config{
		Frequency: 10 * time.Millisecond,
		Timeout:   1 * time.Second,
		Uri:       uri,
		BatchSize: 10000,
		CommonTags: map[string]string{
			"nf.app":     "test",
			"nf.cluster": "test-main",
			"nf.asg":     "test-main-v001",
			"nf.region":  "us-west-1",
		},
		PublishWorkers: 1,
	}
}

func myMeters(registry *spectator.Registry) []spectator.Meter {
	var myMeters []spectator.Meter
	for _, meter := range registry.Meters() {
		if !strings.HasPrefix(meter.MeterId().Name(), "spectator.") {
			myMeters = append(myMeters, meter)
		}
	}
	return myMeters
}
