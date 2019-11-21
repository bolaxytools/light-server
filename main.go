package main

import (
	"github.com/alecthomas/log4go"
	"os"
	"path/filepath"
	"wallet-svc/config"
	"wallet-svc/domain"
	"wallet-svc/persist/mysql"
	controller "wallet-svc/router"
)

func Getwd() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	dir = "/Users/rudy/gospace/wallet-svc"
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

	follower := domain.NewBlockFollower()

	go follower.FollowBlockChain()

	controller.InitRouter()
}
