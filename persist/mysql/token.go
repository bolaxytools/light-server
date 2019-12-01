package mysql

import (
	"fmt"
	"github.com/alecthomas/log4go"
	"github.com/jmoiron/sqlx"
	"wallet-svc/dto/resp"
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
	log4go.Debug("INSERT INTO `token` result=%d\n", lid)

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

func (dao *TokenDao) CheckContractExists(contract string) bool {
	sql := "SELECT " +
		"count(1) from token WHERE contract=?"
	var count int64
	err := dao.db.Get(&count, sql)

	if err != nil {
		return false
	}

	log4go.Debug("query sql=%s,rows=%d\n", sql, count)

	return count>0
}



func (dao *TokenDao) QueryTokenByAddr(addr string, page, pageSize int32) ([]*model.Asset, error) {
	sql := "SELECT " +
		"t.symbol,f.balance,t.contract,t.logo,t.desc,t.decimals from follow f,token t where f.contract!='BUSD' and f.contract=t.contract and f.wallet = ? limit ?,?"
	var assets []*model.Asset
	er := dao.db.Select(&assets, sql, addr, (page-1)*pageSize, pageSize)
	if er != nil {
		return nil, er
	}
	return assets, nil
}

func (dao *TokenDao) QueryTokenByAddrAndContract(addr, contract string) (*model.Asset, error) {
	sql := "SELECT " +
		"t.symbol,f.balance,t.contract,t.logo,t.desc,t.decimals from follow f,token t where f.contract=t.contract and f.contract=? and f.wallet = ? limit 0,1"
	assets := new(model.Asset)
	er := dao.db.Get(assets, sql, contract, addr)
	if er != nil {
		return nil, er
	}
	return assets, nil
}

func (dao *TokenDao) QueryAllTokens(page, pageSize int32) ([]*resp.AssetInfo, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 5 {
		pageSize = 5
	}
	sql := "SELECT " +
		"t.symbol as 'symbol',t.desc as 'name',t.quantity as 'quantity',t.contract as 'contract',t.logo as 'logo' from token t limit ?,?"
	var assets []*resp.AssetInfo
	er := dao.db.Select(&assets, sql, (page-1)*pageSize, pageSize)
	if er != nil {
		return nil, er
	}
	return assets, nil
}

func (dao *TokenDao) QueryTokenByAddrForExplore(page, pageSize int32, addr string) ([]*resp.AssetInfo, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 5 {
		pageSize = 5
	}
	sql := "SELECT " +
		" t.symbol as 'symbol',t.`desc` as 'name',f.balance as 'quantity',t.contract,t.logo as 'logo' from follow f,token t where f.contract=t.contract and f.wallet = ? limit ?,?"
	var assets []*resp.AssetInfo
	er := dao.db.Select(&assets, sql, addr, (page-1)*pageSize, pageSize)
	if er != nil {
		return nil, er
	}
	return assets, nil
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

	var txs []*model.Token

	er := dao.db.Select(&txs, sql, content, like, (page-1)*pageSize, pageSize)
	if er != nil {
		return nil, er
	}

	var ctrcts []string

	for _, t := range txs {

		ctrcts = append(ctrcts, t.Contract)

	}

	var addrs []string

	if len(ctrcts) > 0 {
		sql2 := "SELECT f.contract from  follow f  WHERE f.wallet = ? and f.followed=true and f.contract in (?)"

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
