[![Go Reference](https://pkg.go.dev/badge/github.com/Netflix/spectator-go.svg)](https://pkg.go.dev/github.com/Netflix/spectator-go-runtime-metrics)
[![Snapshot](https://github.com/Netflix/spectator-go/actions/workflows/snapshot.yml/badge.svg)](https://github.com/Netflix/spectator-go-runtime-metrics/actions/workflows/snapshot.yml)
[![Release](https://github.com/Netflix/spectator-go/actions/workflows/release.yml/badge.svg)](https://github.com/Netflix/spectator-go-runtime-metrics/actions/workflows/release.yml)

# Spectator-go Runtime Metrics

Library to collect runtime metrics for Golang applications using [spectator-go](https://github.com/Netflix/spectator-go).

## Instrumenting Code

```go
package main

import (
	"github.com/Netflix/spectator-go-runtime-metrics/runmetrics"
	"github.com/Netflix/spectator-go/v2/spectator"
)

func main() {
	config, _ := spectator.NewConfig("", nil, nil)
	registry, _ := spectator.NewRegistry(config)
	defer registry.Close()

	runmetrics.CollectRuntimeMetrics(registry)
}
```
