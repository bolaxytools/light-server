package mysql

import (
	"fmt"
	"github.com/alecthomas/log4go"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"wallet-svc/config"
)

var db *sqlx.DB

func InitMySQL() {
	var err error
	msq := config.Cfg.MySQL
	dsn := fmt.Sprintf("%s:%s@%s(%s:%s)/%s?parseTime=true&charset=utf8", msq.User, msq.Pasw, msq.Prot, msq.Host, msq.Port, msq.Dbnm)

	db, err = sqlx.Open("mysql", dsn)

	if err != nil {
		panic(err)
	}
	//defer gerror.DeferError(db.Close,"db.Close")

	err = db.Ping()

	if err != nil {
		panic(err)
	}
	log4go.Info("mysql init success")

}

func GetDb() *sqlx.DB {
	return db
}
