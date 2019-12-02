package mysql

import (
	"fmt"
	"testing"
	"wallet-svc/model"
)

func TestNewTokenDao(t *testing.T) {
	dao := NewTokenDao()
	tkn := &model.Token{
		Contract:"contract10001",
		Symbol:"sbl10001",
		Logo:"http://wwww.imgage.com/b.png",

	}
	e := dao.Add(tkn)
	if e != nil {
		t.Error(e)
		return
	}
	fmt.Printf("success\n")

}

func TestNewSliceScan(t *testing.T) {
	dao := NewFollowDao()

	tkns,e := dao.QueryTokenByContract(1,100,"smb","0x9c7D0b2F633C7896f07B64f1F5fe71E748169Bf4")
	if e != nil {
		t.Error(e)
		return
	}
	fmt.Printf("tkns=%v\n",tkns)
	fmt.Printf("success\n")

}

