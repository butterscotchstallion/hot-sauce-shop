package lib

import (
	"os"
	"time"

	"github.com/BurntSushi/toml"
)

const ConfigFilename = "../config.toml"

type ConfigDatabase struct {
	Dsn string `toml:"dsn"`
}

type ConfigTestUsers struct {
	BoardAdminUsername   string `toml:"boardAdminUsername"`
	BoardAdminPassword   string `toml:"boardAdminPassword"`
	UnprivilegedUsername string `toml:"unprivilegedUsername"`
	UnprivilegedPassword string `toml:"unprivilegedPassword"`
	AdminUsername        string `toml:"adminUsername"`
	AdminPassword        string `toml:"adminPassword"`
}

type ConfigServer struct {
	Address             string `toml:"address"`
	AddressWithProtocol string `toml:"addressWithProtocol"`
	TimeZone            string `toml:"timeZone"`
}

type ConfigCache struct {
	DefaultCacheTime time.Duration `toml:"default"`
	PostList         time.Duration `toml:"postList"`
}

type HotSauceShopConfig struct {
	Server    ConfigServer    `toml:"server"`
	Database  ConfigDatabase  `toml:"database"`
	TestUsers ConfigTestUsers `toml:"testUsers"`
	Cache     ConfigCache     `toml:"cache"`
}

func ReadConfig(filename string) (HotSauceShopConfig, error) {
	fileBytes, readFileErr := os.ReadFile(filename)
	if readFileErr != nil {
		return HotSauceShopConfig{}, readFileErr
	}
	fileContents := string(fileBytes)
	var hotSauceShopConfig HotSauceShopConfig
	if _, err := toml.Decode(fileContents, &hotSauceShopConfig); err != nil {
		return HotSauceShopConfig{}, err
	}
	return hotSauceShopConfig, nil
}
