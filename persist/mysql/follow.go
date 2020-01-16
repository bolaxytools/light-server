package mysql

import (
	"github.com/alecthomas/log4go"
	"github.com/jmoiron/sqlx"
	"wallet-svc/model"
)

type FollowDao struct {
	db *sqlx.DB
}

func NewFollowDao() *FollowDao {
	return &FollowDao{
		db: GetDb(),
	}

}

func (dao *FollowDao) Add(gd *model.Follow) error {
	sql := "INSERT INTO `follow`(`contract`, `wallet`, `balance`,`followed`) " +
		"VALUES " +
		"(:contract, :wallet, :balance,:followed) ON DUPLICATE KEY UPDATE `balance`=:balance,`followed`=:followed"

	log4go.Info("FollowDao.Add.sql=%s,contract=%s,wallet=%s,followed=%t",sql,gd.Contract,gd.Wallet,gd.Followed)
	re, err := dao.db.NamedExec(sql, gd)

	if err != nil {
		return err
	}

	lid, er := re.LastInsertId()
	if er != nil {
		return er
	}
	log4go.Debug("INSERT INTO `follow` result=%d\n", lid)

	return nil
}

func (dao *FollowDao) Update(gd *model.Follow) error {
	sql := "update `follow` set  `balance`=:balance where `contract`=:contract and `wallet`=:wallet"

	log4go.Info("FollowDao.Add.sql=%s,contract=%s,wallet=%s,followed=%t",sql,gd.Contract,gd.Wallet,gd.Followed)
	re, err := dao.db.NamedExec(sql, gd)

	if err != nil {
		return err
	}

	lid, er := re.RowsAffected()
	if er != nil {
		return er
	}
	log4go.Debug("INSERT INTO `follow` RowsAffected=%d\n", lid)

	return nil
}

func (dao *FollowDao) QueryCount() (int64, error) {
	sql := "select " +
		"count(1) from follow"
	var count int64
	err := dao.db.Get(&count, sql)

	if err != nil {
		return 0, err
	}

	log4go.Debug("query sql=%s,rows=%d\n", sql, count)

	return count, nil
}

func (dao *FollowDao) QueryBrcs(page, pageSize int32, addr string) ([]*model.Asset, error) {

	if page < 1 {
		page = 1
	}
	if pageSize < 5 {
		pageSize = 5
	}

	sql := "select " +
		"* FROM follow f LEFT JOIN token t on t.contract = f.contract where f.wallet=? and f.followed=true"

	rows, err := dao.db.Queryx(sql, addr, addr, (page-1)*pageSize, pageSize)

	if err != nil {
		return nil, err
	}

	var txs []*model.Asset

	for rows.Next() {
		tx := new(model.Asset)
		er := rows.StructScan(tx)
		if er != nil {
			return nil, er
		}
		txs = append(txs, tx)
	}

	log4go.Debug("query sql=%s,rows=%d\n", sql, len(txs))

	return txs, nil

}
