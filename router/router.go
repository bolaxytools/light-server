
package controller

import (
	"fmt"
	"github.com/alecthomas/log4go"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"wallet-svc/config"
)

var (
	engine *gin.Engine
)

func InitRouter() {
	engine = gin.Default()


	engine.LoadHTMLGlob("resources/*/*.html")
	files, err := ioutil.ReadDir("resources/pages")
	if err != nil {
		log4go.Info("os.Open resources/templates error=%v\n")
		os.Exit(0)
	}

	for _, f := range files {
		fn := f.Name()
		relfn := strings.Replace(fn, ".html", "", 1)
		engine.GET(relfn, func(context *gin.Context) {
			context.HTML(http.StatusOK, fn, gin.H{})
		})
	}

	engine.Static("img", "resources/static/img")
	engine.Static("css", "resources/static/css")

	initExplorerRouter()
	initAssetRouter()
	initTxRouter()
	initLeagueRouter()
	servAddr := fmt.Sprintf(":%d", config.Cfg.Global.Port)
	log4go.Info("InitRouter: [%v]", servAddr)

	err = engine.Run(servAddr)
	if err != nil {
		log4go.Info("gin start http service err=%v\n", err)
	}
	log4go.Info("api初始化完毕，原地待命")
}
