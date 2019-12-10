package domain

import (
	"fmt"
	"github.com/syndtr/goleveldb/leveldb/util"
	"testing"
	"wallet-svc/config"
	"wallet-svc/mock"
)

var (
	flr *BlockFollower
)

func TestMain(m *testing.M) {
	config.LoadConfig(mock.Getwd())
	flr = NewBlockFollower()
	m.Run()
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
	fmt.Println(0x1==0x01)
}