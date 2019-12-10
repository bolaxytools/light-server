package mysql

import (
	"github.com/alecthomas/log4go"
	"github.com/jmoiron/sqlx"
	"wallet-svc/model"
)

type AddressDao struct {
	db *sqlx.DB
}

func NewAddressDao() *AddressDao {
	return &AddressDao{
		db: GetDb(),
	}

}

func (dao *AddressDao) Add(gd *model.Address) error {
	sql := "INSERT INTO `address`(`addr`, `add_time`, `update_time`) " +
		"VALUES " +
		"(:addr, :add_time, :update_time) ON DUPLICATE KEY UPDATE `update_time`=:update_time"
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

func (dao *AddressDao) QueryCount() (int64, error) {
	sql := "select " +
		"count(1) from address"
	var count int64
	err := dao.db.Get(&count, sql)

	if err != nil {
		return 0, err
	}

	log4go.Debug("query sql=%s,rows=%d\n", sql, count)

	return count, nil
}
