package configo

import (
	"bytes"
	"time"

	"github.com/BurntSushi/toml"
	"fmt"
)

type Duration struct {
	time.Duration
}

func (d *Duration) UnmarshalText(text []byte) error {
	var err error
	d.Duration, err = time.ParseDuration(string(text))
	return err
}

func Unify(obj interface{}, conf map[string]interface{}) error {
	buf := new(bytes.Buffer)
	if err := toml.NewEncoder(buf).Encode(conf); err != nil {
		return err
	}
	_, err := toml.Decode(buf.String(), obj)
	return err
}

func LoadFromFile(fpath string, env string, v interface{}) error {
	var environments map[string]map[string]interface{}
	if _, err := toml.DecodeFile(fpath, &environments); err != nil {
		return err
	}

	value, ok := environments[env]
	if !ok {
		return fmt.Errorf("Environment %v doesn't exist in the configuration file.", env)
	}

	return Unify(v, value)
}


// assumes configuration filepath is `pwd`/config.toml
func Load(env string, v interface{}) error {
	return LoadFromFile("application.toml", env, v)
}
