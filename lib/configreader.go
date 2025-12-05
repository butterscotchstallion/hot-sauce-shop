package lib

import (
	"os"

	"github.com/BurntSushi/toml"
)

type ConfigDatabase struct {
	Dsn string `toml:"dsn"`
}

type ConfigTestUsers struct {
	BoardAdminUsername   string `toml:"boardAdminUsername"`
	BoardAdminPassword   string `toml:"boardAdminPassword"`
	UnprivilegedUsername string `toml:"unprivilegedUsername"`
	UnprivilegedPassword string `toml:"unprivilegedPassword"`
}

type ConfigServer struct {
	Address             string
	AddressWithProtocol string
}

type HotSauceShopConfig struct {
	Server    ConfigServer    `toml:"server"`
	Database  ConfigDatabase  `toml:"database"`
	TestUsers ConfigTestUsers `toml:"testUsers"`
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
