package mysql

import (
	"fmt"
	"github.com/boxproject/bolaxy/cmd/sdk"
	"testing"
	"time"
	"wallet-svc/model"
)

func TestGoodsDao_AddAsset(t *testing.T) {
	dao := NewTxDao()
	asset := &model.Tx{
		AddrFrom:"bx01",
		AddrTo:"bx01",
		Amount:"34",
		MinerFee:"232",
		TxHash:"txhash001",
		BlockHeight:1,
		TxTime:time.Now().UnixNano()/1e6,
		Memo:"memo001",
		TxType:0,
	}
	er := dao.Add(asset)
	if er != nil {
		t.Error(er)
		return
	}
	fmt.Printf("add asset success \n")
}


func TestGoodsDao_QueryTxs(t *testing.T) {
	dao := NewTxDao()
	asset := "bx0001"
	txs,er := dao.Query(asset,1,5)
	if er != nil {
		t.Error(er)
		return
	}
	for _,tx := range txs {
		fmt.Printf("%+v\n",tx)
	}
}


func TestGoodsDao_BashSave(t *testing.T) {
	dao := NewTxDao()
	asset1 := &sdk.Transaction{
		Hash:"hash20001",
		From:"from201",
		To:"to201",
		Data:[]byte("dt1"),
	}

	asset2 := &sdk.Transaction{
		Hash:"hash20002",
		From:"from202",
		To:"to202",
		Data:[]byte("dt2"),
	}

	txs := []*sdk.Transaction{asset1,asset2}

	er := dao.BashSave(txs,7,time.Now().Unix())
	if er != nil {
		t.Error(er)
		return
	}
	fmt.Printf("add asset success \n")
}

func TestAddressDao_QueryCount(t *testing.T) {
	dao := NewAddressDao()
	c,e := dao.QueryCount()
	if e != nil {
		t.Error(e)
		return
	}
	fmt.Printf("c=%d\n",c)
}


func TestAddressDao_Add(t *testing.T) {
	dao := NewAddressDao()

	addm := &model.Address{
		Addr:"addr10004",
		AddTime:time.Now().UnixNano()/1e6,
		UpdateTime:time.Now().UnixNano()/1e6,
	}

	e := dao.Add(addm)
	if e != nil {
		t.Error(e)
		return
	}
	fmt.Printf("success!!!")
}

func TestAddressDao_Get(t *testing.T) {
	dao := NewTxDao()


	tx,e := dao.GetTxByHash("0x250fb43c0a76d9f8cdbde67c0c97ffa285d9f5622ea7a7d6397c85eecc8a28d3")
	if e != nil {
		t.Error(e)
		return
	}
	fmt.Printf("success,tx=%+v\n!!!",tx)
}
