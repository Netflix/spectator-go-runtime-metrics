[![Go Reference](https://pkg.go.dev/badge/github.com/Netflix/spectator-go.svg)](https://pkg.go.dev/github.com/Netflix/spectator-go)
[![Snapshot](https://github.com/Netflix/spectator-go/actions/workflows/snapshot.yml/badge.svg)](https://github.com/Netflix/spectator-go/actions/workflows/snapshot.yml)
[![Release](https://github.com/Netflix/spectator-go/actions/workflows/release.yml/badge.svg)](https://github.com/Netflix/spectator-go/actions/workflows/release.yml)

# Spectator-go Runtime Metrics

> :warning: Experimental

Library to collect runtime metrics for Golang applications using spectator-go.

## Instrumenting Code

```go
package main

import (
	"github.com/Netflix/spectator-go"
	"time"
)

func main() {
	config := &spectator.Config{
		Frequency: 5 * time.Second,
		Timeout:   1 * time.Second,
		Uri:       "http://example.org/api/v1/publish",
	}
	registry := spectator.NewRegistry(config)

	// optionally set custom logger (it must implement Debugf, Infof, Errorf)
	// registry.SetLogger(logger)
	registry.Start()
	defer registry.Stop()

	// collect memory and file descriptor metrics
	spectator.CollectRuntimeMetrics(registry)
}

```

## Logging

Logging is implemented with the standard Golang [log package](https://pkg.go.dev/log). The logger
defines interfaces for [Debugf, Infof, and Errorf](./logger.go#L10-L14) which means that under
normal operation, you will see log messages for all of these levels. There are
[useful messages](https://github.com/Netflix/spectator-go/blob/master/registry.go#L268-L273)
implemented at the Debug level which can help diagnose the metric publishing workflow. If you do
not see any of these messages, then it is an indication that the Registry may not be started.

If you do not wish to see debug log messages from spectator-go, then you should configure a custom
logger which implements the Logger interface. A library such as [Zap](https://github.com/uber-go/zap)
can provide this functionality, which will then allow for log level control at the command line
with the `--log-level=debug` flag.

## Debugging Metric Payloads

Set the following environment variable to enumerate the metrics payloads which
are sent to the backend. This is useful for debugging metric publishing issues.

```
export SPECTATOR_DEBUG_PAYLOAD=1
```