package mysql

import (
	"fmt"
	"github.com/alecthomas/log4go"
	"github.com/jmoiron/sqlx"
	"wallet-svc/model"
)

type TokenDao struct {
	db *sqlx.DB
}

func NewTokenDao() *TokenDao {
	return &TokenDao{
		db: GetDb(),
	}

}

func (dao *TokenDao) Add(gd *model.Token) error {
	sql := "INSERT INTO `token`(`contract`, `symbol`, `logo`) " +
		"VALUES " +
		"(:contract, :symbol, :logo) ON DUPLICATE KEY UPDATE `logo`=:logo"
	re, err := dao.db.NamedExec(sql, gd)

	if err != nil {
		return err
	}

	lid, er := re.LastInsertId()
	if er != nil {
		return er
	}
	log4go.Info("INSERT INTO `token` result=%d\n", lid)

	return nil
}

func (dao *TokenDao) QueryCount() (int64, error) {
	sql := "select " +
		"count(1) from token"
	var count int64
	err := dao.db.Get(&count, sql)

	if err != nil {
		return 0, err
	}

	log4go.Debug("query sql=%s,rows=%d\n", sql, count)

	return count, nil
}

func (dao *TokenDao) queryTokenByAddr(addr string) {

}

func (dao *FollowDao) QueryTokenByContract(page, pageSize int32, content, addr string) ([]*model.Token, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	sql := "select " +
		"* FROM token t where t.contract=? or t.symbol like ? limit ?,?"

	like := fmt.Sprintf("%s%s%s", "%", content, "%")

	rows, err := dao.db.Queryx(sql, content, like, (page-1)*pageSize, pageSize)

	if err != nil {
		return nil, err
	}

	var txs []*model.Token

	var ctrcts []string

	for rows.Next() {
		tkn := new(model.Token)
		er := rows.StructScan(tkn)
		if er != nil {
			return nil, er
		}

		ctrcts = append(ctrcts, tkn.Contract)

		txs = append(txs, tkn)
	}

	var addrs []string

	if len(ctrcts) > 0 {
		sql2 := "SELECT f.contract from  follow f  WHERE f.wallet = ? and f.contract in (?)"

		query, args, err := sqlx.In(sql2, addr, ctrcts)
		if err != nil {
			return nil, err
		}
		query = dao.db.Rebind(query)

		er := dao.db.Select(&addrs, query, args...)

		if er != nil {
			return nil, er
		}

		for _, ctt := range addrs {
			for _, tknnnn := range txs {
				if ctt == tknnnn.Contract {
					tknnnn.Followed = true
				}
			}
		}

	}

	log4go.Debug("query sql=%s,rows=%d\n", sql, len(txs))

	return txs, nil

}

func (dao *TokenDao) QueryCountByContent(content string) (uint64, error) {
	sql := "select " +
		"count(1) FROM token t where t.contract=? or t.symbol like ?"
	var count uint64
	err := dao.db.Get(&count, sql, content, content)

	if err != nil {
		return 0, err
	}

	log4go.Debug("query count sql=%s,rows=%d\n", sql, count)

	return count, nil
}
