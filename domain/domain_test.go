package domain

import (
	"fmt"
	"github.com/boxproject/bolaxy/cmd/sdk"
	"github.com/syndtr/goleveldb/leveldb/util"
	"testing"
	"wallet-svc/config"
	"wallet-svc/mock"
	"wallet-svc/model"
	"wallet-svc/persist/mysql"
)

var (
	flr *BlockFollower
)

func TestMain(m *testing.M) {
	config.LoadConfig(mock.Getwd())
	mysql.InitMySQL()
	flr = NewBlockFollower()
	m.Run()
}

func TestBalanceOf(t *testing.T) {
	addr := "0x343e7F717d268903F6033766cE1D210a3D82C097"
	contract := "0xa79d70c4a0b3a31043541cd593828170bf037afe"
	str,er := flr.getChildBalance(&model.FollowAsset{
		Address:addr,
		Contract:contract,
	})

	if er != nil {
		t.Error(er)
		return
	}
	fmt.Printf("str=%s\n",str)


}

func TestSendTx(t *testing.T) {
	str,er := flr.SendRawTx("0xf89180861b48eb57e00082753094d586f735c97545927d2a49ca8accd950bf24d53e880de0b6b3a7640000808080a000000000000000000000000000000000000000000000000000000000000000000126a0cb546564dc6ab695ecaab055a0b63f01af41ccd596138865c3d1ddd436fa97c6a072660deac8a47c1d70c1024b8d51adecb318df4baeb2e542965860250cbb232b")
	if er != nil {
		t.Error(er)
		return
	}
	fmt.Printf("str=%s\n",str)
}

func TestPa(t *testing.T)  {
	flr.pa(3,5)
}

func TestBlockFollower_GetCurrentBlockHeight(t *testing.T) {
	hei,er := flr.GetCurrentBlockHeight()
	if er != nil {
		t.Error(er)
		return
	}
	fmt.Printf("hei=%d\n",hei)

}

func TestLvdb(t *testing.T) {
	k1 := []byte("address.1")
	v1 := []byte{0x1}

	er := flr.db.Put(k1,v1,nil)
	if er != nil {
		t.Error(er)
		return
	}

	k2 := []byte("address.2")

	er = flr.db.Put(k2,v1,nil)
	if er != nil {
		t.Error(er)
		return
	}


	si,er := flr.db.SizeOf([]util.Range{
		*util.BytesPrefix([]byte("address")),
	})
	if er != nil {
		t.Error(er)
		return
	}
	fmt.Printf("si=%v\n",si[0])


}

func TestEq(t *testing.T)  {
	jbstr := `{"Body":{"Index":60,"RoundReceived":179,"StateHash":"T0A2Xm7zXNYo9/fZoNBvRZt0yFZKE7Z1PUA0KWQvaNk=","PeersHash":"WFujRruhyw8nX2VNXnAJM+E348cs8HKauJo35embruw=","Transactions":["+M4zheeqnx4Agw9CQJSnnXDEoLOjEENUHNWTgoFwvwN6/oC4RKkFnLsAAAAAAAAAAAAAAACcfQsvYzx4lvB7ZPH1/nHnSBab9QAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAIrHIwSJ6AAAgICgAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAACAJqByHkh1LoI5yAkaJOTOaOsugvBzdlV5uqhKtJIWWB21TqAE/q4dmRGyVCEbI3heCfm5P+Py0hROT0ROrfYuuQuAWw=="],"InternalTransactions":[],"InternalTransactionReceipts":[]},"Signatures":{"0X026FFA1CC91EDDE209B33F013E341DD855B8E7458A140345AA3B812085613228BE":"0x83165d5d5f916fabb55521b8c1849423de790ef316127f88a0d79f76d092e4a72ec62a0af55e62d89fac46a8687c5812992c0a57b2e729ba812e2e092c956c9f01","0X028269510925EFD170A03C07716B7C17392F382A789A53931FD1A6A7A67FAFF2A3":"0xbb077c8a5d80d32f5692833677d7be8de54c0c42fcd37ae2c2aad5477174d1d3035ed2ed15347c49e7f0d3b204077db8866c9840d1af542f1754d23942933d0e00","0X03AD9D8CFFC3E626DD3DB92853CB86E133C73037944E48911C63353B0F37AC11A2":"0xb0d646ecdf56a50598ff6259b27c35f2d1fb6e569335175b463d0a57098649ba38b8d7fdf17c3ffd6f0b4d23ba0affd1b8ec1dcb0db075d5bba98eac2eade82700"}}`
	txs,_,_ := sdk.GetTransactions(jbstr)
	fmt.Printf("txs=%v\n",txs)
}