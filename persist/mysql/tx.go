package mysql

import (
	"bytes"
	"github.com/alecthomas/log4go"
	"github.com/boxproject/bolaxy/cmd/sdk"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"wallet-svc/model"
)

type TxDao struct {
	db *sqlx.DB
}

func NewTxDao() *TxDao {
	return &TxDao{
		db: GetDb(),
	}

}

func (dao *TxDao) Add(gd *model.Tx) error {
	sql := "INSERT INTO `tx`(`tx_hash`, `addr_from`, `addr_to`, `block_height`, `tx_time`, `memo`) " +
		"VALUES " +
		"(:tx_hash, :addr_from, :addr_to, :block_height, :tx_time, :memo)"
	re, err := dao.db.NamedExec(sql, gd)

	if err != nil {
		return err
	}

	lid, er := re.LastInsertId()
	if er != nil {
		return er
	}
	log4go.Debug("INSERT INTO `asset` result=%d\n", lid)

	return nil
}

func (dao *TxDao) BatchSave(gds []*sdk.Transaction, height int64, txTime int64) error {
	if gds == nil || len(gds) < 1 {
		return errors.New("BatchSave empty txs")
	}
	sql := "INSERT INTO `tx`(`tx_hash`, `addr_from`, `addr_to`, `block_height`, `tx_time`, `memo`,`amount`,`miner_fee`) " +
		"VALUES " +
		"(?, ?, ?, ?, ?, ?,?,?)"

	sql = "INSERT INTO `tx`(`tx_hash`, `addr_from`, `addr_to`, `block_height`, `tx_time`, `memo`,`amount`,`miner_fee`) " +
		"VALUES "
	buf := bytes.NewBufferString(sql)

	arrLen := len(gds)
	if arrLen > 0 {

		for i := 0; i < arrLen; i++ {
			if i == (arrLen - 1) {
				buf.WriteString("(?, ?, ?, ?, ?, ?,?,?)")
			} else {
				buf.WriteString("(?, ?, ?, ?, ?, ?,?,?),")
			}
		}

		finalSql := buf.String()

		stmt, err := dao.db.Prepare(finalSql)

		if err != nil {
			return err
		}

		defer stmt.Close()

		args := make([]interface{}, 0)

		for _, rtx := range gds {
			datastr := ""

			_, c, sbf := sdk.UnWrapData(rtx.Data)

			if c == sdk.DataDisplay {
				datastr = string(sbf)
			}

			if len(datastr) > 100 {
				datastr = datastr[0:100]
			}
			args = append(args, rtx.Hash, rtx.From, rtx.To, height, txTime, datastr, rtx.Value, rtx.Gas)
		}
		_, errex := stmt.Exec(args...)

		if errex != nil {
			return errex
		}
		log4go.Info("block height=%d insert %d txs.\n", height, arrLen)
	}

	return nil
}

func (dao *TxDao) Query(addr string, page, pageSize int32) ([]*model.Tx, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 5 {
		pageSize = 5
	}
	sql := "select " +
		"* from tx where addr_to = ? or addr_from = ? order by tx_time desc limit ?,?"
	rows, err := dao.db.Queryx(sql, addr, addr, (page-1)*pageSize, pageSize)

	if err != nil {
		return nil, err
	}

	var txs []*model.Tx

	for rows.Next() {
		tx := new(model.Tx)
		er := rows.StructScan(tx)
		if er != nil {
			return nil, er
		}
		txs = append(txs, tx)
	}

	log4go.Debug("query sql=%s,rows=%d\n", sql, len(txs))

	return txs, nil
}

func (dao *TxDao) QueryLatestTx(page, pageSize int32) ([]*model.Tx, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 5 {
		pageSize = 5
	}
	sql := "select " +
		"* from tx order by tx_time desc limit ?,?"
	rows, err := dao.db.Queryx(sql, (page-1)*pageSize, pageSize)

	if err != nil {
		return nil, err
	}

	var txs []*model.Tx

	for rows.Next() {
		tx := new(model.Tx)
		er := rows.StructScan(tx)
		if er != nil {
			return nil, er
		}
		txs = append(txs, tx)
	}

	log4go.Debug("query sql=%s,rows=%d\n", sql, len(txs))

	return txs, nil
}

func (dao *TxDao) GetTxByHash(txHash string) (*model.Tx, error) {

	sql := "select " +
		"* from `tx` t where t.tx_hash=?"
	tx := new(model.Tx)
	err := dao.db.Get(tx, sql, txHash)

	if err != nil {
		return nil, err
	}
	return tx, nil
}

func (dao *TxDao) QueryCount() (int64, error) {
	sql := "select " +
		"count(1) from tx"
	var count int64
	err := dao.db.Get(&count, sql)

	if err != nil {
		return 0, err
	}

	log4go.Debug("query sql of tx count=%s,rows=%d\n", sql, count)

	return count, nil
}

func (dao *TxDao) QueryContractTxCount(contract string) (int64, error) {
	sql := "select " +
		"count(1) from tx where contract=?"
	var count int64
	err := dao.db.Get(&count, sql, contract)

	if err != nil {
		return 0, err
	}

	log4go.Debug("query sql of tx count=%s,rows=%d\n", sql, count)

	return count, nil
}

func (dao *TxDao) BatchAddTx(gds []*model.Tx) error {
	if gds == nil || len(gds) < 1 {
		return errors.New("BatchAddTx empty txs")
	}
	sql := "INSERT INTO `tx`(`tx_hash`, `addr_from`, `addr_to`, `block_height`, `tx_time`, `memo`,`amount`,`miner_fee`,`tx_type`,`contract`) " +
		"VALUES " +
		"(?, ?, ?, ?, ?, ?,?,?)"

	sql = "INSERT INTO `tx`(`tx_hash`, `addr_from`, `addr_to`, `block_height`, `tx_time`, `memo`,`amount`,`miner_fee`,`tx_type`,`contract`) " +
		"VALUES "
	buf := bytes.NewBufferString(sql)

	arrLen := len(gds)
	if arrLen > 0 {

		for i := 0; i < arrLen; i++ {
			if i == (arrLen - 1) {
				buf.WriteString("(?, ?, ?, ?, ?, ?,?,?,?,?)")
			} else {
				buf.WriteString("(?, ?, ?, ?, ?, ?,?,?,?,?),")
			}
		}

		finalSql := buf.String()

		stmt, err := dao.db.Prepare(finalSql)

		if err != nil {
			return err
		}

		defer stmt.Close()

		args := make([]interface{}, 0)

		for _, rtx := range gds {
			args = append(args, rtx.TxHash, rtx.AddrFrom, rtx.AddrTo, rtx.BlockHeight, rtx.TxTime, rtx.Memo, rtx.Amount, rtx.MinerFee, rtx.TxType, rtx.Contract)
		}
		_, errex := stmt.Exec(args...)

		if errex != nil {
			return errex
		}
		log4go.Info("block insert %d txs.\n", arrLen)
	}

	return nil
}
