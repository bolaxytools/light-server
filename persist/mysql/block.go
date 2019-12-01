package mysql

import (
	"bytes"
	"github.com/alecthomas/log4go"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"wallet-svc/model"
)

type BlockDao struct {
	db *sqlx.DB
}

func NewBlockDao() *BlockDao {
	return &BlockDao{
		db: GetDb(),
	}
}

func (dao *BlockDao) Add(gd *model.Block) error {
	sql := "INSERT INTO `block`(`height`, `hash`, `tx_count`, `block_time`) " +
		"VALUES " +
		"(:height, :hash, :tx_count, :block_time)"
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

func (dao *BlockDao) Query(page, pageSize int32) ([]*model.Block, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 5 {
		pageSize = 5
	}
	sql := "select " +
		"* from block order by block_time desc limit ?,?"
	rows, err := dao.db.Queryx(sql, (page-1)*pageSize, pageSize)

	if err != nil {
		return nil, err
	}

	var txs []*model.Block

	for rows.Next() {
		tx := new(model.Block)
		er := rows.StructScan(tx)
		if er != nil {
			return nil, er
		}
		txs = append(txs, tx)
	}

	log4go.Debug("query sql=%s,rows=%d\n", sql, len(txs))

	return txs, nil
}

func (dao *BlockDao) GetBlockByHeight(height string) (*model.Block, error) {
	blk := new(model.Block)
	sql := "select " +
		"* from `block` b where b.height =?"
	err := dao.db.Get(blk, sql, height)

	if err != nil {
		return nil, err
	}

	return blk, nil
}

func (dao *BlockDao) GetBlockByHeightX(height uint64) (*model.Block, error) {
	blk := new(model.Block)
	sql := "select " +
		"* from `block` b where b.height =?"
	err := dao.db.Get(blk, sql, height)

	if err != nil {
		return nil, err
	}

	sql = "SELECT " +
		"signer_address from block_signer where b_height=?"
	var signers []string
	er := dao.db.Select(&signers,sql,height)
	if er != nil {
		log4go.Info("dao.db.Select.signers error:%v\n",er)
	}else{
		blk.Signers = signers
	}

	return blk, nil
}

func (dao *BlockDao) QueryCount() (int64, error) {
	sql := "select " +
		"count(1) from block"
	var count int64
	err := dao.db.Get(&count, sql)

	if err != nil {
		return 0, err
	}

	log4go.Debug("query sql of tx count=%s,rows=%d\n", sql, count)

	return count, nil
}

func (dao *BlockDao) BatchAddSingers(gds map[string]string, height uint64) error {
	if gds == nil || len(gds) < 1 {
		return errors.New("empty gds")
	}

	sql := "INSERT INTO block_signer (b_height,signer_address) values "
	buf := bytes.NewBufferString(sql)

	arrLen := len(gds)
	if arrLen > 0 { //表明redis中已经有数据

		for i := 0; i < arrLen; i++ {
			if i == (arrLen - 1) {
				buf.WriteString("(?, ?)")
			} else {
				buf.WriteString("(?, ?),")
			}
		}

		finalSql := buf.String()

		stmt, err := dao.db.Prepare(finalSql)

		if err != nil {
			return err
		}

		defer stmt.Close()

		args := make([]interface{}, 0)

		for k, _ := range gds {
			args = append(args, height, k)
		}
		_, errex := stmt.Exec(args...)

		if errex != nil {
			return errex
		}
		log4go.Info("block height=%d add %d signers success.\n", height, arrLen)
	}

	return nil
}
