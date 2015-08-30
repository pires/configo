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

func TestDecodeProductionEnvironmentFromFile(t *testing.T) {
	var config configuration
	err := LoadFromFile("example.toml", "production", &config)
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

func TestDecodeDevelopmentEnvironment(t *testing.T) {
	var config configuration
	err := Load("development", &config)
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

func TestDecodeUnexistingFile(t *testing.T) {
	var config configuration
	err := LoadFromFile("xxx.toml", "xxx", &config)
	if err == nil {
		t.Fatal()
	}
}

func TestDecodeUnexistingEnvironment(t *testing.T) {
	var config configuration
	err := Load("xxx", &config)
	if err == nil {
		t.Fatal()
	}
}
