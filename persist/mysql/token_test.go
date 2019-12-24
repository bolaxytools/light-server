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