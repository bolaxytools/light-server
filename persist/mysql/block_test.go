package mysql

import (
	"fmt"
	"testing"
)

func TestBlockDao_Get(t *testing.T) {
	dao := NewBlockDao()


	tx,e := dao.GetBlockByHeightX(47)
	if e != nil {
		t.Error(e)
		return
	}
	fmt.Printf("success,block=%+v\n!!!",tx)
}
