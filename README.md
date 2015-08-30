# configo

TOML-based, strong-typed configuration for Go

[![Build Status](https://drone.io/github.com/pires/configo/status.png)](https://drone.io/github.com/pires/configo/latest)
[![GoDoc](https://godoc.org/github.com/pires/configo?status.svg)](https://godoc.org/github.com/pires/configo)

## Pre-requisites
```
go get github.com/BurntSushi/toml
github.com/pires/configo
```

## Usage

```toml
[stage]
pubsub = [
  "redisstg1",
  "redisstg2",
  "redisstg3"
]
keepalive = true
heartbeat = "500ms"
timeout = "5s"

[production]
pubsub = [
  "redisprod1",
  "redisprod2"
]
keepalive = true
heartbeat = "200ms"
timeout = "1m"
```

```go
package main

import (
        "time"
        "github.com/pires/configo"
)

type configuration struct {
	Pubsub    []string
	Keepalive bool
	Heartbeat *configo.Duration
	Timeout   *configo.Duration
}

var (
	done = make (chan struct{})
)

func main() {
	var config configuration
	if err := configo.LoadEnvironment("stage", &config); err != nil {
		panic(err)
	}

        for _, pubsub := range config.Pubsub {
		println("Connecting to", pubsub)
	}

	if config.Keepalive {
		println("Enabling keep-alive...")
		heartbeat := time.NewTicker(config.Heartbeat.Duration)
		defer heartbeat.Stop()

		// run timeout
		timeout := time.NewTimer(config.Timeout.Duration)
		go func(){
			<-timeout.C
			close(done)
		}()

		// run heartbeat
		for {
			select {
			case <-done:
				println("Done")
				return
			case <- heartbeat.C:
				println("Received heartbeat")
			}
		}
	} else {
		println("Keep-alive disabled.")
	}
}
```

The output should be:
```
Connecting to redisstg1
Connecting to redisstg2
Connecting to redisstg3
Enabling keep-alive...
Received heartbeat
Received heartbeat
Received heartbeat
Received heartbeat
Received heartbeat
Received heartbeat
Received heartbeat
Received heartbeat
Received heartbeat
Received heartbeat
Done
```

**Attention**: you may have noticed the type `configo.Duration` instead of `time.Duration`. This happens because [a special kind of unmarshalling is needed](https://github.com/BurntSushi/toml#using-the-encodingtextunmarshaler-interface).
Feel free to implement your own types.
