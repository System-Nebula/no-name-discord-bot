package config

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/BurntSushi/toml"
)

type conf struct {
}

// GetConfig -- get configuration values from disk
func GetConfig() toml.MetaData {
	var c conf
	d := getConfigFromFile("config.toml")
	t, err := toml.Decode(d, &c)
	if err != nil {
		fmt.Println("Problems reading toml config file")
	}
	fmt.Println(t)
	return t
}

func getConfigFromFile(file string) string {
	c, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatal("unable to load file, err=", err)
	}

	t := string(c)
	return t
}
