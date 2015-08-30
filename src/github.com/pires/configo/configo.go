/*
	Package configo provides functionality to load an application configuration
	from TOML files.

	It is a shameless wrapper around github.com/BurntSushi/toml that tries to
	to enforce an environment-based configuration style.
 */
package configo

import (
	"bytes"
	"fmt"
	"time"

	"github.com/BurntSushi/toml"
)

// See https://github.com/BurntSushi/toml#using-the-encodingtextunmarshaler-interface
type Duration struct {
	time.Duration
}

// See https://github.com/BurntSushi/toml#using-the-encodingtextunmarshaler-interface
func (d *Duration) UnmarshalText(text []byte) error {
	var err error
	d.Duration, err = time.ParseDuration(string(text))
	return err
}

func load(obj interface{}, conf map[string]interface{}) error {
	buf := new(bytes.Buffer)
	if err := toml.NewEncoder(buf).Encode(conf); err != nil {
		return err
	}
	_, err := toml.Decode(buf.String(), obj)
	return err
}

// Loads a specific environment from a specific configuratino file.
func LoadEnvironmentFromFile(fpath string, env string, conf interface{}) error {
	var environments map[string]map[string]interface{}
	if _, err := toml.DecodeFile(fpath, &environments); err != nil {
		return err
	}

	// if env is empty, try first available environment
	value, ok := environments[env]
	if !ok {
		return fmt.Errorf("Environment %v doesn't exist in the configuration file.", env)
	}
	return load(conf, value)
}


// Loads a specific environment from config.toml file located in the working directory.
func LoadEnvironment(env string, v interface{}) error {
	return LoadEnvironmentFromFile("application.toml", env, v)
}
