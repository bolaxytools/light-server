package domain

import (
	"bytes"
	"encoding/binary"
	"github.com/alecthomas/log4go"
	"github.com/syndtr/goleveldb/leveldb/errors"
)

/*
	add an address to whitelist
*/
func (flr *BlockFollower) AddToWhiteList(addr string) error {

	key := buildWhiteKey(addr)
	v := []byte{0x01}
	er := flr.db.Put(key, v, nil)
	return er
}

/*
	delete an address from blacklist
*/
func (flr *BlockFollower) RemBlackList(addr string) error {

	key := buildBlackKey(addr)
	er := flr.db.Delete(key, nil)

	return er
}

/*
	删除白名单列表
*/
func (flr *BlockFollower) RemWhiteList(addr string) error {

	key := buildWhiteKey(addr)
	er := flr.db.Delete(key, nil)
	return er
}

/*
	是否存在于黑名单列表
*/
func (flr *BlockFollower) CheckBlackList(addr string) bool {

	key := buildBlackKey(addr)
	_, er := flr.db.Get(key, nil)

	return er == nil
}

/*
	是否存在于白名单列表
*/
func (flr *BlockFollower) CheckWhiteList(addr string) bool {

	key := buildWhiteKey(addr)
	_, er := flr.db.Get(key, nil)
	return er == nil
}

func (flr *BlockFollower) SetTotalGasCost(gas uint64) error {
	k1 := gasKey()
	v1 := make([]byte, 16)

	binary.LittleEndian.PutUint64(v1, gas)

	er := flr.db.Put(k1, v1, nil)
	if er != nil {
		return er
	}
	return nil
}

func (flr *BlockFollower) GeTotalGasCost() (int64, error) {
	vf, er := flr.db.Get(gasKey(), nil)
	if er != nil {
		if er == errors.ErrNotFound {
			return -1, nil
		}
		return 0, er
	}

	p := binary.LittleEndian.Uint64(vf)
	return int64(p), nil
}

func (flr *BlockFollower) increaseGasTotal(gas uint64) {
	total, e := flr.GeTotalGasCost()
	if e != nil {
		log4go.Info("flr.GeTotalGasCost error=%v\n", e)
		return
	}
	if total == -1 {
		total = int64(gas)
	} else {
		total = total + int64(gas)
	}

	e = flr.SetTotalGasCost(uint64(total))
	if e != nil {
		log4go.Info("flr.SetTotalGasCost error=%v\n", e)
	}
}

/*
	将某地址添加到黑名单列表
*/
func (flr *BlockFollower) AddToBlackList(addr string) error {

	key := buildBlackKey(addr)
	v := []byte{0x01}
	er := flr.db.Put(key, v, nil)
	return er
}

func onlyKey() []byte {
	k := []byte{0x01}
	return k
}

func gasKey() []byte {
	k := []byte{0x04}
	return k
}

func buildWhiteKey(addr string) []byte {
	k := []byte{0x03}
	buf := bytes.NewBuffer(k)
	buf.WriteString(addr)
	return buf.Bytes()
}

func buildBlackKey(addr string) []byte {
	k := []byte{0x02}
	buf := bytes.NewBuffer(k)
	buf.WriteString(addr)
	return buf.Bytes()
}

func (flr *BlockFollower) getDealtBlockHeight() (int64, error) {
	vf, er := flr.db.Get(onlyKey(), nil)
	if er != nil {
		if er == errors.ErrNotFound {
			return -1, nil
		}
		return 0, er
	}

	p := binary.LittleEndian.Uint64(vf)
	return int64(p), nil
}

func (flr *BlockFollower) setDealtBlockHeight(hei int64) error {
	k1 := onlyKey()
	v1 := make([]byte, 16)

	binary.LittleEndian.PutUint64(v1, uint64(hei))

	er := flr.db.Put(k1, v1, nil)
	if er != nil {
		return er
	}
	return nil
}