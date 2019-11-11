package config

import (
	"github.com/alecthomas/log4go"
	"testing"
)

func TestInc(t *testing.T) {
	var i uint8 = 0
	for {
		i++
		log4go.Info("i=%d\n", i)
		if i == 0 {
			break
		}
	}
}

func TestLoadConfig(t *testing.T) {

}
