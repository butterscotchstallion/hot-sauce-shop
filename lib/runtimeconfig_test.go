package lib

import (
	"testing"
	"time"
)

func TestSetDynamicConfigProperty(t *testing.T) {
	cfg, cfgErr := ReadConfig(ConfigFilename)
	if cfgErr != nil {
		t.Fatal(cfgErr)
	}
	SetRuntimeConfig(cfg)
	newCacheValue := time.Duration(42)
	setErr := SetDynamicConfigProperty("Cache.DefaultCacheTime", newCacheValue)
	if setErr != nil {
		t.Fatal(setErr)
	}
	updatedCfg := GetRuntimeConfig()
	if updatedCfg.Cache.DefaultCacheTime != newCacheValue {
		t.Fatal("failed to update config property")
	}
}
