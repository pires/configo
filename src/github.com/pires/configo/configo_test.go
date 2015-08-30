package configo

import (
	"testing"
	"time"
)

type configuration struct {
	A         string
	Pubsub    []string
	Keepalive bool
	Heartbeat *Duration
}

func TestDecodeEnvironment(t *testing.T) {
	var config configuration
	err := LoadFromFile("example.toml", "production", &config)
	if err != nil {
		t.Fatal(err)
	}

	if config.A != "a" {
		t.Fail()
	}

	if !config.Keepalive {
		t.Fail()
	}

	if 200 * time.Millisecond != config.Heartbeat.Duration {
		t.Fail()
	}
}
