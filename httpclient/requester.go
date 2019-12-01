package httpclient

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/alecthomas/log4go"
	"io/ioutil"
	"net/http"
	"strings"
)

type Requester struct {
	BaseUrl string
}

func (requester *Requester) RequestHttpByGet(endpoint string, params map[string]string) ([]byte, error) {
	return requester.RequestHttp("GET", endpoint, params)
}

func (requester *Requester) RequestHttpByPost(endpoint string, params map[string]string) ([]byte, error) {
	return requester.RequestHttp("POST", endpoint, params)
}

func (requester *Requester) RequestHttp(method string, endpoint string, params map[string]string) ([]byte, error) {
	transport := &http.Transport{}
	client := &http.Client{
		Transport: transport,
	}

	url := fmt.Sprintf("%s/%s", requester.BaseUrl, endpoint)
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("build request error=%s", err.Error()))
	}
	req.WithContext(context.Background())

	q := req.URL.Query()
	for key, val := range params {
		q.Add(key, val)
	}

	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, errors.New("request alpha ping error")
	}

	return buf, nil
}

func (requester *Requester) PostString(endpoint string, reqstr string) ([]byte, error) {
	resp, err := http.Post(fmt.Sprintf("%s/%s", requester.BaseUrl,endpoint), "text/plain", strings.NewReader(reqstr))
	if err != nil {
		return nil, err
	}

	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, errors.New("request alpha ping error")
	}

	return buf, nil
}

func (requester *Requester) PostJson(endpoint string, request interface{}) ([]byte, error) {
	fullurl := fmt.Sprintf("%s/%s", requester.BaseUrl, endpoint)
	reqbuf, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}
	log4go.Info("postJson请求url=%s,数据：%s\n", fullurl, reqbuf)

	buffer := bytes.NewBuffer(reqbuf)
	req, err := http.NewRequest("POST", fullurl, buffer)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json;charset=UTF-8")
	client := http.Client{}
	respbody, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if respbody.StatusCode != 200 {
		return nil, errors.New("request alpha ping error")
	}

	buf, err := ioutil.ReadAll(respbody.Body)
	if err != nil {
		return nil, err
	}

	reqsult := string(buf)
	log4go.Debug("response string-->%s\n", reqsult)

	return buf, nil

}
