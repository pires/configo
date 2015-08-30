package configo

import (
	"testing"
	"time"
)

type configuration struct {
	Pubsub    []string
	Keepalive bool
	Heartbeat *Duration
}

func TestLoadProductionEnvironmentFromFile(t *testing.T) {
	var config configuration
	err := LoadEnvironmentFromFile("example.toml", "production", &config)
	if err != nil {
		t.Fatal(err)
	}

	if !config.Keepalive {
		t.Fail()
	}

	if 200 * time.Millisecond != config.Heartbeat.Duration {
		t.Fail()
	}
}

func TestLoadDevelopmentEnvironment(t *testing.T) {
	var config configuration
	err := LoadEnvironment("development", &config)
	if err != nil {
		t.Fatal(err)
	}

	if config.Keepalive {
		t.Fail()
	}

	if config.Pubsub[0] != "localhost:7788" {
		t.Fail()
	}

	if 1 * time.Second != config.Heartbeat.Duration {
		t.Fail()
	}
}

func TestLoadUnexistingFile(t *testing.T) {
	var config configuration
	err := LoadEnvironmentFromFile("xxx.toml", "xxx", &config)
	if err == nil {
		t.Fatal()
	}
}

func TestLoadUnexistingEnvironment(t *testing.T) {
	var config configuration
	err := LoadEnvironment("xxx", &config)
	if err == nil {
		t.Fatal()
	}
}

func TestLoadNoEnvironment(t *testing.T) {
	var config configuration
	err := LoadEnvironmentFromFile("invalid.toml", "xxx", &config)
	if err == nil {
		t.Fatal()
	}
}
