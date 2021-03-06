package controller

import (
	"encoding/json"
	"github.com/alecthomas/log4go"
	"github.com/gin-gonic/gin"
	"strings"
	"wallet-svc/domain"
	"wallet-svc/dto/req"
	"wallet-svc/dto/resp"
	"wallet-svc/model"
	"wallet-svc/werror"

	"net/http"
)

func initAssetRouter() {
	grp := engine.Group("asset", func(context *gin.Context) {})
	grp.POST("getbalance", getbalance)
	grp.POST("getnonce", getNonce)
	grp.POST("followtoken", followToken)
	grp.POST("searchtoken", searchToken)
	grp.POST("tokeninfo", getTokenInfo)
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

	asts, er := flr.QueryAddrAssets(1, 100, inner.Addr)
	if er != nil {
		log4go.Info("flr.QueryAddrAssets error=%v\n", asts)
	}

	coinbox := &resp.AssetBox{
		MainCoin:    &model.Asset{Symbol: "BUSD", Balance: n.Balance.String(), Logo: "https://cdn.mytoken.org/Frdw6OBZGQhL5WaU2zvJEBgrh3FK", Desc: "BUSD", Decimals: 18,Bap:21001},
		ExtCoinList: asts,
	}

	c.JSON(http.StatusOK, resp.NewSuccessResp(coinbox))
}

/*
	获取子token余额
*/
func getTokenInfo(c *gin.Context) {
	reqdata := new(req.ReqData)
	err := c.BindJSON(reqdata)
	if err != nil {
		c.JSON(http.StatusOK, resp.BindJsonErrorResp(err.Error()))
		return
	}
	inner := new(req.ReqTokenInfo)
	err = reqdata.Reverse(inner)
	if err != nil {
		c.JSON(http.StatusOK, resp.BindJsonErrorResp(err.Error()))
		return
	}

	flr := domain.NewBlockFollower()

	asts, er := flr.QueryAddrContractAsset(inner.Contract, inner.Addr)
	if er != nil {
		log4go.Info("flr.QueryAddrAssets error=%v\n", asts)
	}

	ct,txs, err := domain.GetTokenHistory(inner.Contract,inner.Addr, inner.Page, inner.PageSize)
	if err != nil {
		c.JSON(http.StatusOK, resp.NewErrorResp(werror.QueryError, err.Error()))
		return
	}


	txlist := resp.NewTxHistory(txs, uint64(ct))

	resq := resp.NewChildInfo(asts, txlist)

	c.JSON(http.StatusOK, resp.NewSuccessResp(resq))
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

	log4go.Info("from:%s,nonce=%d\n", inner.Addr, n.Nonce)

	coinbox := &resp.NonceObj{
		Nonce: n.Nonce,
	}

	c.JSON(http.StatusOK, resp.NewSuccessResp(coinbox))
}

/*
	关注token种
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
	搜索token种
*/
func searchToken(c *gin.Context) {
	reqdata := new(req.ReqData)
	err := c.BindJSON(reqdata)
	if err != nil {
		c.JSON(http.StatusOK, resp.BindJsonErrorResp(err.Error()))
		return
	}

	jbf,_ := json.Marshal(reqdata)
	log4go.Info("search data=%s\n",jbf)

	inner := new(req.ReqSearch)
	err = reqdata.Reverse(inner)
	if err != nil {
		c.JSON(http.StatusOK, resp.BindJsonErrorResp(err.Error()))
		return
	}

	flr := domain.NewBlockFollower()



	tkns, r := flr.SearchToken(inner.Content, strings.ToLower(inner.Addr))
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
