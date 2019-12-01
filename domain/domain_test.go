package domain

import (
	"fmt"
	"github.com/syndtr/goleveldb/leveldb/util"
	"testing"
	"wallet-svc/config"
	"wallet-svc/mock"
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

func TestPa(t *testing.T)  {
	flr.pa(0,3)
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
	fmt.Println(len("0X03573555DCF1B4518816DF6F1EDFF2C16B4D29F64CF6D9BDA78ECB6F50826CCA0B"))
}