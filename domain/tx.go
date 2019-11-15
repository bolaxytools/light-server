package domain

import (
	"wallet-svc/model"
	"wallet-svc/persist/mysql"
)

func GetHistory(addr string,page,size int32) ([]*model.Tx,error)  {
	dao := mysql.NewTxDao()
	txs,err := dao.Query(addr)
	return txs,err
}