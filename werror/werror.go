package werror

import "github.com/alecthomas/log4go"

type ErrorG int32

const(
	empty ErrorG = iota+10000
	Success
	BindJsonError//json解析错误
	QueryError//查询错误
	RequestChainError//请求公链错误
)

func DeferError(errFunc func() error,logPosition string) {
	err := errFunc()
	if err != nil {
		log4go.Info("%s error=%v\n",logPosition,err)
	}
}

