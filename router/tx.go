package controller

import (
	"github.com/gin-gonic/gin"
	"wallet-svc/domain"
	"wallet-svc/dto/req"
	"wallet-svc/dto/resp"
	"wallet-svc/werror"

	"net/http"
)

func initTxRouter() {
	grp := engine.Group("tx", func(context *gin.Context) {})
	grp.POST("sendtx", sendTx)
	grp.POST("gethistory", getHistory)
}

/*
	发送交易
*/

func sendTx(c *gin.Context) {
	reqdata := new(req.ReqData)
	err := c.BindJSON(reqdata)
	if err != nil {
		c.JSON(http.StatusOK, resp.BindJsonErrorResp(err.Error()))
		return
	}
	inner := new(req.ReqSendTx)
	err = reqdata.Reverse(inner)
	if err != nil {
		c.JSON(http.StatusOK, resp.BindJsonErrorResp(err.Error()))
		return
	}

	//TODO send tx to block chain
	flr := domain.NewBlockFollower()

	txhash,err := flr.SendRawTx(inner.SignedTx)
	if err != nil {
		c.JSON(http.StatusOK, resp.BindJsonErrorResp(err.Error()))
		return
	}

	rsp := &resp.SendTxResp{
		TxHash:txhash,
	}



	c.JSON(http.StatusOK, resp.NewSuccessResp(rsp))
}


/*
	获取余额
*/
// swagger:route POST /asset/getbalance users UpdateUserResponseWrapper
func getHistory(c *gin.Context) {
	reqdata := new(req.ReqData)
	err := c.BindJSON(reqdata)
	if err != nil {
		c.JSON(http.StatusOK, resp.BindJsonErrorResp(err.Error()))
		return
	}
	inner := new(req.ReqListBase)
	err = reqdata.Reverse(inner)
	if err != nil {
		c.JSON(http.StatusOK, resp.BindJsonErrorResp(err.Error()))
		return
	}

	txs,err := domain.GetHistory(inner.Addr,inner.Page,inner.PageSize)
	if err != nil {
		c.JSON(http.StatusOK, resp.NewErrorResp(werror.QueryError,err.Error()))
		return
	}

	c.JSON(http.StatusOK, resp.NewSuccessResp(resp.NewTxHistory(txs)))
}