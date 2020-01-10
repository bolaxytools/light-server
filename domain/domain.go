package domain

import (
	"github.com/alecthomas/log4go"
	"github.com/boxproject/bolaxy/cmd/sdk"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/errors"
	"github.com/syndtr/goleveldb/leveldb/filter"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"time"
	"wallet-svc/config"
	"wallet-svc/dto/resp"
	"wallet-svc/httpclient"
	"wallet-svc/model"
	"wallet-svc/persist/mysql"
)

var (
	bfr *BlockFollower
)

/*
	block follower
*/
type BlockFollower struct {
	requester  *httpclient.Requester
	db         *leveldb.DB
	txDao      *mysql.TxDao
	addressDao *mysql.AddressDao
	blockDao   *mysql.BlockDao
	tokenDao   *mysql.TokenDao
	followDao  *mysql.FollowDao
	chan_txs   chan []*sdk.Transaction
	chan_fa    chan *model.FollowAsset
}

/*
	get the count of addresses
*/
func (bf *BlockFollower) GetAddressCount() uint64 {
	c, e := bf.addressDao.QueryCount()
	if e != nil {
		return 0
	}
	return uint64(c)
}

/*
	process transactions
*/
func (bf *BlockFollower) procTxs() {
	for {
		select {
		case txs := <-bf.chan_txs:
			for _, tx := range txs {
				from := tx.From
				to := tx.To
				now := time.Now().UnixNano() / 1e6
				if len(from) > 0 {
					fa := &model.Address{
						Addr:       from,
						AddTime:    now,
						UpdateTime: now,
					}
					er := bf.addressDao.Add(fa)
					if er != nil {
						log4go.Info("bf.addressDao.Add from error=%v\n", er)
					}
				}
				if len(to) > 0 {
					ta := &model.Address{
						Addr:       to,
						AddTime:    now,
						UpdateTime: now,
					}
					er := bf.addressDao.Add(ta)
					if er != nil {
						log4go.Info("bf.addressDao.Add to error=%v\n", er)
					}
				}
			}

		}
	}
}

/*
	if an address followed a brc coin,we should keep its balance
*/
func (bf *BlockFollower) procFollowAssetBalance() {

	for {
		select {
		case fa := <-bf.chan_fa:

			balance, er := bf.getChildBalance(fa)
			if er != nil {
				log4go.Info(" bf.getChildBalance error=%v\n", er)
				continue
			}

			flw := &model.Follow{
				Contract: fa.Contract,
				Wallet:   fa.Address,
				Balance:  balance,
			}

			err := bf.followDao.Add(flw)
			if err != nil {
				log4go.Info("bf.followDao.Add error=%v\n", err)
			}

		}
	}
}


/*
	instantiate level database
*/
func NewLDBDatabase() (*leveldb.DB, error) {

	// Ensure we have some minimal caching and file guarantees
	var (
		cache   = 16
		handles = 16
	)

	// Open the db and recover any potential corruptions
	db, err := leveldb.OpenFile(config.Cfg.Global.LevelDbFilePath, &opt.Options{
		OpenFilesCacheCapacity: handles,
		BlockCacheCapacity:     cache / 2 * opt.MiB,
		WriteBuffer:            cache / 4 * opt.MiB, // Two of these are used internally
		Filter:                 filter.NewBloomFilter(10),
	})

	if _, corrupted := err.(*errors.ErrCorrupted); corrupted {
		db, err = leveldb.RecoverFile(config.Cfg.Global.LevelDbFilePath, nil)
	}

	if err != nil {
		return nil, err
	}

	return db, nil
}

/*
instantiate Block Follower
*/
func NewBlockFollower() *BlockFollower {

	if bfr != nil {
		return bfr
	}

	tmp, er := NewLDBDatabase()
	if er != nil {
		panic(er)
	}
	bfr = &BlockFollower{
		requester: &httpclient.Requester{
			BaseUrl: config.Cfg.Global.BolaxyNodeUrl,
		},
		db:         tmp,
		txDao:      mysql.NewTxDao(),
		addressDao: mysql.NewAddressDao(),
		blockDao:   mysql.NewBlockDao(),
		tokenDao:   mysql.NewTokenDao(),
		followDao:  mysql.NewFollowDao(),
	}
	bfr.chan_txs = make(chan []*sdk.Transaction, 1024)
	bfr.chan_fa = make(chan *model.FollowAsset, 1024)

	go bfr.procTxs()
	go bfr.procFollowAssetBalance()
	return bfr
}

/*
	watch the block chain,follow it when it raise
*/
func (flr *BlockFollower) FollowBlockChain() {
	for {
		lastDealHei, er := flr.getDealtBlockHeight()
		if er != nil {
			log4go.Info("getDealtBlockHeight error=%v\n", er)
			time.Sleep(time.Second * 30)
			continue
		}

		curHei, er := flr.GetCurrentBlockHeight()

		if er != nil {
			log4go.Info("flr.GetCurrentBlockHeight error=%v\n", er)
			time.Sleep(time.Second * 60)
			continue
		}

		flr.pa(lastDealHei, curHei)

		time.Sleep(time.Second * 2)

	}
}

/*
	follow block chain by step
	*/
func (flr *BlockFollower) pa(lh, ch int64) bool {
	log4go.Info("已处理到高度：%d，最新区块高度:%d\n", lh, ch)
	if lh < ch {

		next := lh + 1

		blk, signerMap, txs, err := flr.GetBlockTxs(next)
		if err != nil {
			log4go.Info("flr.GetBlockTxs error=%v\n", err)
			er := flr.setDealtBlockHeight(next)
			if er != nil {
				log4go.Info("flr.setDealtBlockHeight to %d\n", next)
			}
			time.Sleep(time.Second * 30)
			return false
		}

		flr.chan_txs <- txs //把爬到的交易扔到channel里去处理

		nowmills := time.Now().UnixNano() / 1e6

		if lh == 4 || lh == 5 {
			log4go.Info("a")
		}

		/******begin******/
		mdxs, er := flr.translate(txs, uint64(next), nowmills)
		if er != nil {
			log4go.Info("flr.translate error=%v\n", er)
		}

		err = flr.txDao.BatchAddTx(mdxs)
		if err != nil {
			log4go.Info("flr.txDao.BatchAddTx error=%v,blockHeight=%d,txs.len=%d\n", err, lh, len(txs))
			er := flr.setDealtBlockHeight(next)
			if er != nil {
				log4go.Info("flr.setDealtBlockHeight to %d\n", next)
			}
		}
		/******end******/

		er = flr.blockDao.Add(blk)
		if er != nil {
			log4go.Info("flr.blockDao.Add error=%v\n", er)
		} else {
			er := flr.blockDao.BatchAddSingers(signerMap, blk.Height)
			if er != nil {
				log4go.Info("flr.blockDao.BatchAddSingers error:%v\n", er)
			}
		}

		log4go.Info("process block=%d success\n", next)
		er = flr.setDealtBlockHeight(next)
		if er != nil {
			log4go.Info("flr.setDealtBlockHeight to %d\n", next)
		}

	} else {
		time.Sleep(time.Second * 10)
	}
	return true
}

/*
	check weather a contract address exists in the db
*/
func (flr *BlockFollower) checkExists(fooAddr string) bool {
	return flr.tokenDao.CheckContractExists(fooAddr)
}

/*
	follow a token,add it to wallet
*/
func (bf *BlockFollower) FollowToken(contract, addr, balance string) error {
	flw := &model.Follow{
		Contract: contract,
		Wallet:   addr,
		Balance:  balance,
		Followed: true,
	}

	fa := &model.FollowAsset{
		Contract: contract,
		Address:  addr,
	}

	bf.chan_fa <- fa

	return bf.followDao.Add(flw)
}

/*
	search a token,the content which can be a contract address or symbol
*/
func (bf *BlockFollower) SearchToken(content, addr string) ([]*model.Token, error) {
	return bf.followDao.QueryTokenByContract(1, 100, content, addr)
}

/*
	list checked-in tokens info
*/
func (bf *BlockFollower) QueryAddrAssets(page, pageSize int32, addr string) ([]*model.Asset, error) {
	return bf.tokenDao.QueryTokenByAddr(addr, page, pageSize)
}

/*
	query one's specific contract token info
*/
func (bf *BlockFollower) QueryAddrContractAsset(contract, addr string) (*model.Asset, error) {
	return bf.tokenDao.QueryTokenByAddrAndContract(addr, contract)
}

/*
	list all checked-in tokens
*/
func (bf *BlockFollower) QueryAllTokens(page, pageSize int32) ([]*resp.AssetInfo, error) {
	return bf.tokenDao.QueryAllTokens(page, pageSize)
}

/*
	query account's tokens
*/
func (bf *BlockFollower) QueryAccountTokens(page, pageSize int32, addr string) ([]*resp.AssetInfo, error) {
	return bf.tokenDao.QueryTokenByAddrForExplore(page, pageSize, addr)
}

/*
	count token total
*/
func (bf *BlockFollower) QueryTokenCount() (int64, error) {
	return bf.tokenDao.QueryCount()
}

/*
	search token by content which could be symbol or contract address
*/
func (bf *BlockFollower) QuerySearchTokenCount(content string) (uint64, error) {
	return bf.tokenDao.QueryCountByContent(content)
}
