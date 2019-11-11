package controller

import (
	"github.com/gin-gonic/gin"
	"wallet-service/dto/req"
	"wallet-service/dto/resp"

	"net/http"
)

func initGoodsRouter() {
	grp := engine.Group("asset", func(context *gin.Context) {})
	grp.POST("getbalance", getbalance)
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
