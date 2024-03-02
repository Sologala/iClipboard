package main

import (
	"fmt"
	"io/ioutil"

	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v2"
)

const ConfigFile = "config.yaml"
const LogFolder = "log"

var GLogPath string

// Config represents configuration for applicaton
var Gconfig Config = DefaultConfig

// type Config struct {
// 	port                    string `yaml:"port"`
// }

type Config struct {
	Port                    string `yaml:"port"`
	Authkey                 string `yaml:"authkey"`
	Authkey_expired_timeout int64  `yaml:"authkey_expired_timeout"`
}

var DefaultConfig = Config{
	Port:                    "8086",
	Authkey_expired_timeout: 30,
	Authkey:                 "",
}

func (this Config) LoadConfig() error {
	if !IsExistFile(ConfigFile) {
		DefaultConfig.WriteConfig(ConfigFile)
	}
	data, err := ioutil.ReadFile(ConfigFile)
	if err != nil {
		log.Fatal().Msgf("Cannot read config file %s : %v", ConfigFile, err)
		return err
	}
	err = yaml.Unmarshal(data, &this)
	if err != nil {
		log.Fatal().Msgf("Parse YAML file %s faild: %v", ConfigFile, err)
	}
	return nil
}

func (this Config) WriteConfig(path string) error {
	CreateFileIfNotExist(path)
	data, err := yaml.Marshal(&DefaultConfig)
	fmt.Printf("%s", data)
	if err != nil {
		log.Fatal().Msg("config to yaml faild ")
	}
	err = ioutil.WriteFile("config.yaml", data, 0755)
	if err != nil {
		log.Fatal().Msgf("write to %s faild", "config.yaml")
	}
	return err
}
