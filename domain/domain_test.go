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

func TestQueryAddrContractAsset(t *testing.T) {
	addr := "0x9c7D0b2F633C7896f07b64f1F5FE71e748169bF5"
	contract := "0xA79D70c4a0B3A31043541Cd593828170Bf037aFE"
	str,er := flr.QueryAddrContractAsset(contract,addr)

	if er != nil {
		t.Error(er)
		return
	}
	fmt.Printf("str=%v\n",str)


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
	flr.pa(62,64)
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
	jbstr := `{"Body":{"Index":63,"RoundReceived":188,"StateHash":"e+aEixuIPI+Hq0B+heKczfp9tDH4cYFzAaFLbcxE29s=","PeersHash":"WFujRruhyw8nX2VNXnAJM+E348cs8HKauJo35embruw=","Transactions":["+JEIhhtI61fgAIJ1MJTk1KGebBhXIlRBisEL4WKRZBkJcYgbwW1nTsgAAICAgKAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAEmoKZ+wg1ErQxxWJAx2PXrL6TgqHE87rdHyg7N9vXHG819oH2ui+ZL0E/4FrvDm6uW1mrzFUizr7NHFnALrMCS8FPv"],"InternalTransactions":[],"InternalTransactionReceipts":[]},"Signatures":{"0X026FFA1CC91EDDE209B33F013E341DD855B8E7458A140345AA3B812085613228BE":"0xf39e5c71022ba97fcc9b09d4769099bdb0df5ac15469e26f4ef8fd82156fb3137ef51684a7f7f07215244d0377e1cc6e6b8215c47a163f8e329a906a289f651c00","0X028269510925EFD170A03C07716B7C17392F382A789A53931FD1A6A7A67FAFF2A3":"0xbf8617b989dcfa4b0becbb46a24a0179a91d9e7b1da68e2a269fa7c15c90918459d78592f2a9977603b2e655182a021de41805d92f7bc9650d5cfa67f122738500","0X03AD9D8CFFC3E626DD3DB92853CB86E133C73037944E48911C63353B0F37AC11A2":"0xfff2883b8e83f94defedc3938c836197d7326c33d7a29357a0b764fcc5dfb1b15c734610f2a22a48925e7f122ed8ab6f033bd381fb1aff5ea8e368e5c615bba500"}}`
	txs,_,_ := sdk.GetTransactions(jbstr)
	fmt.Printf("txs=%v\n",txs)
}