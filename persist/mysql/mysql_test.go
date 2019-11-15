package mysql

import (
	"github.com/alecthomas/log4go"
	"testing"
	"wallet-svc/config"
	"wallet-svc/mock"
)

func TestMain(m *testing.M) {
	defer log4go.Close()
	//config.LoadConfig(mock.Getwd())
	config.LoadConfig(mock.Getwd())
	InitMySQL()
	m.Run()
}