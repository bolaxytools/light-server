package mysql

import (
	"fmt"
	"testing"
	"time"
	"wallet-service/model"
)

func TestGoodsDao_AddAsset(t *testing.T) {
	dao := NewAssetsDao()
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
	dao := NewAssetsDao()
	asset := "bx0001"
	txs,er := dao.Query(asset)
	if er != nil {
		t.Error(er)
		return
	}
	for _,tx := range txs {
		fmt.Printf("%+v\n",tx)
	}
}
