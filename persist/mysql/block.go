package mysql

import (
	"github.com/alecthomas/log4go"
	"github.com/jmoiron/sqlx"
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
		"(:height, :hash, :addr_to, :tx_count, :block_time)"
	re, err := dao.db.NamedExec(sql, gd)

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

	return blk, nil
}
