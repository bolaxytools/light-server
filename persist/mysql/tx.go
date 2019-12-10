package mysql

import (
	"github.com/alecthomas/log4go"
	"github.com/boxproject/bolaxy/cmd/sdk"
	"github.com/jmoiron/sqlx"
	"wallet-svc/model"
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

func (dao *TxDao) BashSave(gds []*sdk.Transaction,height int64,txTime int64) error {
	if gds == nil || len(gds) < 1 {
		return nil
	}
	sql := "INSERT INTO `tx`(`tx_hash`, `addr_from`, `addr_to`, `block_height`, `tx_time`, `memo`,`amount`,`miner_fee`) " +
		"VALUES " +
		"(?, ?, ?, ?, ?, ?,?,?)"

	tx,err := db.Beginx()

	if err != nil {
		return err
	}


	 //vss := make([][]interface{},len(gds))

	stt,er := tx.Preparex(db.Rebind(sql))


	if er != nil {
		return er
	}

	defer stt.Close()

	for _,rtx := range gds  {
		tmp := make([]interface{},8)
		tmp[0]=rtx.Hash
		tmp[1]=rtx.From
		tmp[2]=rtx.To
		tmp[3]=height
		tmp[4]=txTime
		tmp[5]=string(rtx.Data)
		tmp[6]=rtx.Value
		tmp[7]=rtx.Gas
		_, er := stt.Exec(tmp...)
		if er != nil {
			return er
		}

	}


	err = tx.Commit()
	if err != nil {
		return err

	}

	//log4go.Info("INSERT INTO `asset` result=%d\n", lid)

	return nil
}

func (dao *TxDao) Query(addr string,page,pageSize int32) ([]*model.Tx,error) {
	if page <1  {
		page = 1
	}
	if pageSize<5 {
		pageSize=5
	}
	sql := "select " +
		"* from tx where addr_to = ? or addr_from = ? order by tx_time desc limit ?,?"
	rows,err := dao.db.Queryx(sql,addr,addr,(page-1)*pageSize,pageSize)

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

func (dao *TxDao) QueryLatestTx(page,pageSize int32) ([]*model.Tx,error) {
	if page <1  {
		page = 1
	}
	if pageSize<5 {
		pageSize=5
	}
	sql := "select " +
		"* from tx order by tx_time desc limit ?,?"
	rows,err := dao.db.Queryx(sql,(page-1)*pageSize,pageSize)

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

func (dao *TxDao) GetTxByHash(txHash string) (*model.Tx,error) {

	sql := "select " +
		"* from `tx` t where t.tx_hash=?"
	tx  := new(model.Tx)
	err := dao.db.Get(tx,sql,txHash)


	if err != nil {
		return nil,err
	}


	return tx,nil
}