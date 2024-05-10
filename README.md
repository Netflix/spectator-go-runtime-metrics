[![Go Reference](https://pkg.go.dev/badge/github.com/Netflix/spectator-go.svg)](https://pkg.go.dev/github.com/Netflix/spectator-go)
[![Snapshot](https://github.com/Netflix/spectator-go/actions/workflows/snapshot.yml/badge.svg)](https://github.com/Netflix/spectator-go/actions/workflows/snapshot.yml)
[![Release](https://github.com/Netflix/spectator-go/actions/workflows/release.yml/badge.svg)](https://github.com/Netflix/spectator-go/actions/workflows/release.yml)

# Spectator-go Runtime Metrics

Library to collect runtime metrics for Golang applications using [spectator-go](https://github.com/Netflix/spectator-go).

## Instrumenting Code

```go
package main

import (
	"github.com/Netflix/spectator-go-runtime-metrics/runtime-metrics"
	"github.com/Netflix/spectator-go/spectator"
)

func main() {
	config := &spectator.Config{}
	registry, _ := spectator.NewRegistry(config)
	defer registry.Close()

	runtime_metrics.CollectRuntimeMetrics(registry)
}
```