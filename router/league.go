package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"wallet-svc/config"
	"wallet-svc/domain"
	"wallet-svc/dto/req"
	"wallet-svc/dto/resp"
)

func initLeagueRouter() {
	grp := engine.Group("league", func(context *gin.Context) {})
	grp.POST("checkjoin", checkJoin)
	grp.POST("getdefaultleague", getDefaultLeague)
}

func checkJoin(c *gin.Context) {

	reqdata := new(req.ReqData)
	err := c.BindJSON(reqdata)
	if err != nil {
		c.JSON(http.StatusOK, resp.BindJsonErrorResp(err.Error()))
		return
	}
	inner := new(req.ReqBase)
	err = reqdata.Reverse(inner)
	if err != nil {
		c.JSON(http.StatusOK, resp.BindJsonErrorResp(err.Error()))
		return
	}

	flr := domain.NewBlockFollower()

	white := flr.CheckWhiteList(inner.Addr)
	black := flr.CheckWhiteList(inner.Addr)

	//暂时先返回true，允许所有人拉取信息
	rsp := &resp.CheckJoinResp{
		Allow: !black || white || (!white && !black),
	}

	c.JSON(http.StatusOK, resp.NewSuccessResp(rsp))
}

func getDefaultLeague(c *gin.Context) {

	reqdata := new(req.ReqData)
	err := c.BindJSON(reqdata)
	if err != nil {
		c.JSON(http.StatusOK, resp.BindJsonErrorResp(err.Error()))
		return
	}
	inner := new(req.ReqBase)
	err = reqdata.Reverse(inner)
	if err != nil {
		c.JSON(http.StatusOK, resp.BindJsonErrorResp(err.Error()))
		return
	}

	glb := config.Cfg.Global

	//暂时先返回true，允许所有人拉取信息
	rsp := resp.NewLeagueItem(glb.DefChainId, glb.DefName, glb.DefHost, glb.DefDesc, glb.DefPort)

	c.JSON(http.StatusOK, resp.NewSuccessResp(rsp))
}
