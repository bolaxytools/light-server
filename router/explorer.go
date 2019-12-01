package controller

import (
	"github.com/gin-gonic/gin"
	"math/big"
	"net/http"
	"wallet-svc/domain"
	"wallet-svc/dto/req"
	"wallet-svc/dto/resp"
	"wallet-svc/werror"
)

const (
	addr_len = 42
	hash_len = 66
)

func initExplorerRouter() {
	group := engine.Group("explore", func(context *gin.Context) {})
	group.POST("index", index)
	group.POST("txlist", getTxHistory)
	group.POST("getblock", getHistoryBlock)
	group.POST("getassets", getAssetInfo)
	group.POST("search", search)
	group.POST("getblockbyid", getBlockById)
}

/*
	首页信息
*/
func index(c *gin.Context) {
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
	hei,err := flr.GetCurrentBlockHeight()
	if err != nil {
		c.JSON(http.StatusOK, resp.BindJsonErrorResp(err.Error()))
		return
	}


	txtotal,er := domain.GetTxTotal()
	if er != nil {
		c.JSON(http.StatusOK, resp.BindJsonErrorResp(er.Error()))
		return
	}

	ret := &resp.IndexRet{
		ChainId:       "chainId10011",
		BlockCount:    uint64(hei),
		AddressCount:  flr.GetAddressCount(),
		MainCoinCount: 72774,
		TxCount:       txtotal,
		CrossMax:      100000,
		GasCostCount:  29929229,
	}

	txs, err := domain.GetLatestTx(1, 5)
	if err != nil {
		c.JSON(http.StatusOK, resp.NewErrorResp(werror.QueryError, err.Error()))
		return
	}
	ret.Txs = txs

	blocks, err := domain.GetHistoryBlock(1, 5)
	if err != nil {
		c.JSON(http.StatusOK, resp.NewErrorResp(werror.QueryError, err.Error()))
		return
	}
	ret.Blocks = blocks

	c.JSON(http.StatusOK, resp.NewSuccessResp(ret))
}

/*交易列表*/
func getTxHistory(c *gin.Context) {
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

	txs, err := domain.GetLatestTx(inner.Page, inner.PageSize)
	if err != nil {
		c.JSON(http.StatusOK, resp.NewErrorResp(werror.QueryError, err.Error()))
		return
	}

	c.JSON(http.StatusOK, resp.NewSuccessResp(resp.NewTxHistory(txs)))
}

/*
	获取历史区块
*/
func getHistoryBlock(c *gin.Context) {
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

	blk, err := domain.GetHistoryBlock(inner.Page, inner.PageSize)
	if err != nil {
		c.JSON(http.StatusOK, resp.NewErrorResp(werror.QueryError, err.Error()))
		return
	}

	c.JSON(http.StatusOK, resp.NewSuccessResp(resp.NewBlockHistory(blk)))
}

/*
	获取资产列表
*/
func getAssetInfo(c *gin.Context) {
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

	assets := make([]*resp.AssetInfo, 2)
	assets[0] = &resp.AssetInfo{
		Name:     "酒财币",
		Contract: "0xaaafffbbbccceeed0002223",
		Type:     "积分币",
		Symbol:   "JCB",
		Quantity: 23323244,
	}

	assets[1] = &resp.AssetInfo{
		Name:     "二哈币",
		Contract: "0xaaafffbbbccceeed0002224",
		Type:     "BRCn",
		Symbol:   "RHB",
		Quantity: 23323245,
	}

	c.JSON(http.StatusOK, resp.NewSuccessResp(resp.NewAssetList(assets)))
}

func search(c *gin.Context) {
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

	lenth := len(inner.Content)

	if lenth == addr_len { //搜地址
		assets := make([]*resp.AssetInfo, 2)
		assets[0] = &resp.AssetInfo{
			Name:     "酒财币",
			Contract: "0xaaafffbbbccceeed0002223",
			Type:     "积分币",
			Symbol:   "JCB",
			Quantity: 23323244,
		}

		assets[1] = &resp.AssetInfo{
			Name:     "二哈币",
			Contract: "0xaaafffbbbccceeed0002224",
			Type:     "BRCn",
			Symbol:   "RHB",
			Quantity: 23323245,
		}
		ret := resp.NewSearchRet(resp.Ret_Addr,resp.NewAssetList(assets))
		c.JSON(http.StatusOK, resp.NewSuccessResp(ret))
		return
	} else if lenth == hash_len { //搜交易hash
		tx, err := domain.GetTxById(inner.Content)
		if err != nil {
			c.JSON(http.StatusOK, resp.NewErrorResp(werror.QueryError, err.Error()))
			return
		}
		ret := resp.NewSearchRet(resp.Ret_Hash,tx)
		c.JSON(http.StatusOK, resp.NewSuccessResp(ret))
		return
	} else { //搜区块高度

		heit,_ := big.NewInt(0).SetString(inner.Content,0)

		blk, err := domain.GetBlockById(heit.Uint64())
		if err != nil {
			c.JSON(http.StatusOK, resp.NewErrorResp(werror.QueryError, err.Error()))
			return
		}

		ret := resp.NewSearchRet(resp.Ret_Hash,blk)
		c.JSON(http.StatusOK, resp.NewSuccessResp(ret))
		return
	}
}

/*
	获取指定高度区块
*/
func getBlockById(c *gin.Context) {
	reqdata := new(req.ReqData)
	err := c.BindJSON(reqdata)
	if err != nil {
		c.JSON(http.StatusOK, resp.BindJsonErrorResp(err.Error()))
		return
	}
	inner := new(req.ReqBlockHeight)
	err = reqdata.Reverse(inner)
	if err != nil {
		c.JSON(http.StatusOK, resp.BindJsonErrorResp(err.Error()))
		return
	}

	blk, err := domain.GetBlockById(inner.Height)
	if err != nil {
		c.JSON(http.StatusOK, resp.NewErrorResp(werror.QueryError, err.Error()))
		return
	}

	c.JSON(http.StatusOK, resp.NewSuccessResp(blk))
}