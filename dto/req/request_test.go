package req

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestJson(t *testing.T) {

	mp := make(map[string]interface{})
	mp["code"]="jsksl"

	req := ReqData{
		Data:mp,
		Sign:"ojbk",
	}

	jf,er := json.Marshal(req)
	if er != nil {
		t.Error(er)
		return
	}

	fmt.Printf("json=%s\n",jf)


	rev := new(ReqData)

	er = json.Unmarshal(jf,rev)
	if er != nil {
		t.Error(er)
		return
	}
	lg := new(Login)
	rev.Reverse(lg)
	fmt.Printf("lg.code=%s\n",lg.Code)

}
