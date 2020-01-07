package controller

import (
	"github.com/alecthomas/log4go"
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
	grp.POST("gethistory", getLatestTx)
	grp.POST("gettxbyhash", getTxById)
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
	err = reqdata.ReflectReverse(inner)
	if err != nil {
		c.JSON(http.StatusOK, resp.BindJsonErrorResp(err.Error()))
		return
	}

	flr := domain.NewBlockFollower()

	log4go.Info("交易发送者地址=%s，签过名的tx=%s\n",inner.Addr,inner.SignedTx)

	txhash, err := flr.SendRawTx(inner.SignedTx)
	if err != nil {
		c.JSON(http.StatusOK, resp.BindJsonErrorResp(err.Error()))
		return
	}

	rsp := &resp.SendTxResp{
		TxHash: txhash,
	}

	c.JSON(http.StatusOK, resp.NewSuccessResp(rsp))
}

/*
	获取历史交易
*/
func getLatestTx(c *gin.Context) {
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

	log4go.Info("tx gethistory page=%d,page_size=%d\n",inner.Page,inner.PageSize)

	ct,txs, err := domain.GetHistory(inner.Addr,inner.Page, inner.PageSize)
	if err != nil {
		c.JSON(http.StatusOK, resp.NewErrorResp(werror.QueryError, err.Error()))
		return
	}

	c.JSON(http.StatusOK, resp.NewSuccessResp(resp.NewTxHistory(txs,ct)))
}

/*
	获取指定txHash交易
*/
func getTxById(c *gin.Context) {
	reqdata := new(req.ReqData)
	err := c.BindJSON(reqdata)
	if err != nil {
		c.JSON(http.StatusOK, resp.BindJsonErrorResp(err.Error()))
		return
	}
	inner := new(req.ReqTxHash)
	err = reqdata.Reverse(inner)
	if err != nil {
		c.JSON(http.StatusOK, resp.BindJsonErrorResp(err.Error()))
		return
	}

	txs, err := domain.GetTxById(inner.Txnash)
	if err != nil {
		c.JSON(http.StatusOK, resp.NewErrorResp(werror.QueryError, err.Error()))
		return
	}

	c.JSON(http.StatusOK, resp.NewSuccessResp(txs))
}



