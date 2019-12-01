package domain

import (
	"wallet-svc/model"
	"wallet-svc/persist/mysql"
)

func GetHistory(addr string,page,size int32) ([]*model.Tx,error)  {
	dao := mysql.NewTxDao()
	txs,err := dao.Query(addr,page,size)
	return txs,err
}

func GetLatestTx(page,size int32) ([]*model.Tx,error)  {
	dao := mysql.NewTxDao()
	txs,err := dao.QueryLatestTx(page,size)
	return txs,err
}

func GetTxById(txHash string) (*model.Tx,error)  {
	dao := mysql.NewTxDao()
	txs,err := dao.GetTxByHash(txHash)
	return txs,err
}

func GetHistoryBlock(page,size int32) ([]*model.Block,error)  {
	dao := mysql.NewBlockDao()
	txs,err := dao.Query(page,size)
	return txs,err
}

func GetBlockTotal() (int64,error)  {
	dao := mysql.NewBlockDao()
	txs,err := dao.QueryCount()
	return txs,err
}

func GetBlockById(txHash string) (*model.Block,error)  {
	dao := mysql.NewBlockDao()
	txs,err := dao.GetBlockByHeight(txHash)
	return txs,err
}

func GetBlockByHeight(height uint64) (*model.Block,error)  {
	dao := mysql.NewBlockDao()
	blk,err := dao.GetBlockByHeightX(height)
	return blk,err
}

func GetTxTotal() (uint64,error)  {
	dao := mysql.NewTxDao()
	c,err := dao.QueryCount()
	return uint64(c),err
}