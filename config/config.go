package config

import (
	"os"
	"path/filepath"
)
import "github.com/spf13/viper"

type NacosStrcut struct {
	Ip          string
	Port        uint64
	Scheme      string
	ContextPath string
}

type ConfigStruct struct {
	Nacos NacosStrcut
}

var Config *ConfigStruct

func GetConf() *ConfigStruct {
	return Config
}

func ParseConf() {
	v := viper.New()
	pwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	environment := os.Getenv("environment")
	absolutePath := filepath.Join(pwd, "config", environment+".yaml")
	v.SetConfigFile(absolutePath)
	v.SetConfigType("yaml")
	if err := v.ReadInConfig(); err != nil {
		panic(err)
	}
	var config *ConfigStruct
	err = v.Unmarshal(&config)
	if err != nil {
		panic(err)
	}
	Config = config
}
