package config

import (
	"github.com/BurntSushi/toml"
	"github.com/alecthomas/log4go"
)

var (
	Cfg *Config
)

/*
	Configuration
*/
type Config struct {
	MySQL  MySQL
	Global Global
}

/*
	Load Configuration
*/
func LoadConfig(dir string) {
	Cfg = new(Config)
	if _, err := toml.DecodeFile(dir+"/config.toml", Cfg); err != nil {
		log4go.Exit(err)
	}
}

/*
	MySQL Configuration
*/
type MySQL struct {
	User string `toml:"user"`	//user name
	Pasw string `toml:"pasw"`	//password
	Prot string `toml:"prot"`	//protocol
	Host string `toml:"host"`	//host name or ip
	Port string `toml:"port"`	//port
	Dbnm string `toml:"dbnm"`	//database name
}

/*
	Global Configuration
*/
type Global struct {
	Port            int32  `toml:"port"`
	ServerOrigin    string `toml:"server_origin"`
	LevelDbFilePath string `toml:"level_db_file_path"`
	BolaxyNodeUrl   string `toml:"bolaxy_node_url"`
	DefHost         string `toml:"def_host"`
	DefPort         uint32 `toml:"def_port"`
	DefChainId      string `toml:"def_chain_id"`
	DefName         string `toml:"def_name"`
	DefDesc         string `toml:"def_desc"`
}
