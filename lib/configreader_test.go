package lib

import "testing"

func TestReadExampleConfig(t *testing.T) {
	configData := ReadConfig("config.example.toml")
}
