package util

import (
	"fmt"
	"testing"
)

func TestAddrLen(t *testing.T)  {
	fmt.Println(len("0x3bD6361959306B1b50797D3ff82B9A43541c3e47"))
}


func TestHashLen(t *testing.T)  {
	fmt.Println(len("0xaf4a11b7bf6dc0d734de90b04306af3198a197adacd8d2efcd0929d2b7d5b200"))
}

func TestIsNumber(t *testing.T) {
	str := "aaa"
	b := IsNumber(str)
	fmt.Printf("b=%t\n",b)
}

func TestCutAddress(t *testing.T) {
	la := "0x0000000000000000000000003bd6361959306b1b50797d3ff82b9a43541c3e47"
	sa := CutAddress(la)
	fmt.Printf("addr=%s\n",sa)
	fmt.Printf("addr=%s\n","0x3bd6361959306b1b50797d3ff82b9a43541c3e47")
}

func TestDataToValue(t *testing.T) {
	data := "0x0000000000000000000000000000000000000000000000000000000000000457"
	v := DataToValue(data)
	fmt.Printf("%s\n",v)
}