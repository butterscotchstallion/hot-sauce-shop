package lib

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

/**
 * The primary motivation for this set of functions is to dynamically
 * set caching values during testing so that they don't interfere. This
 * also prevents reading the file multiple times, as we are not passing
 * the config everywhere by reference.
 *
 * The flow works like this:
 * - Read the config from the file system when the server runs
 * - Set the runtime config using that value
 * - Update it in memory as necessary
 */
var runtimeConfig HotSauceShopConfig

func SetRuntimeConfig(config HotSauceShopConfig) {
	runtimeConfig = config
}

func GetRuntimeConfig() HotSauceShopConfig {
	return runtimeConfig
}

// SetDynamicConfigProperty - can set nested properties as well
func SetDynamicConfigProperty(path string, value interface{}) error {
	cfg := GetRuntimeConfig()

	val := reflect.ValueOf(&cfg).Elem()
	parts := strings.Split(path, ".")
	for i, part := range parts {
		// Handle potential pointer/interface indirection
		for val.Kind() == reflect.Ptr || val.Kind() == reflect.Interface {
			val = val.Elem()
		}

		if val.Kind() != reflect.Struct {
			return fmt.Errorf("cannot traverse %s: not a struct", part)
		}

		field := val.FieldByName(part)
		if !field.IsValid() {
			field = val.FieldByName(cases.Title(language.English).String(part))
			if !field.IsValid() {
				return fmt.Errorf("no such field: %s", part)
			}
		}

		if i == len(parts)-1 {
			if !field.CanSet() {
				return fmt.Errorf("cannot set field: %s", part)
			}

			valToSet := reflect.ValueOf(value)
			if field.Type() != valToSet.Type() {
				return fmt.Errorf(
					"type mismatch for %s: expected %v, got %v", path, field.Type(), valToSet.Type(),
				)
			}

			field.Set(valToSet)
		} else {
			val = field
		}
	}

	SetRuntimeConfig(cfg)
	return nil
}

func DisableCaching() error {
	var err error

	err = SetDynamicConfigProperty("cache.DefaultCacheTime", time.Duration(0))
	if err != nil {
		return err
	}

	err = SetDynamicConfigProperty("cache.PostList", time.Duration(0))
	if err != nil {
		return err
	}

	return nil
}
