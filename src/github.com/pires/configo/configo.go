package configo

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
)

func setField(obj interface{}, name string, value interface{}) error {
	structValue := reflect.ValueOf(obj).Elem()
	structFieldValue := structValue.FieldByNameFunc(func(toFold string) bool {
		return strings.EqualFold(name, toFold)
	})

	if !structFieldValue.IsValid() {
		return fmt.Errorf("\t\tNo such field: %s in obj", name)
	}

	if !structFieldValue.CanSet() {
		return fmt.Errorf("\t\tCannot set %s field value", name)
	}

	structFieldType := structFieldValue.Type()
	val := reflect.ValueOf(value)
	if structFieldType != val.Type() {
		if structFieldType.String() == "time.Duration" {
			d, _ := time.ParseDuration(val.String())
			structFieldValue.Set(reflect.ValueOf(d))
		} else if val.Kind() == reflect.Slice {
			for i := 0; i < val.Len(); i++ {
				structFieldValue.Set(reflect.Append(structFieldValue, val.Index(i).Elem()))
			}
		} else {
			return fmt.Errorf("\t\tProvided value type didn't match obj field type. %v != %v | %v", structFieldType, val.Type(), val.Kind())
		}
	} else {
		structFieldValue.Set(val)
	}
	return nil
}

func processResult(obj interface{}, value interface{}) error {
	for k, v := range value.(map[string]interface{}) {
		setField(obj, k, v)
	}
	return nil
}

func Decode(fpath string, env string, v interface{}) error {
	var environments map[string]interface{}
	if _, err := toml.DecodeFile(fpath, &environments); err != nil {
		return err
	}
	processResult(v, environments[env])
	return nil
}
