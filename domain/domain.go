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

func (flr *BlockFollower) SendRawTx(reqstr string) (string, error) {


	buf, err := flr.requester.PostString("rawtx", reqstr)
	if err != nil {
		return "", err
	}

	tmp := struct {
		TxHash string `json:"txHash"`
	}{}

	er := json.Unmarshal(buf,&tmp)

	if er != nil {
		return "",er
	}

	return tmp.TxHash, nil
}

/*
	将某地址添加到黑名单列表
*/
func (flr *BlockFollower) AddToBlackList(addr string) error {

	key := buildBlackKey(addr)
	v := []byte{0x01}
	er :=flr.db.Put(key,v,nil)
	return er
}


/*
	将某地址添加到白名单列表
*/
func (flr *BlockFollower) AddToWhiteList(addr string) error {

	key := buildWhiteKey(addr)
	v := []byte{0x01}
	er := flr.db.Put(key,v,nil)
	return er
}

/*
	删除黑名单列表
*/
func (flr *BlockFollower) RemBlackList(addr string) error {

	key := buildBlackKey(addr)
	er := flr.db.Delete(key,nil)

	return er
}


/*
	删除白名单列表
*/
func (flr *BlockFollower) RemWhiteList(addr string) error {

	key := buildWhiteKey(addr)
	er := flr.db.Delete(key,nil)
	return er
}


/*
	是否存在于黑名单列表
*/
func (flr *BlockFollower) CheckBlackList(addr string) bool {

	key := buildBlackKey(addr)
	_,er := flr.db.Get(key,nil)

	return er==nil
}


/*
	是否存在于白名单列表
*/
func (flr *BlockFollower) CheckWhiteList(addr string) bool {

	key := buildWhiteKey(addr)
	_,er := flr.db.Get(key,nil)
	return er==nil
}
