package config

import (
	"fmt"

	"github.com/BurntSushi/toml"
)

// GetConfig -- get configuration values from disk
func GetConfig(data string, clazz interface{}) {
	_, err := toml.Decode(data, &clazz)
	if err != nil {
		fmt.Println("Problems reading toml config file")
	}
	// todo: implement logging - see https://github.com/mozilla/sops/blob/master/logging/logging.go
}
