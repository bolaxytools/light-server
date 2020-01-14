package domain

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/alecthomas/log4go"
	sdk "github.com/bolaxytools/tool-sdk"
	"math/big"
	"strconv"
	"time"
	"wallet-svc/model"
	"wallet-svc/util"
)

/*
	get transactions from point height
*/
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

/*
	get account info
*/
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

/*
	publish raw transaction to block chain
*/
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

/*
	get one's balance of brc20 in a league
*/
func (flr *BlockFollower) getChildBalance(fa *model.FollowAsset) (string, error) {

	buf, err := flr.requester.PostJson("balanceof", fa)
	if err != nil {
		return "", err
	}

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
	get total of main coin in a league
*/
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

/*
	translate sdk.transactions to model.transactions
*/
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

/*
	rebuild a transaction to model.Tx
*/
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

/*
	get current block height
*/
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
