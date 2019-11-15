package domain

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"github.com/alecthomas/log4go"
	"github.com/boxproject/bolaxy/cmd/sdk"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/errors"
	"github.com/syndtr/goleveldb/leveldb/filter"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"strconv"
	"time"
	"wallet-svc/config"
	"wallet-svc/httpclient"
	"wallet-svc/model"
	"wallet-svc/persist/mysql"
)

var bfr *BlockFollower

type BlockFollower struct {
	requester *httpclient.Requester
	db        *leveldb.DB
	txDao     *mysql.TxDao
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

		db:    tmp,
		txDao: mysql.NewTxDao(),
	}

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

		log4go.Info("已处理到高度：%d，最新区块高度:%d\n",lastDealHei,curHei)
		if lastDealHei < curHei {

			next := lastDealHei + 1

			txs, err := flr.GetBlockTxs(next)
			if err != nil {
				log4go.Info("flr.GetBlockTxs error=%v\n", err)
				er := flr.setDealtBlockHeight(next)
				if er != nil {
					log4go.Info("flr.setDealtBlockHeight to %d\n", next)
				}
				time.Sleep(time.Second*30)
				continue
			}

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

			log4go.Info("process block=%d success\n", next)
			er := flr.setDealtBlockHeight(next)
			if er != nil {
				log4go.Info("flr.setDealtBlockHeight to %d\n", next)
			}

		}

		time.Sleep(time.Second*10)

	}
}

func onlyKey() []byte {
	k := []byte{0x01}
	return k
}

func (flr *BlockFollower) getDealtBlockHeight() (int64, error) {
	vf, er := flr.db.Get(onlyKey(), nil)
	if er != nil {
		if er == errors.ErrNotFound {
			return -1,nil
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

func (flr *BlockFollower) GetBlockTxs(hei int64) ([]*sdk.Transaction, error) {

	endpoint := fmt.Sprintf("block/%d", hei)

	buf, err := flr.requester.RequestHttpByGet(endpoint, nil)
	if err != nil {
		return nil, err
	}

	jb := string(buf)

	txs, err := sdk.GetTransactions(jb)
	if err != nil {
		return nil, err
	}

	return txs, nil
}
//0xcfce91565f523b0625677d8ba872edc35d1550f6bff7bf17f0c602a19d271a41
//0x7a4ddce8ac9be67627c2582fde3fbdd61f44b31ecf42cdb7513e9539322e9c91

func (flr *BlockFollower) GetNonce(addr string) (uint64, error) {

	endpoint := fmt.Sprintf("account/%s", addr)

	buf, err := flr.requester.RequestHttpByGet(endpoint, nil)
	if err != nil {
		return 0, err
	}

	tmp := struct {
		Nonce uint64 `json:"nonce"`
	}{}

	er := json.Unmarshal(buf,&tmp)

	if er != nil {
		return 0,er
	}

	return tmp.Nonce, nil
}