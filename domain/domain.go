package domain

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"github.com/alecthomas/log4go"
	"github.com/boxproject/bolaxy/cmd/sdk"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/errors"
	"github.com/syndtr/goleveldb/leveldb/filter"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"math/big"
	"strconv"
	"time"
	"wallet-svc/config"
	"wallet-svc/httpclient"
	"wallet-svc/model"
	"wallet-svc/persist/mysql"
	"wallet-svc/util"
)

var bfr *BlockFollower

type BlockFollower struct {
	requester  *httpclient.Requester
	db         *leveldb.DB
	txDao      *mysql.TxDao
	addressDao *mysql.AddressDao
	blockDao   *mysql.BlockDao
	tokenDao   *mysql.TokenDao
	followDao  *mysql.FollowDao
	chan_txs   chan []*sdk.Transaction
}

func (bf *BlockFollower) GetAddressCount() uint64 {
	c, e := bf.addressDao.QueryCount()
	if e != nil {
		return 0
	}
	return uint64(c)
}

func (bf *BlockFollower) procTxs() {
	for {
		select {
		case txs := <-bf.chan_txs:
			for _, tx := range txs {
				from := tx.From
				to := tx.To
				now := time.Now().UnixNano() / 1e6
				fa := &model.Address{
					Addr:       from,
					AddTime:    now,
					UpdateTime: now,
				}
				er := bf.addressDao.Add(fa)
				log4go.Info("bf.addressDao.Add from error=%v\n", er)
				ta := &model.Address{
					Addr:       to,
					AddTime:    now,
					UpdateTime: now,
				}
				er = bf.addressDao.Add(ta)
				log4go.Info("bf.addressDao.Add to error=%v\n", er)
			}

		}
	}
}

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

	go bfr.procTxs()
	return bfr
}

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

		log4go.Info("已处理到高度：%d，最新区块高度:%d\n", lastDealHei, curHei)
		if lastDealHei < curHei {

			next := lastDealHei + 1

			blk, txs, err := flr.GetBlockTxs(next)
			if err != nil {
				log4go.Info("flr.GetBlockTxs error=%v\n", err)
				er := flr.setDealtBlockHeight(next)
				if er != nil {
					log4go.Info("flr.setDealtBlockHeight to %d\n", next)
				}
				time.Sleep(time.Second * 30)
				continue
			}

			flr.chan_txs <- txs //把爬到的交易扔到channel里去处理

			nowmills := time.Now().UnixNano() / 1e6

			err = flr.txDao.BashSave(txs, next, nowmills)
			if err != nil {
				log4go.Info("flr.txDao.BashSave error=%v\n", err)
				er := flr.setDealtBlockHeight(next)
				if er != nil {
					log4go.Info("flr.setDealtBlockHeight to %d\n", next)
				}
				time.Sleep(time.Second * 10)
				continue
			}

			er := flr.blockDao.Add(blk)
			if er != nil {
				log4go.Info("flr.blockDao.Add error=%v\n", er)
			}

			log4go.Info("process block=%d success\n", next)
			er = flr.setDealtBlockHeight(next)
			if er != nil {
				log4go.Info("flr.setDealtBlockHeight to %d\n", next)
			}

		}

		time.Sleep(time.Second * 10)

	}
}

func (flr *BlockFollower) Translate(txs []*sdk.Transaction) error {
	for _, tx := range txs {
		bigint, _ := big.NewInt(0).SetString(tx.Value, 0)
		if bigint.Int64() == 0 && flr.checkExists(tx.To) { //TODO 判断to的地址是不是在登记的池里，如果在就是智能合约转账
			return flr.remakeTx(tx)
		}
	}
	return nil
}

func (flr *BlockFollower) checkExists(fooAddr string) bool {
	return false
}

func (flr *BlockFollower) remakeTx(tx *sdk.Transaction) error {
	endpoint := fmt.Sprintf("tx/%s", tx.Hash)

	buf, err := flr.requester.RequestHttpByGet(endpoint, nil)
	if err != nil {
		return err
	}

	tmp := new(model.ReceiptLog)

	er := json.Unmarshal(buf, tmp)
	if er != nil {
		return er
	}

	if len(tmp.Topics) != 3 {
		return errors.New("not brc20 transfer")
	}

	tx.To = util.CutAddress(tmp.Topics[2])
	tx.Value = util.DataToValue(tmp.Data)

	return nil
}

func onlyKey() []byte {
	k := []byte{0x01}
	return k
}

func buildWhiteKey(addr string) []byte {
	k := []byte{0x03}
	buf := bytes.NewBuffer(k)
	buf.WriteString(addr)
	return buf.Bytes()
}

func buildBlackKey(addr string) []byte {
	k := []byte{0x02}
	buf := bytes.NewBuffer(k)
	buf.WriteString(addr)
	return buf.Bytes()
}

func (flr *BlockFollower) getDealtBlockHeight() (int64, error) {
	vf, er := flr.db.Get(onlyKey(), nil)
	if er != nil {
		if er == errors.ErrNotFound {
			return -1, nil
		}
		return 0, er
	}

	p := binary.LittleEndian.Uint64(vf)
	return int64(p), nil
}

func (flr *BlockFollower) setDealtBlockHeight(hei int64) error {
	k1 := onlyKey()
	v1 := make([]byte, 16)

	binary.LittleEndian.PutUint64(v1, uint64(hei))

	er := flr.db.Put(k1, v1, nil)
	if er != nil {
		return er
	}
	return nil
}

func (flr *BlockFollower) GetCurrentBlockHeight() (int64, error) {
	cf := new(model.ChainInfo)

	buf, err := flr.requester.RequestHttpByGet("info", nil)
	if err != nil {
		return 0, err
	}

	err = json.Unmarshal(buf, cf)

	if err != nil {
		return 0, err
	}

	hei, er := strconv.ParseInt(cf.LastBlockIndex, 10, 64)
	if er != nil {
		return 0, er
	}

	return hei, nil
}

func (flr *BlockFollower) GetBlockTxs(hei int64) (*model.Block, []*sdk.Transaction, error) {

	endpoint := fmt.Sprintf("block/%d", hei)

	buf, err := flr.requester.RequestHttpByGet(endpoint, nil)
	if err != nil {
		return nil, nil, err
	}

	tmp := struct {
		Index     uint64 `json:"Index"`
		StateHash string `json:"StateHash"`
	}{}

	er := json.Unmarshal(buf, &tmp)
	if er != nil {
		return nil, nil, er
	}

	jb := string(buf)

	txs, err := sdk.GetTransactions(jb)
	if err != nil {
		return nil, nil, err
	}

	blk := &model.Block{
		Height:    tmp.Index,
		Hash:      tmp.StateHash,
		TxCount:   int32(len(txs)),
		BlockTime: time.Now().UnixNano() / 1e6,
	}

	return blk, txs, nil
}

func (flr *BlockFollower) GetAccount(addr string) (*model.Account, error) {

	endpoint := fmt.Sprintf("account/%s", addr)

	buf, err := flr.requester.RequestHttpByGet(endpoint, nil)
	if err != nil {
		return nil, err
	}

	tmp := new(model.Account)

	er := json.Unmarshal(buf, tmp)

	if er != nil {
		return nil, er
	}

	return tmp, nil
}

func (flr *BlockFollower) SendRawTx(reqstr string) (string, error) {

	buf, err := flr.requester.PostString("rawtx", reqstr)
	if err != nil {
		return "", err
	}

	tmp := struct {
		TxHash string `json:"txHash"`
	}{}

	er := json.Unmarshal(buf, &tmp)

	if er != nil {
		return "", er
	}

	return tmp.TxHash, nil
}

/*
	将某地址添加到黑名单列表
*/
func (flr *BlockFollower) AddToBlackList(addr string) error {

	key := buildBlackKey(addr)
	v := []byte{0x01}
	er := flr.db.Put(key, v, nil)
	return er
}

/*
	将某地址添加到白名单列表
*/
func (flr *BlockFollower) AddToWhiteList(addr string) error {

	key := buildWhiteKey(addr)
	v := []byte{0x01}
	er := flr.db.Put(key, v, nil)
	return er
}

/*
	删除黑名单列表
*/
func (flr *BlockFollower) RemBlackList(addr string) error {

	key := buildBlackKey(addr)
	er := flr.db.Delete(key, nil)

	return er
}

/*
	删除白名单列表
*/
func (flr *BlockFollower) RemWhiteList(addr string) error {

	key := buildWhiteKey(addr)
	er := flr.db.Delete(key, nil)
	return er
}

/*
	是否存在于黑名单列表
*/
func (flr *BlockFollower) CheckBlackList(addr string) bool {

	key := buildBlackKey(addr)
	_, er := flr.db.Get(key, nil)

	return er == nil
}

/*
	是否存在于白名单列表
*/
func (flr *BlockFollower) CheckWhiteList(addr string) bool {

	key := buildWhiteKey(addr)
	_, er := flr.db.Get(key, nil)
	return er == nil
}

func (bf *BlockFollower) FollowToken(contract, addr, balance string) error {
	flw := &model.Follow{
		Contract: contract,
		Wallet:   addr,
		Balance:  balance,
	}
	return bf.followDao.Add(flw)
}

func (bf *BlockFollower) SearchToken(content string) ([]*model.Token, error) {
	return bf.followDao.QueryTokenByContract(1, 100, content)
}

func (bf *BlockFollower) QuerySearchTokenCount(content string) (uint64, error) {
	return bf.tokenDao.QueryCountByContent(content)
}
