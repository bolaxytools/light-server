package controller

import (
	"github.com/alecthomas/log4go"
	"github.com/gin-gonic/gin"
	"wallet-svc/domain"
	"wallet-svc/dto/req"
	"wallet-svc/dto/resp"
	"wallet-svc/model"

	"net/http"
)

func initAssetRouter() {
	grp := engine.Group("asset", func(context *gin.Context) {})
	grp.POST("getbalance", getbalance)
	grp.POST("getnonce", getNonce)
	grp.POST("followtoken", followToken)
	grp.POST("searchtoken", searchToken)
}

/*
	获取余额
*/
func getbalance(c *gin.Context) {
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

	n, r := flr.GetAccount(inner.Addr)
	if r != nil {
		c.JSON(http.StatusOK, resp.BindJsonErrorResp(r.Error()))
		return
	}

	coinbox := &resp.AssetBox{
		MainCoin:    &model.Asset{Symbol: "BUSD", Balance: n.Balance.String()},
		ExtCoinList: []*model.Asset{&model.Asset{Symbol: "Brc1", Balance: "100000"}, &model.Asset{Symbol: "Brc5", Balance: "900000"}},
	}

	c.JSON(http.StatusOK, resp.NewSuccessResp(coinbox))
}

/*
	获取nonce值
*/
func getNonce(c *gin.Context) {
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

	n, r := flr.GetAccount(inner.Addr)
	if r != nil {
		c.JSON(http.StatusOK, resp.BindJsonErrorResp(r.Error()))
		return
	}

	log4go.Info("from:%s,nonce=%d\n",inner.Addr,n.Nonce)

	coinbox := &resp.NonceObj{
		Nonce: n.Nonce,
	}

	c.JSON(http.StatusOK, resp.NewSuccessResp(coinbox))
}

/*
	关注币种
*/
func followToken(c *gin.Context) {
	reqdata := new(req.ReqData)
	err := c.BindJSON(reqdata)
	if err != nil {
		c.JSON(http.StatusOK, resp.BindJsonErrorResp(err.Error()))
		return
	}
	inner := new(req.ReqFollow)
	err = reqdata.Reverse(inner)
	if err != nil {
		c.JSON(http.StatusOK, resp.BindJsonErrorResp(err.Error()))
		return
	}

	flr := domain.NewBlockFollower()

	r := flr.FollowToken(inner.Contract, inner.Addr, "0")
	if r != nil {
		c.JSON(http.StatusOK, resp.BindJsonErrorResp(r.Error()))
		return
	}

	c.JSON(http.StatusOK, resp.NewSuccessResp(nil))
}

/*
	搜索币种
*/
func searchToken(c *gin.Context) {
	reqdata := new(req.ReqData)
	err := c.BindJSON(reqdata)
	if err != nil {
		c.JSON(http.StatusOK, resp.BindJsonErrorResp(err.Error()))
		return
	}
	inner := new(req.ReqSearch)
	err = reqdata.Reverse(inner)
	if err != nil {
		c.JSON(http.StatusOK, resp.BindJsonErrorResp(err.Error()))
		return
	}

	flr := domain.NewBlockFollower()

	tkns, r := flr.SearchToken(inner.Content,inner.Addr)
	if r != nil {
		c.JSON(http.StatusOK, resp.BindJsonErrorResp(r.Error()))
		return
	}

	total, r := flr.QuerySearchTokenCount(inner.Content)

	if r != nil {
		c.JSON(http.StatusOK, resp.BindJsonErrorResp(r.Error()))
		return
	}

	ret := resp.NewSearchTokenRet(tkns, total)
	c.JSON(http.StatusOK, resp.NewSuccessResp(ret))
}
