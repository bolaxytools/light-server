package mysql

import (
	"testing"
	"wallet-svc/model"
)

func TestFollowDao_Add(t *testing.T) {
	fd := NewFollowDao()
	flo := &model.Follow{
		Contract:"contract10001",
		Wallet:"address10001",
		Balance:"13000",
	}
	e := fd.Add(flo)
	if e != nil {
		t.Error(e)
		return
	}

}
