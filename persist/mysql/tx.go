package mysql

import (
	"github.com/alecthomas/log4go"
	"github.com/jmoiron/sqlx"
	"wallet-service/model"
)

type TxDao struct {
	db *sqlx.DB
}

func NewTxDao() *TxDao {
	return &TxDao{
		db:GetDb(),
	}

}

func (dao *TxDao) Add(gd *model.Tx) error {
	sql := "INSERT INTO `tx`(`tx_hash`, `addr_from`, `addr_to`, `block_height`, `tx_time`, `memo`) " +
		"VALUES " +
		"(:tx_hash, :addr_from, :addr_to, :block_height, :tx_time, :memo)"
	re,err := dao.db.NamedExec(sql,gd)

	if err != nil {
		return err
	}

	lid, er := re.LastInsertId()
	if er != nil {
		return er
	}
	log4go.Info("INSERT INTO `asset` result=%d\n", lid)

	return nil
}

func (dao *TxDao) Query(addr string) ([]*model.Tx,error) {
	sql := "select " +
		"* from tx where addr_to = ? or addr_from = ? order by tx_time desc"
	rows,err := dao.db.Queryx(sql,addr,addr)

	if err != nil {
		return nil,err
	}

	var txs []*model.Tx


	for rows.Next() {
		tx := new(model.Tx)
		er := rows.StructScan(tx)
		if er != nil {
			return nil,er
		}
		txs = append(txs, tx)
	}

	log4go.Debug("query sql=%s,rows=%d\n", sql,len(txs))

	return txs,nil
}