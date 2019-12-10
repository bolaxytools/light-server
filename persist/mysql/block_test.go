package mysql

import (
	"fmt"
	"testing"
)

func TestBlockDao_Get(t *testing.T) {
	dao := NewBlockDao()


	tx,e := dao.GetBlockByHeight("0")
	if e != nil {
		t.Error(e)
		return
	}
	fmt.Printf("success,block=%+v\n!!!",tx)
}
