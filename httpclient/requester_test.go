package httpclient

import (
	"fmt"
	"testing"
	"wallet-svc/model"
)

func TestRequester_RequestHttp(t *testing.T) {
	rster := &Requester{BaseUrl:"http://192.168.10.189:8080"}

	bf,err := rster.RequestHttp("GET","/block/0",nil)

	if err != nil {
		fmt.Printf("error=%v\n",err)
		return
	}

	fmt.Printf("resp=%s\n",bf)

}


func TestRequester_RequestJson(t *testing.T) {
	rster := &Requester{BaseUrl:"http://192.168.9.128:8082"}

	bf,err := rster.PostJson("balanceof",&model.FollowAsset{
		Address:"0x2112D0B5EA38BC1A5cAa113dB8393E1BeCaaC6b1",
		Contract:"0xd9f607a27b909858e623b0bd06ec7f46265e7199",
	})

	if err != nil {
		fmt.Printf("error=%v\n",err)
		return
	}

	fmt.Printf("resp=%s\n",bf)

}


