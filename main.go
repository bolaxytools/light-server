package main

import (
	"github.com/alecthomas/log4go"
	"os"
	"path/filepath"
	"wallet-service/config"
	"wallet-service/persist/mysql"
	controller "wallet-service/router"
)

func Getwd() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	//dir = "/Users/rudy/gospace/wallet-service"
	if err != nil {
		panic("No caller information")
	}
	return dir
}

func main() {
	defer log4go.Close()

	log4go.Info("钱包服务初始化开始......")
	config.LoadConfig(Getwd())
	mysql.InitMySQL()
	controller.InitRouter()
}
