module wallet-svc

require (
	github.com/BurntSushi/toml v0.3.1
	github.com/alecthomas/log4go v0.0.0-20180109082532-d146e6b86faa
	github.com/boxproject/bolaxy v0.0.0
	github.com/gin-contrib/sse v0.1.0 // indirect
	github.com/gin-gonic/gin v1.3.0
	github.com/go-sql-driver/mysql v1.4.1
	github.com/jmoiron/sqlx v1.2.0
	github.com/json-iterator/go v1.1.9 // indirect
	github.com/mitchellh/mapstructure v1.1.2
	github.com/syndtr/goleveldb v1.0.0
	golang.org/x/crypto v0.0.0-20190308221718-c2843e01d9a2
	google.golang.org/appengine v1.6.5 // indirect
	gopkg.in/go-playground/assert.v1 v1.2.1 // indirect
	gopkg.in/go-playground/validator.v8 v8.18.2 // indirect
)

replace github.com/boxproject/bolaxy v0.0.0 => ../bolaxy

go 1.13
