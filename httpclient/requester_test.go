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
	rster := &Requester{BaseUrl:"http://192.168.10.189:8080"}

	bf,err := rster.PostJson("balanceof",&model.FollowAsset{
		Address:"0x343e7F717d268903F6033766cE1D210a3D82C097",
		Contract:"0xA79D70c4a0B3A31043541Cd593828170Bf037aFE",
	})

	if err != nil {
		fmt.Printf("error=%v\n",err)
		return
	}

	fmt.Printf("resp=%s\n",bf)

}


