package config

import (
	"github.com/BurntSushi/toml"
	"github.com/alecthomas/log4go"
)

var (
	Cfg *Config
)

type Config struct {
	MySQL  MySQL
	Global Global
}

func LoadConfig(dir string) {
	Cfg = new(Config)
	if _, err := toml.DecodeFile(dir+"/config.toml", Cfg); err != nil {
		log4go.Exit(err)
	}
}

type MySQL struct {
	User string `toml:"user"`
	Pasw string `toml:"pasw"`
	Prot string `toml:"prot"`
	Host string `toml:"host"`
	Port string `toml:"port"`
	Dbnm string `toml:"dbnm"`
}

type Global struct {
	Port         int32
	ServerOrigin string
}
