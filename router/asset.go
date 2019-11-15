package controller

import (
	"github.com/gin-gonic/gin"
	"wallet-svc/domain"
	"wallet-svc/dto/req"
	"wallet-svc/dto/resp"

	"net/http"
)

func initAssetRouter() {
	grp := engine.Group("asset", func(context *gin.Context) {})
	grp.POST("getbalance", getbalance)
	grp.POST("getnonce",  getNonce)
}

/*
	获取余额
*/
// swagger:route POST /asset/getbalance users UpdateUserResponseWrapper
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

	coinbox := &resp.AssetBox{
		MainCoin:&resp.Asset{Symbol:"Box",Balance:"100"},
		ExtCoinList:[]*resp.Asset{&resp.Asset{Symbol:"Brc1",Balance:"100000"},&resp.Asset{Symbol:"Brc5",Balance:"900000"}},
	}

	c.JSON(http.StatusOK, resp.NewSuccessResp(coinbox))
}


/*
	获取nonce值
*/
// swagger:route POST /asset/getbalance users UpdateUserResponseWrapper
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

	n,r := flr.GetNonce(inner.Addr)
	if r != nil {
		c.JSON(http.StatusOK, resp.BindJsonErrorResp(r.Error()))
		return
	}

	coinbox := &resp.NonceObj{
		Nonce:n,
	}

	c.JSON(http.StatusOK, resp.NewSuccessResp(coinbox))
}