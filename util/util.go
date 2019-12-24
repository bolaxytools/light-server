package util

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"regexp"
)


const HashLength = 32

func IsNumber(content string) bool  {
	pattern := "\\d+" //反斜杠要转义
	result,e := regexp.MatchString(pattern,content)
	if e != nil {
		return false
	}
	return result
}

func CutAddress(longAddr string) string {
	//0x0000000000000000000000003bd6361959306b1b50797d3ff82b9a43541c3e47
	//0x3bd6361959306b1b50797d3ff82b9a43541c3e47
	return fmt.Sprintf("0x%s",longAddr[26:])
}

func DataToValue(data string) string {
	hash := HexToHash(data)
	bt := new(big.Int).SetBytes(hash[:])
	return bt.String()
}

func HexToHash(s string) Hash { return BytesToHash(FromHex(s)) }

func Hex2Bytes(str string) []byte {
	h, _ := hex.DecodeString(str)
	return h
}

func FromHex(s string) []byte {
	if len(s) > 1 {
		if s[0:2] == "0x" || s[0:2] == "0X" {
			s = s[2:]
		}
	}
	if len(s)%2 == 1 {
		s = "0" + s
	}
	return Hex2Bytes(s)
}


func (h Hash) Big() *big.Int { return new(big.Int).SetBytes(h[:]) }

func BytesToHash(b []byte) Hash {
	var h Hash
	h.SetBytes(b)
	return h
}

type Hash [HashLength]byte
func (h *Hash) SetBytes(b []byte) {
	if len(b) > len(h) {
		b = b[len(b)-HashLength:]
	}

	copy(h[HashLength-len(b):], b)
}