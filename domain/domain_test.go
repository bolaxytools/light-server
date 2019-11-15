package domain

import (
	"encoding/binary"
	"fmt"
	"testing"
	"wallet-svc/config"
	"wallet-svc/mock"
)

var (
	flr *BlockFollower
)

func TestMain(m *testing.M) {
	config.LoadConfig(mock.Getwd())
	flr = NewBlockFollower("http://192.168.10.189:8080")
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
	k1 := []byte("k1")
	v1 := make([]byte,16)

	binary.LittleEndian.PutUint64(v1,4)


	er := flr.db.Put(k1,v1,nil)
	if er != nil {
		t.Error(er)
		return
	}


	vf,er := flr.db.Get(k1,nil)
	if er != nil {
		t.Error(er)
		return
	}

	p := binary.LittleEndian.Uint64(vf)
	fmt.Printf("get p = %d\n",p)


}

func TestEq(t *testing.T)  {
	fmt.Println(0x1==0x01)
}