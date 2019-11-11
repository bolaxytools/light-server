package mysql

import (
	"github.com/alecthomas/log4go"
	"testing"
	"wallet-service/config"
	"wallet-service/mock"
)

func TestMain(m *testing.M) {
	defer log4go.Close()
	//config.LoadConfig(mock.Getwd())
	config.LoadConfig(mock.Getwd())
	InitMySQL()
	m.Run()
}