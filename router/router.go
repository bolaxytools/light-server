
// 这是一份用于实验GoSwagger的学习代码
//
//	Schemes: http, https
//	Host: 192.168.10.153
//	BasePath: v1
//	Version: 0.0.1
//
// swagger:meta
package controller

import (
	"fmt"
	"github.com/alecthomas/log4go"
	"github.com/gin-gonic/gin"
	"wallet-svc/config"
)

var (
	engine *gin.Engine
)

func InitRouter() {
	engine = gin.Default()
	initAssetRouter()
	initTxRouter()
	initLeagueRouter()
	servAddr := fmt.Sprintf(":%d", config.Cfg.Global.Port)
	log4go.Info("InitRouter: [%v]", servAddr)

	err := engine.Run(servAddr)
	if err != nil {
		log4go.Info("gin start http service err=%v\n", err)
	}
	log4go.Info("api初始化完毕，原地待命")
}
