package main

import (
	"fmt"
	"os"
	"reflect"
	"strings"
)

// options, loaded from environment
var opts struct {
	Port        string `required:"true"`
	DatabaseURL string `required:"true"`
}

// optsFromEnv loads the options struct values from the environment.
// It returns an error if any required option is not set
func optsFromEnv() error {

	v := reflect.ValueOf(opts)
	for i := 0; i < v.NumField(); i++ {
		field := v.Type().Field(i)
		name := strings.ToUpper(field.Name)
		envVal := os.Getenv(name)
		required := field.Tag.Get("required")
		if envVal == "" && required == "true" {
			return fmt.Errorf("%s not set in environment", name)
		}
		reflect.ValueOf(&opts).Elem().Field(i).SetString(envVal)
	}

	return nil
}
