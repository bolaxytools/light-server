package mysql

import (
	"bytes"
	"encoding/hex"
	"github.com/alecthomas/log4go"
	"github.com/boxproject/bolaxy/cmd/sdk"
	"github.com/boxproject/bolaxy/common"
	"github.com/boxproject/bolaxy/rlp"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"golang.org/x/crypto/sha3"
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
		return errors.New("empty txs")
	}
	sql := "INSERT INTO `tx`(`tx_hash`, `addr_from`, `addr_to`, `block_height`, `tx_time`, `memo`,`amount`,`miner_fee`) " +
		"VALUES " +
		"(?, ?, ?, ?, ?, ?,?,?)"

	sql = "INSERT INTO `tx`(`tx_hash`, `addr_from`, `addr_to`, `block_height`, `tx_time`, `memo`,`amount`,`miner_fee`) " +
		"VALUES "
	buf := bytes.NewBufferString(sql)

	arrLen := len(gds)
	if arrLen > 0 { //表明redis中已经有数据

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



		for _, redistick := range gds {
			datastr := hex.EncodeToString(redistick.Data)
			if len(datastr)>100 {
				datastr = datastr[0:100]
			}
			args = append(args, redistick.Hash, redistick.From, redistick.To, height, txTime, datastr, redistick.Value, redistick.Value)
		}
		_, errex := stmt.Exec(args...)

		if errex != nil {
			return errex
		}
		log4go.Info("block height=%d insert %d txs.\n",height,arrLen)
	}

	return nil
}

func rlpHash(x interface{}) (h common.Hash) {
	hw := sha3.NewLegacyKeccak256()
	rlp.Encode(hw, x)
	hw.Sum(h[:0])
	return h
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

/*
["+IUDgIJSCJScfQsvYzx4lvB7ZPH1/nHnSBab9IInEICAgKAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAIAloNnnwdkb0/N9utWHiqefailLWf0GKPda5ysnCXil1/QsoHYVbAB36VLeTdFka1Qb69Sj26+bIsLOwVhHDlWio7JJ"],"InternalTransactions":[],"InternalTransactionReceipts":[]},"Signatures":{"0X03573555DCF1B4518816DF6F1EDFF2C16B4D29F64CF6D9BDA78ECB6F50826CCA0B":"0x2799d38bd4879df08e9cc2b0f87aac79506f11ccd3d94c78615d2453085e4c3a318a43dfd31593aef136072642d6a794a7fb8c744c1415a839ce5ac30d043a8401"}}
["+IUEgIJSCJScfQsvYzx4lvB7ZPH1/nHnSBab9IInEICAgKAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAIAloFDeqURUbeJM/yYrM6ma8eFYcDISGwMZIJW0lrfB+tPFoEzJV6dwxj9RnLCdx7H31HV1rmy7rztX+4XRbaBnZB3c"],"InternalTransactions":[],"InternalTransactionReceipts":[]},"Signatures":{"0X03573555DCF1B4518816DF6F1EDFF2C16B4D29F64CF6D9BDA78ECB6F50826CCA0B":"0x6c18d4ddb4b086ba469a212a51281bd71436fe50729037215941d796a4419d6d13fb99ef3ce92ac11231775d07ad9d46c696beb62593db329c9165af481dbccd01"}}
*/
