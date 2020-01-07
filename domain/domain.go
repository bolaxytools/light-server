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
	"wallet-svc/dto/resp"
	"wallet-svc/httpclient"
	"wallet-svc/model"
	"wallet-svc/persist/mysql"
	"wallet-svc/util"
)

var (
	bfr *BlockFollower
)

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
	bfr.chan_fa = make(chan *model.FollowAsset, 1024)

	go bfr.procTxs()
	go bfr.procFollowAssetBalance()
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

		flr.pa(lastDealHei, curHei)

		time.Sleep(time.Second * 2)

	}
}

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

func (flr *BlockFollower) translate(txs []*sdk.Transaction, height uint64, txTime int64) ([]*model.Tx, error) {
	var rt []*model.Tx
	for _, tx := range txs {
		//bigint, _ := big.NewInt(0).SetString(tx.Value, 0)
		//if bigint.Int64() == 0 {
		tkx, e := flr.remakeTx(tx, height, txTime)
		if e != nil {
			log4go.Info("flr.remakeTx error=%v\n", e)
			continue
		}
		rt = append(rt, tkx)
		//}
	}
	return rt, nil
}

func (flr *BlockFollower) checkExists(fooAddr string) bool {
	return flr.tokenDao.CheckContractExists(fooAddr)
}

func (flr *BlockFollower) remakeTx(tx *sdk.Transaction, height uint64, txTime int64) (*model.Tx, error) {
	endpoint := fmt.Sprintf("tx/%s", tx.Hash)

	buf, err := flr.requester.RequestHttpByGet(endpoint, nil)
	if err != nil {
		return nil, err
	}

	type Tmp struct {
		Data *model.TxReceipt
		Err  string
	}

	tmp := new(Tmp)

	er := json.Unmarshal(buf, tmp)
	if er != nil {
		return nil, er
	}

	_, c, sbf := sdk.UnWrapData(tx.Data)

	datastr := ""
	if c == sdk.DataDisplay {
		datastr = string(sbf)
	}

	if len(datastr) > 100 {
		datastr = datastr[0:100]
	}
	mtx := &model.Tx{
		TxType:      0,
		AddrFrom:    tx.From,
		AddrTo:      tx.To,
		Amount:      tx.Value,
		MinerFee:    fmt.Sprintf("%d", tx.Gas),
		TxHash:      tx.Hash,
		BlockHeight: height,
		TxTime:      txTime,
		Memo:        datastr,
		Status:      tmp.Data.Status,
	}

	if tmp.Data.ContractAddress != "0x0000000000000000000000000000000000000000" {
		mtx.TxType = 1
	}

	gasPriceBig, _ := big.NewInt(0).SetString(tx.GasPrice, 0)
	gused := big.NewInt(int64(tmp.Data.GasUsed))

	realCost := gasPriceBig.Mul(gasPriceBig, gused)

	flr.increaseGasTotal(realCost.Uint64())

	mtx.MinerFee = realCost.String()

	if len(tmp.Data.Logs) > 0 && len(tmp.Data.Logs[0].Topics) == 3 { //如果是智能合约的转账就把目的地址和交易额提出来
		mtx.Contract = tmp.Data.Logs[0].Address
		mtx.AddrTo = util.CutAddress(tmp.Data.Logs[0].Topics[2])
		mtx.Amount = util.DataToValue(tmp.Data.Logs[0].Data)

		fa_fr := &model.FollowAsset{
			Contract: mtx.Contract,
			Address:  mtx.AddrFrom,
		}

		flr.chan_fa <- fa_fr

		fa_to := &model.FollowAsset{
			Contract: mtx.Contract,
			Address:  mtx.AddrTo,
		}
		flr.chan_fa <- fa_to

	} else {
		mtx.Contract = "BUSD"
	}
	return mtx, nil
}

func onlyKey() []byte {
	k := []byte{0x01}
	return k
}

func gasKey() []byte {
	k := []byte{0x04}
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

	type TMP struct {
		Data *model.ChainInfo
		Err  string
	}

	cf := new(TMP)

	buf, err := flr.requester.RequestHttpByGet("info", nil)
	if err != nil {
		return 0, err
	}

	//fmt.Printf("fuckbuf=%s\n",buf)

	err = json.Unmarshal(buf, cf)

	if err != nil {
		return 0, err
	}

	if cf.Err != "" {
		return 0, errors.New(cf.Err)
	}

	hei, er := strconv.ParseInt(cf.Data.LastBlockIndex, 10, 64)
	if er != nil {
		return 0, er
	}

	return hei, nil
}

func (flr *BlockFollower) GetTotalBUSD() (string, error) {

	cf := struct {
		Data struct {
			Data string `json:"data"`
		}
		Err string
	}{}

	buf, err := flr.requester.RequestHttpByGet("totalbusd", nil)
	if err != nil {
		return "0", err
	}

	err = json.Unmarshal(buf, &cf)

	if err != nil {
		return "0", err
	}

	if cf.Err != "" {
		return "", errors.New(cf.Err)
	}

	return cf.Data.Data, nil
}

func (flr *BlockFollower) GetBlockTxs(hei int64) (*model.Block, map[string]string, []*sdk.Transaction, error) {

	endpoint := fmt.Sprintf("block/%d", hei)

	buf, err := flr.requester.RequestHttpByGet(endpoint, nil)
	if err != nil {
		return nil, nil, nil, err
	}

	//fmt.Printf("return buf=%s\n", buf)

	type bk struct {
		Index     uint64 `json:"Index"`
		StateHash string `json:"StateHash"`
	}

	type bd struct {
		Body bk `json:"body"`
	}

	tmp := new(model.AlphaPing)

	er := json.Unmarshal(buf, tmp)
	if er != nil {
		return nil, nil, nil, er
	}

	//jb := string(buf)
	jbuf, e := json.Marshal(tmp.Data)
	if e != nil {
		log4go.Info("json.Marshal(tmp.Data) error=%v\n", e)
		return nil, nil, nil, er
	}
	jb := string(jbuf)
	//fmt.Printf("to=%s\n", jb)

	txs, signerMap, err := sdk.GetTransactions(jb)
	if err != nil {
		return nil, nil, nil, err
	}

	//mdmap := tmp.Data.(map[string]interface{})

	pmt := new(bd)

	er = json.Unmarshal(jbuf, pmt)
	if er != nil {
		return nil, nil, nil, er
	}

	blk := &model.Block{
		Height:    pmt.Body.Index,
		Hash:      pmt.Body.StateHash,
		TxCount:   int32(len(txs)),
		BlockTime: time.Now().UnixNano() / 1e6,
	}

	return blk, signerMap, txs, nil
}

func (flr *BlockFollower) GetAccount(addr string) (*model.Account, error) {

	endpoint := fmt.Sprintf("account/%s", addr)

	buf, err := flr.requester.RequestHttpByGet(endpoint, nil)
	if err != nil {
		return nil, err
	}

	type TMP struct {
		Data *model.Account
		Err  string
	}

	tmp := new(TMP)

	er := json.Unmarshal(buf, tmp)

	if er != nil {
		return nil, er
	}

	if tmp.Err != "" {
		return nil, errors.New(tmp.Err)
	}

	return tmp.Data, nil
}

func (flr *BlockFollower) SendRawTx(reqstr string) (string, error) {

	buf, err := flr.requester.PostString("rawtx", reqstr)
	if err != nil {
		return "", err
	}

	tmp := struct {
		Data struct {
			TxHash string `json:"txHash"`
		}
		Err string
	}{}

	er := json.Unmarshal(buf, &tmp)

	if er != nil {
		return "", er
	}

	if tmp.Err != "" {
		log4go.Info("send raw tx =%s\n", tmp.Err)
		return "", errors.New(tmp.Err)
	}

	return tmp.Data.TxHash, nil
}

func (flr *BlockFollower) getChildBalance(fa *model.FollowAsset) (string, error) {

	buf, err := flr.requester.PostJson("balanceof", fa)
	if err != nil {
		return "", err
	}

	fmt.Printf("postresult=%s\n", buf)

	tmp := struct {
		Data struct {
			Data string `json:"data"`
		}
		Err string
	}{}

	er := json.Unmarshal(buf, &tmp)

	if er != nil {
		return "", er
	}

	if tmp.Err != "" {
		return "", errors.New(tmp.Err)
	}

	return tmp.Data.Data, nil
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
		Followed: true,
	}

	fa := &model.FollowAsset{
		Contract: contract,
		Address:  addr,
	}

	bf.chan_fa <- fa

	return bf.followDao.Add(flw)
}

func (bf *BlockFollower) SearchToken(content, addr string) ([]*model.Token, error) {
	return bf.followDao.QueryTokenByContract(1, 100, content, addr)
}

func (bf *BlockFollower)  QueryAddrAssets(page, pageSize int32, addr string) ([]*model.Asset, error) {
	return bf.tokenDao.QueryTokenByAddr(addr, page, pageSize)
}

func (bf *BlockFollower) QueryAddrContractAsset(contract, addr string) (*model.Asset, error) {
	return bf.tokenDao.QueryTokenByAddrAndContract(addr, contract)
}

func (bf *BlockFollower) QueryAllTokens(page, pageSize int32) ([]*resp.AssetInfo, error) {
	return bf.tokenDao.QueryAllTokens(page, pageSize)
}

func (bf *BlockFollower) QueryAccountTokens(page, pageSize int32, addr string) ([]*resp.AssetInfo, error) {
	return bf.tokenDao.QueryTokenByAddrForExplore(page, pageSize, addr)
}

func (bf *BlockFollower) QueryTokenCount() (int64, error) {
	return bf.tokenDao.QueryCount()
}

func (bf *BlockFollower) QuerySearchTokenCount(content string) (uint64, error) {
	return bf.tokenDao.QueryCountByContent(content)
}

func (flr *BlockFollower) increaseGasTotal(gas uint64) {
	total, e := flr.GeTotalGasCost()
	if e != nil {
		log4go.Info("flr.GeTotalGasCost error=%v\n", e)
		return
	}
	if total == -1 {
		total = int64(gas)
	} else {
		total = total + int64(gas)
	}

	e = flr.SetTotalGasCost(uint64(total))
	if e != nil {
		log4go.Info("flr.SetTotalGasCost error=%v\n", e)
	}

}

func (flr *BlockFollower) SetTotalGasCost(gas uint64) error {
	k1 := gasKey()
	v1 := make([]byte, 16)

	binary.LittleEndian.PutUint64(v1, gas)

	er := flr.db.Put(k1, v1, nil)
	if er != nil {
		return er
	}
	return nil
}

func (flr *BlockFollower) GeTotalGasCost() (int64, error) {
	vf, er := flr.db.Get(gasKey(), nil)
	if er != nil {
		if er == errors.ErrNotFound {
			return -1, nil
		}
		return 0, er
	}

	p := binary.LittleEndian.Uint64(vf)
	return int64(p), nil
}
