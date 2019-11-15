package httpclient

import (
	"fmt"
	"testing"
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